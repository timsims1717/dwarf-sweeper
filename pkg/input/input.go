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
	updateConsume(win, i.Joystick)
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
		modePressed := Any
		repeated := false
		if i.joyConn && !set.noJoy && i.Mode != KeyboardMouse {
			for _, g := range set.Buttons {
				if c, ok := consumeGamepad[g]; !ok || !c {
					nowPressed = win.JoystickPressed(i.Joystick, g) || nowPressed
					if win.JoystickPressed(i.Joystick, g) {
						modePressed = Gamepad
					}
					if i.StickD {
						if g == pixelgl.ButtonDpadLeft && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftX) < -Deadzone {
							nowPressed = true
							modePressed = Gamepad
						} else if g == pixelgl.ButtonDpadRight && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftX) > Deadzone {
							nowPressed = true
							modePressed = Gamepad
						}
						if g == pixelgl.ButtonDpadUp && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftY) < -Deadzone {
							nowPressed = true
							modePressed = Gamepad
						} else if g == pixelgl.ButtonDpadDown && win.JoystickAxis(i.Joystick, pixelgl.AxisLeftY) > Deadzone {
							nowPressed = true
							modePressed = Gamepad
						}
					}
				}
			}
			if set.AxisV != 0 &&
				((win.JoystickAxis(i.Joystick, set.Axis) > Deadzone && set.AxisV > 0) ||
					(win.JoystickAxis(i.Joystick, set.Axis) < -Deadzone && set.AxisV < 0)) {
				nowPressed = true
				modePressed = Gamepad
			}
		}
		if i.Mode != Gamepad {
			for _, s := range set.Keys {
				if c, ok := consumeKey[s]; !ok || !c {
					nowPressed = win.Pressed(s) || nowPressed
					repeated = win.Repeated(s) || repeated
					if win.Pressed(s) {
						modePressed = KeyboardMouse
					}
				}
			}
			if set.Scroll != 0 {
				if (win.MouseScroll().Y > 0. && set.Scroll > 0) || (win.MouseScroll().Y < 0. && set.Scroll < 0) {
					nowPressed = true
					modePressed = KeyboardMouse
				}
			}
		}
		set.Button.justPressed = nowPressed && !wasPressed
		set.Button.pressed = nowPressed
		set.Button.justReleased = !nowPressed && wasPressed
		set.Button.repeated = repeated
		set.LastMode = Mode(modePressed)
	}
}

func New(n pixelgl.Button, g pixelgl.GamepadButton) *ButtonSet {
	return &ButtonSet{
		Keys:    []pixelgl.Button{n},
		Buttons: []pixelgl.GamepadButton{g},
	}
}

func NewJoyless(n pixelgl.Button) *ButtonSet {
	return &ButtonSet{
		Keys:  []pixelgl.Button{n},
		noJoy: true,
	}
}

func (in *Input) FirstKey(s string) string {
	if in.Get(s) == nil {
		return ""
	}
	if in.Mode != Gamepad && len(in.Get(s).Keys) > 0{
		return in.Get(s).Keys[0].String()
	} else if in.Mode != KeyboardMouse && len(in.Get(s).Buttons) > 0 {
		return GamepadString(in.Get("up").Buttons[0])
	}
	return ""
}

func (i *Input) AnyJustPressed(consume bool) (bool, Mode) {
	for _, b := range i.Buttons {
		if b.JustPressed() {
			if consume {
				b.Consume()
			}
			return true, b.LastMode
		}
	}
	return false, Any
}

func (i *Input) Get(s string) *ButtonSet {
	if b, ok := i.Buttons[s]; ok {
		return b
	}
	return nil
}

type Button struct {
	justPressed  bool
	pressed      bool
	justReleased bool
	repeated     bool
}

func (bs *ButtonSet) JustPressed() bool {
	if bs == nil {
		return false
	}
	return bs.Button.justPressed
}

func (bs *ButtonSet) Pressed() bool {
	if bs == nil {
		return false
	}
	return bs.Button.pressed
}

func (bs *ButtonSet) JustReleased() bool {
	if bs == nil {
		return false
	}
	return bs.Button.justReleased
}

func (bs *ButtonSet) Repeated() bool {
	if bs == nil {
		return false
	}
	return bs.Button.repeated
}

func (bs *ButtonSet) Consume() {
	if bs == nil {
		return
	}
	for _, b := range bs.Buttons {
		consumeGamepad[b] = true
	}
	for _, k := range bs.Keys {
		consumeKey[k] = true
	}
}

type AxisSet struct {
	F float64
	A pixelgl.GamepadAxis
}

type ButtonSet struct {
	Button   Button
	Keys     []pixelgl.Button        `toml:"keys"`
	Scroll   int                     `toml:"scroll"`
	Buttons  []pixelgl.GamepadButton `toml:"buttons"`
	Axis     pixelgl.GamepadAxis     `toml:"axis"`
	AxisV    int                     `toml:"axis_v"`
	noJoy    bool
	LastMode Mode
}