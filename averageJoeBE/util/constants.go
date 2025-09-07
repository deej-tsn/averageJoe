package util

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
	Voting RoundState = iota
	Calculating
	RoundFinished
)

var RoundStateMap = map[RoundState]string{
	0: "Voting",
	1: "Calculating",
	2: "Finished",
}

func (rs RoundState) String() string {
	return RoundStateMap[rs]
}
