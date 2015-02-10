package websocket

import (
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"code.google.com/p/go-uuid/uuid"
	"github.com/gorilla/websocket"
)

type filter struct {
	prefixRegex, dataRegex string
}

type Conn struct {
	*websocket.Conn
	mutex   *sync.Mutex
	filters []filter
	id      string
}

var websocketTable = struct {
	sync.RWMutex
	conns map[string]*Conn
}{sync.RWMutex{}, map[string]*Conn{}}

func addWebsock(conn *Conn) {
	websocketTable.Lock()
	defer websocketTable.Unlock()
	if _, ok := websocketTable.conns[conn.id]; ok {
		log.Printf("Tried to add websock to websocketTable, key was not empty %s", conn.id)
	} else {
		websocketTable.conns[conn.id] = conn
	}
}

func dropWebsock(conn *Conn) {
	websocketTable.Lock()
	defer websocketTable.Unlock()
	if _, ok := websocketTable.conns[conn.id]; ok {
		delete(websocketTable.conns, conn.id)
	} else {
		log.Printf("Tried to drop websock from websocketTable, key was empty %s", conn.id)

	}
}

type websocketBody struct {
	Event string        `json:"event"`
	Args  []interface{} `json:"args"`
}

func StartServer(addr string) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	http.HandleFunc("/websocket", func(w http.ResponseWriter, req *http.Request) {
		orig, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println(err)
			return
		}
		conn := &Conn{Conn: orig, mutex: &sync.Mutex{}, id: uuid.New()}
		addWebsock(conn)
		conn.readLoop()
	})

	http.Handle("/", http.FileServer(http.Dir("./pipelisten")))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func (conn *Conn) readLoop() {
	for {
		input := &websocketBody{}
		err := conn.ReadJSON(input)
		if err != nil {
			_, isOpErr := err.(*net.OpError)
			knownErr := strings.Contains(err.Error(), "EOF") || isOpErr
			if !knownErr {
				log.Printf("Websocket disconnected for an unknown reason: %#v", err)
			}
			dropWebsock(conn)
			conn.Close()
			return
		}
		go conn.handle(input)
	}
}
