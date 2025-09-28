package util

import "time"

type GameState int

const (
	Lobby GameState = iota
	Running
	GameFinished
)

var GameStateMap = map[GameState]string{
	0: "Lobby",
	1: "Running",
	2: "Finished",
}

func (gs GameState) String() string {
	return GameStateMap[gs]
}

type RoundState int

const (
	RoundFinished RoundState = iota
	Voting
	Calculating
)

var RoundStateMap = map[RoundState]string{
	0: "Voting",
	1: "Calculating",
	2: "Finished",
}

func (rs RoundState) String() string {
	return RoundStateMap[rs]
}

// 10 second round duration
const ROUND_DURATION = time.Second * 10

// 1 second broadcast
const ROUND_CLICK_RATE = time.Second
