package main

import (
	"flag"

	"github.com/Sean-Der/pipeshout/pipe"
	"github.com/Sean-Der/pipeshout/websocket"
)

var addr = flag.String("addr", ":8080", "http addr")
var pipePath = flag.String("pipe", "./pipeshout.pipe", "Path to pipeshout's pipe")

func main() {
	flag.Parse()

	go websocket.StartServer(*addr)
	pipe.StartHandler(*pipePath)
}
