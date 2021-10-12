package menu

import (
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/input"
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
	TransformEffect *transform.Effect
	
	Show bool
	Cam  *camera.Camera

	Mask        color.RGBA
	ColorEffect *transform.ColorEffect
}

func NewMenu(rect pixel.Rect, cam *camera.Camera) *Menu {
	tran := transform.NewTransform()
	tran.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	tran.SetRect(rect)
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

func (m *Menu) Update(world pixel.Vec, clicked *input.ButtonSet) {
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
			m.Transform.UIZoom = m.Cam.GetZoomScale()
			m.Transform.UIPos = m.Cam.APos
		}
		m.Transform.Update()
		for _, item := range m.Items {
			point := m.Transform.Mat.Unproject(world)
			point.X += m.Canvas.Bounds().W() * 0.5
			point.Y += m.Canvas.Bounds().H() * 0.5
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