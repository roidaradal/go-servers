package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func RunClient(host string, port int) {
	address := fmt.Sprintf("%s:%d", host, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatal("ResolveTCPAddr failed:", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("DialTCP failed:", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		cleanLine := strings.TrimSpace(line)
		if cleanLine == "exit" || cleanLine == "" {
			fmt.Println("exiting...")
			break
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatal("Server write failed:", err)
		}

		reply := make([]byte, 8)
		_, err = conn.Read(reply)
		if err != nil {
			log.Fatal("Server reply failed:", err)
		}

		fmt.Println("TCP:", string(reply))
	}
}
