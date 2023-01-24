package data

import (
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
	pxginput "github.com/timsims1717/pixel-go-input"
)

type Puzzle interface {
	Create(*pixel.Vec, int)
	IsOpen() bool
	IsClosed() bool
	Open(*Player, string)
	Close()
	Update(*pxginput.Input)
	Draw(pixel.Target)
	Solved() bool
	Failed() bool
	OnSolve()
	OnFail()
	SetOnSolve(func())
	SetOnFail(func())
	SetTimer(*timing.Timer)
}