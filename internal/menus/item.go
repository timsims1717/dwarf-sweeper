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
	SRaw    string
	SText   *text.Text
	clickFn func()
	leftFn  func()
	rightFn func()

	Transform  *transform.Transform
	STransform *transform.Transform
	RightMarg  float64

	TextColor color.RGBA

	Hovered  bool
	Disabled bool
	NoHover  bool
	hovered  bool
	disabled bool
	clicked  bool
}

func NewItem(key, raw, sRaw string) *Item {
	tran := transform.NewTransform()
	tran.Anchor = transform.Anchor{
		H: transform.Left,
		V: transform.Bottom,
	}
	tran.Scalar = DefaultSize
	stran := transform.NewTransform()
	stran.Anchor = transform.Anchor{
		H: transform.Left,
		V: transform.Bottom,
	}
	stran.Scalar = DefaultSize
	tex := text.New(pixel.ZV, typeface.BasicAtlas)
	tex.LineHeight *= 1.5
	tex2 := text.New(pixel.ZV, typeface.BasicAtlas)
	tex2.LineHeight *= 1.5
	return &Item{
		Key:         key,
		Raw:         raw,
		Text:        tex,
		SRaw:        sRaw,
		SText:       tex2,
		Transform:   tran,
		STransform:  stran,
		TextColor:   DefaultColor,
	}
}

func (i *Item) Update() {
	if i.Disabled && !i.disabled {
		i.disabled = true
		i.hovered = false
		i.clicked = false
		i.TextColor = DisabledColor
		i.Transform.Scalar = DefaultSize
		i.STransform.Scalar = DefaultSize
	} else if !i.Disabled && i.disabled {
		i.disabled = false
		i.TextColor = DefaultColor
		i.Transform.Scalar = DefaultSize
		i.STransform.Scalar = DefaultSize
	}
	if !i.disabled {
		if i.Hovered && !i.hovered {
			i.hovered = true
			i.TextColor = HoverColor
			i.Transform.Scalar = HoverSize
			i.STransform.Scalar = HoverSize
		} else if !i.Hovered && i.hovered {
			i.hovered = false
			i.TextColor = DefaultColor
			i.Transform.Scalar = DefaultSize
			i.STransform.Scalar = DefaultSize
		}
	}
	i.Text.Clear()
	i.Text.Color = i.TextColor
	fmt.Fprintln(i.Text, i.Raw)

	if i.SRaw != "" {
		i.SText.Clear()
		i.SText.Color = i.TextColor
		i.SText.Dot.X -= i.Text.BoundsOf(i.SRaw).W()
		fmt.Fprintln(i.SText, i.SRaw)
		i.STransform.UIZoom = camera.Cam.GetZoomScale()
		i.STransform.UIPos = camera.Cam.APos
		i.STransform.Update()
	}

	i.Transform.UIZoom = camera.Cam.GetZoomScale()
	i.Transform.UIPos = camera.Cam.APos
	i.Transform.Update()
}

func (i *Item) Draw(target pixel.Target) {
	if i.Text != nil {
		i.Text.Draw(target, i.Transform.Mat)
	}
	if i.SText != nil && i.SRaw != "" {
		i.SText.Draw(target, i.STransform.Mat)
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