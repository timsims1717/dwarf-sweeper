package typeface

import (
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/transform"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"image/color"
)

type Text struct {
	Raw     string
	Text    *text.Text
	Color   color.RGBA
	Align   Alignment
	Symbols []symbolHandle
	NoShow  bool

	Increment bool
	CurrPos   int
	Width     float64
	Height    float64
	MaxWidth  float64
	MaxHeight float64
	MaxLines  int

	Cam          *camera.Camera
	RelativeSize float64
	SymbolSize   float64
	Transform    *transform.Transform

	rawLines   []string
	lineWidths []float64
	fullHeight float64
}

func New(cam *camera.Camera, atlas string, align Alignment, lineHeight, relativeSize, maxWidth, maxHeight float64) *Text {
	tex := text.New(pixel.ZV, Atlases[atlas])
	tex.LineHeight *= lineHeight
	return &Text{
		Text:         tex,
		Align:        align,
		Color:        colornames.White,
		Width:        maxWidth,
		Height:       maxHeight,
		MaxWidth:     maxWidth,
		MaxHeight:    maxHeight,
		MaxLines:     int(maxHeight / (tex.LineHeight * relativeSize)),
		Cam:          cam,
		RelativeSize: relativeSize,
		SymbolSize:   1.,
		Transform:    transform.New(),
	}
}

func (item *Text) Update() {
	item.Transform.Scalar = pixel.V(item.RelativeSize, item.RelativeSize)
	if item.Cam != nil {
		item.Transform.UIZoom = item.Cam.GetZoomScale()
		item.Transform.UIPos = item.Cam.APos
	}
	item.Transform.Update()
	for _, sym := range item.Symbols {
		if item.Cam != nil {
			sym.trans.UIZoom = item.Cam.GetZoomScale()
			sym.trans.UIPos = item.Cam.APos
		}
		sym.trans.Update()
	}
}

func (item *Text) Draw(target pixel.Target) {
	if !item.NoShow {
		item.Text.Draw(target, item.Transform.Mat)
		for _, sym := range item.Symbols {
			sym.symbol.spr.Draw(target, sym.trans.Mat)
		}
	}
}

func (item *Text) SetWidth(width float64) {
	item.MaxWidth = width
	item.SetText(item.Raw)
}

func (item *Text) SetHeight(height float64) {
	item.MaxHeight = height
	item.SetText(item.Raw)
}

func (item *Text) SetColor(col color.RGBA) {
	item.Color = col
	item.updateText()
}

func (item *Text) SetSize(size float64) {
	item.RelativeSize = size
	item.SetText(item.Raw)
}

func (item *Text) SetPos(pos pixel.Vec) {
	item.Transform.Pos = pos
	item.updateText()
}

func (item *Text) IncrementTextPos() {
	if item.Increment {

	}
}

func (item *Text) SkipIncrement() {
	if item.Increment {

	}
}

func (item *Text) PrintLines() {
	for _, line := range item.rawLines {
		fmt.Println(line)
	}
}