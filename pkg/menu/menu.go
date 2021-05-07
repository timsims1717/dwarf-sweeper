package menu

import (
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/color"
)

type Menu struct {
	Items map[string]*Item
	
	Transform       *transform.Transform
	Canvas          *pixelgl.Canvas
	TransformEffect *transform.TransformEffect
	
	Show bool
	Cam  *camera.Camera

	Mask        color.RGBA
	ColorEffect *transform.ColorEffect
}

func NewMenu(rect pixel.Rect, cam *camera.Camera) *Menu {
	tran := transform.NewTransform(true)
	tran.Anchor = transform.Anchor{
		H: transform.Left,
		V: transform.Bottom,
	}
	tran.Rect = rect
	return &Menu{
		Items:     map[string]*Item{},
		Transform: tran,
		Canvas:    pixelgl.NewCanvas(rect),
		Show:      true,
		Cam:       cam,
		Mask:      colornames.White,
	}
}

func (m *Menu) AddItem(key string, item *Item) {
	m.Items[key] = item
}

func (m *Menu) Update(world pixel.Vec, clicked input.Toggle) {
	if m.Show {
		if m.TransformEffect != nil {
			m.TransformEffect.Update()
			if m.TransformEffect.IsDone() {
				m.TransformEffect = nil
			}
		}
		if m.ColorEffect != nil {
			m.ColorEffect.Update()
			if m.ColorEffect.IsDone() {
				m.ColorEffect = nil
			}
		}
		if m.Cam != nil {
			m.Canvas.SetBounds(pixel.R(0,0, m.Cam.Width, m.Cam.Height))
			m.Transform.Rect = m.Canvas.Bounds()
		}
		m.Transform.Update(pixel.Rect{})
		if m.Cam != nil {
			m.Transform.Mat = m.Cam.UITransform(m.Transform.RPos, m.Transform.Scalar, m.Transform.Rot)
		}
		for _, item := range m.Items {
			point := m.Transform.Mat.Unproject(world)
			if m.Transform.OCenter {
				point.X += m.Canvas.Bounds().W() * 0.5
				point.Y += m.Canvas.Bounds().H() * 0.5
			}
			item.Update(m.Canvas.Bounds(), point, clicked)
		}
	}
}

func (m *Menu) Draw(target pixel.Target) {
	m.Canvas.Clear(color.RGBA{})
	if m.Show {
		for _, item := range m.Items {
			item.Draw(m.Canvas)
		}
	}
	m.Canvas.Draw(target, m.Transform.Mat)
}