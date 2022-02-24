package data

import (
	"dwarf-sweeper/pkg/input"
	"github.com/faiface/pixel/pixelgl"
)

var (
	GameInputP1 = &input.Input{
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
)
