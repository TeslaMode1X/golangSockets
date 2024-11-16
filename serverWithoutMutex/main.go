package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка при получении токена", err)
		return
	}
	defer conn.Close()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Ошибка чтения сообщения:", err)
			return
		}
		fmt.Println("Сообщение от клиента", string(msg))

		err = conn.WriteMessage(messageType, []byte(fmt.Sprintf("Сервер получил: %s", msg)))
		if err != nil {
			log.Println("Ошибка отправки сообщения:", err)
			break
		}
	}
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("WebSocket сервер запущен на :8080")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
