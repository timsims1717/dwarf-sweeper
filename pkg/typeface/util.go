package typeface

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

var (
	SymbolMarker = '^'
	SymbolItem   = "^"
)

func SetText(txt *text.Text, raw string, maxWidth float64, align Alignment) []pixel.Vec {
	var symbols []pixel.Vec
	b := 0
	e := 0
	cut := false
	space := false
	var syms []int
	for i, r := range raw {
		switch r {
		case '\n':
			cut = true
			e = i
		case ' ', '\t':
			space = true
			e = i
		case SymbolMarker:
			syms = append(syms, i)
			continue
		}
		if maxWidth > 0. && txt.BoundsOf(raw[b:i]).W() > maxWidth && space {
			cut = true
		}
		if cut {
			if b >= e || e < 0 {
				fmt.Fprintln(txt, "")
			} else {
				if align.H == Center {
					txt.Dot.X -= txt.BoundsOf(raw[b:e]).W() * 0.5
				} else if align.H == Right {
					txt.Dot.X -= txt.BoundsOf(raw[b:e]).W()
				}
				if len(syms) > 0 {
					nb := b
					for j := 0; j < len(syms); j++ {
						if syms[j] >= nb {
							fmt.Fprint(txt, raw[nb:syms[j]])
							symbols = append(symbols, pixel.V(txt.Dot.X+txt.BoundsOf(SymbolItem).W()*0.5, txt.Dot.Y))
							txt.Dot.X += txt.BoundsOf(SymbolItem).W()
							nb = syms[j] + 1
						}
					}
					fmt.Fprintf(txt, "%s\n", raw[nb:e])
				} else {
					fmt.Fprintf(txt, "%s\n", raw[b:e])
				}
			}
			cut = false
			space = false
			b = e + 1
			syms = []int{}
		}
	}
	if align.H == Center {
		txt.Dot.X -= txt.BoundsOf(raw[b:]).W() * 0.5
	} else if align.H == Right {
		txt.Dot.X -= txt.BoundsOf(raw[b:]).W()
	}
	if len(syms) > 0 {
		nb := b
		for j := 0; j < len(syms); j++ {
			if syms[j] >= nb {
				fmt.Fprint(txt, raw[nb:syms[j]])
				symbols = append(symbols, pixel.V(txt.Dot.X+txt.BoundsOf(SymbolItem).W()*0.5, txt.Dot.Y))
				txt.Dot.X += txt.BoundsOf(SymbolItem).W()
				nb = syms[j] + 1
			}
		}
		fmt.Fprintln(txt, raw[nb:])
	} else {
		fmt.Fprintln(txt, raw[b:])
	}
	return symbols
}

func RawLines(txt *text.Text, raw string, maxWidth float64) []string {
	var lines []string
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
		if maxWidth > 0. && txt.BoundsOf(raw[b:i]).W() > maxWidth && space {
			cut = true
		}
		if cut {
			if b >= e || e < 0 {
				lines = append(lines, "")
			} else {
				lines = append(lines, raw[b:e])
			}
			cut = false
			space = false
			b = e + 1
		}
	}
	lines = append(lines, raw[b:])
	return lines
}
