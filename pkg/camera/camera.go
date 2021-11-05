package camera

import (
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"math/rand"
	"time"
)

var (
	Cam *Camera
)

type Camera struct {
	Height float64
	Width  float64
	Mat    pixel.Matrix
	Pos    pixel.Vec
	APos   pixel.Vec
	Zoom   float64
	zStep  float64
	Opt    Options
	Mask   color.RGBA
	Effect *transform.ColorEffect
	IsWin  bool
	iLock  bool

	interX *gween.Tween
	interY *gween.Tween
	interZ *gween.Tween
	shakeX *gween.Tween
	shakeY *gween.Tween
	lock   bool
	random *rand.Rand
}

type Options struct {
	ScrollSpeed float64
	ZoomStep    float64
	ZoomSpeed   float64
	WindowScale float64
}

func New(isWin bool) *Camera {
	return &Camera{
		Mat:   pixel.IM,
		Pos:   pixel.ZV,
		Zoom:  1.0,
		zStep: 1.0,
		Opt: Options{
			ScrollSpeed: 40.0,
			ZoomStep:    1.2,
			ZoomSpeed:   0.2,
			WindowScale: 900.,
		},
		Mask:   colornames.White,
		IsWin:  isWin,
		random: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (c *Camera) SetSize(width, height float64) {
	c.Width = width
	c.Height = height
}

func (c *Camera) GetZoomScale() float64 {
	return 1 / c.Zoom
}

func (c *Camera) Moving() bool {
	return c.lock
}

func (c *Camera) Restrict(bl, tr pixel.Vec) {
	world := c.Pos
	if bl.X <= tr.X {
		if bl.X > world.X {
			world.X = bl.X
		} else if tr.X < world.X {
			world.X = tr.X
		}
	}
	if bl.Y <= tr.Y {
		if bl.Y > world.Y {
			world.Y = bl.Y
		} else if tr.Y < world.Y {
			world.Y = tr.Y
		}
	}
	c.Pos = world
}

func (c *Camera) Update(win *pixelgl.Window) {
	if c.IsWin {
		c.SetSize(win.Bounds().W(), win.Bounds().H())
	}
	fin := true
	if c.interX != nil {
		x, finX := c.interX.Update(timing.DT)
		c.Pos.X = x
		if finX {
			c.interX = nil
		} else {
			fin = false
		}
	}
	if c.interY != nil {
		y, finY := c.interY.Update(timing.DT)
		c.Pos.Y = y
		if finY {
			c.interY = nil
		} else {
			fin = false
		}
	}
	if c.interZ != nil {
		z, finZ := c.interZ.Update(timing.DT)
		c.Zoom = z
		if finZ {
			c.interZ = nil
		} else {
			fin = false
		}
	}
	if fin && c.lock {
		c.lock = false
	}
	if c.Effect != nil {
		c.Effect.Update()
		if c.Effect.IsDone() {
			c.Effect = nil
		}
	}
	c.APos = c.Pos
	if c.shakeX != nil {
		x, finSX := c.shakeX.Update(timing.DT)
		c.APos.X += x
		if finSX {
			c.shakeX = nil
		}
	}
	if c.shakeY != nil {
		y, finSY := c.shakeY.Update(timing.DT)
		c.APos.Y += y
		if finSY {
			c.shakeY = nil
		}
	}
	if c.iLock {
		c.APos.X = math.Round(c.APos.X)
		c.APos.Y = math.Round(c.APos.Y)
	}
	c.Mat = pixel.IM.Scaled(c.APos, c.Height / c.Opt.WindowScale).Scaled(c.APos, c.Zoom).Moved(win.Bounds().Center().Sub(c.APos))
	win.SetMatrix(c.Mat)
	win.SetColorMask(c.Mask)
}

func (c *Camera) Stop() {
	c.lock = false
	c.interX = nil
	c.interY = nil
}

func (c *Camera) SnapTo(v pixel.Vec) {
	if !c.lock {
		c.Pos.X = v.X
		c.Pos.Y = v.Y
	}
}

func (c *Camera) StayWithin(v pixel.Vec, d float64) {
	if !c.lock {
		if c.Pos.X >= v.X + d {
			c.Pos.X = v.X + d
		} else if c.Pos.X <= v.X - d {
			c.Pos.X = v.X - d
		}
		if c.Pos.Y >= v.Y + d {
			c.Pos.Y = v.Y + d
		} else if c.Pos.Y <= v.Y - d {
			c.Pos.Y = v.Y - d
		}
	}
}

func (c *Camera) MoveTo(v pixel.Vec, dur float64, lock bool) {
	if !c.lock {
		c.interX = gween.New(c.Pos.X, v.X, dur, ease.InOutQuad)
		c.interY = gween.New(c.Pos.Y, v.Y, dur, ease.InOutQuad)
		c.lock = lock
	}
}

func (c *Camera) Follow(v pixel.Vec, spd float64) {
	if !c.lock {
		c.Pos.X += spd * timing.DT * (v.X - c.Pos.X)
		c.Pos.Y += spd * timing.DT * (v.Y - c.Pos.Y)
	}
}

func (c *Camera) CenterOn(points []pixel.Vec) {
	if !c.lock {
		if points == nil || len(points) == 0 {
			return
		} else if len(points) == 1 {
			c.Pos = points[0]
		} else {
			// todo: center on multiple points + change zoom
		}
	}
}

func (c *Camera) Left() {
	if !c.lock {
		c.Pos.X -= c.Opt.ScrollSpeed * timing.DT
	}
}

func (c *Camera) Right() {
	if !c.lock {
		c.Pos.X += c.Opt.ScrollSpeed * timing.DT
	}
}

func (c *Camera) Down() {
	if !c.lock {
		c.Pos.Y -= c.Opt.ScrollSpeed * timing.DT
	}
}

func (c *Camera) Up() {
	if !c.lock {
		c.Pos.Y += c.Opt.ScrollSpeed * timing.DT
	}
}

func (c *Camera) SetZoom(zoom float64) {
	c.Zoom = zoom
	c.zStep = zoom
}

func (c *Camera) ZoomIn(zoom float64) {
	if !c.lock {
		c.zStep *= math.Pow(c.Opt.ZoomStep, zoom)
		c.interZ = gween.New(c.Zoom, c.zStep, c.Opt.ZoomSpeed, ease.OutQuad)
	}
}

// UITransform returns a pixel.Matrix that can move the center of a pixel.Rect
// to the center of the screen.
func (c *Camera) UITransform(pos, scalar pixel.Vec, rot float64) pixel.Matrix {
	zoom := c.GetZoomScale()
	mat := pixel.IM
	mat = mat.ScaledXY(pixel.ZV, scalar.Scaled(zoom))
	mat = mat.Rotated(pixel.ZV, rot)
	mat = mat.Moved(pixel.V(c.Pos.X, c.Pos.Y))
	mat = mat.Moved(pos.Scaled(zoom))
	return mat
}

func (c *Camera) SetILock(b bool) {
	c.iLock = b
}

func (c *Camera) GetColor() color.RGBA {
	return c.Mask
}

func (c *Camera) SetColor(col color.RGBA) {
	c.Mask = col
}

func (c *Camera) Shake() {
	c.shakeX = gween.New((rand.Float64() - 0.5) * 8., 0., 0.5, Sine)
	c.shakeY = gween.New((rand.Float64() - 0.5) * 8., 0., 0.5, Sine)
}

func Sine(t, b, c, d float64) float64 {
	return b * math.Pow(math.E, -math.Abs(c) * t) * math.Sin(10. * math.Pi * t)
}