package config

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/pkg/input"
	"github.com/faiface/pixel/pixelgl"
)

var (
	DefaultInput = inputs{
		Gamepad:      -1,
		AimDedicated: true,
		DigOnRelease: false,
		Deadzone:     0.25,
		LeftStick:    true,
		Left:         input.New(pixelgl.KeyA, pixelgl.ButtonDpadLeft),
		Right:        input.New(pixelgl.KeyD, pixelgl.ButtonDpadRight),
		Up:           input.New(pixelgl.KeyW, pixelgl.ButtonDpadUp),
		Down:         input.New(pixelgl.KeyS, pixelgl.ButtonDpadDown),
		Jump:         input.New(pixelgl.KeySpace, pixelgl.ButtonA),
		Dig: &input.ButtonSet{
			Keys:    []pixelgl.Button{pixelgl.MouseButtonLeft},
			Buttons: []pixelgl.GamepadButton{pixelgl.ButtonX},
			Axis:    pixelgl.AxisRightTrigger,
			AxisV:   1,
		},
		Flag: &input.ButtonSet{
			Keys:  []pixelgl.Button{pixelgl.MouseButtonRight},
			Axis:  pixelgl.AxisLeftTrigger,
			AxisV: 1,
		},
		Use:      input.New(pixelgl.KeyE, pixelgl.ButtonB),
		Interact: input.New(pixelgl.KeyQ, pixelgl.ButtonY),
		Prev: &input.ButtonSet{
			Buttons: []pixelgl.GamepadButton{pixelgl.ButtonLeftBumper},
			Scroll:  -1,
		},
		Next: &input.ButtonSet{
			Buttons: []pixelgl.GamepadButton{pixelgl.ButtonRightBumper},
			Scroll:  1,
		},
	}
)

