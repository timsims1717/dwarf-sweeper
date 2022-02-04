package menus

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"image/color"
)

const (
	VStep    = 300.
	HStep    = 400.
)

var (
	DefaultColor  color.RGBA
	HoverColor    color.RGBA
	DisabledColor color.RGBA
	SymbolScalar  float64

	corner *pixel.Sprite
	sideV  *pixel.Sprite
	sideH  *pixel.Sprite
	inner  *pixel.Sprite
	arrow  *pixel.Sprite
	hintA  *pixel.Sprite
)

func Initialize() {
	corner = img.Batchers[constants.MenuSprites].Sprites["menu_corner"]
	sideV = img.Batchers[constants.MenuSprites].Sprites["menu_side_v"]
	sideH = img.Batchers[constants.MenuSprites].Sprites["menu_side_h"]
	inner = img.Batchers[constants.MenuSprites].Sprites["menu_inner"]
	arrow = img.Batchers[constants.MenuSprites].Sprites["menu_arrow"]
	hintA = img.Batchers[constants.MenuSprites].Sprites["menu_side_entry"]
	DefaultColor = color.RGBA{
		R: 74,
		G: 84,
		B: 98,
		A: 255,
	}
	HoverColor = colornames.Mediumblue
	DisabledColor = colornames.Darkgray
	SymbolScalar = 0.8
	DefaultDist = world.TileSize * 4.
}

type MenuBox struct {
	Pos    pixel.Vec
	Center *transform.Transform
	CTUL   *transform.Transform
	CTUR   *transform.Transform
	CTDR   *transform.Transform
	CTDL   *transform.Transform
	STU    *transform.Transform
	STR    *transform.Transform
	STD    *transform.Transform
	STL    *transform.Transform

	Rect pixel.Rect
	Cam  *camera.Camera

	closed  bool
	closing bool
	opened  bool
	StepV   float64
	StepH   float64
}

func NewBox(cam *camera.Camera) *MenuBox {
	Center := transform.NewTransform()
	CTUL := transform.NewTransform()
	CTUR := transform.NewTransform()
	CTDR := transform.NewTransform()
	CTDL := transform.NewTransform()
	STU := transform.NewTransform()
	STR := transform.NewTransform()
	STD := transform.NewTransform()
	STL := transform.NewTransform()
	CTUR.Flip = true
	CTDR.Flip = true
	CTDR.Flop = true
	CTDL.Flop = true
	STR.Flip = true
	STD.Flop = true

	return &MenuBox{
		Center:  Center,
		CTUL:    CTUL,
		CTUR:    CTUR,
		CTDR:    CTDR,
		CTDL:    CTDL,
		STU:     STU,
		STR:     STR,
		STD:     STD,
		STL:     STL,
		Rect:    pixel.R(0., 0., 16., 16.),
		Cam:     cam,
		StepV:   16.,
		StepH:   16.,
	}
}

func (mb *MenuBox) IsOpen() bool {
	return mb.opened
}

func (mb *MenuBox) IsClosed() bool {
	return mb.closed
}

func (mb *MenuBox) Open() {
	mb.closed = false
	mb.closing = false
	mb.opened = false
}

func (mb *MenuBox) Close() {
	mb.closing = true
	mb.opened = false
}

func (mb *MenuBox) CloseInstant() {
	mb.closing = true
	mb.closed = true
	mb.opened = false
	mb.StepV = 16.
	mb.StepH = 16.
}

func (mb *MenuBox) SetSize(r pixel.Rect) {
	mb.Rect = r
}

