package model

import (
	"log"

	"github.com/gorilla/websocket"
)

type Connections map[*websocket.Conn]bool

func (conns Connections) BroadcastJSON(messageType string, jsonData any) {
	for client, _ := range conns {
		if err := client.WriteJSON(map[string]any{"messageType": messageType, "data": jsonData}); err != nil {
			log.Default().Printf("ERROR : %s", err.Error())
			client.Close()
			delete(conns, client)
		}
	}
}
