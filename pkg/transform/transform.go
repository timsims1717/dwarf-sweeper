package transform

import (
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
	"github.com/google/uuid"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
)

type Alignment int

const (
	Left = iota
	Center
	Right
	Top    = Left
	Bottom = Right
)

type Anchor struct {
	H Alignment
	V Alignment
}

type Transform struct {
	ID uuid.UUID

	Anchor  Anchor
	Rect    pixel.Rect
	oRect   pixel.Rect
	Parent  pixel.Rect
	oPRect  pixel.Rect
	Mat     pixel.Matrix
	Pos     pixel.Vec
	Offset  pixel.Vec
	APos    pixel.Vec
	LastPos pixel.Vec
	Rot     float64
	Scalar  pixel.Vec
	Flip    bool
	Flop    bool

	Hide bool
	Mask color.RGBA

	UIPos  pixel.Vec
	UIZoom float64
}

func New() *Transform {
	return &Transform{
		ID: uuid.New(),
		Scalar: pixel.Vec{
			X: 1.,
			Y: 1.,
		},
		UIZoom: 1.,
		Mask:   colornames.White,
	}
}

func (t *Transform) SetRect(r pixel.Rect) {
	t.Rect = r
	t.oRect = r
}

func (t *Transform) SetParent(r pixel.Rect) {
	t.Parent = r
	t.oPRect = r
}

func (t *Transform) Update() {
	t.APos = t.Pos
	if t.Anchor.H == Left {
		t.APos.X += t.Rect.W() * t.Scalar.X / 2.
	} else if t.Anchor.H == Center {
		t.APos.X += t.Parent.W() / 2.
	} else if t.Anchor.H == Right {
		t.APos.X += t.Parent.W()
		t.APos.X -= t.Rect.W() * t.Scalar.X / 2.
	}
	if t.Anchor.V == Bottom {
		t.APos.Y += t.Rect.H() * t.Scalar.Y / 2.
	} else if t.Anchor.V == Center {
		t.APos.Y += t.Parent.H() / 2.
	} else if t.Anchor.V == Top {
		t.APos.Y += t.Parent.H()
		t.APos.Y -= t.Rect.H() * t.Scalar.Y / 2.
	}
	t.APos.X += t.Offset.X
	t.APos.Y += t.Offset.Y
	t.APos.X = math.Round(t.APos.X)
	t.APos.Y = math.Round(t.APos.Y)
	t.Mat = pixel.IM
	if t.Flip && t.Flop {
		t.Mat = t.Mat.Scaled(pixel.ZV, -1.)
	} else if t.Flip {
		t.Mat = t.Mat.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
	} else if t.Flop {
		t.Mat = t.Mat.ScaledXY(pixel.ZV, pixel.V(1., -1.))
	}
	if t.oRect.W() > 0. && t.oRect.H() > 0. {
		t.Mat = t.Mat.ScaledXY(pixel.ZV, pixel.V(t.Rect.W()/t.oRect.W(), t.Rect.H()/t.oRect.H()))
	}
	if t.oPRect.W() > 0. && t.oPRect.H() > 0. {
		t.Mat = t.Mat.ScaledXY(pixel.ZV, pixel.V(t.Parent.W()/t.oPRect.W(), t.Parent.H()/t.oPRect.H()))
	}
	t.Mat = t.Mat.ScaledXY(pixel.ZV, t.Scalar.Scaled(t.UIZoom))
	t.Mat = t.Mat.Rotated(pixel.ZV, math.Pi*t.Rot)
	t.Mat = t.Mat.Moved(t.APos.Scaled(t.UIZoom))
	t.Mat = t.Mat.Moved(t.UIPos)
}

type Effect struct {
	target  *Transform
	interX  *gween.Tween
	interY  *gween.Tween
	interR  *gween.Tween
	interSX *gween.Tween
	interSY *gween.Tween
	isDone  bool
}

func (e *Effect) Update() {
	isDone := true
	pos := e.target.Pos
	rot := e.target.Rot
	sca := e.target.Scalar
	if e.interX != nil {
		x, fin := e.interX.Update(timing.DT)
		pos.X = x
		if fin {
			e.interX = nil
		} else {
			isDone = false
		}
	}
	if e.interY != nil {
		y, fin := e.interY.Update(timing.DT)
		pos.Y = y
		if fin {
			e.interY = nil
		} else {
			isDone = false
		}
	}
	if e.interR != nil {
		r, fin := e.interR.Update(timing.DT)
		rot = r
		if fin {
			e.interR = nil
		} else {
			isDone = false
		}
	}
	if e.interSX != nil {
		x, fin := e.interSX.Update(timing.DT)
		sca.X = x
		if fin {
			e.interSX = nil
		} else {
			isDone = false
		}
	}
	if e.interSY != nil {
		y, fin := e.interSY.Update(timing.DT)
		sca.Y = y
		if fin {
			e.interSY = nil
		} else {
			isDone = false
		}
	}
	e.target.Pos = pos
	e.target.Rot = rot
	e.target.Scalar = sca
	e.isDone = isDone
}

func (e *Effect) IsDone() bool {
	return e.isDone
}

type TransformBuilder struct {
	Transform *Transform
	InterX    *gween.Tween
	InterY    *gween.Tween
	InterR    *gween.Tween
	InterSX   *gween.Tween
	InterSY   *gween.Tween
}

func (b *TransformBuilder) Build() *Effect {
	return &Effect{
		target:  b.Transform,
		interX:  b.InterX,
		interY:  b.InterY,
		interR:  b.InterR,
		interSX: b.InterSX,
		interSY: b.InterSY,
	}
}
