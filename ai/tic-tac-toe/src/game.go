package src

type Player struct {
	name string
}

type Board struct {
}

type Game struct {
	p1    Player
	p2    Player
	state Board
}

func (g Game) getState() string {
	return "state"
}
