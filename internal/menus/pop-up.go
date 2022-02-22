package menus

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"github.com/faiface/pixel"
)

const Offset = 20.

var DefaultDist float64

type PopUp struct {
	Raw  string
	Text *typeface.Text
	Dist float64
	Box  *MenuBox

	Tran     *transform.Transform
	Display  bool
	MaxWidth float64
}

func NewPopUp(raw string) *PopUp {
	tran := transform.New()
	tex := typeface.New(nil, "main", typeface.NewAlign(typeface.Center, typeface.Center), 1.2, constants.ActualHintSize, 120., 0.)
	tex.SetColor(DefaultColor)
	tex.SetText(raw)
	box := NewBox(nil, 1.0)
	box.SetSize(pixel.R(0., 0., tex.Width, tex.Height))
	box.SetEntry(Bottom)
	return &PopUp{
		Raw:      raw,
		Text:     tex,
		Dist:     DefaultDist,
		Box:      box,
		Tran:     tran,
		MaxWidth: DefaultMax,
	}
}

func (p *PopUp) SetText(raw string) {
	p.Text.SetText(raw)
	p.Box.SetSize(pixel.R(0., 0., p.Text.Width, p.Text.Height))
}

func (p *PopUp) Update() {
	if p.Display {
		if !p.Box.IsOpen() {
			p.Box.Open()
		}
	} else {
		if p.Box.IsOpen() {
			p.Box.Close()
		}
	}
	p.Tran.Pos.Y += Offset + p.Text.Height*0.5
	p.Tran.Update()
	p.Box.Pos = p.Tran.Pos
	p.Box.Update()
	p.Text.SetPos(p.Tran.Pos)
	p.Text.Update()
}

func (p *PopUp) Draw(target pixel.Target) {
	p.Box.Draw(target)
	if !p.Box.IsClosed() && p.Box.IsOpen() {
		p.Text.Draw(target)
	}
}
