package player

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
			Keys:    []pixelgl.Button{pixelgl.MouseButtonLeft, pixelgl.KeyLeftShift},
			Buttons: []pixelgl.GamepadButton{pixelgl.ButtonX},
			Axis:    pixelgl.AxisRightTrigger,
			AxisV:   1,
		},
		"mark": {
			Keys:    []pixelgl.Button{pixelgl.MouseButtonRight, pixelgl.KeyLeftControl},
			Buttons: []pixelgl.GamepadButton{pixelgl.ButtonY},
			Axis:    pixelgl.AxisLeftTrigger,
			AxisV:   1,
		},
		"left":  input.New(pixelgl.KeyA, pixelgl.ButtonDpadLeft),
		"right": input.New(pixelgl.KeyD, pixelgl.ButtonDpadRight),
		"up":    input.New(pixelgl.KeyW, pixelgl.ButtonDpadUp),
		"down":  input.New(pixelgl.KeyS, pixelgl.ButtonDpadDown),
		"jump":  input.New(pixelgl.KeySpace, pixelgl.ButtonA),
		"use":   input.New(pixelgl.KeyF, pixelgl.ButtonB),
		"prev": {
			Keys:    []pixelgl.Button{pixelgl.KeyQ},
			Buttons: []pixelgl.GamepadButton{pixelgl.ButtonLeftBumper},
			Scroll:  -1,
		},
		"next": {
			Keys:    []pixelgl.Button{pixelgl.KeyE},
			Buttons: []pixelgl.GamepadButton{pixelgl.ButtonRightBumper},
			Scroll:  1,
		},
	},
	StickD: true,
	Mode: input.KeyboardMouse,
}