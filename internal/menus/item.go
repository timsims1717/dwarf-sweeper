package menus

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"image/color"
)

type Item struct {
	Key     string
	Raw     string
	Hint    string
	Text    *text.Text
	Symbols []string
	SymMats []pixel.Matrix

	clickFn   func()
	leftFn    func()
	rightFn   func()
	hoverFn   func()
	unHoverFn func()

	Transform  *transform.Transform

	TextColor color.RGBA

	Right    bool
	Hovered  bool
	Disabled bool
	NoHover  bool
	Ignore   bool
	NoDraw   bool
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
			t.Pos.Y += pos.Y + i.Text.LineHeight * 0.5
			t.Update()
			i.SymMats = append(i.SymMats, t.Mat)
		}
	}
	i.Transform.UIZoom = camera.Cam.GetZoomScale()
	i.Transform.UIPos = camera.Cam.APos
	i.Transform.Update()
}

func (i *Item) Draw(target pixel.Target) {
	if i.Text != nil && !i.Ignore && !i.noShowT && !i.NoDraw {
		i.Text.Draw(target, i.Transform.Mat)
		if len(i.SymMats) == len(i.Symbols) {
			for j := 0; j < len(i.Symbols); j++ {
				sym := img.Batchers[constants.MenuSprites].Sprites[i.Symbols[j]]
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