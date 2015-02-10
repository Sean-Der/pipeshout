package websocket

import (
	"errors"
	"log"
	"reflect"
)

var events = map[string]interface{}{
	"setFilters": func(conn *Conn, filters []filter) {

	},
}

func (conn *Conn) handle(body *websocketBody) {
	err, funcType, funcValue := eventFunc(body.Event)
	if err != nil {
		log.Printf("Error getting func for event: %s err:%v", body.Event, err.Error())
		return
	}
	if (funcType.NumIn() - 1) != len(body.Args) {
		log.Printf("Arg mismatch for event %s: args: %v", body.Event, body.Args)
		return
	}
	funcArgs := []reflect.Value{reflect.ValueOf(conn)}
	if funcType.NumIn() > 1 {
		for _, value := range body.Args {
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
