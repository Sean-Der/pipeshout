package websocket

import (
	"fmt"
	"time"
)

func lineMatchRegexes(prefix, line string, regexes []regex) bool {
	if len(regexes) == 0 {
		return true
	}

	for _, regex := range regexes {
		if regex.prefixCompiled == nil || regex.lineCompiled == nil {
			fmt.Println("nil regex")
			continue
		}

		if regex.prefixCompiled.MatchString(prefix) && regex.lineCompiled.MatchString(line) {
			return true
		}
	}

	return false
}

func EmitAddLine(start time.Time, prefix, line string) {
	body := newWebsocketBody("addLine", []interface{}{start, prefix, line})

	websocketTable.RLock()
	defer websocketTable.RUnlock()
	for _, conn := range websocketTable.conns {
		conn.mutex.RLock()
		if !lineMatchRegexes(prefix, line, conn.regexes) {
			continue
		}
		conn.mutex.RUnlock()
		conn.emit(body)
	}

}

func (conn *Conn) emit(body *websocketBody) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.WriteJSON(body)
}
