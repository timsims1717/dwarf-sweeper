package menu

import (
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"image/color"
)

type ItemText struct {
	Text  *text.Text
	Raw   string
	Align TextAlign
	// todo: add VAlign

	Transform       *transform.Transform
	DefaultSize     pixel.Vec
	HoverSize       pixel.Vec
	TransformEffect *transform.TransformEffect

	TextColor    color.RGBA
	DefaultColor color.RGBA
	HoverColor   color.RGBA
	ColorEffect  *transform.ColorEffect
}

type TextAlign int

const (
	Left = iota
	Center
	Right
)

func NewItemText(raw string, color color.RGBA, size pixel.Vec, centered bool) *ItemText {
	tran := transform.NewTransform(false)
	tran.Scalar = size
	var align TextAlign
	if centered {
		tran.Anchor = transform.Anchor{
			H: transform.Center,
			V: transform.Center,
		}
		align = Center
	} else {
		tran.Anchor = transform.Anchor{
			H: transform.Left,
			V: transform.Bottom,
		}
		align = Left
	}
	return &ItemText{
		Text:         text.New(pixel.ZV, typeface.BasicAtlas),
		Raw:          raw,
		TextColor:    color,
		DefaultColor: color,
		HoverColor:   color,
		DefaultSize:  size,
		HoverSize:    size,
		Transform:    tran,
		Align:        align,
	}
}

func (t *ItemText) Update(r pixel.Rect) {
	if t.TransformEffect != nil {
		t.TransformEffect.Update()
		if t.TransformEffect.IsDone() {
			t.TransformEffect = nil
		}
	}
	if t.ColorEffect != nil {
		t.ColorEffect.Update()
		if t.ColorEffect.IsDone() {
			t.ColorEffect = nil
		}
	}
	t.Text.Clear()
	if t.Align == Center {
		t.Text.Dot.X -= t.Text.BoundsOf(t.Raw).W() / 2.
	} else if t.Align == Right {
		t.Text.Dot.X -= t.Text.BoundsOf(t.Raw).W()
	}
	t.Text.Color = t.TextColor
	fmt.Fprintf(t.Text, t.Raw)
	t.Transform.Rect = t.Text.Bounds()
	t.Transform.Update(r)
}

func (t *ItemText) Draw(target pixel.Target) {
	t.Text.Draw(target, t.Transform.Mat)
}

func (t *ItemText) GetColor() color.RGBA {
	return t.TextColor
}

func (t *ItemText) SetColor(c color.RGBA) {
	t.TextColor = c
}