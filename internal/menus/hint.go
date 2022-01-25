package menus

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"math"
)

const (
	DefaultMax = 96.
)

type HintBox struct {
	Raw  string
	Text *text.Text

	Symbols []string
	SymMats []pixel.Matrix

	Tran     *transform.Transform
	Rect     pixel.Rect
	closing  bool
	opened   bool
	Closed   bool
	StepV    float64
	StepH    float64
	Cam      *camera.Camera
	MaxWidth float64

	TTran  *transform.Transform
	Center *transform.Transform
	CTUL   *transform.Transform
	CTUR   *transform.Transform
	CTDR   *transform.Transform
	CTDL   *transform.Transform
	STU    *transform.Transform
	STR    *transform.Transform
	STD    *transform.Transform
	STL    *transform.Transform
	EntryT *transform.Transform
}

func NewHint(cam *camera.Camera) *HintBox {
	tran := transform.NewTransform()
	tran.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	tTran := transform.NewTransform()
	tTran.Scalar = pixel.V(constants.ActualHintSize, constants.ActualHintSize)
	Center := transform.NewTransform()
	CTUL := transform.NewTransform()
	CTUR := transform.NewTransform()
	CTDR := transform.NewTransform()
	CTDL := transform.NewTransform()
	STU := transform.NewTransform()
	STR := transform.NewTransform()
	STD := transform.NewTransform()
	STL := transform.NewTransform()
	EntryT := transform.NewTransform()
	CTUR.Flip = true
	CTDR.Flip = true
	CTDR.Flop = true
	CTDL.Flop = true
	STR.Flip = true
	STD.Flop = true
	tex := text.New(pixel.ZV, typeface.Atlases["main"])
	tex.LineHeight *= 1.2
	return &HintBox{
		Text:     tex,
		Tran:     tran,
		Rect:     pixel.R(0., 0., 16., 16.),
		Closed:   true,
		StepV:    8.,
		StepH:    8.,
		Cam:      cam,
		MaxWidth: DefaultMax,
		TTran:    tTran,
		Center:   Center,
		CTUL:     CTUL,
		CTUR:     CTUR,
		CTDR:     CTDR,
		CTDL:     CTDL,
		STU:      STU,
		STR:      STR,
		STD:      STD,
		STL:      STL,
		EntryT:   EntryT,
	}
}

func (h *HintBox) Update() {
	h.UpdateTransforms()
	h.Text.Clear()
	h.Text.Color = DefaultColor
	if h.Raw != "" {
		symPos := typeface.SetText(h.Text, h.Raw, h.MaxWidth / constants.ActualHintSize, typeface.DefaultAlign)
		if len(symPos) > 0 {
			t := transform.NewTransform()
			t.Scalar = h.Tran.Scalar.Scaled(SymbolScalar)
			h.SymMats = []pixel.Matrix{}
			for _, pos := range symPos {
				t.Pos = h.TTran.APos
				t.Pos.X += pos.X
				t.Pos.Y += pos.Y + 2.
				t.Update()
				h.SymMats = append(h.SymMats, t.Mat)
			}
		}
	}
}

func (h *HintBox) UpdateSize() {
	h.closing = h.Raw == ""
	if !h.closing {
		h.Closed = false
		var width, height float64
		fullW := h.Text.BoundsOf(h.Raw).W() * constants.ActualHintSize
		if fullW < h.MaxWidth {
			width = fullW
			height = h.Text.LineHeight * constants.ActualHintSize
		} else {
			width = h.MaxWidth
			height = math.Ceil(fullW/h.MaxWidth) * h.Text.LineHeight * constants.ActualHintSize
		}
		h.Rect = pixel.R(0., 0., width, height)
		if h.StepV < h.Rect.H()*0.5 {
			h.StepV += timing.DT * VStep
		}
		if h.StepV > h.Rect.H()*0.5 {
			h.StepV = h.Rect.H() * 0.5
		}
		if h.StepH < h.Rect.W()*0.5 {
			h.StepH += timing.DT * HStep
		}
		if h.StepH > h.Rect.W()*0.5 {
			h.StepH = h.Rect.W() * 0.5
		}
		if h.StepH >= h.Rect.W()*0.5 && h.StepV >= h.Rect.H()*0.5 {
			h.opened = true
		} else {
			h.opened = false
		}
	} else {
		h.opened = false
		if h.StepV > 8. {
			h.StepV -= timing.DT * VStep
		}
		if h.StepV < 8. {
			h.StepV = 8.
		}
		if h.StepH > 8. {
			h.StepH -= timing.DT * HStep
		}
		if h.StepH < 8. {
			h.StepH = 8.
		}
		if h.StepH < 10. && h.StepV < 10. {
			h.Closed = true
		}
	}
}

