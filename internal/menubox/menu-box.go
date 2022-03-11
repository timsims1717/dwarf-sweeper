package menubox

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
)

type EntryDir int

const (
	Left = iota
	Top
	Right
	Bottom
)

const (
	VStep = 250.
	HStep = 300.
	stretch = 0.25
)

var (
	corner *pixel.Sprite
	sideV  *pixel.Sprite
	sideH  *pixel.Sprite
	inner  *pixel.Sprite
	entry  *pixel.Sprite
)

func Initialize() {
	corner = img.Batchers[constants.MenuSprites].Sprites["menu_corner"]
	sideV = img.Batchers[constants.MenuSprites].Sprites["menu_side_v"]
	sideH = img.Batchers[constants.MenuSprites].Sprites["menu_side_h"]
	inner = img.Batchers[constants.MenuSprites].Sprites["menu_inner"]
	entry = img.Batchers[constants.MenuSprites].Sprites["menu_side_entry"]
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

	EntryT   *transform.Transform
	EntryDir EntryDir

	Rect   pixel.Rect
	Scale  float64
	Parent *pixel.Vec

	closed  bool
	closing bool
	opened  bool
	StepV   float64
	StepH   float64
}

func NewBox(parent *pixel.Vec, scale float64) *MenuBox {
	Center := transform.New()
	CTUL := transform.New()
	CTUR := transform.New()
	CTDR := transform.New()
	CTDL := transform.New()
	STU := transform.New()
	STR := transform.New()
	STD := transform.New()
	STL := transform.New()
	CTUR.Flip = true
	CTDR.Flip = true
	CTDR.Flop = true
	CTDL.Flop = true
	STR.Flip = true
	STD.Flop = true

	return &MenuBox{
		Center: Center,
		CTUL:   CTUL,
		CTUR:   CTUR,
		CTDR:   CTDR,
		CTDL:   CTDL,
		STU:    STU,
		STR:    STR,
		STD:    STD,
		STL:    STL,
		Rect:   pixel.R(0., 0., 4., 4.),
		Parent: parent,
		Scale:  scale,
		StepV:  4.,
		StepH:  4.,
		closed: true,
	}
}

