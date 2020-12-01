package main

import (
	"fmt"

	"github.com/paulyeo21/go-playground/ai/tic-tac-toe/src"
)

func main() {
	p1 := src.Player{name: "human"}
	p2 := src.Player{name: "computer"}
	board := src.Board{}
	g := src.Game{
		p1:    p1,
		p2:    p2,
		state: board,
	}
	fmt.Sprintf("game %s\n", g.getState())
}
