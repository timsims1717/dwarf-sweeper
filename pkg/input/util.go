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
	for next < int(pixelgl.JoystickLast)+1 {
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
		return "ButtonA"
	case pixelgl.ButtonB:
		return "ButtonB"
	case pixelgl.ButtonX:
		return "ButtonX"
	case pixelgl.ButtonY:
		return "ButtonY"
	case pixelgl.ButtonLeftBumper:
		return "ButtonLeftBumper"
	case pixelgl.ButtonRightBumper:
		return "ButtonRightBumper"
	case pixelgl.ButtonBack:
		return "ButtonBack"
	case pixelgl.ButtonStart:
		return "ButtonStart"
	case pixelgl.ButtonLeftThumb:
		return "ButtonLeftThumb"
	case pixelgl.ButtonRightThumb:
		return "ButtonRightThumb"
	case pixelgl.ButtonDpadUp:
		return "ButtonDpadUp"
	case pixelgl.ButtonDpadRight:
		return "ButtonDpadRight"
	case pixelgl.ButtonDpadDown:
		return "ButtonDpadDown"
	case pixelgl.ButtonDpadLeft:
		return "ButtonDpadLeft"
	}
	return "unknown"
}

func AxisString(a pixelgl.GamepadAxis) string {
	switch a {
	case pixelgl.AxisLeftX:
		return "AxisLeftX"
	case pixelgl.AxisLeftY:
		return "AxisLeftY"
	case pixelgl.AxisRightX:
		return "AxisRightX"
	case pixelgl.AxisRightY:
		return "AxisRightY"
	case pixelgl.AxisLeftTrigger:
		return "AxisLeftTrigger"
	case pixelgl.AxisRightTrigger:
		return "AxisRightTrigger"
	}
	return "unknown"
}

func AxisDirString(a pixelgl.GamepadAxis, pos bool) string {
	if pos {
		switch a {
		case pixelgl.AxisLeftX:
			return "AxisLeftXPos"
		case pixelgl.AxisLeftY:
			return "AxisLeftYPos"
		case pixelgl.AxisRightX:
			return "AxisRightXPos"
		case pixelgl.AxisRightY:
			return "AxisRightYPos"
		case pixelgl.AxisLeftTrigger:
			return "AxisLeftTrigger"
		case pixelgl.AxisRightTrigger:
			return "AxisRightTrigger"
		}
	} else {
		switch a {
		case pixelgl.AxisLeftX:
			return "AxisLeftXNeg"
		case pixelgl.AxisLeftY:
			return "AxisLeftYNeg"
		case pixelgl.AxisRightX:
			return "AxisRightXNeg"
		case pixelgl.AxisRightY:
			return "AxisRightYNeg"
		}
	}
	return "unknown"
}
