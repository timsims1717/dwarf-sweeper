package menus

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
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
	Symbols []string
	SymMats []pixel.Matrix

	clickFn   func()
	leftFn    func()
	rightFn   func()
	hoverFn   func()
	unHoverFn func()

	Transform  *transform.Transform
	HTransform *transform.Transform

	TextColor color.RGBA

	Right    bool
	Hovered  bool
	Disabled bool
	NoHover  bool
	NoShow   bool
	hovered  bool
	disabled bool
	noShowT  bool
	CurrLine int
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
	align := typeface.DefaultAlign
	if i.Right {
		align.H = typeface.Right
	}
	symPos := typeface.SetText(i.Text, i.Raw, 0., align)
	if len(symPos) > 0 {
		t := transform.NewTransform()
		t.Scalar = i.Transform.Scalar.Scaled(SymbolScalar)
		t.UIZoom = camera.Cam.GetZoomScale()
		t.UIPos = camera.Cam.APos
		i.SymMats = []pixel.Matrix{}
		for _, pos := range symPos {
			t.Pos = i.Transform.APos
			t.Pos.X += pos.X
			t.Pos.Y += pos.Y
			t.Update()
			i.SymMats = append(i.SymMats, t.Mat)
		}
	}
	i.Transform.UIZoom = camera.Cam.GetZoomScale()
	i.Transform.UIPos = camera.Cam.APos
	i.Transform.Update()
}

func (i *Item) Draw(target pixel.Target) {
	if i.Text != nil && !i.NoShow && !i.noShowT {
		i.Text.Draw(target, i.Transform.Mat)
		if len(i.SymMats) == len(i.Symbols) {
			for j := 0; j < len(i.Symbols); j++ {
				sym := img.Batchers[cfg.MenuSprites].Sprites[i.Symbols[j]]
				if sym != nil {
					sym.Draw(target, i.SymMats[j])
				}
			}
		}
	}
}

func (i *Item) SetHoverFn(fn func()) {
	i.hoverFn = fn
}

func (i *Item) SetUnhoverFn(fn func()) {
	i.unHoverFn = fn
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