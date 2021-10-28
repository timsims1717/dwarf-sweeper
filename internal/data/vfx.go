package data

import "dwarf-sweeper/pkg/timing"

type VFX struct {
	Effects []interface{}
}

type Blink struct{
	Timer      *timing.FrameTimer
	BlinkTimer *timing.FrameTimer
	Blink      bool
}

var (
	BlinkSec = 0.18
	ShowSec  = 0.22
)

func NewBlink(blinkSec float64) *Blink {
	return &Blink{
		Timer:      timing.New(blinkSec),
		BlinkTimer: timing.New(BlinkSec),
		Blink:      true,
	}
}