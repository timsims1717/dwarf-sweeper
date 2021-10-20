package debug

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
)

var (
	imd *imdraw.IMDraw
)

func InitializeLines() {
	imd = imdraw.New(nil)
}

func DrawLines(win *pixelgl.Window) {
	imd.Draw(win)
	imd.Clear()
}

func AddLine(color color.RGBA, shape imdraw.EndShape, a, b pixel.Vec, thickness float64) {
	imd.Color = color
	imd.EndShape = shape
	imd.Push(a, b)
	imd.Polygon(thickness)
}

func AddCircle(color color.RGBA, c pixel.Vec, r float64) {
	imd.Color = color
	imd.EndShape = imdraw.NoEndShape
	imd.Push(c)
	imd.Circle(r, 0.)
}

func AddRect(color color.RGBA, c pixel.Vec, r pixel.Rect) {
	imd.Color = color
	imd.EndShape = imdraw.NoEndShape
	nr := r.Moved(c).Moved(pixel.V(r.W()*-0.5, r.H()*-0.5))
	vt := nr.Vertices()
	imd.Push(vt[0])
	imd.Push(vt[1])
	imd.Push(vt[2])
	imd.Push(vt[3])
	imd.Polygon(0.)
}