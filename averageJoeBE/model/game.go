package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"

	"time"

	"github.com/deej-tsn/averageJoe/util"
	"github.com/gorilla/websocket"
)

type GameMgr struct {
	Games map[string]*Game
}

type Game struct {
	GameID         string
	CurrentRound   *Round
	PreviousRounds []*Round
	State          util.GameState
	Connections    Connections
	RoundTimer     *RoundTimer
}

type Player struct {
	PlayerID string
	Conn     *websocket.Conn
}

type Round struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Votes    []int    `json:"votes"`
	hasVoted map[*websocket.Conn]int
	State    util.RoundState `json:"state"`
}

type RoundTimer struct {
	timer  *time.Timer
	ticker *time.Ticker
}

type roundInfoData struct {
	RoundID   int    `json:"roundID"`
	RoundData *Round `json:"roundData"`
}

type Data []Round

func NewGM() *GameMgr {
	return &GameMgr{
		Games: make(map[string]*Game),
	}
}

func (gm *GameMgr) NewGameFromCode(gameCode string, round *Round) *Game {
	roundTimer := &RoundTimer{
		timer:  nil,
		ticker: nil,
	}
	game := &Game{
		GameID:         gameCode,
		CurrentRound:   round,
		State:          util.Lobby,
		PreviousRounds: make([]*Round, 0),
		Connections:    make(map[*websocket.Conn]bool),
		RoundTimer:     roundTimer,
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

func (gm *GameMgr) StartGame(gameID string, data *Data) error {
	game, exists := gm.Games[gameID]
	if !exists {
		return fmt.Errorf("game not found")
	}
	if game.State != util.Lobby {
		return fmt.Errorf("game already started or finished")
	}
	game.State = util.Running
	game.Connections.BroadcastJSON("Game Started", map[string]string{gameID: "running"})
	log.Default().Println(len(game.PreviousRounds))
	log.Default().Println("Test")
	go func() {
		for {
			if len(game.Connections) == 0 {
				return
			}
			if len(game.Connections) == 1 && len(game.PreviousRounds) > 0 {
				game.Connections.BroadcastJSON("GAME-END", map[string]any{})
				for player := range game.Connections {
					player.Close()
				}
				return
			}
			log.Default().Println("Checking if to start new round")
			if game.CurrentRound.State == util.RoundFinished {
				gm.StartRound(game.GameID, data)
				log.Default().Println("Waiting for round to end")
				time.Sleep(10 * time.Second)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}

func (gm *GameMgr) StartRound(gameID string, data *Data) error {
	game, exists := gm.Games[gameID]
	if !exists {
		return fmt.Errorf("game not found")
	}
	if game.State != util.Running {
		return fmt.Errorf("game has not started yet")
	}
	if game.CurrentRound.State != util.RoundFinished {
		return fmt.Errorf("current round has not ended yet")
	}

	// we can start a new round now
	round := data.GetRandomRound()
	round.State = util.Voting
	log.Default().Println("Starting new round")
	// TODO remove generated random Responses
	// generates random response
	log.Default().Println("Simulating Random Votes")
	for index := range len(round.Votes) {
		round.Votes[index] = rand.IntN(100)
	}

	game.PreviousRounds = append(game.PreviousRounds, game.CurrentRound)
	game.CurrentRound = round
	log.Default().Println("Broadcast new round")
	game.BroadcastRound()
	game.RoundTimer.startRoundTimer(round, game.Connections)

	return nil
}

func (rT *RoundTimer) startRoundTimer(round *Round, conns Connections) {
	rT.timer = time.NewTimer(util.ROUND_DURATION)
	rT.ticker = time.NewTicker(util.ROUND_CLICK_RATE)
	endTime := time.Now().Add(util.ROUND_DURATION)
	go func() {
		for {
			select {
			case <-rT.ticker.C:
				remaining := time.Until(endTime)
				fmt.Printf("Tick : %v\n", remaining)
				if remaining >= time.Second {
					conns.BroadcastJSON("Round-Timer", map[string]time.Duration{"time-left": time.Duration(remaining.Seconds())})
				}
			case <-rT.timer.C:
				fmt.Println("Timer Done")
				round.State = util.Calculating
				conns.BroadcastJSON("Round-End", map[string]string{"round-state": "end"})
				round.calculateRound(conns)
				return
			}
		}
	}()
}

func roundEndResponse(roundData map[string]string, outcomeText string) map[string]any {
	response := map[string]any{
		"messageType": "Round-Results",
		"data": map[string]any{
			"outcome":   outcomeText,
			"roundData": roundData,
		},
	}
	return response
}

func incorrectOption(player *websocket.Conn, response map[string]any, conns Connections) {
	player.WriteJSON(response)
	player.WriteMessage(websocket.TextMessage, []byte("Thanks for playing"))
	player.Close()
	delete(conns, player)
}

func (r *Round) calculateRound(conns Connections) {
	var wG sync.WaitGroup
	output := map[string]string{}
	maxCount := r.Votes[0]
	maxIndex := 0

	for index, option := range r.Options {
		output[option] = fmt.Sprint(r.Votes[index])
		if r.Votes[index] > maxCount {
			maxCount = r.Votes[index]
			maxIndex = index
		}
	}

	incorrectOptionResponse := roundEndResponse(output, "Unlucky")
	correctOptionResponse := roundEndResponse(output, "Lucky")
	didntVoteResponse := roundEndResponse(output, "Next Time, try voting faster")
	for player := range conns {
		wG.Go(func() {
			voteIndex, ok := r.hasVoted[player]
			if !ok {
				incorrectOption(player, didntVoteResponse, conns)
			}
			if voteIndex != maxIndex {
				incorrectOption(player, incorrectOptionResponse, conns)
			} else {
				player.WriteJSON(correctOptionResponse)
			}
		})
	}
	wG.Wait()
	r.State = util.RoundFinished
}

func (game *Game) BroadcastRound() {
	round := game.CurrentRound

	bc := &roundInfoData{
		RoundID:   len(game.PreviousRounds),
		RoundData: round,
	}
	log.Default().Println("BOARDCASTING")
	fmt.Println(bc)
	game.Connections.BroadcastJSON("ROUND-START", bc)
}

// List active games with states
func (gm *GameMgr) ListGames() map[string]string {
	keys := make(map[string]string)
	for k, v := range gm.Games {
		keys[k] = v.State.String()
	}
	return keys
}

func (round *Round) VoteInRound(player *Player, optionIndex int) error {
	if round.State != util.Voting {
		return errors.New("round not accepting vote, current state of round " + round.State.String())
	}
	if optionIndex < 0 || optionIndex >= len(round.Options) {
		return errors.New("invalid Option Index")
	}
	if _, ok := round.hasVoted[player.Conn]; ok {
		return fmt.Errorf("player %s has already voted", player.PlayerID)
	}
	round.Votes[optionIndex] += 1
	round.hasVoted[player.Conn] = optionIndex
	client := player.Conn
	if err := client.WriteJSON(map[string]any{"messageType": "vote-received", "data": optionIndex}); err != nil {
		client.Close()
	}
	return nil
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
	round := (*d)[index]
	roundDeepCopy := &Round{
		Question: round.Question,
		Options:  round.Options,
		Votes:    make([]int, len(round.Options)),
		hasVoted: make(map[*websocket.Conn]int),
		State:    util.RoundFinished,
	}
	return roundDeepCopy
}
