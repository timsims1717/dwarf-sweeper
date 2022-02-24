package input

import "github.com/faiface/pixel/pixelgl"

func GetAllJustPressedMouse(win *pixelgl.Window) []pixelgl.Button {
	var result []pixelgl.Button
	if win.JustPressed(pixelgl.MouseButton1) {
		result = append(result, pixelgl.MouseButton1)
	}
	if win.JustPressed(pixelgl.MouseButton2) {
		result = append(result, pixelgl.MouseButton2)
	}
	if win.JustPressed(pixelgl.MouseButton3) {
		result = append(result, pixelgl.MouseButton3)
	}
	if win.JustPressed(pixelgl.MouseButton4) {
		result = append(result, pixelgl.MouseButton4)
	}
	if win.JustPressed(pixelgl.MouseButton5) {
		result = append(result, pixelgl.MouseButton5)
	}
	if win.JustPressed(pixelgl.MouseButton6) {
		result = append(result, pixelgl.MouseButton6)
	}
	if win.JustPressed(pixelgl.MouseButton7) {
		result = append(result, pixelgl.MouseButton7)
	}
	if win.JustPressed(pixelgl.MouseButton8) {
		result = append(result, pixelgl.MouseButton8)
	}
	if win.JustPressed(pixelgl.MouseButtonLast) {
		result = append(result, pixelgl.MouseButtonLast)
	}
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		result = append(result, pixelgl.MouseButtonLeft)
	}
	if win.JustPressed(pixelgl.MouseButtonRight) {
		result = append(result, pixelgl.MouseButtonRight)
	}
	if win.JustPressed(pixelgl.MouseButtonMiddle) {
		result = append(result, pixelgl.MouseButtonMiddle)
	}
	return result
}

