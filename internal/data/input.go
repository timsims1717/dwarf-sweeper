package data

import (
	"dwarf-sweeper/pkg/input"
	"github.com/faiface/pixel/pixelgl"
)

var GameInput = &input.Input{
	Axes: map[string]*input.AxisSet{
		"targetX": {
			A: pixelgl.AxisRightX,
		},
		"targetY": {
			A: pixelgl.AxisRightY,
		},
	},
	Buttons: map[string]*input.ButtonSet{
		"dig": {
			Keys:    []pixelgl.Button{pixelgl.MouseButtonLeft},
			Buttons: []pixelgl.GamepadButton{pixelgl.ButtonX},
			Axis:    pixelgl.AxisRightTrigger,
			AxisV:   1,
		},
		"flag": {
			Keys:  []pixelgl.Button{pixelgl.MouseButtonRight},
			Axis:  pixelgl.AxisLeftTrigger,
			AxisV: 1,
		},
		"left":   input.New(pixelgl.KeyA, pixelgl.ButtonDpadLeft),
		"right":  input.New(pixelgl.KeyD, pixelgl.ButtonDpadRight),
		"up":     input.New(pixelgl.KeyW, pixelgl.ButtonDpadUp),
		"down":   input.New(pixelgl.KeyS, pixelgl.ButtonDpadDown),
		"jump":   input.New(pixelgl.KeySpace, pixelgl.ButtonA),
		"pickUp": input.New(pixelgl.KeyQ, pixelgl.ButtonY),
		"use":    input.New(pixelgl.KeyE, pixelgl.ButtonB),
		"prev": {
			Buttons: []pixelgl.GamepadButton{pixelgl.ButtonLeftBumper},
			Scroll:  -1,
		},
		"next": {
			Buttons: []pixelgl.GamepadButton{pixelgl.ButtonRightBumper},
			Scroll:  1,
		},
	},
	StickD: true,
	Mode:   input.KeyboardMouse,
}
