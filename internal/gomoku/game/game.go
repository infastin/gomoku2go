package game

import "fmt"

type GameState uint
type FieldType uint
type PlayerType uint

const (
	NotFinished GameState = iota
	FirstPlayerWin
	SecondPlayerWin
	NobodyWins
)

const (
	EmptyField FieldType = iota
	FirstPlayerField
	SecondPlayerField
)

const (
	FirstPlayer PlayerType = iota
	SecondPlayer
)

const (
	MinSize = 3
	MaxSize = 20
)

type Game struct {
	state     GameState
	fields    [][]FieldType
	players   []*Player
	curplayer PlayerType
	size      uint
	winCond   uint
}

type Field struct {
	X, Y uint
	Ft   FieldType
}

func NewGame(p1, p2 *Player, size, winCond uint) (*Game, error) {
	if size < MinSize || size > MaxSize {
		return nil, fmt.Errorf("The size of the board can only be in the range of 3x3 to 20x20 squares")
	}

	if winCond < MinSize || winCond > size {
		return nil, fmt.Errorf("The number of markers in a consecutive row to win can only be in the range of 3 to the length of the board")
	}

	result := &Game{
		state:     NotFinished,
		fields:    make([][]FieldType, size),
		players:   []*Player{p1, p2},
		curplayer: 0,
		size:      size,
		winCond:   winCond,
	}

	for i := uint(0); i < size; i++ {
		result.fields[i] = make([]FieldType, size)
	}

	return result, nil
}

func (g *Game) SetField(x, y uint) (bool, error) {
	if x >= g.size || y >= g.size {
		return false, fmt.Errorf("Out of board bounds")
	}

	if g.fields[x][y] == EmptyField {
		g.fields[x][y] = FieldType(g.curplayer) + 1
		return true, nil
	}

	return false, nil
}

func (g *Game) Field(x, y uint) (FieldType, error) {
	if x >= g.size || y >= g.size {
		return EmptyField, fmt.Errorf("Out of board bounds")
	}

	return g.fields[x][y], nil
}

func (g *Game) Size() uint {
	return g.size
}

func (g *Game) CurrentPlayer() PlayerType {
	return g.curplayer
}

func (g *Game) ChangePlayer() {
	g.curplayer = (g.curplayer + 1) % 2
}

func (g *Game) State() GameState {
	return g.state
}

func (g *Game) NotEmptyFields() []Field {
	var res []Field

	for i := uint(0); i < g.size; i++ {
		for j := uint(0); j < g.size; j++ {
			if g.fields[i][j] != EmptyField {
				res = append(res, Field{
					X:  i,
					Y:  j,
					Ft: g.fields[i][j],
				})
			}
		}
	}

	return res
}
