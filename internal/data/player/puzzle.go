package player

import (
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type Puzzle interface {
	Create(*pixel.Vec, int)
	IsOpen() bool
	IsClosed() bool
	Open(*Player, string)
	Close()
	Update(*input.Input)
	Draw(pixel.Target)
	Solved() bool
	Failed() bool
	OnSolve()
	OnFail()
	SetOnSolve(func())
	SetOnFail(func())
	SetTimer(*timing.Timer)
}