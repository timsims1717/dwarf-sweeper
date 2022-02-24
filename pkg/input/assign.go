package input

import "github.com/faiface/pixel/pixelgl"

// CheckAssign checks all inputs and assigns any new inputs to the provided
// Input. Returns true if an assignment was made.
func CheckAssign(win *pixelgl.Window, in *Input, key string) bool {
	assigned := false
	if in.Mode != Gamepad {
		keys := GetAllJustPressedKey(win)
		if len(keys) > 0 {
			in.Buttons[key].Keys = append(in.Buttons[key].Keys, keys[0])
			assigned = true
		} else {
			mKeys := GetAllJustPressedMouse(win)
			if len(mKeys) > 0 {
				in.Buttons[key].Keys = append(in.Buttons[key].Keys, mKeys[0])
				assigned = true
			} else if in.ScrollV > 0. {
				in.Buttons[key].Scroll = 1
				assigned = true
			} else if in.ScrollV < 0. {
				in.Buttons[key].Scroll = -1
				assigned = true
			}
		}
	}
	if !assigned && in.Mode != KeyboardMouse {
		buttons := GetAllJustPressedGamepad(win, in.Joystick)
		if len(buttons) > 0 {
			in.Buttons[key].Buttons = append(in.Buttons[key].Buttons, buttons[0])
			assigned = true
		} else {
			axes := GetAllAxisGamepad(win, in)
			if len(axes) > 0 {
				in.Buttons[key].Axis = axes[0].Axis
				if axes[0].Dir > 0. {
					in.Buttons[key].AxisV = 1
					assigned = true
				} else if axes[0].Dir < 0. {
					in.Buttons[key].AxisV = -1
					assigned = true
				}
			}
		}
	}
	return assigned
}

func ClearInput(in *Input, key string) {
	if in.Mode != Gamepad {
		in.Buttons[key].Keys = []pixelgl.Button{}
		in.Buttons[key].Scroll = 0
	}
	if in.Mode != KeyboardMouse {
		in.Buttons[key].Buttons = []pixelgl.GamepadButton{}
		in.Buttons[key].AxisV = 0
	}
}