func (mb *MenuBox) SetEntry(dir EntryDir) {
	mb.EntryT = transform.New()
	switch dir {
	case Top:
		mb.EntryT.Rot = -0.5
	case Right:
		mb.EntryT.Flip = true
	case Bottom:
		mb.EntryT.Rot = 0.5
	}
	mb.EntryDir = dir
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
	mb.StepV = 4.
	mb.StepH = 4.
	mb.Update()
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
		if mb.StepH < 32. && mb.StepV < 32. {
			mb.closed = true
		}
	}
	if mb.Parent != nil {
		//mb.CTUL.UIZoom = mb.Cam.GetZoomScale()
		mb.CTUL.UIPos = *mb.Parent
		//mb.CTUR.UIZoom = mb.Cam.GetZoomScale()
		mb.CTUR.UIPos = *mb.Parent
		//mb.CTDR.UIZoom = mb.Cam.GetZoomScale()
		mb.CTDR.UIPos = *mb.Parent
		//mb.CTDL.UIZoom = mb.Cam.GetZoomScale()
		mb.CTDL.UIPos = *mb.Parent
		//mb.STU.UIZoom = mb.Cam.GetZoomScale()
		mb.STU.UIPos = *mb.Parent
		//mb.STR.UIZoom = mb.Cam.GetZoomScale()
		mb.STR.UIPos = *mb.Parent
		//mb.STD.UIZoom = mb.Cam.GetZoomScale()
		mb.STD.UIPos = *mb.Parent
		//mb.STL.UIZoom = mb.Cam.GetZoomScale()
		mb.STL.UIPos = *mb.Parent
		//mb.Center.UIZoom = mb.Cam.GetZoomScale()
		mb.Center.UIPos = *mb.Parent
		if mb.EntryT != nil {
			//mb.EntryT.UIZoom = mb.Cam.GetZoomScale()
			mb.EntryT.UIPos = *mb.Parent
		}
	}
	mb.CTUL.Pos = pixel.V(mb.Pos.X-mb.StepH, mb.Pos.Y+mb.StepV)
	mb.CTUL.Scalar = pixel.V(mb.Scale, mb.Scale)
	mb.CTUL.Update()
	mb.CTUR.Pos = pixel.V(mb.Pos.X+mb.StepH, mb.Pos.Y+mb.StepV)
	mb.CTUR.Scalar = pixel.V(mb.Scale, mb.Scale)
	mb.CTUR.Update()
	mb.CTDR.Pos = pixel.V(mb.Pos.X+mb.StepH, mb.Pos.Y-mb.StepV)
	mb.CTDR.Scalar = pixel.V(mb.Scale, mb.Scale)
	mb.CTDR.Update()
	mb.CTDL.Pos = pixel.V(mb.Pos.X-mb.StepH, mb.Pos.Y-mb.StepV)
	mb.CTDL.Scalar = pixel.V(mb.Scale, mb.Scale)
	mb.CTDL.Update()
	mb.STU.Pos = pixel.V(mb.Pos.X, mb.Pos.Y+mb.StepV)
	mb.STU.Scalar = pixel.V(mb.StepH*stretch, mb.Scale)
	mb.STU.Update()
	mb.STR.Pos = pixel.V(mb.Pos.X+mb.StepH, mb.Pos.Y)
	mb.STR.Scalar = pixel.V(mb.Scale, mb.StepV*stretch)
	mb.STR.Update()
	mb.STD.Pos = pixel.V(mb.Pos.X, mb.Pos.Y-mb.StepV)
	mb.STD.Scalar = pixel.V(mb.StepH*stretch, mb.Scale)
	mb.STD.Update()
	mb.STL.Pos = pixel.V(mb.Pos.X-mb.StepH, mb.Pos.Y)
	mb.STL.Scalar = pixel.V(mb.Scale, mb.StepV*stretch)
	mb.STL.Update()
	mb.Center.Pos = mb.Pos
	mb.Center.Scalar = pixel.V(mb.StepH*stretch, mb.StepV*stretch)
	mb.Center.Update()
	if mb.EntryT != nil {
		switch mb.EntryDir {
		case Left:
			mb.EntryT.Pos = pixel.V(mb.Pos.X-mb.StepH-entry.Frame().W()*5/6, mb.Pos.Y)
		case Top:
			mb.EntryT.Pos = pixel.V(mb.Pos.X, mb.Pos.Y+mb.StepV+entry.Frame().W()*5/6)
		case Right:
			mb.EntryT.Pos = pixel.V(mb.Pos.X+mb.StepH+entry.Frame().W()*5/6, mb.Pos.Y)
		case Bottom:
			mb.EntryT.Pos = pixel.V(mb.Pos.X, mb.Pos.Y-mb.StepV-entry.Frame().W()*5/6)
		}
		mb.EntryT.Scalar = pixel.V(mb.Scale, mb.Scale)
		mb.EntryT.Update()
	}
}

func (mb *MenuBox) Draw(target pixel.Target) {
	if !mb.closed {
		inner.Draw(target, mb.Center.Mat)
		sideH.Draw(target, mb.STU.Mat)
		sideV.Draw(target, mb.STR.Mat)
		sideH.Draw(target, mb.STD.Mat)
		sideV.Draw(target, mb.STL.Mat)
		corner.Draw(target, mb.CTUL.Mat)
		corner.Draw(target, mb.CTUR.Mat)
		corner.Draw(target, mb.CTDR.Mat)
		corner.Draw(target, mb.CTDL.Mat)
		if mb.EntryT != nil {
			entry.Draw(target, mb.EntryT.Mat)
		}
	}
}