package menu

import (
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/pkg/animation"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/util"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"image/color"
)

var (
	Start = MenuItem{
		raw: "Start Game",
		text: text.New(pixel.V(0., 0), typeface.BasicAtlas),
		transform: &animation.Transform{
			Pos:    pixel.V(0., -20.),
			Scalar: pixel.V(2.35, 2.35),
		},
		color: colornames.Aliceblue,
	}
	Exit = MenuItem{
		raw: "Exit",
		text: text.New(pixel.V(0., 0), typeface.BasicAtlas),
		transform: &animation.Transform{
			Pos:    pixel.V(0., -60.),
			Scalar: pixel.V(2.35, 2.35),
		},
		color: colornames.Aliceblue,
	}
	Credits = MenuItem{
		raw: "Credits",
		text: text.New(pixel.V(0., 0), typeface.BasicAtlas),
		transform: &animation.Transform{
			Pos:    pixel.V(0., -100.),
			Scalar: pixel.V(2.35, 2.35),
		},
		color: colornames.Aliceblue,
	}
)

type MenuItem struct {
	raw       string
	text      *text.Text
	transform *animation.Transform
	color     color.RGBA
	Clicked   bool
	hovered   bool
}

func (m *MenuItem) Update() {
	m.Clicked = false
	mat := m.transform.Mat.Moved(pixel.V(-m.text.BoundsOf(m.raw).W()*0.5, m.text.BoundsOf(m.raw).H()*0.5).Scaled(3.))
	if util.PointInside(input.Input.World, m.text.BoundsOf(m.raw), mat) {
		if !m.hovered {
			sfx.SoundPlayer.PlaySound("click", 2.0)
			m.hovered = true
		}
		m.color = colornames.Darkblue
		m.transform.Scalar = pixel.V(3., 3.)
		if input.Input.Click {
			m.Clicked = true
		}
	} else {
		m.hovered = false
		m.color = colornames.Aliceblue
		m.transform.Scalar = pixel.V(2.35, 2.35)
	}
	m.text.Clear()
	m.text.Color = m.color
	m.text.Dot.X -= m.text.BoundsOf(m.raw).W() * 0.5
	fmt.Fprintf(m.text, m.raw)
	m.transform.Update(pixel.Rect{})
}

func (m *MenuItem) Draw(win *pixelgl.Window) {
	m.text.Draw(win, m.transform.Mat)
}