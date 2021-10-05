package input

import "github.com/faiface/pixel/pixelgl"

// NextGamepad returns the next integer representation of a valid pixelgl.Joystick.
// -1 is returned if none is found and pixelgl.JoystickLast is reached.
func NextGamepad(win *pixelgl.Window, curr int) int {
	var next int
	if curr < -1 {
		next = 0
	} else {
		next = curr + 1
	}
	for next < int(pixelgl.JoystickLast) + 1 {
		jsN := pixelgl.Joystick(next)
		if win.JoystickPresent(jsN) {
			return next
		}
		next++
	}
	return -1
}

// PrevGamepad returns the previous integer representation of a valid pixelgl.Joystick.
// -1 is returned if none is found and -1 is reached.
func PrevGamepad(win *pixelgl.Window, curr int) int {
	var prev int
	if curr == -1 {
		prev = int(pixelgl.JoystickLast)
	} else {
		prev = curr - 1
	}
	for prev > -1 {
		jsP := pixelgl.Joystick(prev)
		if win.JoystickPresent(jsP) {
			return prev
		}
		prev--
	}
	return -1
}

func GamepadString(b pixelgl.GamepadButton) string {
	switch b {
	case pixelgl.ButtonA:
		return "A"
	case pixelgl.ButtonB:
		return "B"
	case pixelgl.ButtonX:
		return "X"
	case pixelgl.ButtonY:
		return "Y"
	case pixelgl.ButtonLeftBumper:
		return "LBump"
	case pixelgl.ButtonRightBumper:
		return "RBump"
	case pixelgl.ButtonBack:
		return "Back"
	case pixelgl.ButtonStart:
		return "Start"
	case pixelgl.ButtonGuide:
		return "Guide"
	case pixelgl.ButtonLeftThumb:
		return "LThumb"
	case pixelgl.ButtonRightThumb:
		return "RThumb"
	case pixelgl.ButtonDpadUp:
		return "DUp"
	case pixelgl.ButtonDpadRight:
		return "DRight"
	case pixelgl.ButtonDpadDown:
		return "DDown"
	case pixelgl.ButtonDpadLeft:
		return "DLeft"
	}
	return "unknown"
}