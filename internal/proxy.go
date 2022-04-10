package internal

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/oriser/regroup"
)

const (
	HTTPS_START_MESSAGE = "HTTP/1.0 200 Connection established\\r\\nProxy-agent: Kproxy\\r\\n\\r\\n"
)

var (
	wh     = regexp.MustCompile("\r\n")
	cp     = regexp.MustCompile("[Pp]roxy-[Cc]onnection:")
	tp     = regexp.MustCompile("(?:GET|POST|HEAD|OPTIONS).*")
	secure = regexp.MustCompile("CONNECT")
	hp0    = regexp.MustCompile("(?P<name>Host):(?P<host>.*):(?P<port>:(\\d{1,3}))?")
	hp     = regroup.MustCompile("(?P<name>Host):(?P<host>.*)(?P<port>:(\\d{1,3}))?")
)

type Proxy struct {
	Pool *pgxpool.Pool
	//tls later?
}

//main runner
func (p Proxy) Run() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept() // принимаем TCP-соединение от клиента и создаем новый сокет
		if err != nil {
			continue
		}
		go p.execConnection(conn) // обрабатываем запросы клиента в отдельной го-рутине
	}
}

func (p Proxy) execConnection(conn net.Conn) {
	defer conn.Close()

	//read full message
	full, _ := p.readMessage(conn)

	//here we need to check if this is secure connection
	r := bytes.NewReader(full)
	reader := bufio.NewReader(r)
	req, _ := http.ReadRequest(reader)
	log.Println(req)

	if req.Method == http.MethodConnect {
		p.parseSecureConnection(conn, req.Host, string(full))
	} else {
		p.parseConnection(conn, req, string(full))
	}

}

func (p Proxy) readMessage(conn net.Conn) ([]byte, error) {

	full := make([]byte, 0, 0)
	b := make([]byte, 100)

	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 1))

	for {
		n, err := conn.Read(b)
		if n == 0 {
			break
		}
		if err != nil {
			return nil, err
		}
		full = append(full, b[:n]...)
	}
	full = append(full, []byte("\r\n\r\n")...)
	return full, nil
}

func writeMessage(conn net.Conn, msg []byte) {
	size := len(msg)
	sentBytes := 0
	for sentBytes < size {
		n, err := conn.Write(msg)
		if err != nil {
			fmt.Println("some error in sending data", n)
		}
		sentBytes += n
	}
	//answer := GetRequest(conn)
	//return answer.FullMsg
}

func (p *Proxy) parseSecureConnection(conn net.Conn, host, message string) error {

	//answer immediately
	writeMessage(conn, []byte(HTTPS_START_MESSAGE))
	//get current path
	dir, err := filepath.Abs("")
	if err != nil {
		return err
	}
	//generate certificates
	err = exec.Command(dir+"/gen_cert.sh", "mail.ru", strconv.Itoa(rand.Int())).Run()
	if err != nil {
		return err
	}

	cert, err := tls.LoadX509KeyPair(dir+"/mitm.crt", dir+"/cert.key")
	if err != nil {
		panic(err)
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	serv := tls.Server(conn, cfg)

	//read message
	bt, err := p.readMessage(serv)

	tlsconnTo, err := tls.Dial("tcp", host, cfg)
	log.Println(tlsconnTo, err)

	SendMessage(conn, []byte(message))
	log.Println(string(bt), err)

	answer, err := p.readMessage(tlsconnTo)
	if err != nil {
		panic(err)
	}

	log.Println(answer)
	if strings.LastIndex(string(answer), "Transfer-Encoding: chunked") != -1 {
		answer = bytes.Replace(answer, []byte("Transfer-Encoding: chunked"), []byte(""), -1)
	}

	defer serv.Close()
	defer conn.Close()

	return nil
}


func SendMessage(conn net.Conn, msg []byte) error {
	bytesSent := 0
	for bytesSent < len(msg) {
		n, err := conn.Write(msg)
		if err != nil {
			return err
		}
		bytesSent += n
	}
	return nil
}


func (p *Proxy) parseConnection(conn net.Conn, req *http.Request, message string) error {

	if !strings.Contains(req.Host, ":"){
		req.Host += ":80"
	}
	dial, err := net.Dial("tcp", req.Host)
	if err != nil {
		return err
	}
	defer dial.Close()
	dial.SetReadDeadline(time.Now().Add(time.Second))
	writeMessage(dial, []byte(message))

	full, err := p.readMessage(dial)
	if err != nil {
		return err
	}
	writeMessage(conn, full)
	p.writeRequest(message)
	p.getRequest(2)
	return nil
}