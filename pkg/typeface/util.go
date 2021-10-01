package typeface

import "github.com/faiface/pixel/text"

func SetText(txt *text.Text, raw string, maxWidth float64) []string {
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
			b = e+1
		}
	}
	lines = append(lines, raw[b:])
	return lines
}