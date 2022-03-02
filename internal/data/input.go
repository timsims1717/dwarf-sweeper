package data

import (
	"dwarf-sweeper/pkg/input"
	"github.com/faiface/pixel/pixelgl"
)

var (
	GameInputP1 = &input.Input{
		Key:  "p1",
		Axes: map[string]*input.AxisSet{
			"targetX": {
				A: pixelgl.AxisRightX,
			},
			"targetY": {
				A: pixelgl.AxisRightY,
			},
		},
		Buttons: map[string]*input.ButtonSet{
			"pause": input.New(pixelgl.KeyEscape, pixelgl.ButtonStart),
		},
	}
	GameInputP2 = &input.Input{
		Key:  "p2",
		Axes: map[string]*input.AxisSet{
			"targetX": {
				A: pixelgl.AxisRightX,
			},
			"targetY": {
				A: pixelgl.AxisRightY,
			},
		},
		Buttons: map[string]*input.ButtonSet{
			"pause": input.New(pixelgl.KeyEscape, pixelgl.ButtonStart),
		},
	}
	GameInputP3 = &input.Input{
		Key:  "p3",
		Axes: map[string]*input.AxisSet{
			"targetX": {
				A: pixelgl.AxisRightX,
			},
			"targetY": {
				A: pixelgl.AxisRightY,
			},
		},
		Buttons: map[string]*input.ButtonSet{
			"pause": input.New(pixelgl.KeyEscape, pixelgl.ButtonStart),
		},
	}
	GameInputP4 = &input.Input{
		Key:  "p4",
		Axes: map[string]*input.AxisSet{
			"targetX": {
				A: pixelgl.AxisRightX,
			},
			"targetY": {
				A: pixelgl.AxisRightY,
			},
		},
		Buttons: map[string]*input.ButtonSet{
			"pause": input.New(pixelgl.KeyEscape, pixelgl.ButtonStart),
		},
	}
	CurrInput = GameInputP1
)
