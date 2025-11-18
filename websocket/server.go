package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunServer(host string, port int) {
	address := fmt.Sprintf("%s:%d", host, port)
	http.HandleFunc("/post", handleInput)

	fmt.Println("Running WebSocket server at", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func handleInput(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			continue
		}
		fmt.Printf("Received: %s\n", string(message))

		if err := conn.WriteMessage(websocket.TextMessage, []byte("OK")); err != nil {
			fmt.Println("Error writing message:", err)
			continue
		}
	}
}
