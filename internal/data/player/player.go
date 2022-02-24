package player

import (
	"dwarf-sweeper/internal/puzzles"
	"dwarf-sweeper/pkg/input"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Player struct {
	Code      string
	Attr      Attributes
	Stats     Stats
	Canvas    *pixelgl.Canvas
	CamPos    pixel.Vec
	CanvasPos pixel.Vec
	CamTar    pixel.Vec
	Input     *input.Input
	Inventory *Inventory
	Puzzle    puzzles.Puzzle
}

func New(code string, in *input.Input) *Player {
	return &Player{
		Code:      code,
		Attr:      DefaultAttr(),
		Canvas:    pixelgl.NewCanvas(pixel.R(0., 0., 10., 10.)),
		Input:     in,
		Inventory: &Inventory{},
	}
}

func (p *Player) StartPuzzle(puzz puzzles.Puzzle) bool {
	if p.Puzzle != nil {
		return false
	}
	p.Puzzle = puzz
	p.Puzzle.Open()
	return true
}