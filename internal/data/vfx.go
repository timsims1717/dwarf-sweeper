package data

import (
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/timing"
	"image/color"
)

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

type FadeOut struct{
	InterR *gween.Tween
	InterG *gween.Tween
	InterB *gween.Tween
	InterA *gween.Tween
}

func NewFadeOut(col color.RGBA, dur float64) *FadeOut {
	return &FadeOut{
		InterR: gween.New(float64(col.R), 0, dur, ease.Linear),
		InterG: gween.New(float64(col.G), 0, dur, ease.Linear),
		InterB: gween.New(float64(col.B), 0, dur, ease.Linear),
		InterA: gween.New(float64(col.A), 0, dur, ease.Linear),
	}
}