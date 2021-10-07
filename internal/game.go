package gomoku

import "fmt"

type GameState uint
type FieldType uint
type PlayerType uint

const (
	NotFinished GameState = iota
	FirstPlayerWin
	SecondPlayerWin
	Draw
)

const (
	EmptyField FieldType = iota
	FirstPlayerField
	SecondPlayerField
)

type Player struct {
	name string
}

type Game struct {
	state     GameState
	fields    [][]FieldType
	players   []*Player
	curplayer uint
	size      uint
	winCond   uint
}

func NewPlayer(name string) *Player {
	return &Player{
		name: name,
	}
}

func NewGame(p1, p2 *Player, size, winCond uint) (*Game, error) {
	if size < 3 || size > 20 {
		return nil, fmt.Errorf("The size of the board can only be in the range of 3x3 to 20x20 squares")
	}

	if winCond < 3 || winCond > size {
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

func (g *Game) SetField(x, y uint) bool {
	if g.fields[x][y] == EmptyField {
		g.fields[x][y] = FieldType(g.curplayer) + 1
		return true
	}

	return false
}

func (g *Game) GetField(x, y uint) FieldType {
	return g.fields[x][y]
}

func (g *Game) GetPlayerName() string {
	return g.players[g.curplayer].name
}

func (g *Game) ChangePlayer() {
	g.curplayer = (g.curplayer + 1) % 2
}

func (g *Game) CheckState() GameState {

	return NotFinished
}
