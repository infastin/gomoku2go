package main

import (
	"fmt"
	"lab2go/internal"
)

func main() {
	p1 := gomoku.NewPlayer("Alex")
	p2 := gomoku.NewPlayer("Marti")

	game := gomoku.NewGame(p1, p2)

	game.SetField(0, 0)
	game.SetField(0, 1)
	game.SetField(0, 2)

	state := game.CheckState()


	fmt.Println(state)
}
