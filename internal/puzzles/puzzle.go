package puzzles

import (
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/input"
	"github.com/faiface/pixel"
)

type Puzzle interface {
	Create(*camera.Camera, int)
	IsOpen() bool
	IsClosed() bool
	Open(string)
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