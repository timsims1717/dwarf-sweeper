package transform

import (
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/timing"
	"golang.org/x/image/colornames"
	"image/color"
)

type Colorable interface {
	GetColor() color.RGBA
	SetColor(color.RGBA)
}

type ColorEffect struct {
	target Colorable
	interR *gween.Tween
	interG *gween.Tween
	interB *gween.Tween
	interA *gween.Tween
	isDone bool
}

func (e *ColorEffect) Update() {
	isDone := true
	col := e.target.GetColor()
	if e.interR != nil {
		r, fin := e.interR.Update(timing.DT)
		col.R = uint8(r)
		if fin {
			e.interR = nil
		} else {
			isDone = false
		}
	}
	if e.interG != nil {
		g, fin := e.interG.Update(timing.DT)
		col.G = uint8(g)
		if fin {
			e.interG = nil
		} else {
			isDone = false
		}
	}
	if e.interB != nil {
		b, fin := e.interB.Update(timing.DT)
		col.B = uint8(b)
		if fin {
			e.interB = nil
		} else {
			isDone = false
		}
	}
	if e.interA != nil {
		a, fin := e.interA.Update(timing.DT)
		col.A = uint8(a)
		if fin {
			e.interA = nil
		} else {
			isDone = false
		}
	}
	e.target.SetColor(col)
	e.isDone = isDone
}

func (e *ColorEffect) IsDone() bool {
	return e.isDone
}

func FadeIn(target Colorable, dur float64) *ColorEffect {
	start := colornames.Black
	start.A = 0
	end := target.GetColor()
	return &ColorEffect{
		target: target,
		interR: gween.New(float64(start.R), float64(end.R), dur, ease.Linear),
		interG: gween.New(float64(start.G), float64(end.G), dur, ease.Linear),
		interB: gween.New(float64(start.B), float64(end.B), dur, ease.Linear),
		interA: gween.New(float64(start.A), float64(end.A), dur, ease.Linear),
		isDone: false,
	}
}

func FadeOut(target Colorable, dur float64) *ColorEffect {
	start := target.GetColor()
	end := colornames.Black
	end.A = 0
	return &ColorEffect{
		target: target,
		interR: gween.New(float64(start.R), float64(end.R), dur, ease.Linear),
		interG: gween.New(float64(start.G), float64(end.G), dur, ease.Linear),
		interB: gween.New(float64(start.B), float64(end.B), dur, ease.Linear),
		interA: gween.New(float64(start.A), float64(end.A), dur, ease.Linear),
		isDone: false,
	}
}

func FadeFrom(target Colorable, col color.RGBA, dur float64) *ColorEffect {
	start := col
	end := target.GetColor()
	return &ColorEffect{
		target: target,
		interR: gween.New(float64(start.R), float64(end.R), dur, ease.Linear),
		interG: gween.New(float64(start.G), float64(end.G), dur, ease.Linear),
		interB: gween.New(float64(start.B), float64(end.B), dur, ease.Linear),
		interA: gween.New(float64(start.A), float64(end.A), dur, ease.Linear),
		isDone: false,
	}
}

func FadeTo(target Colorable, col color.RGBA, dur float64) *ColorEffect {
	start := target.GetColor()
	end := col
	return &ColorEffect{
		target: target,
		interR: gween.New(float64(start.R), float64(end.R), dur, ease.Linear),
		interG: gween.New(float64(start.G), float64(end.G), dur, ease.Linear),
		interB: gween.New(float64(start.B), float64(end.B), dur, ease.Linear),
		interA: gween.New(float64(start.A), float64(end.A), dur, ease.Linear),
		isDone: false,
	}
}

func Reset(target Colorable, dur float64) *ColorEffect {
	return FadeTo(target, colornames.White, dur)
}

type ColorBuilder struct {
	Target Colorable
	InterR *gween.Tween
	InterG *gween.Tween
	InterB *gween.Tween
	InterA *gween.Tween
}

func (b *ColorBuilder) Build() *ColorEffect {
	return &ColorEffect{
		target: b.Target,
		interR: b.InterR,
		interG: b.InterG,
		interB: b.InterB,
		interA: b.InterA,
	}
}
