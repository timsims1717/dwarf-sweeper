package data

import (
	gween "dwarf-sweeper/pkg/gween64"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
	"strings"
)

type Player struct {
	Code  string
	Attr  Attributes
	Stats Stats
	Gems  int

	Canvas     *pixelgl.Canvas
	CanvasPos  pixel.Vec
	CamPos     pixel.Vec
	PostCamPos pixel.Vec
	RelX       float64
	CamTar     pixel.Vec
	CamVel     pixel.Vec
	InterX     *gween.Tween
	InterY     *gween.Tween
	ShakeX     *gween.Tween
	ShakeY     *gween.Tween
	Lock       bool

	Inventory *Inventory
	Puzzle    Puzzle
	Messages  []*RawMessage
	Message   *Message
	Input     *pxginput.Input
}

func New(code string, in *pxginput.Input) *Player {
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
