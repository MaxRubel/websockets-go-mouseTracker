package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	X      uint16
	Y      uint16
	Player uint8
}

type MouseMessage struct {
	Player uint8
}

type ChoosePlayerMessage struct {
	PlayerChoice uint8
}

type playerHasBeenChosen struct {
	MessageType uint8
	PlayerNo    uint8
}

type GameData struct {
	Player1Score  int
	Player2Score  int
	Player1       bool
	Player2       bool
	GameIsRunning bool
}

type MousePositions struct {
	P1x int
	P1y int
	P2x int
	P2y int
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var gameData GameData
var mousePos MousePositions
var posF string
var choosePlayerMessage ChoosePlayerMessage

func choosePlayer(no uint8) playerHasBeenChosen {
	var m playerHasBeenChosen
	if no == 1 {
		if !gameData.Player1 {
			gameData.Player1 = true
			m.MessageType, m.PlayerNo = 2, 1
		}
	}

	if no == 2 {
		if !gameData.Player2 {
			gameData.Player2 = true
			m.MessageType, m.PlayerNo = 2, 2

		}
	}
	return m
}

func mouseClicked(m MouseMessage) {
	fmt.Println("this player scored:", m.Player)
	switch m.Player {
	case 1:
		gameData.Player1Score++
	case 2:
		gameData.Player2Score++
	default:
		fmt.Println("Unable to process mouse click")
	}
	fmt.Println("player1:", gameData.Player1Score)
	fmt.Println("player2:", gameData.Player2Score)
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

		messageType := incoming[0]

		switch messageType {
		case 0:
			var message Message
			err = binary.Read(bytes.NewReader(incoming[1:]), binary.LittleEndian, &message)
			if err != nil {
				fmt.Println("Error decoding mouse message:", err)
				continue
			}

			if message.Player == 1 {
				mousePos.P1x = int(message.X)
				mousePos.P2y = int(message.Y)
			}
			if message.Player == 2 {
				mousePos.P2x = int(message.X)
				mousePos.P2y = int(message.Y)
			}

		case 1:
			var mouseMessage MouseMessage
			err = binary.Read(bytes.NewReader(incoming[1:]), binary.LittleEndian, &mouseMessage)
			if err != nil {
				fmt.Println("Error decoding mouse message:", err)
				continue
			}

			mouseClicked(mouseMessage)

			il := strconv.Itoa(gameData.Player1Score)
			i2 := strconv.Itoa(gameData.Player2Score)
			is := "n" + "/" + il + "/" + i2

			err = conn.WriteMessage(websocket.TextMessage, []byte(is))
			if err != nil {
				log.Println("Write:", err)
				return
			}

		case 2:
			err = binary.Read(bytes.NewReader(incoming[1:]), binary.LittleEndian, &choosePlayerMessage)
			if err != nil {
				fmt.Println("Error decoding mouse message:", err)
				continue
			}
			p := choosePlayer(choosePlayerMessage.PlayerChoice)
			pF := fmt.Sprintf("%d%d", p.MessageType, p.PlayerNo)
			if p.MessageType == 2 {
				err = conn.WriteMessage(websocket.TextMessage, []byte(pF))
			}

		default:
			fmt.Println("Unknown message type:", messageType)
		}
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		select {
		case <-ticker.C:
			posF = fmt.Sprintf("m/%d/%d/%d/%d", mousePos.P1x, mousePos.P1y, mousePos.P2x, mousePos.P2y)
			conn.WriteMessage(websocket.TextMessage, []byte(posF))
		}
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
