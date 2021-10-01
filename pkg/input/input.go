package input

import (
	"dwarf-sweeper/pkg/camera"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var Deadzone = 0.25

type Mode int

const (
	Any = iota
	KeyboardMouse
	Gamepad
)

func (m Mode) String() string {
	switch m {
	case Any:
		return "Any"
	case KeyboardMouse:
		return "Keyboard & Mouse"
	case Gamepad:
		return "Gamepad"
	default:
		return ""
	}
}

type Input struct {
	Cursor     pixel.Vec
	World      pixel.Vec
	MouseMoved bool
	// todo: add mouse axes
	ScrollV    float64
	ScrollH    float64
	Axes       map[string]*AxisSet
	Buttons    map[string]*ButtonSet
	Joystick   pixelgl.Joystick
	StickD     bool
	Mode       Mode
	joyConn    bool
}

func (i *Input) Update(win *pixelgl.Window) {
	i.Cursor = win.MousePosition()
	i.World = camera.Cam.Mat.Unproject(win.MousePosition())
	i.ScrollV = win.MouseScroll().Y
	i.ScrollH = win.MouseScroll().X
	i.MouseMoved = !win.MousePreviousPosition().Eq(win.MousePosition())
	i.joyConn = win.JoystickPresent(i.Joystick)

	if i.joyConn && i.Mode != KeyboardMouse {
		for _, set := range i.Axes {
			f := win.JoystickAxis(i.Joystick, set.A)
			if f > Deadzone || f < -Deadzone {
				set.F = f
			} else {
				set.F = 0.
			}
		}
	}

	for _, set := range i.Buttons {
		wasPressed := set.Button.pressed
		nowPressed := false
		if i.joyConn && !set.noJoy && i.Mode != KeyboardMouse {
			for _, g := range set.GPKey {
				nowPressed = win.JoystickPressed(i.Joystick, g) || nowPressed
				if i.StickD {
					if g == pixelgl.ButtonDpadLeft && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftX) < -Deadzone {
						nowPressed = true
					} else if g == pixelgl.ButtonDpadRight && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftX) > Deadzone {
						nowPressed = true
					}
					if g == pixelgl.ButtonDpadUp && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftY) < -Deadzone {
						nowPressed = true
					} else if g == pixelgl.ButtonDpadDown && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftY) > Deadzone {
						nowPressed = true
					}
				}
			}
			if set.GP != 0 &&
				((win.JoystickAxis(i.Joystick, set.Axis) > Deadzone && set.GP > 0) ||
					(win.JoystickAxis(i.Joystick, set.Axis) < -Deadzone && set.GP < 0)) {
				nowPressed = true
			}
		}
		if i.Mode != Gamepad {
			for _, s := range set.Key {
				nowPressed = win.Pressed(s) || nowPressed
			}
			if set.Scroll != 0 {
				if (win.MouseScroll().Y > 0. && set.Scroll > 0) || (win.MouseScroll().Y < 0. && set.Scroll < 0) {
					set.Button.pressed = true
				}
			}
		}
		set.Button.justPressed = nowPressed && !wasPressed
		set.Button.pressed = nowPressed
		set.Button.justReleased = !nowPressed && wasPressed
		set.Button.consumed = set.Button.consumed && (set.Button.justPressed || set.Button.pressed || set.Button.justReleased)
	}
}

func New(n pixelgl.Button, g pixelgl.GamepadButton) *ButtonSet {
	return &ButtonSet{
		Key:   []pixelgl.Button{n},
		GPKey: []pixelgl.GamepadButton{g},
	}
}

func NewJoyless(n pixelgl.Button) *ButtonSet {
	return &ButtonSet{
		Key:   []pixelgl.Button{n},
		noJoy: true,
	}
}

func (i *Input) AnyJustPressed(consume bool) bool {
	for _, b := range i.Buttons {
		if b.Button.JustPressed() {
			if consume {
				b.Button.Consume()
			}
			return true
		}
	}
	return false
}

func (i *Input) Get(s string) *Button {
	if b, ok := i.Buttons[s]; ok {
		return &b.Button
	}
	return &Button{}
}

type Button struct {
	justPressed  bool
	pressed      bool
	justReleased bool
	consumed     bool
}

func (t *Button) JustPressed() bool {
	return t.justPressed && !t.consumed
}

func (t *Button) Pressed() bool {
	return t.pressed && !t.consumed
}

func (t *Button) JustReleased() bool {
	return t.justReleased && !t.consumed
}

func (t *Button) Consume() {
	t.consumed = true
}

type AxisSet struct {
	F float64
	A pixelgl.GamepadAxis
}

type ButtonSet struct {
	Button Button
	Key    []pixelgl.Button
	Scroll int
	GPKey  []pixelgl.GamepadButton
	Axis   pixelgl.GamepadAxis
	GP     int
	noJoy  bool
}