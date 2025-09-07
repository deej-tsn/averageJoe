package routes

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// POST("connect-to-game")
func (gmc *GameMgrController) WS_handler(c echo.Context) error {

	gameID := c.QueryParam("gameID")
	if gameID == "" {
		c.Echo().StdLogger.Printf("ERROR: could not bind request body to controller\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "could not read request"})
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	// Save connection to game
	game, ok := gmc.gm.Games[gameID]
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "game not found"})
	}
	game.Connections[conn] = true

	// Example: listen for messages
	go func() {
		defer conn.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				delete(game.Connections, conn)
				return
			}
			// Handle incoming messages (answers, ready, etc.)
			println("Message:", string(msg))
			gmc.broadcast(gameID, websocket.TextMessage, msg)
		}
	}()

	return nil
}

func (gmc *GameMgrController) broadcast(gameID string, messageType int, message []byte) {
	for client, _ := range gmc.gm.Games[gameID].Connections {
		if err := client.WriteMessage(messageType, message); err != nil {
			client.Close()
			delete(gmc.gm.Games[gameID].Connections, client)
		}
	}
}
