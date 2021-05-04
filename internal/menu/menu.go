package menu

import (
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/camera"
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
		raw:   "Start Game",
		text:  text.New(pixel.V(0., 0), typeface.BasicAtlas),
		color: colornames.Aliceblue,
		small: pixel.V(8., 8.),
		big:   pixel.V(9., 9.),
	}
	Exit = MenuItem{
		raw:   "Exit",
		text:  text.New(pixel.V(0., 0), typeface.BasicAtlas),
		color: colornames.Aliceblue,
		small: pixel.V(8., 8.),
		big:   pixel.V(9., 9.),
	}
	Credits = MenuItem{
		raw:   "Credits",
		text:  text.New(pixel.V(0., 0), typeface.BasicAtlas),
		color: colornames.Aliceblue,
		small: pixel.V(8., 8.),
		big:   pixel.V(9., 9.),
	}
	Retry = MenuItem{
		raw:   "Retry",
		text:  text.New(pixel.V(0., 0), typeface.BasicAtlas),
		color: colornames.Aliceblue,
		small: pixel.V(3., 3.),
		big:   pixel.V(3.5, 3.5),
	}
	Menu = MenuItem{
		raw:   "Menu",
		text:  text.New(pixel.V(0., 0), typeface.BasicAtlas),
		color: colornames.Aliceblue,
		small: pixel.V(3., 3.),
		big:   pixel.V(3.5, 3.5),
	}
)

type MenuItem struct {
	raw       string
	text      *text.Text
	Transform *transform.Transform
	color     color.RGBA
	Clicked   bool
	hovered   bool
	small     pixel.Vec
	big       pixel.Vec
}

func (m *MenuItem) Update() {
	m.Clicked = false
	offset := m.Transform.RPos
	offset.X -= m.text.BoundsOf(m.raw).W() * m.big.X * 0.5
	offset.Y += m.text.BoundsOf(m.raw).H() * m.big.Y * 0.5
	mat := camera.Cam.UITransform(offset, m.big, m.Transform.Rot)
	if util.PointInside(input.Input.World, m.text.BoundsOf(m.raw), mat) {
		if !m.hovered {
			sfx.SoundPlayer.PlaySound("click", 2.0)
			m.hovered = true
		}
		m.color = colornames.Darkblue
		m.Transform.Scalar = m.big
		if input.Input.Click {
			m.Clicked = true
		}
	} else {
		m.hovered = false
		m.color = colornames.Aliceblue
		m.Transform.Scalar = m.small
	}
	m.text.Clear()
	m.text.Color = m.color
	m.text.Dot.X -= m.text.BoundsOf(m.raw).W() * 0.5
	fmt.Fprintf(m.text, m.raw)
	m.Transform.Update(pixel.Rect{})
	m.Transform.Mat = camera.Cam.UITransform(m.Transform.RPos, m.Transform.Scalar, m.Transform.Rot)
}

func (m *MenuItem) Draw(win *pixelgl.Window) {
	m.text.Draw(win, m.Transform.Mat)
}