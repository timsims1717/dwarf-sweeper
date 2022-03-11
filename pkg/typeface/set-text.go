package typeface

import (
	"bytes"
	"dwarf-sweeper/pkg/transform"
	"fmt"
	"github.com/faiface/pixel"
	"math"
)

const (
	OpenMarker  = '{'
	CloseMarker = '}'
	DivMarker   = ':'
	OpenItem    = "{"
	CloseItem   = "}"
)

var (
	RoundDot = false
)

func (item *Text) SetText(raw string) {
	item.Raw = raw
	item.rawLines = []string{}
	item.lineWidths = []float64{}
	b := 0
	bb := 0
	e := 0
	cut := false
	wSpace := false
	inBrackets := false
	widthMod := 0.
	maxLineWidth := 0.
	mode := ""
	buf := bytes.NewBuffer(nil)
	for i, r := range item.Raw {
		if !inBrackets {
			switch r {
			case '\n':
				cut = true
				e = i
			case ' ', '\t':
				wSpace = true
				e = i
			case OpenMarker:
				inBrackets = true
				bb = i
				continue
			case CloseMarker:
				fmt.Printf("extra closing bracket in text at position %d\n", i)
			}
			lineWidth := item.Text.BoundsOf(item.Raw[b:i]).W() - widthMod
			lineWidthRelative := lineWidth * item.RelativeSize
			if item.MaxWidth > 0. && lineWidthRelative > item.MaxWidth && wSpace {
				cut = true
			}
			if cut {
				if b >= e || e < 0 {
					item.rawLines = append(item.rawLines, "")
					item.lineWidths = append(item.lineWidths, 0.)
				} else {
					item.rawLines = append(item.rawLines, raw[b:e])
					item.lineWidths = append(item.lineWidths, lineWidth)
					if maxLineWidth < lineWidthRelative {
						maxLineWidth = lineWidthRelative
					}
				}
				cut = false
				wSpace = false
				widthMod = 0.
				b = e + 1
			}
		} else {
			switch r {
			case '\n':
				fmt.Printf("new line in bracketed text at position %d\n", i)
			case ' ', '\t':
				continue
			case OpenMarker:
				fmt.Printf("extra opening bracket at position %d\n", i)
				continue
			case CloseMarker:
				switch mode {
				case "symbol":
					if sym, ok := theSymbols[buf.String()]; ok {
						widthMod -= sym.spr.Frame().W() * item.SymbolSize * sym.sca / item.RelativeSize
					}
				}
				widthMod += item.Text.BoundsOf(item.Raw[bb:i+1]).W()
				inBrackets = false
				mode = ""
				buf.Reset()
			case DivMarker:
				mode = buf.String()
				buf.Reset()
			default:
				buf.WriteRune(r)
			}
		}
	}
	item.rawLines = append(item.rawLines, raw[b:])
	lineWidth := item.Text.BoundsOf(item.Raw[b:]).W() - widthMod
	lineWidthRelative := lineWidth * item.RelativeSize
	item.lineWidths = append(item.lineWidths, lineWidth)
	item.fullHeight = float64(len(item.rawLines)) * item.Text.LineHeight
	if maxLineWidth < lineWidthRelative {
		maxLineWidth = lineWidthRelative
	}
	maxX := maxLineWidth
	maxY := item.MaxHeight
	if maxY == 0. {
		maxY = item.fullHeight * item.RelativeSize
	}
	item.Width = maxX
	item.Height = maxY
	item.updateText()
}

func (item *Text) updateText() {
	item.Text.Clear()
	item.Text.Color = item.Color
	//var colorStack []color.RGBA
	item.Symbols = []symbolHandle{}
	inBrackets := false
	mode := ""
	buf := bytes.NewBuffer(nil)
	item.Text.Dot.Y -= item.Text.LineHeight
	if item.Align.V == Center {
		item.Text.Dot.Y += item.fullHeight * 0.5
	} else if item.Align.V == Bottom {
		item.Text.Dot.Y += item.fullHeight
	}
	for li, line := range item.rawLines {
		b := 0
		inBrackets = false
		if item.Align.H == Center {
			item.Text.Dot.X -= item.lineWidths[li] * 0.5
		} else if item.Align.H == Right {
			item.Text.Dot.X -= item.lineWidths[li]
		}
		for i, r := range line {
			if !inBrackets {
				switch r {
				case OpenMarker:
					item.roundDot()
					fmt.Fprintf(item.Text, "%s", line[b:i])
					inBrackets = true
				}
			} else {
				switch r {
				case CloseMarker:
					switch mode {
					case "symbol":
						if sym, ok := theSymbols[buf.String()]; ok {
							item.roundDot()
							trans := transform.New()
							trans.Scalar = pixel.V(item.SymbolSize, item.SymbolSize).Scaled(sym.sca)
							trans.Pos = item.Transform.Pos
							trans.Pos = trans.Pos.Add(item.Text.Dot.Scaled(item.RelativeSize))
							trans.Pos = trans.Pos.Add(pixel.V(sym.spr.Frame().W() * 0.5, sym.spr.Frame().H() * 0.25).Scaled(item.SymbolSize * sym.sca))
							item.Symbols = append(item.Symbols, symbolHandle{
								symbol: sym,
								trans:  trans,
							})
							item.Text.Dot.X += sym.spr.Frame().W() * item.SymbolSize * sym.sca / item.RelativeSize
						}
					}
					b = i+1
					inBrackets = false
					mode = ""
					buf.Reset()
				case DivMarker:
					mode = buf.String()
					buf.Reset()
				default:
					buf.WriteRune(r)
				}
			}
		}
		item.roundDot()
		fmt.Fprintf(item.Text, "%s\n", line[b:])
	}
}

func (item *Text) roundDot() {
	if RoundDot {
		item.Text.Dot.X = math.Floor(item.Text.Dot.X)
		item.Text.Dot.Y = math.Floor(item.Text.Dot.Y)
	}
}