func (h *HintBox) UpdateTransforms() {
	if h.Cam != nil {
		h.Tran.UIZoom = h.Cam.GetZoomScale()
		h.Tran.UIPos = h.Cam.APos
		h.TTran.UIZoom = h.Cam.GetZoomScale()
		h.TTran.UIPos = h.Cam.APos
		h.CTUL.UIZoom = h.Cam.GetZoomScale()
		h.CTUL.UIPos = h.Cam.APos
		h.CTUR.UIZoom = h.Cam.GetZoomScale()
		h.CTUR.UIPos = h.Cam.APos
		h.CTDR.UIZoom = h.Cam.GetZoomScale()
		h.CTDR.UIPos = h.Cam.APos
		h.CTDL.UIZoom = h.Cam.GetZoomScale()
		h.CTDL.UIPos = h.Cam.APos
		h.STU.UIZoom = h.Cam.GetZoomScale()
		h.STU.UIPos = h.Cam.APos
		h.STR.UIZoom = h.Cam.GetZoomScale()
		h.STR.UIPos = h.Cam.APos
		h.STD.UIZoom = h.Cam.GetZoomScale()
		h.STD.UIPos = h.Cam.APos
		h.STL.UIZoom = h.Cam.GetZoomScale()
		h.STL.UIPos = h.Cam.APos
		h.Center.UIZoom = h.Cam.GetZoomScale()
		h.Center.UIPos = h.Cam.APos
		h.EntryT.UIZoom = h.Cam.GetZoomScale()
		h.EntryT.UIPos = h.Cam.APos
	}
	h.Tran.Update()
	h.CTUL.Pos = pixel.V(h.Tran.Pos.X-h.StepH, h.Tran.Pos.Y+h.StepV)
	h.CTUL.Scalar = pixel.V(1.4, 1.4)
	h.CTUL.Update()
	h.CTUR.Pos = pixel.V(h.Tran.Pos.X+h.StepH, h.Tran.Pos.Y+h.StepV)
	h.CTUR.Scalar = pixel.V(1.4, 1.4)
	h.CTUR.Update()
	h.CTDR.Pos = pixel.V(h.Tran.Pos.X+h.StepH, h.Tran.Pos.Y-h.StepV)
	h.CTDR.Scalar = pixel.V(1.4, 1.4)
	h.CTDR.Update()
	h.CTDL.Pos = pixel.V(h.Tran.Pos.X-h.StepH, h.Tran.Pos.Y-h.StepV)
	h.CTDL.Scalar = pixel.V(1.4, 1.4)
	h.CTDL.Update()
	h.STU.Pos = pixel.V(h.Tran.Pos.X, h.Tran.Pos.Y+h.StepV)
	h.STU.Scalar = pixel.V(1.4*h.StepH*0.1735, 1.4)
	h.STU.Update()
	h.STR.Pos = pixel.V(h.Tran.Pos.X+h.StepH, h.Tran.Pos.Y)
	h.STR.Scalar = pixel.V(1.4, 1.4*h.StepV*0.1735)
	h.STR.Update()
	h.STD.Pos = pixel.V(h.Tran.Pos.X, h.Tran.Pos.Y-h.StepV)
	h.STD.Scalar = pixel.V(1.4*h.StepH*0.1735, 1.4)
	h.STD.Update()
	h.STL.Pos = pixel.V(h.Tran.Pos.X-h.StepH, h.Tran.Pos.Y)
	h.STL.Scalar = pixel.V(1.4, 1.4*h.StepV*0.1735)
	h.STL.Update()
	h.Center.Pos = h.Tran.Pos
	h.Center.Scalar = pixel.V(1.4*h.StepH*0.1735, 1.4*h.StepV*0.1735)
	h.Center.Update()
	h.EntryT.Pos = pixel.V(h.Tran.Pos.X-h.StepH-hintA.Frame().W()*7/6, h.Tran.Pos.Y)
	h.EntryT.Scalar = pixel.V(1.4, 1.4)
	h.EntryT.Update()
	h.TTran.Pos = pixel.V(h.Tran.Pos.X-h.Rect.W()*0.5, h.Tran.Pos.Y+(h.Rect.H()-h.Text.BoundsOf(h.Raw).H()*constants.ActualHintSize)*0.5)
	h.TTran.Update()
}

func (h *HintBox) Draw(target pixel.Target) {
	if !h.Closed {
		inner.Draw(target, h.Center.Mat)
		sideH.Draw(target, h.STU.Mat)
		sideV.Draw(target, h.STR.Mat)
		sideH.Draw(target, h.STD.Mat)
		sideV.Draw(target, h.STL.Mat)
		corner.Draw(target, h.CTUL.Mat)
		corner.Draw(target, h.CTUR.Mat)
		corner.Draw(target, h.CTDR.Mat)
		corner.Draw(target, h.CTDL.Mat)
		hintA.Draw(target, h.EntryT.Mat)
		if !h.closing && h.opened && h.Raw != "" {
			h.Text.Draw(target, h.TTran.Mat)
			if len(h.SymMats) == len(h.Symbols) {
				for j := 0; j < len(h.Symbols); j++ {
					sym := img.Batchers[constants.MenuSprites].Sprites[h.Symbols[j]]
					if sym != nil {
						sym.Draw(target, h.SymMats[j])
					}
				}
			}
		}
	}
}