func GetAllJustPressedKey(win *pixelgl.Window) []pixelgl.Button {
	var result []pixelgl.Button
	if win.JustPressed(pixelgl.KeySpace) {
		result = append(result, pixelgl.KeySpace)
	}
	if win.JustPressed(pixelgl.KeyApostrophe) {
		result = append(result, pixelgl.KeyApostrophe)
	}
	if win.JustPressed(pixelgl.KeyComma) {
		result = append(result, pixelgl.KeyComma)
	}
	if win.JustPressed(pixelgl.KeyMinus) {
		result = append(result, pixelgl.KeyMinus)
	}
	if win.JustPressed(pixelgl.KeyPeriod) {
		result = append(result, pixelgl.KeyPeriod)
	}
	if win.JustPressed(pixelgl.KeySlash) {
		result = append(result, pixelgl.KeySlash)
	}
	if win.JustPressed(pixelgl.Key0) {
		result = append(result, pixelgl.Key0)
	}
	if win.JustPressed(pixelgl.Key1) {
		result = append(result, pixelgl.Key1)
	}
	if win.JustPressed(pixelgl.Key2) {
		result = append(result, pixelgl.Key2)
	}
	if win.JustPressed(pixelgl.Key3) {
		result = append(result, pixelgl.Key3)
	}
	if win.JustPressed(pixelgl.Key4) {
		result = append(result, pixelgl.Key4)
	}
	if win.JustPressed(pixelgl.Key5) {
		result = append(result, pixelgl.Key5)
	}
	if win.JustPressed(pixelgl.Key6) {
		result = append(result, pixelgl.Key6)
	}
	if win.JustPressed(pixelgl.Key7) {
		result = append(result, pixelgl.Key7)
	}
	if win.JustPressed(pixelgl.Key8) {
		result = append(result, pixelgl.Key8)
	}
	if win.JustPressed(pixelgl.Key9) {
		result = append(result, pixelgl.Key9)
	}
	if win.JustPressed(pixelgl.KeySemicolon) {
		result = append(result, pixelgl.KeySemicolon)
	}
	if win.JustPressed(pixelgl.KeyEqual) {
		result = append(result, pixelgl.KeyEqual)
	}
	if win.JustPressed(pixelgl.KeyA) {
		result = append(result, pixelgl.KeyA)
	}
	if win.JustPressed(pixelgl.KeyB) {
		result = append(result, pixelgl.KeyB)
	}
	if win.JustPressed(pixelgl.KeyC) {
		result = append(result, pixelgl.KeyC)
	}
	if win.JustPressed(pixelgl.KeyD) {
		result = append(result, pixelgl.KeyD)
	}
	if win.JustPressed(pixelgl.KeyE) {
		result = append(result, pixelgl.KeyE)
	}
	if win.JustPressed(pixelgl.KeyF) {
		result = append(result, pixelgl.KeyF)
	}
	if win.JustPressed(pixelgl.KeyG) {
		result = append(result, pixelgl.KeyG)
	}
	if win.JustPressed(pixelgl.KeyH) {
		result = append(result, pixelgl.KeyH)
	}
	if win.JustPressed(pixelgl.KeyI) {
		result = append(result, pixelgl.KeyI)
	}
	if win.JustPressed(pixelgl.KeyJ) {
		result = append(result, pixelgl.KeyJ)
	}
	if win.JustPressed(pixelgl.KeyK) {
		result = append(result, pixelgl.KeyK)
	}
	if win.JustPressed(pixelgl.KeyL) {
		result = append(result, pixelgl.KeyL)
	}
	if win.JustPressed(pixelgl.KeyM) {
		result = append(result, pixelgl.KeyM)
	}
	if win.JustPressed(pixelgl.KeyN) {
		result = append(result, pixelgl.KeyN)
	}
	if win.JustPressed(pixelgl.KeyO) {
		result = append(result, pixelgl.KeyO)
	}
	if win.JustPressed(pixelgl.KeyP) {
		result = append(result, pixelgl.KeyP)
	}
	if win.JustPressed(pixelgl.KeyQ) {
		result = append(result, pixelgl.KeyQ)
	}
	if win.JustPressed(pixelgl.KeyR) {
		result = append(result, pixelgl.KeyR)
	}
	if win.JustPressed(pixelgl.KeyS) {
		result = append(result, pixelgl.KeyS)
	}
	if win.JustPressed(pixelgl.KeyT) {
		result = append(result, pixelgl.KeyT)
	}
	if win.JustPressed(pixelgl.KeyU) {
		result = append(result, pixelgl.KeyU)
	}
	if win.JustPressed(pixelgl.KeyV) {
		result = append(result, pixelgl.KeyV)
	}
	if win.JustPressed(pixelgl.KeyW) {
		result = append(result, pixelgl.KeyW)
	}
	if win.JustPressed(pixelgl.KeyX) {
		result = append(result, pixelgl.KeyX)
	}
	if win.JustPressed(pixelgl.KeyY) {
		result = append(result, pixelgl.KeyY)
	}
	if win.JustPressed(pixelgl.KeyZ) {
		result = append(result, pixelgl.KeyZ)
	}
	if win.JustPressed(pixelgl.KeyLeftBracket) {
		result = append(result, pixelgl.KeyLeftBracket)
	}
	if win.JustPressed(pixelgl.KeyBackslash) {
		result = append(result, pixelgl.KeyBackslash)
	}
	if win.JustPressed(pixelgl.KeyRightBracket) {
		result = append(result, pixelgl.KeyRightBracket)
	}
	if win.JustPressed(pixelgl.KeyGraveAccent) {
		result = append(result, pixelgl.KeyGraveAccent)
	}
	if win.JustPressed(pixelgl.KeyWorld1) {
		result = append(result, pixelgl.KeyWorld1)
	}
	if win.JustPressed(pixelgl.KeyWorld2) {
		result = append(result, pixelgl.KeyWorld2)
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		result = append(result, pixelgl.KeyEscape)
	}
	if win.JustPressed(pixelgl.KeyEnter) {
		result = append(result, pixelgl.KeyEnter)
	}
	if win.JustPressed(pixelgl.KeyTab) {
		result = append(result, pixelgl.KeyTab)
	}
	if win.JustPressed(pixelgl.KeyBackspace) {
		result = append(result, pixelgl.KeyBackspace)
	}
	if win.JustPressed(pixelgl.KeyInsert) {
		result = append(result, pixelgl.KeyInsert)
	}
	if win.JustPressed(pixelgl.KeyDelete) {
		result = append(result, pixelgl.KeyDelete)
	}
	if win.JustPressed(pixelgl.KeyRight) {
		result = append(result, pixelgl.KeyRight)
	}
	if win.JustPressed(pixelgl.KeyLeft) {
		result = append(result, pixelgl.KeyLeft)
	}
	if win.JustPressed(pixelgl.KeyDown) {
		result = append(result, pixelgl.KeyDown)
	}
	if win.JustPressed(pixelgl.KeyUp) {
		result = append(result, pixelgl.KeyUp)
	}
	if win.JustPressed(pixelgl.KeyPageUp) {
		result = append(result, pixelgl.KeyPageUp)
	}
	if win.JustPressed(pixelgl.KeyPageDown) {
		result = append(result, pixelgl.KeyPageDown)
	}
	if win.JustPressed(pixelgl.KeyHome) {
		result = append(result, pixelgl.KeyHome)
	}
	if win.JustPressed(pixelgl.KeyEnd) {
		result = append(result, pixelgl.KeyEnd)
	}
	if win.JustPressed(pixelgl.KeyCapsLock) {
		result = append(result, pixelgl.KeyCapsLock)
	}
	if win.JustPressed(pixelgl.KeyScrollLock) {
		result = append(result, pixelgl.KeyScrollLock)
	}
	if win.JustPressed(pixelgl.KeyNumLock) {
		result = append(result, pixelgl.KeyNumLock)
	}
	if win.JustPressed(pixelgl.KeyPrintScreen) {
		result = append(result, pixelgl.KeyPrintScreen)
	}
	if win.JustPressed(pixelgl.KeyPause) {
		result = append(result, pixelgl.KeyPause)
	}
	if win.JustPressed(pixelgl.KeyF1) {
		result = append(result, pixelgl.KeyF1)
	}
	if win.JustPressed(pixelgl.KeyF2) {
		result = append(result, pixelgl.KeyF2)
	}
	if win.JustPressed(pixelgl.KeyF3) {
		result = append(result, pixelgl.KeyF3)
	}
	if win.JustPressed(pixelgl.KeyF4) {
		result = append(result, pixelgl.KeyF4)
	}
	if win.JustPressed(pixelgl.KeyF5) {
		result = append(result, pixelgl.KeyF5)
	}
	if win.JustPressed(pixelgl.KeyF6) {
		result = append(result, pixelgl.KeyF6)
	}
	if win.JustPressed(pixelgl.KeyF7) {
		result = append(result, pixelgl.KeyF7)
	}
	if win.JustPressed(pixelgl.KeyF8) {
		result = append(result, pixelgl.KeyF8)
	}
	if win.JustPressed(pixelgl.KeyF9) {
		result = append(result, pixelgl.KeyF9)
	}
	if win.JustPressed(pixelgl.KeyF10) {
		result = append(result, pixelgl.KeyF10)
	}
	if win.JustPressed(pixelgl.KeyF11) {
		result = append(result, pixelgl.KeyF11)
	}
	if win.JustPressed(pixelgl.KeyF12) {
		result = append(result, pixelgl.KeyF12)
	}
	if win.JustPressed(pixelgl.KeyF13) {
		result = append(result, pixelgl.KeyF13)
	}
	if win.JustPressed(pixelgl.KeyF14) {
		result = append(result, pixelgl.KeyF14)
	}
	if win.JustPressed(pixelgl.KeyF15) {
		result = append(result, pixelgl.KeyF15)
	}
	if win.JustPressed(pixelgl.KeyF16) {
		result = append(result, pixelgl.KeyF16)
	}
	if win.JustPressed(pixelgl.KeyF17) {
		result = append(result, pixelgl.KeyF17)
	}
	if win.JustPressed(pixelgl.KeyF18) {
		result = append(result, pixelgl.KeyF18)
	}
	if win.JustPressed(pixelgl.KeyF19) {
		result = append(result, pixelgl.KeyF19)
	}
	if win.JustPressed(pixelgl.KeyF20) {
		result = append(result, pixelgl.KeyF20)
	}
	if win.JustPressed(pixelgl.KeyF21) {
		result = append(result, pixelgl.KeyF21)
	}
	if win.JustPressed(pixelgl.KeyF22) {
		result = append(result, pixelgl.KeyF22)
	}
	if win.JustPressed(pixelgl.KeyF23) {
		result = append(result, pixelgl.KeyF23)
	}
	if win.JustPressed(pixelgl.KeyF24) {
		result = append(result, pixelgl.KeyF24)
	}
	if win.JustPressed(pixelgl.KeyF25) {
		result = append(result, pixelgl.KeyF25)
	}
	if win.JustPressed(pixelgl.KeyKP0) {
		result = append(result, pixelgl.KeyKP0)
	}
	if win.JustPressed(pixelgl.KeyKP1) {
		result = append(result, pixelgl.KeyKP1)
	}
	if win.JustPressed(pixelgl.KeyKP2) {
		result = append(result, pixelgl.KeyKP2)
	}
	if win.JustPressed(pixelgl.KeyKP3) {
		result = append(result, pixelgl.KeyKP3)
	}
	if win.JustPressed(pixelgl.KeyKP4) {
		result = append(result, pixelgl.KeyKP4)
	}
	if win.JustPressed(pixelgl.KeyKP5) {
		result = append(result, pixelgl.KeyKP5)
	}
	if win.JustPressed(pixelgl.KeyKP6) {
		result = append(result, pixelgl.KeyKP6)
	}
	if win.JustPressed(pixelgl.KeyKP7) {
		result = append(result, pixelgl.KeyKP7)
	}
	if win.JustPressed(pixelgl.KeyKP8) {
		result = append(result, pixelgl.KeyKP8)
	}
	if win.JustPressed(pixelgl.KeyKP9) {
		result = append(result, pixelgl.KeyKP9)
	}
	if win.JustPressed(pixelgl.KeyKPDecimal) {
		result = append(result, pixelgl.KeyKPDecimal)
	}
	if win.JustPressed(pixelgl.KeyKPDivide) {
		result = append(result, pixelgl.KeyKPDivide)
	}
	if win.JustPressed(pixelgl.KeyKPMultiply) {
		result = append(result, pixelgl.KeyKPMultiply)
	}
	if win.JustPressed(pixelgl.KeyKPSubtract) {
		result = append(result, pixelgl.KeyKPSubtract)
	}
	if win.JustPressed(pixelgl.KeyKPAdd) {
		result = append(result, pixelgl.KeyKPAdd)
	}
	if win.JustPressed(pixelgl.KeyKPEnter) {
		result = append(result, pixelgl.KeyKPEnter)
	}
	if win.JustPressed(pixelgl.KeyKPEqual) {
		result = append(result, pixelgl.KeyKPEqual)
	}
	if win.JustPressed(pixelgl.KeyLeftShift) {
		result = append(result, pixelgl.KeyLeftShift)
	}
	if win.JustPressed(pixelgl.KeyLeftControl) {
		result = append(result, pixelgl.KeyLeftControl)
	}
	if win.JustPressed(pixelgl.KeyLeftAlt) {
		result = append(result, pixelgl.KeyLeftAlt)
	}
	if win.JustPressed(pixelgl.KeyLeftSuper) {
		result = append(result, pixelgl.KeyLeftSuper)
	}
	if win.JustPressed(pixelgl.KeyRightShift) {
		result = append(result, pixelgl.KeyRightShift)
	}
	if win.JustPressed(pixelgl.KeyRightControl) {
		result = append(result, pixelgl.KeyRightControl)
	}
	if win.JustPressed(pixelgl.KeyRightAlt) {
		result = append(result, pixelgl.KeyRightAlt)
	}
	if win.JustPressed(pixelgl.KeyRightSuper) {
		result = append(result, pixelgl.KeyRightSuper)
	}
	if win.JustPressed(pixelgl.KeyMenu) {
		result = append(result, pixelgl.KeyMenu)
	}
	return result
}

