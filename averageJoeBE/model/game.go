package model

import (
	"encoding/json"
	"math/rand/v2"
)

type Game struct {
	GameID       string
	NumOfPlayers int
	CurrentRound *Round
}

type Round struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Votes    []int
}

type Data []Round

func NewGame(round *Round) *Game {
	return &Game{
		GameID:       "ABCDE1",
		NumOfPlayers: 0,
		CurrentRound: round,
	}
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
