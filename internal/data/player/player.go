package player

import (
	"dwarf-sweeper/pkg/input"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"strings"
)

type Player struct {
	Code      string
	Attr      Attributes
	Stats     Stats
	Flags     Flags
	Canvas    *pixelgl.Canvas
	CanvasPos pixel.Vec
	CamPos    pixel.Vec
	RelX      float64
	CamTar    pixel.Vec
	CamVel    pixel.Vec
	Input     *input.Input
	Inventory *Inventory
	Puzzle    Puzzle
	Messages  []*RawMessage
	Message   *Message
}

func New(code string, in *input.Input) *Player {
	p := &Player{
		Code:      code,
		Attr:      DefaultAttr(),
		Canvas:    pixelgl.NewCanvas(pixel.R(0., 0., 10., 10.)),
		Input:     in,
		Inventory: &Inventory{},
	}
	p.Message = NewMessage(p)
	return p
}

func (p *Player) StartPuzzle(puzz Puzzle) bool {
	if p.Puzzle != nil {
		return false
	}
	p.Puzzle = puzz
	p.Puzzle.Open(p, p.Code)
	return true
}

func (p *Player) GiveMessage(raw string, fn func()) {
	p.Messages = append(p.Messages, &RawMessage{
		Raw:     strings.Replace(raw, "player", p.Code, -1),
		OnClose: fn,
	})
}