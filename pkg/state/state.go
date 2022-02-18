package state

import (
	"github.com/faiface/pixel/pixelgl"
)

type State interface {
	Unload()
	Load(chan struct{})
	Update(*pixelgl.Window)
	Draw(*pixelgl.Window)
	SetAbstract(*AbstractState)
}

type AbstractState struct {
	State
	LoadPrc float64
}

func New(state State) *AbstractState {
	aState := &AbstractState{
		State:   state,
	}
	state.SetAbstract(aState)
	return aState
}