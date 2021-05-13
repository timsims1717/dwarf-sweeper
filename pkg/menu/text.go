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
	Text     *text.Text
	Raw      string
	lines    []string
	HAlign   TextAlign
	VAlign   TextAlign
	MaxWidth float64

	Transform       *transform.Transform
	DefaultSize     pixel.Vec
	HoverSize       pixel.Vec
	TransformEffect *transform.Effect

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
	Bottom = Left
	Top    = Right
)

func NewItemText(raw string, color color.RGBA, size pixel.Vec, hAlign, vAlign TextAlign) *ItemText {
	tran := transform.NewTransform()
	if hAlign == Left {
		tran.Anchor.H = transform.Left
	} else if hAlign == Center {
		tran.Anchor.H = transform.Center
	} else {
		tran.Anchor.H = transform.Right
	}
	if vAlign == Bottom {
		tran.Anchor.V = transform.Bottom
	} else if vAlign == Center {
		tran.Anchor.V = transform.Center
	} else {
		tran.Anchor.V = transform.Top
	}
	tran.Scalar = size
	item := &ItemText{
		Text:         text.New(pixel.ZV, typeface.BasicAtlas),
		MaxWidth:     0.,
		TextColor:    color,
		DefaultColor: color,
		HoverColor:   color,
		DefaultSize:  size,
		HoverSize:    size,
		Transform:    tran,
		HAlign:       hAlign,
		VAlign:       vAlign,
	}
	item.SetText(raw)
	return item
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
	t.Text.Color = t.TextColor
	for i, s := range t.lines {
		if t.VAlign == Top {
			t.Text.Dot.Y = t.Text.Orig.Y - t.Text.LineHeight*float64(i+1)
		} else if t.VAlign == Center {
			if len(t.lines) % 2 == 0 {
				t.Text.Dot.Y = t.Text.Orig.Y + t.Text.LineHeight*(float64(len(t.lines) / 2-i)-1.)
			} else {
				t.Text.Dot.Y = t.Text.Orig.Y + t.Text.LineHeight*(float64(len(t.lines) / 2-i)-0.5)
			}
		} else {
			t.Text.Dot.Y = t.Text.Orig.Y + t.Text.LineHeight*float64(len(t.lines)-i-1)
		}
		if t.HAlign == Center {
			t.Text.Dot.X -= t.Text.BoundsOf(s).W() / 2.
		} else if t.HAlign == Right {
			t.Text.Dot.X -= t.Text.BoundsOf(s).W()
		}
		fmt.Fprintln(t.Text, s)
	}
	t.Transform.Parent = r
	t.Transform.Update()
}

func (t *ItemText) Draw(target pixel.Target) {
	t.Text.Draw(target, t.Transform.Mat)
}

func (t *ItemText) SetMaxWidth(w float64) {
	t.MaxWidth = w
}

func (t *ItemText) SetText(raw string) {
	t.Raw = raw
	t.lines = []string{}
	b := 0
	e := 0
	cut := false
	space := false
	for i, r := range raw {
		switch r {
		case '\n':
			cut = true
			e = i
		case ' ', '\t':
			space = true
			e = i
		}
		if t.MaxWidth > 0. && t.Text.BoundsOf(raw[b:i]).W() > t.MaxWidth && space {
			cut = true
		}
		if cut {
			if b >= e || e < 0 {
				t.lines = append(t.lines, "")
			} else {
				t.lines = append(t.lines, raw[b:e])
			}
			cut = false
			space = false
			b = e+1
		}
	}
	t.lines = append(t.lines, raw[b:])
}

func (t *ItemText) GetColor() color.RGBA {
	return t.TextColor
}

func (t *ItemText) SetColor(c color.RGBA) {
	t.TextColor = c
}