package model

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/deej-tsn/averageJoe/util"
	"github.com/gorilla/websocket"
)

type GameMgr struct {
	Games map[string]*Game
}

type Game struct {
	GameID         string
	Players        map[string]*PlayerGameRecord
	CurrentRound   *Round
	PreviousRounds []*Round
	State          util.GameState
	Connections    map[*websocket.Conn]bool
}

type Player struct {
	PlayerID string
}

type PlayerGameRecord struct {
	GameID   string
	PlayerID string
	Score    int
	Live     bool
	Answers  []string
}

type Round struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Votes    []int
	State    util.RoundState
}

type Data []Round

func NewGM() *GameMgr {
	return &GameMgr{
		Games: make(map[string]*Game),
	}
}

func (gm *GameMgr) NewGameFromCode(gameCode string, round *Round) *Game {
	game := &Game{
		GameID:       gameCode,
		Players:      make(map[string]*PlayerGameRecord),
		CurrentRound: round,
		State:        util.Lobby,
		Connections:  make(map[*websocket.Conn]bool),
	}

	gm.Games[game.GameID] = game
	return game
}

func (gm *GameMgr) generateValidRoomCode() string {
	var gameCode string
	for {
		gameCode = util.GenerateRoomCode(6)
		if _, exists := gm.Games[gameCode]; !exists {
			break
		}
	}
	return gameCode
}

func (gm *GameMgr) NewGame(playerGameCode string, round *Round) (*Game, error) {
	// check player given code already exists
	if _, exist := gm.Games[playerGameCode]; exist {
		return nil, fmt.Errorf("gameID already in use")
	}

	var game *Game
	if playerGameCode != "" {
		game = gm.NewGameFromCode(playerGameCode, round)
	} else {
		game = gm.NewGameFromCode(gm.generateValidRoomCode(), round)
	}

	return game, nil
}

func (gm *GameMgr) JoinGame(gameID string, playerID string) (*Game, error) {
	game, exists := gm.Games[gameID]
	if !exists {
		return nil, fmt.Errorf("game with id : %s not found", gameID)
	}
	if game.State != util.Lobby {
		return nil, fmt.Errorf("game already started")
	}

	if _, ok := game.Players[playerID]; !ok {
		game.Players[playerID] = &PlayerGameRecord{
			PlayerID: playerID,
			GameID:   gameID,
			Score:    0,
			Live:     true,
			Answers:  make([]string, 0),
		}
	}

	return game, nil
}

func (gm *GameMgr) StartGame(gameID string) error {
	game, exists := gm.Games[gameID]
	if !exists {
		return fmt.Errorf("game not found")
	}
	if game.State != util.Lobby {
		return fmt.Errorf("game already started or finished")
	}
	game.State = util.Running
	return nil
}

// List active games with states
func (gm *GameMgr) ListGames() map[string]string {
	keys := make(map[string]string)
	for k, v := range gm.Games {
		keys[k] = v.State.String()
	}
	return keys
}

func LoadData(fileData []byte) *Data {
	data := new(Data)
	if err := json.Unmarshal(fileData, &data); err != nil {
		panic(err)
	}

	for index, row := range *data {
		(*data)[index].Votes = make([]int, len(row.Options))
	}
	return data
}

func (d *Data) GetRandomRound() *Round {
	index := rand.IntN(len(*d))
	return &(*d)[index]
}
