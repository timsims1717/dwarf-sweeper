package config

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/pkg/input"
	"fmt"
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
		PuzzLeave:    input.New(pixelgl.KeyE, pixelgl.ButtonB),
		PuzzHelp:     input.New(pixelgl.KeyQ, pixelgl.ButtonY),
		MinePuzzBomb: input.New(pixelgl.MouseButtonRight, pixelgl.ButtonA),
		MinePuzzSafe: input.New(pixelgl.MouseButtonLeft, pixelgl.ButtonX),
	}
)

//goland:noinspection GoNilness
func loadInput(conf *config) {
	for i := 0; i < 4; i++ {
		var in *input.Input
		var cf inputs
		switch i {
		case 0:
			in = data.GameInputP1
			cf = conf.InputP1
		case 1:
			in = data.GameInputP2
			cf = conf.InputP2
		case 2:
			in = data.GameInputP3
			cf = conf.InputP3
		case 3:
			in = data.GameInputP4
			cf = conf.InputP4
		}
		if cf.Gamepad < 0 {
			in.Mode = input.KeyboardMouse
		} else {
			in.Mode = input.Gamepad
			in.Joystick = pixelgl.Joystick(cf.Gamepad)
		}
		in.AimDedicated = cf.AimDedicated
		in.DigOnRelease = cf.DigOnRelease
		in.StickD = cf.LeftStick
		in.Deadzone = cf.Deadzone
		if cf.Left != nil {
			in.Buttons["left"] = cf.Left
		} else {
			in.Buttons["left"] = DefaultInput.Left
		}
		if cf.Right != nil {
			in.Buttons["right"] = cf.Right
		} else {
			in.Buttons["right"] = DefaultInput.Right
		}
		if cf.Up != nil {
			in.Buttons["up"] = cf.Up
		} else {
			in.Buttons["up"] = DefaultInput.Up
		}
		if cf.Down != nil {
			in.Buttons["down"] = cf.Down
		} else {
			in.Buttons["down"] = DefaultInput.Down
		}
		if cf.Jump != nil {
			in.Buttons["jump"] = cf.Jump
		} else {
			in.Buttons["jump"] = DefaultInput.Jump
		}
		if cf.Dig != nil {
			in.Buttons["dig"] = cf.Dig
		} else {
			in.Buttons["dig"] = DefaultInput.Dig
		}
		if cf.Flag != nil {
			in.Buttons["flag"] = cf.Flag
		} else {
			in.Buttons["flag"] = DefaultInput.Flag
		}
		if cf.Use != nil {
			in.Buttons["use"] = cf.Use
		} else {
			in.Buttons["use"] = DefaultInput.Use
		}
		if cf.Interact != nil {
			in.Buttons["interact"] = cf.Interact
		} else {
			in.Buttons["interact"] = DefaultInput.Interact
		}
		if cf.Prev != nil {
			in.Buttons["prev"] = cf.Prev
		} else {
			in.Buttons["prev"] = DefaultInput.Prev
		}
		if cf.Next != nil {
			in.Buttons["next"] = cf.Next
		} else {
			in.Buttons["next"] = DefaultInput.Next
		}
		if cf.PuzzLeave != nil {
			in.Buttons["puzz_leave"] = cf.PuzzLeave
		} else {
			in.Buttons["puzz_leave"] = DefaultInput.PuzzLeave
		}
		if cf.PuzzHelp != nil {
			in.Buttons["puzz_help"] = cf.PuzzHelp
		} else {
			in.Buttons["puzz_help"] = DefaultInput.PuzzHelp
		}
		if cf.MinePuzzBomb != nil {
			in.Buttons["mine_puzz_bomb"] = cf.MinePuzzBomb
		} else {
			in.Buttons["mine_puzz_bomb"] = DefaultInput.MinePuzzBomb
		}
		if cf.MinePuzzSafe != nil {
			in.Buttons["mine_puzz_safe"] = cf.MinePuzzSafe
		} else {
			in.Buttons["mine_puzz_safe"] = DefaultInput.MinePuzzSafe
		}
		in.Key = fmt.Sprintf("p%d", i+1)
	}
}

//goland:noinspection GoNilness
func saveInput(conf *config) {
	for i := 0; i < 4; i++ {
		var in *input.Input
		var cf inputs
		switch i {
		case 0:
			in = data.GameInputP1
		case 1:
			in = data.GameInputP2
		case 2:
			in = data.GameInputP3
		case 3:
			in = data.GameInputP4
		}

		if in.Mode == input.KeyboardMouse {
			cf.Gamepad = -1
		} else {
			cf.Gamepad = int(in.Joystick)
		}
		cf.AimDedicated = in.AimDedicated
		cf.DigOnRelease = in.DigOnRelease
		cf.LeftStick = in.StickD
		cf.Deadzone = in.Deadzone
		cf.Left = in.Buttons["left"]
		cf.Right = in.Buttons["right"]
		cf.Up = in.Buttons["up"]
		cf.Down = in.Buttons["down"]
		cf.Jump = in.Buttons["jump"]
		cf.Dig = in.Buttons["dig"]
		cf.Flag = in.Buttons["flag"]
		cf.Use = in.Buttons["use"]
		cf.Interact = in.Buttons["interact"]
		cf.Prev = in.Buttons["prev"]
		cf.Next = in.Buttons["next"]
		cf.PuzzLeave = in.Buttons["puzz_leave"]
		cf.PuzzHelp = in.Buttons["puzz_help"]
		cf.MinePuzzBomb = in.Buttons["mine_puzz_bomb"]
		cf.MinePuzzSafe = in.Buttons["mine_puzz_safe"]
		cf.Key = fmt.Sprintf("p%d", i+1)

		switch i {
		case 0:
			conf.InputP1 = cf
		case 1:
			conf.InputP2 = cf
		case 2:
			conf.InputP3 = cf
		case 3:
			conf.InputP4 = cf
		}
	}
}