func GetAllJustPressedGamepad(win *pixelgl.Window, js pixelgl.Joystick) []pixelgl.GamepadButton {
	var result []pixelgl.GamepadButton
	if win.JoystickPresent(js) {
		if win.JoystickJustPressed(js, pixelgl.ButtonA) {
			result = append(result, pixelgl.ButtonA)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonB) {
			result = append(result, pixelgl.ButtonB)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonX) {
			result = append(result, pixelgl.ButtonX)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonY) {
			result = append(result, pixelgl.ButtonY)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonLeftBumper) {
			result = append(result, pixelgl.ButtonLeftBumper)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonRightBumper) {
			result = append(result, pixelgl.ButtonRightBumper)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonBack) {
			result = append(result, pixelgl.ButtonBack)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonStart) {
			result = append(result, pixelgl.ButtonStart)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonLeftThumb) {
			result = append(result, pixelgl.ButtonLeftThumb)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonRightThumb) {
			result = append(result, pixelgl.ButtonRightThumb)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonDpadUp) {
			result = append(result, pixelgl.ButtonDpadUp)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonDpadRight) {
			result = append(result, pixelgl.ButtonDpadRight)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonDpadDown) {
			result = append(result, pixelgl.ButtonDpadDown)
		}
		if win.JoystickJustPressed(js, pixelgl.ButtonDpadLeft) {
			result = append(result, pixelgl.ButtonDpadLeft)
		}
	}
	return result
}

