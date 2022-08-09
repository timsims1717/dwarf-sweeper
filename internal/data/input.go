package data

import (
	"github.com/faiface/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
)

var (
	GameInputP1 = &pxginput.Input{
		Key:  "p1",
		Axes: map[string]*pxginput.AxisSet{
			"targetX": {
				A: pixelgl.AxisRightX,
			},
			"targetY": {
				A: pixelgl.AxisRightY,
			},
		},
		Buttons: map[string]*pxginput.ButtonSet{
			"pause": pxginput.New(pixelgl.KeyEscape, pixelgl.ButtonStart),
		},
	}
	GameInputP2 = &pxginput.Input{
		Key:  "p2",
		Axes: map[string]*pxginput.AxisSet{
			"targetX": {
				A: pixelgl.AxisRightX,
			},
			"targetY": {
				A: pixelgl.AxisRightY,
			},
		},
		Buttons: map[string]*pxginput.ButtonSet{
			"pause": pxginput.New(pixelgl.KeyEscape, pixelgl.ButtonStart),
		},
	}
	GameInputP3 = &pxginput.Input{
		Key:  "p3",
		Axes: map[string]*pxginput.AxisSet{
			"targetX": {
				A: pixelgl.AxisRightX,
			},
			"targetY": {
				A: pixelgl.AxisRightY,
			},
		},
		Buttons: map[string]*pxginput.ButtonSet{
			"pause": pxginput.New(pixelgl.KeyEscape, pixelgl.ButtonStart),
		},
	}
	GameInputP4 = &pxginput.Input{
		Key:  "p4",
		Axes: map[string]*pxginput.AxisSet{
			"targetX": {
				A: pixelgl.AxisRightX,
			},
			"targetY": {
				A: pixelgl.AxisRightY,
			},
		},
		Buttons: map[string]*pxginput.ButtonSet{
			"pause": pxginput.New(pixelgl.KeyEscape, pixelgl.ButtonStart),
		},
	}
	CurrInput = GameInputP1
)
