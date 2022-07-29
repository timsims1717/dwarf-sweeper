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
	LoadPrc  float64
	ShowLoad bool
}

func New(state State, showLoad bool) *AbstractState {
	aState := &AbstractState{
		State:    state,
		ShowLoad: showLoad,
	}
	state.SetAbstract(aState)
	return aState
}