func loadInput(conf *config) {
	if conf.InputP1.Gamepad < 0 {
		data.GameInputP1.Mode = input.KeyboardMouse
	} else {
		data.GameInputP1.Mode = input.Gamepad
		data.GameInputP1.Joystick = pixelgl.Joystick(conf.InputP1.Gamepad)
	}
	data.GameInputP1.AimDedicated = conf.InputP1.AimDedicated
	data.GameInputP1.DigOnRelease = conf.InputP1.DigOnRelease
	data.GameInputP1.StickD = conf.InputP1.LeftStick
	data.GameInputP1.Deadzone = conf.InputP1.Deadzone
	data.GameInputP1.Buttons["left"] = conf.InputP1.Left
	data.GameInputP1.Buttons["right"] = conf.InputP1.Right
	data.GameInputP1.Buttons["up"] = conf.InputP1.Up
	data.GameInputP1.Buttons["down"] = conf.InputP1.Down
	data.GameInputP1.Buttons["jump"] = conf.InputP1.Jump
	data.GameInputP1.Buttons["dig"] = conf.InputP1.Dig
	data.GameInputP1.Buttons["flag"] = conf.InputP1.Flag
	data.GameInputP1.Buttons["use"] = conf.InputP1.Use
	data.GameInputP1.Buttons["interact"] = conf.InputP1.Interact
	data.GameInputP1.Buttons["prev"] = conf.InputP1.Prev
	data.GameInputP1.Buttons["next"] = conf.InputP1.Next
	data.GameInputP1.Key = "p1"

	if conf.InputP2.Gamepad < 0 {
		data.GameInputP2.Mode = input.KeyboardMouse
	} else {
		data.GameInputP2.Mode = input.Gamepad
		data.GameInputP2.Joystick = pixelgl.Joystick(conf.InputP2.Gamepad)
	}
	data.GameInputP2.AimDedicated = conf.InputP2.AimDedicated
	data.GameInputP2.DigOnRelease = conf.InputP2.DigOnRelease
	data.GameInputP2.StickD = conf.InputP2.LeftStick
	data.GameInputP2.Deadzone = conf.InputP2.Deadzone
	data.GameInputP2.Buttons["left"] = conf.InputP2.Left
	data.GameInputP2.Buttons["right"] = conf.InputP2.Right
	data.GameInputP2.Buttons["up"] = conf.InputP2.Up
	data.GameInputP2.Buttons["down"] = conf.InputP2.Down
	data.GameInputP2.Buttons["jump"] = conf.InputP2.Jump
	data.GameInputP2.Buttons["dig"] = conf.InputP2.Dig
	data.GameInputP2.Buttons["flag"] = conf.InputP2.Flag
	data.GameInputP2.Buttons["use"] = conf.InputP2.Use
	data.GameInputP2.Buttons["interact"] = conf.InputP2.Interact
	data.GameInputP2.Buttons["prev"] = conf.InputP2.Prev
	data.GameInputP2.Buttons["next"] = conf.InputP2.Next
	data.GameInputP2.Key = "p2"

	if conf.InputP3.Gamepad < 0 {
		data.GameInputP3.Mode = input.KeyboardMouse
	} else {
		data.GameInputP3.Mode = input.Gamepad
		data.GameInputP3.Joystick = pixelgl.Joystick(conf.InputP3.Gamepad)
	}
	data.GameInputP3.AimDedicated = conf.InputP3.AimDedicated
	data.GameInputP3.DigOnRelease = conf.InputP3.DigOnRelease
	data.GameInputP3.StickD = conf.InputP3.LeftStick
	data.GameInputP3.Deadzone = conf.InputP3.Deadzone
	data.GameInputP3.Buttons["left"] = conf.InputP3.Left
	data.GameInputP3.Buttons["right"] = conf.InputP3.Right
	data.GameInputP3.Buttons["up"] = conf.InputP3.Up
	data.GameInputP3.Buttons["down"] = conf.InputP3.Down
	data.GameInputP3.Buttons["jump"] = conf.InputP3.Jump
	data.GameInputP3.Buttons["dig"] = conf.InputP3.Dig
	data.GameInputP3.Buttons["flag"] = conf.InputP3.Flag
	data.GameInputP3.Buttons["use"] = conf.InputP3.Use
	data.GameInputP3.Buttons["interact"] = conf.InputP3.Interact
	data.GameInputP3.Buttons["prev"] = conf.InputP3.Prev
	data.GameInputP3.Buttons["next"] = conf.InputP3.Next
	data.GameInputP3.Key = "p3"

	if conf.InputP4.Gamepad < 0 {
		data.GameInputP4.Mode = input.KeyboardMouse
	} else {
		data.GameInputP4.Mode = input.Gamepad
		data.GameInputP4.Joystick = pixelgl.Joystick(conf.InputP4.Gamepad)
	}
	data.GameInputP4.AimDedicated = conf.InputP4.AimDedicated
	data.GameInputP4.DigOnRelease = conf.InputP4.DigOnRelease
	data.GameInputP4.StickD = conf.InputP4.LeftStick
	data.GameInputP4.Deadzone = conf.InputP4.Deadzone
	data.GameInputP4.Buttons["left"] = conf.InputP4.Left
	data.GameInputP4.Buttons["right"] = conf.InputP4.Right
	data.GameInputP4.Buttons["up"] = conf.InputP4.Up
	data.GameInputP4.Buttons["down"] = conf.InputP4.Down
	data.GameInputP4.Buttons["jump"] = conf.InputP4.Jump
	data.GameInputP4.Buttons["dig"] = conf.InputP4.Dig
	data.GameInputP4.Buttons["flag"] = conf.InputP4.Flag
	data.GameInputP4.Buttons["use"] = conf.InputP4.Use
	data.GameInputP4.Buttons["interact"] = conf.InputP4.Interact
	data.GameInputP4.Buttons["prev"] = conf.InputP4.Prev
	data.GameInputP4.Buttons["next"] = conf.InputP4.Next
	data.GameInputP4.Key = "p4"
}

