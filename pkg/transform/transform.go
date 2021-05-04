package transform

import (
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

//type Transformable interface {
//	GetPos() pixel.Vec
//	SetPos(pixel.Vec)
//	GetRot() float64
//	SetRot(float64)
//	GetScaled() pixel.Vec
//	SetScaled(pixel.Vec)
//}

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
	//Cam    *camera.Camera
	Anchor Anchor
	Rect   pixel.Rect
	Mat    pixel.Matrix
	Pos    pixel.Vec
	Offset pixel.Vec
	RPos   pixel.Vec
	Rot    float64
	Scalar pixel.Vec
	ocent  bool
	Flip   bool
	Flop   bool
}

func NewTransform(isOrigCentered bool) *Transform {
	return &Transform{
		Scalar: pixel.Vec{
			X: 1.,
			Y: 1.,
		},
		ocent: isOrigCentered,
	}
}

func (t *Transform) Update(r pixel.Rect) {
	t.RPos = t.Pos
	if t.ocent {
		if t.Anchor.H == Left {
			t.RPos.X += t.Rect.W() * t.Scalar.X / 2.
		} else if t.Anchor.H == Center {
			t.RPos.X += r.W() / 2.
		} else if t.Anchor.H == Right {
			t.RPos.X += r.W()
			t.RPos.X -= t.Rect.W() * t.Scalar.X / 2.
		}
		if t.Anchor.V == Bottom {
			t.RPos.Y += t.Rect.H() * t.Scalar.Y / 2.
		} else if t.Anchor.V == Center {
			t.RPos.Y += r.H() / 2.
		} else if t.Anchor.V == Top {
			t.RPos.Y += r.H()
			t.RPos.Y -= t.Rect.H() * t.Scalar.Y / 2.
		}
	} else {
		if t.Anchor.H == Center {
			t.RPos.X += r.W() / 2.
		} else if t.Anchor.H == Right {
			t.RPos.X += r.W()
		}
		if t.Anchor.V == Center {
			t.RPos.Y += r.H() / 2.
		} else if t.Anchor.V == Top {
			t.RPos.Y += r.H()
		}
	}
	//if t.Anchor.V == Bottom {
	//	t.RPos.Y += t.Rect.H() / 2.
	//} else if t.Anchor.V == Top {
	//	t.RPos.Y -= t.Rect.H() / 2.
	//}
	t.RPos.X += t.Offset.X
	t.RPos.Y += t.Offset.Y
	//if t.Cam != nil {
	//	t.Mat = t.Cam.UITransform(t.RPos, t.Scalar, t.Rot)
	//} else {
	t.Mat = pixel.IM
	if t.Flip && t.Flop {
		t.Mat = t.Mat.Scaled(pixel.ZV, -1.)
	} else if t.Flip {
		t.Mat = t.Mat.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
	} else if t.Flop {
		t.Mat = t.Mat.ScaledXY(pixel.ZV, pixel.V(1., -1.))
	}
	t.Mat = t.Mat.ScaledXY(pixel.ZV, t.Scalar)
	t.Mat = t.Mat.Rotated(pixel.ZV, t.Rot)
	t.Mat = t.Mat.Moved(t.RPos)
	//}
}

type TransformEffect struct {
	target  *Transform
	interX  *gween.Tween
	interY  *gween.Tween
	interR  *gween.Tween
	interSX *gween.Tween
	interSY *gween.Tween
	isDone  bool
}

func (e *TransformEffect) Update() {
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

func (e *TransformEffect) IsDone() bool {
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

func (b *TransformBuilder) Build() *TransformEffect {
	return &TransformEffect{
		target:  b.Transform,
		interX:  b.InterX,
		interY:  b.InterY,
		interR:  b.InterR,
		interSX: b.InterSX,
		interSY: b.InterSY,
	}
}
