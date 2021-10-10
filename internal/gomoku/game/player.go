package game

type Player struct {
	name string
}

func NewPlayer(name string) *Player {
	return &Player{
		name: name,
	}
}

func (p *Player) Name() string {
	return p.name
}
