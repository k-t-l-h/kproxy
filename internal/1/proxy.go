package first

import (
	"fmt"
	"github.com/oriser/regroup"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

var wh = regexp.MustCompile("\r\n")
var cp = regexp.MustCompile("[Pp]roxy-[Cc]onnection:")
var tp = regexp.MustCompile("(?:GET|POST|HEAD|OPTIONS).*")
var hp0 = regexp.MustCompile("(?P<name>Host):(?P<host>.*)(?P<port>:(\\d{1,3}))?")
var hp = regroup.MustCompile("(?P<name>Host):(?P<host>.*)(?P<port>:(\\d{1,3}))?")

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

	message := string(full)
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
