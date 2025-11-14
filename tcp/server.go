package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Server struct {
	host string
	port int
}

type Conn struct {
	id   int
	conn net.Conn
}

func NewServer(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

func (server *Server) Run() {
	address := fmt.Sprintf("%s:%d", server.host, server.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("Running TCP server at", address)

	id := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		id += 1
		client := &Conn{conn: conn, id: id}
		go client.handleRequest()
	}
}

func (client *Conn) handleRequest() {
	fmt.Printf("Client %d connected\n", client.id)
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			return
		}
		fmt.Printf("Received from %d: %s\n", client.id, message)
		client.conn.Write([]byte("OK\n"))
	}
}
