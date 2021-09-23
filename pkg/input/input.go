package input

import (
	"dwarf-sweeper/pkg/camera"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var Deadzone = 0.25

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
	joyConn    bool
}

func (i *Input) Update(win *pixelgl.Window) {
	i.Cursor = win.MousePosition()
	i.World = camera.Cam.Mat.Unproject(win.MousePosition())
	i.ScrollV = win.MouseScroll().Y
	i.ScrollH = win.MouseScroll().X
	i.MouseMoved = !win.MousePreviousPosition().Eq(win.MousePosition())
	i.joyConn = win.JoystickPresent(i.Joystick)

	if i.joyConn {
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
		jp := false
		p := false
		jr := false
		if i.joyConn && !set.noJoy {
			if set.GP == 0 {
				jp = win.JoystickJustPressed(i.Joystick, set.GPKey)
				p = win.JoystickPressed(i.Joystick, set.GPKey)
				jr = win.JoystickJustReleased(i.Joystick, set.GPKey)
				if i.StickD {
					if set.GPKey == pixelgl.ButtonDpadLeft && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftX) < -Deadzone {
						jp = !set.Button.pressed
						p = true
						jr = false
					} else if set.GPKey == pixelgl.ButtonDpadRight && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftX) > Deadzone {
						jp = !set.Button.pressed
						p = true
						jr = false
					}
					if set.GPKey == pixelgl.ButtonDpadUp && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftY) < -Deadzone {
						jp = !set.Button.pressed
						p = true
						jr = false
					} else if set.GPKey == pixelgl.ButtonDpadDown && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftY) > Deadzone {
						jp = !set.Button.pressed
						p = true
						jr = false
					}
				}
			} else {
				if (win.JoystickAxis(i.Joystick, set.Axis) > Deadzone && set.GP > 0) || (win.JoystickAxis(i.Joystick, set.Axis) < -Deadzone && set.GP < 0) {
					jp = !set.Button.pressed
					p = true
					jr = false
				} else {
					jr = set.Button.pressed
					p = false
					jp = false
				}
			}
		}
		if set.Scroll == 0 {
			set.Button.justPressed = win.JustPressed(set.Key) || jp
			set.Button.pressed = win.Pressed(set.Key) || p
			set.Button.justReleased = win.JustReleased(set.Key) || jr
		} else {
			if (win.MouseScroll().Y > 0. && set.Scroll > 0) || (win.MouseScroll().Y < 0. && set.Scroll < 0) {
				set.Button.justPressed = !set.Button.pressed || jp
				set.Button.pressed = true
				set.Button.justReleased = jr
			} else {
				set.Button.justReleased = set.Button.pressed || jr
				set.Button.pressed = p
				set.Button.justPressed = jp
			}
		}
		set.Button.consumed = set.Button.consumed && (set.Button.justPressed || set.Button.pressed || set.Button.justReleased)
	}
}

func New(n pixelgl.Button, g pixelgl.GamepadButton) *ButtonSet {
	return &ButtonSet{
		Key:   n,
		GPKey: g,
	}
}

func NewJoyless(n pixelgl.Button) *ButtonSet {
	return &ButtonSet{
		Key:   n,
		noJoy: true,
	}
}

func (i *Input) SetStandard(s string, n pixelgl.Button) {
	if b, ok := i.Buttons[s]; ok {
		b.Key = n
	} else {
		i.Buttons[s] = &ButtonSet{
			Button: Button{},
			Key:    n,
		}
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
	Key    pixelgl.Button
	Scroll int
	GPKey  pixelgl.GamepadButton
	Axis   pixelgl.GamepadAxis
	GP     int
	noJoy  bool
}