package websocket

import (
	"fmt"

	"github.com/Sean-Der/pipeshout/pipe"
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

func (conn *Conn) EmitSetLines() {
	lines := pipe.GetLineCache()
	goodLines := []pipe.Line{}

	conn.mutex.RLock()
	for _, line := range lines {
		if lineMatchRegexes(line.Prefix, line.Line, conn.regexes) {
			goodLines = append(goodLines, line)
		}
	}
	conn.mutex.RUnlock()

	conn.emit(newWebsocketBody("setLines", []interface{}{goodLines}))
}

func (conn *Conn) emit(body *websocketBody) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.WriteJSON(body)
}

func addLineEmitter() {
	for {
		line := <-pipe.LinesChan
		body := newWebsocketBody("addLine", []interface{}{line})

		websocketTable.RLock()
		for _, conn := range websocketTable.conns {
			conn.mutex.RLock()
			if !lineMatchRegexes(line.Prefix, line.Line, conn.regexes) {
				continue
			}
			conn.mutex.RUnlock()
			conn.emit(body)
		}
		websocketTable.RUnlock()
	}

}
