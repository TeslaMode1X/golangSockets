package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Global list of all active connections and a mutex for synchronization
var connections = make([]*websocket.Conn, 0)
var mu sync.Mutex

// Function to handle WebSocket connections
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка при получении токена", err)
		return
	}
	defer conn.Close()

	// Add the new connection to the global list
	mu.Lock()
	connections = append(connections, conn)
	mu.Unlock()

	for {
		// Read message from the WebSocket
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Ошибка чтения сообщения:", err)
			return
		}
		fmt.Println("Сообщение от клиента", string(msg))

		// Broadcast the message to all connected clients
		mu.Lock()
		for _, c := range connections {
			err = c.WriteMessage(messageType, []byte(fmt.Sprintf("Сервер получил: %s", msg)))
			if err != nil {
				log.Println("Ошибка отправки сообщения:", err)
				c.Close() // Close the connection if there's an error
			}
		}
		mu.Unlock()
	}
}

// Set up HTTP route for WebSocket
func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("WebSocket сервер запущен на :8080")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
