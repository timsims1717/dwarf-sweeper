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