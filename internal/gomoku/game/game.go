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
	empty     uint
}

type Field struct {
	X, Y uint
	Ft   FieldType
}

type Strike struct {
	X0, Y0 uint
	X1, Y1 uint
}

func CheckSettings(p1, p2 *Player, size, winCond uint) error {
	if size < MinSize || size > MaxSize {
		return fmt.Errorf("the size of the board can only be in the range of 3x3 to 20x20 squares")
	}

	if winCond < MinSize || winCond > size {
		return fmt.Errorf("the number of markers in a consecutive row to win can only be in the range of 3 to the length of the board")
	}

	return nil
}

func NewGame(p1, p2 *Player, size, winCond uint) (*Game, error) {
	if size < MinSize || size > MaxSize {
		return nil, fmt.Errorf("the size of the board can only be in the range of 3x3 to 20x20 squares")
	}

	if winCond < MinSize || winCond > size {
		return nil, fmt.Errorf("the number of markers in a consecutive row to win can only be in the range of 3 to the length of the board")
	}

	result := &Game{
		state:     NotFinished,
		fields:    make([][]FieldType, size),
		players:   []*Player{p1, p2},
		curplayer: 0,
		size:      size,
		winCond:   winCond,
		empty:     size*size,
	}

	fieldsSub := make([]FieldType, size*size)

	for i := uint(0); i < size; i++ {
		result.fields[i] = fieldsSub[:size]
		fieldsSub = fieldsSub[size:]
	}

	return result, nil
}

func (g *Game) SetField(x, y uint) (bool, error) {
	if x >= g.size || y >= g.size {
		return false, fmt.Errorf("out of board bounds")
	}

	if g.state != NotFinished {
		return false, fmt.Errorf("game over")
	}

	if g.fields[x][y] == EmptyField {
		g.fields[x][y] = FieldType(g.curplayer) + 1
		g.empty -= 1
		return true, nil
	}

	return false, nil
}

func (g *Game) Field(x, y uint) (FieldType, error) {
	if x >= g.size || y >= g.size {
		return EmptyField, fmt.Errorf("out of board bounds")
	}

	return g.fields[x][y], nil
}

func (g *Game) Size() uint {
	return g.size
}

func (g *Game) CurrentPlayer() PlayerType {
	return g.curplayer
}

func (g *Game) Player(pt PlayerType) *Player {
	return g.players[pt]
}

func (g *Game) ChangePlayer() {
	g.curplayer = (g.curplayer + 1) % 2
}

func (g *Game) checkWinnerAxis(x, y uint, xincr, yincr int) (Strike, bool) {
	count := uint(1)

	x0 := int(x)
	y0 := int(y)

	x0 -= (int(g.winCond) - 1) * xincr
	y0 -= (int(g.winCond) - 1) * yincr

	x1 := x0 + xincr
	y1 := y0 + yincr

	for n := uint(0); n < 2*g.winCond-1; n++ {
		if x0 < 0 || y0 < 0 || x1 < 0 || y1 < 0 ||
			x0 >= int(g.size) || y0 >= int(g.size) ||
			x1 >= int(g.size) || y1 >= int(g.size) {
			goto incr
		}

		if g.fields[x0][y0] == g.fields[x1][y1] && g.fields[x0][y0] != EmptyField {
			count++

			if count == g.winCond {
				return Strike{
					X0: uint(x1 - (int(g.winCond)-1)*xincr),
					Y0: uint(y1 - (int(g.winCond)-1)*yincr),
					X1: uint(x1),
					Y1: uint(y1),
				}, true
			}
		} else {
			count = 1
		}

	incr:
		x0 += xincr
		y0 += yincr
		x1 += xincr
		y1 += yincr
	}

	return Strike{}, false
}

func (g *Game) CheckDraw() bool {
	if g.empty == 0 {
		g.state = NobodyWins
		return true
	}

	return false
}

func (g *Game) CheckWinner(x, y uint) (s Strike, win bool) {
	if s, win = g.checkWinnerAxis(x, y, 1, 0); win {
		g.state = GameState(g.curplayer) + 1
		return
	}

	if s, win = g.checkWinnerAxis(x, y, 0, 1); win {
		g.state = GameState(g.curplayer) + 1
		return
	}

	if s, win = g.checkWinnerAxis(x, y, 1, 1); win {
		g.state = GameState(g.curplayer) + 1
		return
	}

	if s, win = g.checkWinnerAxis(x, y, 1, -1); win {
		g.state = GameState(g.curplayer) + 1
		return
	}

	return
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
