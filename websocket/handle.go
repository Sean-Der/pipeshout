package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"regexp"
)

var events = map[string]interface{}{
	"setRegexes": func(conn *Conn, regexes []regex) {
		for i, regex := range regexes {
			regexes[i].prefixCompiled, _ = regexp.Compile(regex.PrefixRegex)
			regexes[i].lineCompiled, _ = regexp.Compile(regex.LineRegex)
		}
		conn.mutex.Lock()
		conn.regexes = regexes
		conn.mutex.Unlock()
	},
}

func (conn *Conn) handle(body *websocketBody) {
	err, funcType, funcValue := eventFunc(body.Event)
	if err != nil {
		log.Printf("Error getting func for event: %s err:%v", body.Event, err.Error())
		return
	}
	funcArgs := []reflect.Value{reflect.ValueOf(conn)}
	if funcType.NumIn() > 1 {
		serializedJSON := []interface{}{}
		for i := 1; i < funcType.NumIn(); i++ {
			serializedJSON = append(serializedJSON, reflect.New(funcType.In(i)).Interface())
		}
		if err := json.Unmarshal(body.Args, &serializedJSON); err != nil {
			log.Println(err)
			return
		}
		for _, value := range serializedJSON {
			val := reflect.ValueOf(value)
			if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			funcArgs = append(funcArgs, val)
		}
	}
	funcValue.Call(funcArgs)

}
func eventFunc(eventName string) (err error, funcType reflect.Type, funcValue reflect.Value) {
	var (
		event interface{}
		found bool
	)
	if event, found = events[eventName]; !found {
		err = errors.New("ERROR: No func for event: " + eventName)
		return
	}

	funcType = reflect.TypeOf(event)
	funcValue = reflect.ValueOf(event)
	if funcType.Kind() != reflect.Func {
		err = errors.New("ERROR: Nonfunc entry for event: " + eventName)
		return
	}
	return
}
