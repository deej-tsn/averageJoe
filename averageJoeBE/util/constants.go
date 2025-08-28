package util

type GameState int

const (
	Lobby GameState = iota
	Running
	Finished
)

var GameStateMap = map[GameState]string{
	0: "Lobby",
	1: "Running",
	2: "Finished",
}

func (gs GameState) String() string {
	return GameStateMap[gs]
}
