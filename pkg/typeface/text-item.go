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

var (
	padding = 5.
	padOrig = pixel.V(padding, padding)
)

type Item struct {
	Raw     string
	Text    *text.Text
	Color   color.RGBA
	Align   Alignment
	Symbols []symbolHandle

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

func New(cam *camera.Camera, atlas string, align Alignment, lineHeight, relativeSize, maxWidth, maxHeight float64) *Item {
	tex := text.New(pixel.ZV, Atlases[atlas])
	tex.LineHeight *= lineHeight
	return &Item{
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

func (item *Item) Update() {
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

func (item *Item) Draw(target pixel.Target) {
	item.Text.Draw(target, item.Transform.Mat)
	for _, sym := range item.Symbols {
		sym.symbol.spr.Draw(target, sym.trans.Mat)
	}
}

func (item *Item) SetWidth(width float64) {
	item.MaxWidth = width
	item.updateText()
}

func (item *Item) SetHeight(height float64) {
	item.MaxHeight = height
	item.updateText()
}

func (item *Item) SetColor(col color.RGBA) {
	item.Color = col
}

func (item *Item) SetPos(pos pixel.Vec) {
	item.Transform.Pos = pos
	item.updateText()
}

func (item *Item) IncrementTextPos() {
	if item.Increment {

	}
}

func (item *Item) SkipIncrement() {
	if item.Increment {

	}
}

func (item *Item) PrintLines() {
	for _, line := range item.rawLines {
		fmt.Println(line)
	}
}