type AxisReturn struct {
	Axis pixelgl.GamepadAxis
	Dir  float64
}

func GetAllAxisGamepad(win *pixelgl.Window, in *Input) []AxisReturn {
	var result []AxisReturn
	if win.JoystickPresent(in.Joystick) {
		if win.JoystickAxis(in.Joystick, pixelgl.AxisLeftX) > in.Deadzone {
			result = append(result, AxisReturn{
				Axis: pixelgl.AxisLeftX,
				Dir:  1.,
			})
		}
		if win.JoystickAxis(in.Joystick, pixelgl.AxisLeftY) > in.Deadzone {
			result = append(result, AxisReturn{
				Axis: pixelgl.AxisLeftY,
				Dir:  1.,
			})
		}
		if win.JoystickAxis(in.Joystick, pixelgl.AxisRightX) > in.Deadzone {
			result = append(result, AxisReturn{
				Axis: pixelgl.AxisRightX,
				Dir:  1.,
			})
		}
		if win.JoystickAxis(in.Joystick, pixelgl.AxisRightY) > in.Deadzone {
			result = append(result, AxisReturn{
				Axis: pixelgl.AxisRightY,
				Dir:  1.,
			})
		}
		if win.JoystickAxis(in.Joystick, pixelgl.AxisLeftTrigger) > in.Deadzone {
			result = append(result, AxisReturn{
				Axis: pixelgl.AxisLeftTrigger,
				Dir:  1.,
			})
		}
		if win.JoystickAxis(in.Joystick, pixelgl.AxisRightTrigger) > in.Deadzone {
			result = append(result, AxisReturn{
				Axis: pixelgl.AxisRightTrigger,
				Dir:  1.,
			})
		}
		if win.JoystickAxis(in.Joystick, pixelgl.AxisLeftX) < -in.Deadzone {
			result = append(result, AxisReturn{
				Axis: pixelgl.AxisLeftX,
				Dir:  -1.,
			})
		}
		if win.JoystickAxis(in.Joystick, pixelgl.AxisLeftY) < -in.Deadzone {
			result = append(result, AxisReturn{
				Axis: pixelgl.AxisLeftY,
				Dir:  -1.,
			})
		}
		if win.JoystickAxis(in.Joystick, pixelgl.AxisRightX) < -in.Deadzone {
			result = append(result, AxisReturn{
				Axis: pixelgl.AxisRightX,
				Dir:  -1.,
			})
		}
		if win.JoystickAxis(in.Joystick, pixelgl.AxisRightY) < -in.Deadzone {
			result = append(result, AxisReturn{
				Axis: pixelgl.AxisRightY,
				Dir:  -1.,
			})
		}
	}
	return result
}
