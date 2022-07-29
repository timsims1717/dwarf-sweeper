package data

import (
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"image/color"
)

type VFX struct {
	Effects []Effect
}

type Effect interface {
	Update(*transform.Transform)
	IsDone(*transform.Transform) bool
}

type Blink struct {
	Trans      *transform.Transform
	Timer      *timing.Timer
	BlinkTimer *timing.Timer
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

func (blink *Blink) Update(trans *transform.Transform) {
	if blink.BlinkTimer.UpdateDone() {
		blink.Blink = !blink.Blink
		if blink.Blink {
			blink.BlinkTimer = timing.New(BlinkSec)
		} else {
			blink.BlinkTimer = timing.New(ShowSec)
		}
		trans.Hide = blink.Blink
	}
}

func (blink *Blink) IsDone(trans *transform.Transform) bool {
	if blink.Timer.UpdateDone() {
		trans.Hide = false
		return true
	}
	return false
}

type FadeOut struct {
	InterR *gween.Tween
	InterG *gween.Tween
	InterB *gween.Tween
	InterA *gween.Tween
	allFin bool
}

func NewFadeOut(col color.RGBA, dur float64) *FadeOut {
	return &FadeOut{
		InterR: gween.New(float64(col.R), 0, dur, ease.Linear),
		InterG: gween.New(float64(col.G), 0, dur, ease.Linear),
		InterB: gween.New(float64(col.B), 0, dur, ease.Linear),
		InterA: gween.New(float64(col.A), 0, dur, ease.Linear),
	}
}

func NewFadeBlack(col color.RGBA, dur float64) *FadeOut {
	return &FadeOut{
		InterR: gween.New(float64(col.R), 0, dur, ease.Linear),
		InterG: gween.New(float64(col.G), 0, dur, ease.Linear),
		InterB: gween.New(float64(col.B), 0, dur, ease.Linear),
	}
}

func (fade *FadeOut) Update(trans *transform.Transform) {
	allFin := true
	if fade.InterR != nil {
		f, fin := fade.InterR.Update(timing.DT)
		trans.Mask.R = uint8(f)
		if fin {
			fade.InterR = nil
		} else {
			allFin = false
		}
	}
	if fade.InterG != nil {
		f, fin := fade.InterG.Update(timing.DT)
		trans.Mask.G = uint8(f)
		if fin {
			fade.InterG = nil
		} else {
			allFin = false
		}
	}
	if fade.InterB != nil {
		f, fin := fade.InterB.Update(timing.DT)
		trans.Mask.B = uint8(f)
		if fin {
			fade.InterB = nil
		} else {
			allFin = false
		}
	}
	if fade.InterA != nil {
		f, fin := fade.InterA.Update(timing.DT)
		trans.Mask.A = uint8(f)
		if fin {
			fade.InterA = nil
		} else {
			allFin = false
		}
	}
	fade.allFin = allFin
}

func (fade *FadeOut) IsDone(trans *transform.Transform) bool {
	return fade.allFin
}