package menus

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/menubox"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"github.com/faiface/pixel"
)

const (
	DefaultMax = 96.
)

type HintBox struct {
	Text *typeface.Text
	Box  *menubox.MenuBox

	Tran     *transform.Transform
	Display  bool
	MaxWidth float64
}

func NewHintBox(raw string, cam *camera.Camera) *HintBox {
	tran := transform.New()
	tex := typeface.New(cam, "main", typeface.NewAlign(typeface.Left, typeface.Center), 1.2, constants.ActualHintSize, DefaultMax, 0.)
	tex.SetColor(constants.DefaultColor)
	tex.SetText(raw)
	box := menubox.NewBox(nil, 1.0)
	box.Cam = cam
	box.SetSize(pixel.R(0., 0., tex.Width, tex.Height))
	box.SetEntry(menubox.Left)
	return &HintBox{
		Text:     tex,
		Box:      box,
		Tran:     tran,
		MaxWidth: DefaultMax,
	}
}

func (p *HintBox) SetText(raw string) {
	p.Text.SetText(raw)
	p.Box.SetSize(pixel.R(0., 0., p.Text.Width, p.Text.Height))
}

func (p *HintBox) Update() {
	if p.Display {
		if !p.Box.IsOpen() {
			p.Box.Open()
		}
	} else {
		if p.Box.IsOpen() {
			p.Box.Close()
		}
	}
	//p.Tran.Pos.Y += Offset + p.Text.Height*0.5
	p.Tran.Update()
	p.Box.Pos = p.Tran.Pos
	p.Box.Update()
	tPos := p.Tran.Pos
	tPos.X -= p.Box.Rect.W()*0.5
	tPos.Y += 2.
	p.Text.SetPos(tPos)
	p.Text.Update()
}

func (p *HintBox) Draw(target pixel.Target) {
	p.Box.Draw(target)
	if !p.Box.IsClosed() && p.Box.IsOpen() {
		p.Text.Draw(target)
	}
}
