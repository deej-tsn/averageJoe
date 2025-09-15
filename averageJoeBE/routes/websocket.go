package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/deej-tsn/averageJoe/model"
	"github.com/deej-tsn/averageJoe/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
)

type websocketResponse struct {
	MessageType string                 `json:"messageType"`
	Data        map[string]interface{} `json:"data"`
}

type roundRespData struct {
	playerID string
	roundID  string
	choice   string
}

type roundInfoData struct {
	RoundID   string       `json:"roundID"`
	RoundData *model.Round `json:"roundData"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// POST("connect-to-game")
func (gmc *GameMgrController) WS_handler(c echo.Context) error {
	req := c.Request()
	res := c.Response()

	// Extract optional protocol
	protocolHeader := req.Header.Get("Sec-WebSocket-Protocol")
	var selectedProtocol string

	if protocolHeader != "" {
		// Example format: "auth-token.abc123"
		parts := strings.Split(protocolHeader, "-")
		if len(parts) >= 2 {
			token := parts[1]
			log.Println("Extracted token:", token)
			// Prepare to echo the protocol
			selectedProtocol = protocolHeader
		}
	}

	// If a protocol was sent, echo it back in the handshake
	respHeaders := http.Header{}
	if selectedProtocol != "" {
		upgrader.Subprotocols = []string{selectedProtocol}
		respHeaders.Set("Sec-WebSocket-Protocol", selectedProtocol)
	}

	gameID := c.QueryParam("gameID")
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(*util.JWT_CustomClaim)
	log.Printf("userID : %s", claim.Name)
	if gameID == "" {
		c.Echo().StdLogger.Printf("ERROR: could not bind request body to controller\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "could not read request"})
	}

	conn, err := upgrader.Upgrade(res, req, respHeaders)
	if err != nil {
		return err
	}

	// Save connection to game
	game, ok := gmc.gm.Games[gameID]
	if !ok {
		conn.WriteJSON(util.ErrorMessage("Game not found"))
		conn.Close()
		return nil
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
			data, err := parseMsg(msg)
			if err != nil {
				log.Print("invalid format of message")
				log.Print(string(msg))
				errMessage := []byte("invalid message sent")
				conn.WriteMessage(websocket.TextMessage, errMessage)
				continue
			}

			gmc.messageTypeToDataStruct(gameID, data)
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

func (gmc *GameMgrController) broadcastJSON(gameID string, messageType string, jsonData any) {
	for client, _ := range gmc.gm.Games[gameID].Connections {
		if err := client.WriteJSON(map[string]any{"messageType": messageType, "data": jsonData}); err != nil {
			client.Close()
			delete(gmc.gm.Games[gameID].Connections, client)
		}
	}
}

func (gmc *GameMgrController) sendRound(gameID string, roundID string) {
	round := gmc.data.GetRandomRound()

	bc := &roundInfoData{
		RoundID:   roundID,
		RoundData: round,
	}
	gmc.gm.Games[gameID].CurrentRound = round
	gmc.broadcastJSON(gameID, "ROUND", bc)
}

func (gmc *GameMgrController) voteInRound(gameID string, roundID string) {
	gmc.gm.Games[gameID].CurrentRound.Votes
}

func parseMsg(msg []byte) (*websocketResponse, error) {
	var jsonMessage websocketResponse
	if err := json.Unmarshal(msg, &jsonMessage); err != nil {
		return nil, err
	}
	log.Printf("%v", jsonMessage)
	if jsonMessage.MessageType == "" {
		return nil, errors.New("Invalid JSON format, (missing messageType)")
	}
	return &jsonMessage, nil
}

func (gmc *GameMgrController) messageTypeToDataStruct(gameID string, wsR *websocketResponse) error {
	switch wsR.MessageType {
	case "START":
		gmc.sendRound(gameID, "1")
	case "VOTE":
		if err := mapstructure.Decode(wsR.Data, roundRespData{}); err != nil {
			return errors.New("Incorrect Vote format")
		}
	default:
		return errors.New("Unknown messageType found")
	}
	return nil
}

/*
	{
		messageType : AUTH
		playerToken : xyz123
	}
*/
func (gmc *GameMgrController) authenicatePlayer() {
	log.Println("Authorising user")
}

/*
	{
		messageType : CHOICE
		playerToken : xyz123,
		round : a,
		option : 1
	}
*/
func (gmc *GameMgrController) readOptionChoose(gameID string, msg []byte) {

}

func (gmc *GameMgrController) broadcastRound(gameID string) {

}
