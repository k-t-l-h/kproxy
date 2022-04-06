package first

import (
	"crypto/tls"
	"fmt"
	"github.com/oriser/regroup"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	HTTPS_START_MESSAGE = "HTTP/1.0 200 Connection established\\r\\nProxy-agent: Kproxy\\r\\n\\r\\n"
)

var (
	wh = regexp.MustCompile("\r\n")
	cp = regexp.MustCompile("[Pp]roxy-[Cc]onnection:")
	tp = regexp.MustCompile("(?:GET|POST|HEAD|OPTIONS).*")
	secure = regexp.MustCompile("CONNECT")
	hp0 = regexp.MustCompile("(?P<name>Host):(?P<host>.*)(?P<port>:(\\d{1,3}))?")
	hp = regroup.MustCompile("(?P<name>Host):(?P<host>.*)(?P<port>:(\\d{1,3}))?")
)

type Proxy struct {
	//tls later?
}

//main runner
func (p Proxy) Run() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println(err)
	}

	for {
		conn, err := listener.Accept() // принимаем TCP-соединение от клиента и создаем новый сокет
		if err != nil {
			continue
		}
		go parseConnection(conn) // обрабатываем запросы клиента в отдельной го-рутине
	}
}

func parseConnection(conn net.Conn) {
	defer conn.Close()

	//data for connection
	var (
		host string
		port string
	)

	//new message
	full := make([]byte, 0, 0)
	b := make([]byte, 100)

	for {
		n, err := conn.Read(b)
		if err != nil {
			break
		}
		if n < len(b) {
			full = append(full, b[:n]...)
			break
		}
		full = append(full, b...)
	}


	//here we need to check if this is secure connection

	message := string(full)

	if isSecure(message) {

	}

	message = strings.TrimSpace(message)
	messages := wh.Split(message, -1)
	cmessages := make([]string, 0, 0)

	//
	//log.Printf("%+v",messages)
	for i := 0; i < len(messages); i++ {

		//parse host+port
		if hp0.MatchString(messages[i]) {
			groups, _ := hp.Groups(messages[i])
			host = groups["host"]
			host = strings.TrimSpace(host)
			port = groups["port"]
			if port == "" {
				port = ":80"
			}

			port = strings.TrimSpace(port)
			log.Print(host, port)
		}
		if tp.MatchString(messages[i]) {
			//cmessage := strings.Replace(messages[i], "http", "", -1)
			//cmessages = append(cmessages, cmessage)
			//continue
		}

		//delete proxy-connection
		if cp.MatchString(messages[i]) {
			continue
		}
		cmessages = append(cmessages, messages[i])
	}

	dial, err := net.Dial("tcp", host+port)
	if err != nil {
		return
	}
	defer dial.Close()
	dial.SetReadDeadline(time.Now().Add(time.Second))
	writeRequest(dial, []byte(strings.Join(cmessages, "\r\n")))

	d := make([]byte, 10)
	full = nil
	for {
		n, err := dial.Read(d)
		if err != nil || n == 0 {
			break
		}

		full = append(full, d[:n]...)
	}

	log.Println(full)

	writeRequest(conn, full)
}

func writeRequest(conn net.Conn, msg []byte) {
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


func parseSecureConnection(conn net.Conn) error {

	//answer immediately
	writeRequest(conn, []byte(HTTPS_START_MESSAGE))
	//get current path
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	//generate certificates
	err = exec.Command(dir+"/gen_cert.sh", "", strconv.Itoa(rand.Int())).Run()
	if err != nil {
		return err
	}

	cert, err := tls.LoadX509KeyPair(dir+"/mitm.crt", dir+"/cert.key")
	if err != nil {
		panic(err)
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	serv := tls.Server(conn, cfg)
	defer serv.Close()
	defer conn.Close()

	return nil
}

func isSecure(message string) bool {
	return false
}
