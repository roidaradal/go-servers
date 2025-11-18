package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/go-servers/sse"
	"github.com/roidaradal/go-servers/tcp"
	"github.com/roidaradal/go-servers/websocket"
)

const (
	host string = "localhost"
	port int    = 6969
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		displayUsage()
		return
	}

	key := args[0]
	handler, ok := handlers[key]
	if !ok {
		displayUsage()
		return
	}
	handler()
}

var handlers = map[string]func(){
	"tcp-server": func() {
		tcp.NewServer(host, port).Run()
	},
	"tcp-client": func() {
		tcp.RunClient(host, port)
	},
	"sse-server": func() {
		sse.RunServer(host, port)
	},
	"ws-server": func() {
		websocket.RunServer(host, port)
	},
}

func displayUsage() {
	fmt.Println("Usage: go run . <option>")
	options := dict.Keys(handlers)
	slices.Sort(options)
	fmt.Println("Options:", strings.Join(options, ", "))
}
