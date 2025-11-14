package main

import (
	"fmt"
	"os"

	"github.com/roidaradal/go-servers/tcp"
)

const (
	usage string = "Usage: go run . <server|client>"
	host  string = "localhost"
	port  int    = 6969
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println(usage)
		return
	}

	switch args[0] {
	case "server":
		runTCPServer()
	case "client":
		runTCPClient()
	default:
		fmt.Println(usage)
		return
	}
}

func runTCPServer() {
	server := tcp.NewServer(host, port)
	server.Run()
}

func runTCPClient() {
	tcp.Client(host, port)
}
