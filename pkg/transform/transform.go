package transform

import (
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
	ID   uuid.UUID
	Hide bool
	Dead bool

	Pos     pixel.Vec
	Rect    pixel.Rect
	Mat     pixel.Matrix
	Offset  pixel.Vec
	APos    pixel.Vec
	LastPos pixel.Vec
	Rot     float64
	Scalar  pixel.Vec
	Flip    bool
	Flop    bool

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

func (t *Transform) Update() {
	t.APos = t.Pos
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
	t.Mat = t.Mat.ScaledXY(pixel.ZV, t.Scalar.Scaled(t.UIZoom))
	t.Mat = t.Mat.Rotated(pixel.ZV, math.Pi*t.Rot)
	t.Mat = t.Mat.Moved(t.APos.Scaled(t.UIZoom))
	t.Mat = t.Mat.Moved(t.UIPos)
}