package pipe

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"

	"github.com/Sean-Der/pipeshout/websocket"
)

func handleConn(conn net.Conn) {

	reader := bufio.NewReader(conn)
	for {
		startTime := time.Now()
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		separator := strings.Index(line, " ")
		if separator == -1 {
			log.Printf("Line has no separator: %s", line)
			continue
		}
		websocket.EmitAddLine(startTime, line[0:separator], line[separator+1:])
	}
}

func StartHandler(pipePath string) {
	l, err := net.Listen("unix", pipePath)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
