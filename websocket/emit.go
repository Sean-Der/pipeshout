package websocket

import "time"

func EmitAddLine(start time.Time, prefix, line string) {
	body := &websocketBody{Event: "addLine", Args: []interface{}{start, prefix, line}}

	websocketTable.RLock()
	defer websocketTable.RUnlock()
	for _, conn := range websocketTable.conns {
		conn.emit(body)
	}

}

func (conn *Conn) emit(body *websocketBody) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.WriteJSON(body)
}
