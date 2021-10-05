package input

import "github.com/faiface/pixel/pixelgl"

// CheckAssign checks all inputs and assigns any new inputs to the provided
// Input. Returns true if an assignment was made.
func CheckAssign(win *pixelgl.Window, in *Input, key string) bool {
	assigned := false
	if in.Mode != Gamepad {
		keys := GetAllJustPressedKey(win)
		if len(keys) > 0 {
			in.Buttons[key].Key = append(in.Buttons[key].Key, keys[0])
			assigned = true
		} else {
			mKeys := GetAllJustPressedMouse(win)
			if len(mKeys) > 0 {
				in.Buttons[key].Key = append(in.Buttons[key].Key, mKeys[0])
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
			in.Buttons[key].GPKey = append(in.Buttons[key].GPKey, buttons[0])
			assigned = true
		} else {
			axes := GetAllAxisGamepad(win, in.Joystick)
			if len(axes) > 0 {
				in.Buttons[key].Axis = axes[0].Axis
				if axes[0].Dir > 0. {
					in.Buttons[key].GP = 1
					assigned = true
				} else if axes[0].Dir < 0. {
					in.Buttons[key].GP = -1
					assigned = true
				}
			}
		}
	}
	return assigned
}

func ClearInput(in *Input, key string) {
	if in.Mode != Gamepad {
		in.Buttons[key].Key = []pixelgl.Button{}
		in.Buttons[key].Scroll = 0
	}
	if in.Mode != KeyboardMouse {
		in.Buttons[key].GPKey = []pixelgl.GamepadButton{}
		in.Buttons[key].GP = 0
	}
}