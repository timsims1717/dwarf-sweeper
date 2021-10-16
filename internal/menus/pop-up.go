package menus

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"math"
)

const Offset = 20.

var DefaultDist float64

type PopUp struct {
	Raw  string
	Text *text.Text
	Dist float64
	Box  bool

	Symbols []string
	SymMats []pixel.Matrix

	Parent   *transform.Transform
	Tran     *transform.Transform
	Rect     pixel.Rect
	Display  bool
	closing  bool
	opened   bool
	Closed   bool
	StepV    float64
	StepH    float64
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

func NewPopUp(raw string, parent *transform.Transform) *PopUp {
	tran := transform.NewTransform()
	tran.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	tTran := transform.NewTransform()
	tTran.Scalar = HintSize
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
	EntryT.Rot = math.Pi * 0.5
	tex := text.New(pixel.ZV, typeface.BasicAtlas)
	tex.LineHeight *= 1.2
	return &PopUp{
		Raw:      raw,
		Text:     tex,
		Dist:     DefaultDist,
		Box:      true,
		Parent:   parent,
		Tran:     tran,
		Rect:     pixel.R(0., 0., 16., 16.),
		Closed:   true,
		StepV:    8.,
		StepH:    8.,
		MaxWidth: DefaultMax*2.,
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

func (p *PopUp) Update() {
	p.UpdateSize()
	p.UpdateTransforms()
	p.Text.Clear()
	p.Text.Color = DefaultColor
	if p.Raw != "" {
		symPos := typeface.SetText(p.Text, p.Raw, p.MaxWidth, typeface.DefaultAlign)
		if len(symPos) > 0 {
			t := transform.NewTransform()
			t.Scalar = p.Tran.Scalar.Scaled(SymbolScalar)
			p.SymMats = []pixel.Matrix{}
			for _, pos := range symPos {
				t.Pos = p.TTran.APos
				t.Pos.X += pos.X
				t.Pos.Y += pos.Y + 2.
				t.Update()
				p.SymMats = append(p.SymMats, t.Mat)
			}
		}
	}
}

func (p *PopUp) UpdateSize() {
	p.closing = !p.Display
	if !p.closing {
		p.Closed = false
		var width, height float64
		fullW := p.Text.BoundsOf(p.Raw).W()
		if fullW < p.MaxWidth {
			width = fullW*0.8
			height = p.Text.LineHeight*0.8
		} else {
			width = p.MaxWidth*0.8
			height = math.Ceil(fullW/p.MaxWidth)* p.Text.LineHeight*0.8
		}
		p.Rect = pixel.R(0.,0., width, height)
		if p.StepV < p.Rect.H() * 0.5 {
			p.StepV += timing.DT * VStep
		}
		if p.StepV > p.Rect.H() * 0.5 {
			p.StepV = p.Rect.H() * 0.5
		}
		if p.StepH < p.Rect.W() * 0.5 {
			p.StepH += timing.DT * HStep
		}
		if p.StepH > p.Rect.W() * 0.5 {
			p.StepH = p.Rect.W() * 0.5
		}
		if p.StepH >= p.Rect.W() * 0.5 && p.StepV >= p.Rect.H() * 0.5 {
			p.opened = true
		} else {
			p.opened = false
		}
	} else {
		p.opened = false
		if p.StepV > 8. {
			p.StepV -= timing.DT * VStep
		}
		if p.StepV < 8. {
			p.StepV = 8.
		}
		if p.StepH > 8. {
			p.StepH -= timing.DT * HStep
		}
		if p.StepH < 8. {
			p.StepH = 8.
		}
		if p.StepH < 10. && p.StepV < 10. {
			p.Closed = true
		}
	}
}

func (p *PopUp) UpdateTransforms() {
	if p.Parent != nil {
		p.Tran.Pos = p.Parent.Pos
	}
	p.Tran.Pos.Y += Offset + p.Rect.H() * 0.5
	p.Tran.Update()
	p.CTUL.Pos = pixel.V(p.Tran.Pos.X-p.StepH, p.Tran.Pos.Y+p.StepV)
	p.CTUL.Scalar = pixel.V(1.4, 1.4)
	p.CTUL.Update()
	p.CTUR.Pos = pixel.V(p.Tran.Pos.X+p.StepH, p.Tran.Pos.Y+p.StepV)
	p.CTUR.Scalar = pixel.V(1.4, 1.4)
	p.CTUR.Update()
	p.CTDR.Pos = pixel.V(p.Tran.Pos.X+p.StepH, p.Tran.Pos.Y-p.StepV)
	p.CTDR.Scalar = pixel.V(1.4, 1.4)
	p.CTDR.Update()
	p.CTDL.Pos = pixel.V(p.Tran.Pos.X-p.StepH, p.Tran.Pos.Y-p.StepV)
	p.CTDL.Scalar = pixel.V(1.4, 1.4)
	p.CTDL.Update()
	p.STU.Pos = pixel.V(p.Tran.Pos.X, p.Tran.Pos.Y+p.StepV)
	p.STU.Scalar = pixel.V(1.4 * p.StepH * 0.1735, 1.4)
	p.STU.Update()
	p.STR.Pos = pixel.V(p.Tran.Pos.X+p.StepH, p.Tran.Pos.Y)
	p.STR.Scalar = pixel.V(1.4, 1.4 * p.StepV * 0.1735)
	p.STR.Update()
	p.STD.Pos = pixel.V(p.Tran.Pos.X, p.Tran.Pos.Y-p.StepV)
	p.STD.Scalar = pixel.V(1.4 * p.StepH * 0.1735, 1.4)
	p.STD.Update()
	p.STL.Pos = pixel.V(p.Tran.Pos.X-p.StepH, p.Tran.Pos.Y)
	p.STL.Scalar = pixel.V(1.4, 1.4 * p.StepV * 0.1735)
	p.STL.Update()
	p.Center.Pos = p.Tran.Pos
	p.Center.Scalar = pixel.V(1.4 * p.StepH * 0.1735, 1.4 * p.StepV * 0.1735)
	p.Center.Update()
	p.EntryT.Pos = pixel.V(p.Tran.Pos.X, p.Tran.Pos.Y-p.StepV-hintA.Frame().W()*7/6)
	p.EntryT.Scalar = pixel.V(1.4, 1.4)
	p.EntryT.Update()
	p.TTran.Pos = pixel.V(p.Tran.Pos.X-p.Rect.W() * 0.5, p.Tran.Pos.Y+(p.Rect.H()-p.Text.BoundsOf(p.Raw).H()) * 0.5)
	p.TTran.Update()
}

func (p *PopUp) Draw(target pixel.Target) {
	if !p.Closed {
		if p.Box {
			inner.Draw(target, p.Center.Mat)
			sideH.Draw(target, p.STU.Mat)
			sideV.Draw(target, p.STR.Mat)
			sideH.Draw(target, p.STD.Mat)
			sideV.Draw(target, p.STL.Mat)
			corner.Draw(target, p.CTUL.Mat)
			corner.Draw(target, p.CTUR.Mat)
			corner.Draw(target, p.CTDR.Mat)
			corner.Draw(target, p.CTDL.Mat)
			hintA.Draw(target, p.EntryT.Mat)
		}
		if !p.closing && p.opened && p.Raw != "" {
			p.Text.Draw(target, p.TTran.Mat)
			if len(p.SymMats) == len(p.Symbols) {
				for j := 0; j < len(p.Symbols); j++ {
					sym := img.Batchers[constants.MenuSprites].Sprites[p.Symbols[j]]
					if sym != nil {
						sym.Draw(target, p.SymMats[j])
					}
				}
			}
		}
	}
}