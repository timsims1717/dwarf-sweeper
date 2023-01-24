package menus

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/menubox"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"github.com/faiface/pixel"
	"strings"
)

const Offset = 20.

var DefaultDist float64

type PopUp struct {
	Raw   string
	Text1 *typeface.Text
	Text2 *typeface.Text
	Text3 *typeface.Text
	Text4 *typeface.Text
	Dist  float64
	Box1  *menubox.MenuBox
	Box2  *menubox.MenuBox
	Box3  *menubox.MenuBox
	Box4  *menubox.MenuBox

	Tran     *transform.Transform
	MaxWidth float64

	Hide     bool
	Display1 bool
	Display2 bool
	Display3 bool
	Display4 bool

	Player1 *data.Player
	Player2 *data.Player
	Player3 *data.Player
	Player4 *data.Player
}

func NewPopUp(raw string) *PopUp {
	tran := transform.New()
	tex1 := typeface.New(nil, "main", typeface.NewAlign(typeface.Center, typeface.Center), 1.2, constants.ActualHintSize, 120., 0.)
	tex1.SetColor(constants.DefaultColor)
	tex1.SetText(strings.Replace(raw, "player", "p1", -1))
	box1 := menubox.NewBox(nil, 1.0)
	box1.SetSize(pixel.R(0., 0., tex1.Width, tex1.Height))
	box1.SetEntry(menubox.Bottom)
	tex2 := typeface.New(nil, "main", typeface.NewAlign(typeface.Center, typeface.Center), 1.2, constants.ActualHintSize, 120., 0.)
	tex2.SetColor(constants.DefaultColor)
	tex2.SetText(strings.Replace(raw, "player", "p2", -1))
	box2 := menubox.NewBox(nil, 1.0)
	box2.SetSize(pixel.R(0., 0., tex2.Width, tex2.Height))
	box2.SetEntry(menubox.Bottom)
	tex3 := typeface.New(nil, "main", typeface.NewAlign(typeface.Center, typeface.Center), 1.2, constants.ActualHintSize, 120., 0.)
	tex3.SetColor(constants.DefaultColor)
	tex3.SetText(strings.Replace(raw, "player", "p3", -1))
	box3 := menubox.NewBox(nil, 1.0)
	box3.SetSize(pixel.R(0., 0., tex3.Width, tex3.Height))
	box3.SetEntry(menubox.Bottom)
	tex4 := typeface.New(nil, "main", typeface.NewAlign(typeface.Center, typeface.Center), 1.2, constants.ActualHintSize, 120., 0.)
	tex4.SetColor(constants.DefaultColor)
	tex4.SetText(strings.Replace(raw, "player", "p4", -1))
	box4 := menubox.NewBox(nil, 1.0)
	box4.SetSize(pixel.R(0., 0., tex4.Width, tex4.Height))
	box4.SetEntry(menubox.Bottom)
	return &PopUp{
		Raw:      raw,
		Text1:    tex1,
		Text2:    tex2,
		Text3:    tex3,
		Text4:    tex4,
		Dist:     DefaultDist,
		Box1:     box1,
		Box2:     box2,
		Box3:     box3,
		Box4:     box4,
		Tran:     tran,
		MaxWidth: DefaultMax,
	}
}

func (p *PopUp) SetText(raw string) {
	p.Text1.SetText(strings.Replace(raw, "player", "p1", -1))
	p.Box1.SetSize(pixel.R(0., 0., p.Text1.Width, p.Text1.Height))
	p.Text2.SetText(strings.Replace(raw, "player", "p2", -1))
	p.Box2.SetSize(pixel.R(0., 0., p.Text2.Width, p.Text2.Height))
	p.Text3.SetText(strings.Replace(raw, "player", "p3", -1))
	p.Box3.SetSize(pixel.R(0., 0., p.Text3.Width, p.Text3.Height))
	p.Text4.SetText(strings.Replace(raw, "player", "p4", -1))
	p.Box4.SetSize(pixel.R(0., 0., p.Text4.Width, p.Text4.Height))
}

func (p *PopUp) Update() {
	if p.Display1 && !p.Hide {
		if !p.Box1.IsOpen() {
			p.Box1.Open()
		}
	} else {
		if p.Box1.IsOpen() {
			p.Box1.Close()
		}
	}
	p.Tran.Pos.Y += Offset + p.Text1.Height*0.5
	p.Tran.Update()
	p.Box1.Pos = p.Tran.Pos
	p.Box1.Update()
	p.Text1.SetPos(p.Tran.Pos)
	p.Text1.Update()

	if p.Display2 && !p.Hide {
		if !p.Box2.IsOpen() {
			p.Box2.Open()
		}
	} else {
		if p.Box2.IsOpen() {
			p.Box2.Close()
		}
	}
	p.Box2.Pos = p.Tran.Pos
	p.Box2.Update()
	p.Text2.SetPos(p.Tran.Pos)
	p.Text2.Update()

	if p.Display3 && !p.Hide {
		if !p.Box3.IsOpen() {
			p.Box3.Open()
		}
	} else {
		if p.Box3.IsOpen() {
			p.Box3.Close()
		}
	}
	p.Box3.Pos = p.Tran.Pos
	p.Box3.Update()
	p.Text3.SetPos(p.Tran.Pos)
	p.Text3.Update()

	if p.Display4 && !p.Hide {
		if !p.Box4.IsOpen() {
			p.Box4.Open()
		}
	} else {
		if p.Box4.IsOpen() {
			p.Box4.Close()
		}
	}
	p.Box4.Pos = p.Tran.Pos
	p.Box4.Update()
	p.Text4.SetPos(p.Tran.Pos)
	p.Text4.Update()
}

func (p *PopUp) Draw() {
	if !p.Box1.IsClosed() {
		p.Box1.Draw(p.Player1.Canvas)
		if p.Box1.IsOpen() {
			p.Text1.Draw(p.Player1.Canvas)
		}
	}
	if !p.Box2.IsClosed() {
		p.Box2.Draw(p.Player2.Canvas)
		if p.Box2.IsOpen() {
			p.Text2.Draw(p.Player2.Canvas)
		}
	}
	if !p.Box3.IsClosed() {
		p.Box3.Draw(p.Player3.Canvas)
		if p.Box3.IsOpen() {
			p.Text3.Draw(p.Player3.Canvas)
		}
	}
	if !p.Box4.IsClosed() {
		p.Box4.Draw(p.Player4.Canvas)
		if p.Box4.IsOpen() {
			p.Text4.Draw(p.Player4.Canvas)
		}
	}
}
