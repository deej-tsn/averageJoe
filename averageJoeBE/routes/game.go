package routes

import (
	"net/http"

	"github.com/deej-tsn/averageJoe/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ConnectRequest struct {
	GameID   string `json:"gameID"`
	PlayerID string `json:"playerID"`
}

type CreateGameRequest struct {
	PlayerGameCode string `json:"playerGameCode"`
}

type GameMgrController struct {
	gm   *model.GameMgr
	data *model.Data
}

func NewGameMgrController(gm *model.GameMgr, data *model.Data) *GameMgrController {
	return &GameMgrController{
		gm:   gm,
		data: data,
	}
}

// GET("/active-games")
func (gmc *GameMgrController) GET_activeGames(c echo.Context) error {
	games := gmc.gm.ListGames()
	return c.JSON(http.StatusAccepted, games)
}

// POST("/connect-to-game")
func (gmc *GameMgrController) POST_connectToGame(c echo.Context) error {
	var body ConnectRequest
	if err := c.Bind(&body); err != nil {
		c.Echo().StdLogger.Printf("ERROR: could not bind request body to controller\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "could not read request"})
	}

	if _, err := gmc.gm.JoinGame(body.GameID, body.PlayerID); err != nil {
		c.Echo().StdLogger.Printf("ERROR: %s\n", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.String(http.StatusAccepted, body.GameID)
}

// POST("/create-game")
func (gmc *GameMgrController) POST_createGame(c echo.Context) error {
	var body CreateGameRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "could not read request"})
	}

	game, err := gmc.gm.NewGame(body.PlayerGameCode, gmc.data.GetRandomRound())
	if err != nil {
		c.Echo().StdLogger.Printf("ERROR: %s\n", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.String(http.StatusAccepted, game.GameID)
}

func GET_newPlayerUUID(c echo.Context) error {
	playerUUID := uuid.NewString()

	return c.JSON(http.StatusAccepted, map[string]string{"uuid": playerUUID})
}
