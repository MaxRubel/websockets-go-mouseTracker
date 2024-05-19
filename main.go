package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	X      uint16
	Y      uint16
	UserID uint32
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, incoming, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Println(incoming)
		buf := bytes.NewReader(incoming)

		var message Message

		err = binary.Read(buf, binary.LittleEndian, &message)
		if err != nil {
			fmt.Println("Error decoding:", err)
			continue
		}

		fmt.Println("userId", message.UserID)
		fmt.Println("X:", message.X)
		fmt.Println("Y:", message.Y)
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	log.Println("Server starting on port 8080...")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