func saveInput(conf *config) {
	if data.GameInputP1.Mode == input.KeyboardMouse {
		conf.InputP1.Gamepad = -1
	} else {
		conf.InputP1.Gamepad = int(data.GameInputP1.Joystick)
	}
	conf.InputP1.AimDedicated = data.GameInputP1.AimDedicated
	conf.InputP1.DigOnRelease = data.GameInputP1.DigOnRelease
	conf.InputP1.LeftStick = data.GameInputP1.StickD
	conf.InputP1.Deadzone = data.GameInputP1.Deadzone
	conf.InputP1.Left = data.GameInputP1.Buttons["left"]
	conf.InputP1.Right = data.GameInputP1.Buttons["right"]
	conf.InputP1.Up = data.GameInputP1.Buttons["up"]
	conf.InputP1.Down = data.GameInputP1.Buttons["down"]
	conf.InputP1.Jump = data.GameInputP1.Buttons["jump"]
	conf.InputP1.Dig = data.GameInputP1.Buttons["dig"]
	conf.InputP1.Flag = data.GameInputP1.Buttons["flag"]
	conf.InputP1.Use = data.GameInputP1.Buttons["use"]
	conf.InputP1.Interact = data.GameInputP1.Buttons["interact"]
	conf.InputP1.Prev = data.GameInputP1.Buttons["prev"]
	conf.InputP1.Next = data.GameInputP1.Buttons["next"]
	conf.InputP1.Key = "p1"

	if data.GameInputP2.Mode == input.KeyboardMouse {
		conf.InputP2.Gamepad = -1
	} else {
		conf.InputP2.Gamepad = int(data.GameInputP2.Joystick)
	}
	conf.InputP2.AimDedicated = data.GameInputP2.AimDedicated
	conf.InputP2.DigOnRelease = data.GameInputP2.DigOnRelease
	conf.InputP2.LeftStick = data.GameInputP2.StickD
	conf.InputP2.Deadzone = data.GameInputP2.Deadzone
	conf.InputP2.Left = data.GameInputP2.Buttons["left"]
	conf.InputP2.Right = data.GameInputP2.Buttons["right"]
	conf.InputP2.Up = data.GameInputP2.Buttons["up"]
	conf.InputP2.Down = data.GameInputP2.Buttons["down"]
	conf.InputP2.Jump = data.GameInputP2.Buttons["jump"]
	conf.InputP2.Dig = data.GameInputP2.Buttons["dig"]
	conf.InputP2.Flag = data.GameInputP2.Buttons["flag"]
	conf.InputP2.Use = data.GameInputP2.Buttons["use"]
	conf.InputP2.Interact = data.GameInputP2.Buttons["interact"]
	conf.InputP2.Prev = data.GameInputP2.Buttons["prev"]
	conf.InputP2.Next = data.GameInputP2.Buttons["next"]
	conf.InputP2.Key = "p2"

	if data.GameInputP3.Mode == input.KeyboardMouse {
		conf.InputP3.Gamepad = -1
	} else {
		conf.InputP3.Gamepad = int(data.GameInputP3.Joystick)
	}
	conf.InputP3.AimDedicated = data.GameInputP3.AimDedicated
	conf.InputP3.DigOnRelease = data.GameInputP3.DigOnRelease
	conf.InputP3.LeftStick = data.GameInputP3.StickD
	conf.InputP3.Deadzone = data.GameInputP3.Deadzone
	conf.InputP3.Left = data.GameInputP3.Buttons["left"]
	conf.InputP3.Right = data.GameInputP3.Buttons["right"]
	conf.InputP3.Up = data.GameInputP3.Buttons["up"]
	conf.InputP3.Down = data.GameInputP3.Buttons["down"]
	conf.InputP3.Jump = data.GameInputP3.Buttons["jump"]
	conf.InputP3.Dig = data.GameInputP3.Buttons["dig"]
	conf.InputP3.Flag = data.GameInputP3.Buttons["flag"]
	conf.InputP3.Use = data.GameInputP3.Buttons["use"]
	conf.InputP3.Interact = data.GameInputP3.Buttons["interact"]
	conf.InputP3.Prev = data.GameInputP3.Buttons["prev"]
	conf.InputP3.Next = data.GameInputP3.Buttons["next"]
	conf.InputP3.Key = "p3"

	if data.GameInputP4.Mode == input.KeyboardMouse {
		conf.InputP4.Gamepad = -1
	} else {
		conf.InputP4.Gamepad = int(data.GameInputP4.Joystick)
	}
	conf.InputP4.AimDedicated = data.GameInputP4.AimDedicated
	conf.InputP4.DigOnRelease = data.GameInputP4.DigOnRelease
	conf.InputP4.LeftStick = data.GameInputP4.StickD
	conf.InputP4.Deadzone = data.GameInputP4.Deadzone
	conf.InputP4.Left = data.GameInputP4.Buttons["left"]
	conf.InputP4.Right = data.GameInputP4.Buttons["right"]
	conf.InputP4.Up = data.GameInputP4.Buttons["up"]
	conf.InputP4.Down = data.GameInputP4.Buttons["down"]
	conf.InputP4.Jump = data.GameInputP4.Buttons["jump"]
	conf.InputP4.Dig = data.GameInputP4.Buttons["dig"]
	conf.InputP4.Flag = data.GameInputP4.Buttons["flag"]
	conf.InputP4.Use = data.GameInputP4.Buttons["use"]
	conf.InputP4.Interact = data.GameInputP4.Buttons["interact"]
	conf.InputP4.Prev = data.GameInputP4.Buttons["prev"]
	conf.InputP4.Next = data.GameInputP4.Buttons["next"]
	conf.InputP4.Key = "p4"
}