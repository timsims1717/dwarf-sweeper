package debug

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
)

var (
	imd *imdraw.IMDraw
)

func InitializeLines() {
	imd = imdraw.New(nil)
}

func DrawLines(target pixel.Target) {
	if Debug {
		imd.Draw(target)
	}
}

func AddLine(color color.RGBA, shape imdraw.EndShape, a, b pixel.Vec, thickness float64) {
	imd.Color = color
	imd.EndShape = shape
	imd.Push(a, b)
	imd.Line(thickness)
}

func AddCircle(color color.RGBA, c pixel.Vec, r, t float64) {
	imd.Color = color
	imd.EndShape = imdraw.NoEndShape
	imd.Push(c)
	imd.Circle(r, t)
}

func AddRect(color color.RGBA, c pixel.Vec, r pixel.Rect, t float64) {
	imd.Color = color
	imd.EndShape = imdraw.NoEndShape
	nr := r.Moved(c).Moved(pixel.V(r.W()*-0.5, r.H()*-0.5))
	vt := nr.Vertices()
	imd.Push(vt[0])
	imd.Push(vt[1])
	imd.Push(vt[2])
	imd.Push(vt[3])
	imd.Polygon(t)
}
