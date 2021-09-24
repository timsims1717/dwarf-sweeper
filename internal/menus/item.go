package menus

import (
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"image/color"
)

type Item struct {
	Key     string
	Raw     string
	Text    *text.Text
	clickFn func()
	leftFn  func()
	rightFn func()
	Right   bool

	Transform  *transform.Transform

	TextColor color.RGBA

	Hovered  bool
	Disabled bool
	NoHover  bool
	hovered  bool
	disabled bool
}

func NewItem(key, raw string) *Item {
	tran := transform.NewTransform()
	tran.Anchor = transform.Anchor{
		H: transform.Left,
		V: transform.Bottom,
	}
	tran.Scalar = DefaultSize
	tex := text.New(pixel.ZV, typeface.BasicAtlas)
	tex.LineHeight *= 1.5
	return &Item{
		Key:        key,
		Raw:        raw,
		Text:       tex,
		Transform:  tran,
		TextColor:  DefaultColor,
	}
}

func (i *Item) Update() {
	if i.Disabled && !i.disabled {
		i.disabled = true
		i.hovered = false
		i.TextColor = DisabledColor
		i.Transform.Scalar = DefaultSize
	} else if !i.Disabled && i.disabled {
		i.disabled = false
		i.TextColor = DefaultColor
		i.Transform.Scalar = DefaultSize
	}
	if !i.disabled {
		if i.Hovered && !i.hovered {
			i.hovered = true
			i.TextColor = HoverColor
			i.Transform.Scalar = HoverSize
		} else if !i.Hovered && i.hovered {
			i.hovered = false
			i.TextColor = DefaultColor
			i.Transform.Scalar = DefaultSize
		}
	}
	i.Text.Clear()
	i.Text.Color = i.TextColor
	if i.Right {
		i.Text.Dot.X -= i.Text.BoundsOf(i.Raw).W()
	}
	fmt.Fprintln(i.Text, i.Raw)
	i.Transform.UIZoom = camera.Cam.GetZoomScale()
	i.Transform.UIPos = camera.Cam.APos
	i.Transform.Update()
}

func (i *Item) Draw(target pixel.Target) {
	if i.Text != nil {
		i.Text.Draw(target, i.Transform.Mat)
	}
}

func (i *Item) SetClickFn(fn func()) {
	i.clickFn = fn
}

func (i *Item) SetLeftFn(fn func()) {
	i.leftFn = fn
}

func (i *Item) SetRightFn(fn func()) {
	i.rightFn = fn
}