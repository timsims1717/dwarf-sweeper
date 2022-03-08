package transform

import (
	"github.com/faiface/pixel"
	"github.com/google/uuid"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"math/rand"
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
	Load bool

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

	Shaking bool
	ShakeI  int
	ShakeE  int
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
	if t.Shaking {
		switch t.ShakeI {
		case 0,1,7:
			t.APos.Y += 1.
		case 3,4,5:
			t.APos.Y -= 1.
		}
		switch t.ShakeI {
		case 1,2,3:
			t.APos.X += 1.
		case 5,6,7:
			t.APos.X -= 1.
		}
		t.ShakeI++
		t.ShakeI %= 8
		if t.ShakeI == t.ShakeE {
			t.Shaking = false
		}
	}
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

func (t *Transform) Shake(rando *rand.Rand) {
	t.Shaking = true
	t.ShakeI = rando.Intn(8)
	t.ShakeE = (t.ShakeI + 7) % 8
}