func (mb *MenuBox) Update() {
	if !mb.closing {
		if mb.StepV < mb.Rect.H()*0.5 {
			mb.StepV += timing.DT * VStep
		}
		if mb.StepV > mb.Rect.H()*0.5 {
			mb.StepV = mb.Rect.H() * 0.5
		}
		if mb.StepH < mb.Rect.W()*0.5 {
			mb.StepH += timing.DT * HStep
		}
		if mb.StepH > mb.Rect.W()*0.5 {
			mb.StepH = mb.Rect.W() * 0.5
		}
		if mb.StepH >= mb.Rect.W()*0.5 && mb.StepV >= mb.Rect.H()*0.5 {
			mb.opened = true
		}
	} else {
		if mb.StepV > 16. {
			mb.StepV -= timing.DT * VStep
		}
		if mb.StepV < 16. {
			mb.StepV = 16.
		}
		if mb.StepH > 16. {
			mb.StepH -= timing.DT * HStep
		}
		if mb.StepH < 16. {
			mb.StepH = 16.
		}
		if mb.StepH < 20. && mb.StepV < 20. {
			mb.closed = true
		}
	}
	if mb.Cam != nil {
		mb.CTUL.UIZoom = mb.Cam.GetZoomScale()
		mb.CTUL.UIPos = mb.Cam.APos
		mb.CTUR.UIZoom = mb.Cam.GetZoomScale()
		mb.CTUR.UIPos = mb.Cam.APos
		mb.CTDR.UIZoom = mb.Cam.GetZoomScale()
		mb.CTDR.UIPos = mb.Cam.APos
		mb.CTDL.UIZoom = mb.Cam.GetZoomScale()
		mb.CTDL.UIPos = mb.Cam.APos
		mb.STU.UIZoom = mb.Cam.GetZoomScale()
		mb.STU.UIPos = mb.Cam.APos
		mb.STR.UIZoom = mb.Cam.GetZoomScale()
		mb.STR.UIPos = mb.Cam.APos
		mb.STD.UIZoom = mb.Cam.GetZoomScale()
		mb.STD.UIPos = mb.Cam.APos
		mb.STL.UIZoom = mb.Cam.GetZoomScale()
		mb.STL.UIPos = mb.Cam.APos
		mb.Center.UIZoom = mb.Cam.GetZoomScale()
		mb.Center.UIPos = mb.Cam.APos
	}
	mb.CTUL.Pos = pixel.V(mb.Pos.X-mb.StepH, mb.Pos.Y+mb.StepV)
	mb.CTUL.Scalar = pixel.V(1.4, 1.4)
	mb.CTUL.Update()
	mb.CTUR.Pos = pixel.V(mb.Pos.X+mb.StepH, mb.Pos.Y+mb.StepV)
	mb.CTUR.Scalar = pixel.V(1.4, 1.4)
	mb.CTUR.Update()
	mb.CTDR.Pos = pixel.V(mb.Pos.X+mb.StepH, mb.Pos.Y-mb.StepV)
	mb.CTDR.Scalar = pixel.V(1.4, 1.4)
	mb.CTDR.Update()
	mb.CTDL.Pos = pixel.V(mb.Pos.X-mb.StepH, mb.Pos.Y-mb.StepV)
	mb.CTDL.Scalar = pixel.V(1.4, 1.4)
	mb.CTDL.Update()
	mb.STU.Pos = pixel.V(mb.Pos.X, mb.Pos.Y+mb.StepV)
	mb.STU.Scalar = pixel.V(1.4*mb.StepH*0.1735, 1.4)
	mb.STU.Update()
	mb.STR.Pos = pixel.V(mb.Pos.X+mb.StepH, mb.Pos.Y)
	mb.STR.Scalar = pixel.V(1.4, 1.4*mb.StepV*0.1735)
	mb.STR.Update()
	mb.STD.Pos = pixel.V(mb.Pos.X, mb.Pos.Y-mb.StepV)
	mb.STD.Scalar = pixel.V(1.4*mb.StepH*0.1735, 1.4)
	mb.STD.Update()
	mb.STL.Pos = pixel.V(mb.Pos.X-mb.StepH, mb.Pos.Y)
	mb.STL.Scalar = pixel.V(1.4, 1.4*mb.StepV*0.1735)
	mb.STL.Update()
	mb.Center.Pos = mb.Pos
	mb.Center.Scalar = pixel.V(1.4*mb.StepH*0.1735, 1.4*mb.StepV*0.1735)
	mb.Center.Update()
}

func (mb *MenuBox) Draw(target pixel.Target) {
	inner.Draw(target, mb.Center.Mat)
	sideH.Draw(target, mb.STU.Mat)
	sideV.Draw(target, mb.STR.Mat)
	sideH.Draw(target, mb.STD.Mat)
	sideV.Draw(target, mb.STL.Mat)
	corner.Draw(target, mb.CTUL.Mat)
	corner.Draw(target, mb.CTUR.Mat)
	corner.Draw(target, mb.CTDR.Mat)
	corner.Draw(target, mb.CTDL.Mat)
}