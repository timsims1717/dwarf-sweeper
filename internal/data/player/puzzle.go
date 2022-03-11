package player

import (
	"dwarf-sweeper/pkg/input"
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
}