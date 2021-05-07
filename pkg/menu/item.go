package menu

import (
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/color"
)

type Item struct {
	Text *ItemText
	Show bool

	Transform       *transform.Transform
	Spr             *pixel.Sprite
	Canvas          *pixelgl.Canvas
	TransformEffect *transform.TransformEffect

	Mask        color.RGBA
	ColorEffect *transform.ColorEffect

	hovered        bool
	clicked        bool
	Disabled       bool
	disabled       bool
	HoverDefault   bool
	UnHoverDefault bool
	hoverFn        func()
	onHoverFn      func()
	unHoverFn      func()
	clickFn        func()
	unClickFn      func()
	onDisabledFn   func()
	disabledFn     func()
	onEnabledFn    func()
}

func NewItem(t *ItemText, rect pixel.Rect) *Item {
	tran := transform.NewTransform(true)
	tran.Anchor = transform.Anchor{
		H: transform.Left,
		V: transform.Bottom,
	}
	tran.Rect = rect
	item := &Item{
		Text:      t,
		Transform: tran,
		Canvas:    pixelgl.NewCanvas(rect),
		Mask:      colornames.White,
		Show:      true,

		HoverDefault:   true,
		UnHoverDefault: true,
	}
	return item
}

func (i *Item) Update(r pixel.Rect, cursor pixel.Vec, clicked input.Toggle) {
	if i.Show {
		if i.Disabled {
			if i.disabled && i.disabledFn != nil {
				i.disabledFn()
			} else if !i.disabled && i.onDisabledFn != nil {
				i.onDisabledFn()
			}
			i.disabled = true
			if i.hovered {
				i.unHoverFn()
			}
			i.hovered = false
		} else {
			if i.disabled && i.onEnabledFn != nil {
				i.onEnabledFn()
			}
			i.disabled = false
			mouseOver := i.PointInside(cursor)
			if mouseOver && !i.hovered {
				if i.HoverDefault {
					i.Text.TextColor = i.Text.HoverColor
					i.Text.Transform.Scalar = i.Text.HoverSize
				}
				if i.onHoverFn != nil {
					i.onHoverFn()
				}
			} else if !mouseOver && i.hovered {
				if i.UnHoverDefault {
					i.Text.TextColor = i.Text.DefaultColor
					i.Text.Transform.Scalar = i.Text.DefaultSize
				}
				if i.unHoverFn != nil {
					i.unHoverFn()
				}
			} else if i.hovered && i.hoverFn != nil {
				i.hoverFn()
			}
			i.hovered = mouseOver
			if clicked.JustPressed() && i.hovered {
				clicked.Consume()
				if i.clickFn != nil {
					i.clickFn()
				}
				i.clicked = true
			} else {
				if i.clicked && i.unClickFn != nil {
					i.unClickFn()
				}
				i.clicked = false
			}
		}
		if i.Text != nil {
			i.Text.Update(i.Canvas.Bounds())
		}
		if i.TransformEffect != nil {
			i.TransformEffect.Update()
			if i.TransformEffect.IsDone() {
				i.TransformEffect = nil
			}
		}
		if i.ColorEffect != nil {
			i.ColorEffect.Update()
			if i.ColorEffect.IsDone() {
				i.ColorEffect = nil
			}
		}
		i.Transform.Rect = i.Canvas.Bounds()
		i.Transform.Update(r)
	}
}

func (i *Item) Draw(target pixel.Target) {
	if i.Show {
		i.Canvas.Clear(color.RGBA{})
		if i.Spr != nil {
			r := i.Spr.Frame()
			i.Spr.Draw(i.Canvas, pixel.IM.Moved(pixel.V(r.W()*0.5, r.H()*0.5)))
		}
		if i.Text != nil {
			i.Text.Draw(i.Canvas)
		}
		//i.Canvas.DrawColorMask(target, i.Mat, i.Mask)
		i.Canvas.Draw(target, i.Transform.Mat)
		//DebugDraw.Draw(target)
	}
}

func (i *Item) GetColor() color.RGBA {
	return i.Mask
}

func (i *Item) SetColor(mask color.RGBA) {
	i.Mask = mask
}

func (i *Item) IsHovered() bool {
	return i.hovered
}

func (i *Item) OnHover() {
	i.onHoverFn()
}

func (i *Item) Hover() {
	i.hoverFn()
}

func (i *Item) OnUnHover() {
	i.unHoverFn()
}

func (i *Item) Click() {
	i.clicked = true
}

func (i *Item) IsClicked() bool {
	return i.clicked
}

func (i *Item) OnClick() {
	i.clickFn()
}

func (i *Item) OnUnClick() {
	i.unClickFn()
}

func (i *Item) SetHoverFn(fn func()) {
	i.hoverFn = fn
}

func (i *Item) SetOnHoverFn(fn func()) {
	i.onHoverFn = fn
}

func (i *Item) SetUnHoverFn(fn func()) {
	i.unHoverFn = fn
}

func (i *Item) SetClickFn(fn func()) {
	i.clickFn = fn
}

func (i *Item) SetUnClickFn(fn func()) {
	i.unClickFn = fn
}

func (i *Item) SetOnDisabledFn(fn func()) {
	i.onDisabledFn = fn
}

func (i *Item) SetDisabledFn(fn func()) {
	i.disabledFn = fn
}

func (i *Item) SetEnabledFn(fn func()) {
	i.onEnabledFn = fn
}

func (i *Item) PointInside(point pixel.Vec) bool {
	return util.PointInside(point, i.Canvas.Bounds(), i.Transform.Mat)
}