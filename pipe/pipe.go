package pipe

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Line struct {
	Start  time.Time `json:"start"`
	Prefix string    `json:"prefix"`
	Line   string    `json:"line"`
	Id     string    `json:"id"`
}

func handleConn(conn net.Conn) {

	reader := bufio.NewReader(conn)
	for {
		startTime := time.Now()
		rawLine, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		separator := strings.Index(rawLine, " ")
		if separator == -1 {
			log.Printf("Line has no separator: %s", rawLine)
			continue
		}
		addCacheLine(Line{Start: startTime, Prefix: rawLine[0:separator], Line: rawLine[separator+1:], Id: uuid.New()})
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
