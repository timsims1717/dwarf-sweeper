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
	Hint    string
	Text    *text.Text
	HText   *text.Text
	clickFn func()
	leftFn  func()
	rightFn func()
	Right   bool

	Transform  *transform.Transform
	HTransform *transform.Transform

	TextColor color.RGBA

	Hovered  bool
	Disabled bool
	NoHover  bool
	NoShow   bool
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
	htran := transform.NewTransform()
	htran.Anchor = transform.Anchor{
		H: transform.Left,
		V: transform.Bottom,
	}
	htran.Scalar = HintSize
	tex := text.New(pixel.ZV, typeface.BasicAtlas)
	tex.LineHeight *= 1.5
	htex := text.New(pixel.ZV, typeface.BasicAtlas)
	htex.LineHeight *= 1.2
	return &Item{
		Key:        key,
		Raw:        raw,
		Text:       tex,
		HText:      htex,
		Transform:  tran,
		HTransform: htran,
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
	if i.Hovered && !i.Disabled && i.Hint != "" {
		i.HTransform.Pos = i.Transform.Pos
		i.HText.Clear()
		i.HText.Color = DefaultColor
		fmt.Fprintln(i.HText, i.Hint)
		i.HTransform.UIZoom = camera.Cam.GetZoomScale()
		i.HTransform.UIPos = camera.Cam.APos
		i.HTransform.Update()
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
	if i.Text != nil && !i.NoShow {
		i.Text.Draw(target, i.Transform.Mat)
	}
	//if i.HText != nil && !i.NoShow && i.Hovered && !i.Disabled && i.Hint != "" {
	//	i.HText.Draw(target, i.HTransform.Mat)
	//}
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