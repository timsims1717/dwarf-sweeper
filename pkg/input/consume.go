package input

import "github.com/faiface/pixel/pixelgl"

var (
	consumeGamepad map[pixelgl.GamepadButton]bool
	consumeKey     map[pixelgl.Button]bool
)

func init() {
	consumeGamepad = make(map[pixelgl.GamepadButton]bool)
	consumeKey = make(map[pixelgl.Button]bool)
}

func updateConsume(win *pixelgl.Window, js pixelgl.Joystick) {
	for g, b := range consumeGamepad {
		consumeGamepad[g] = b && win.JoystickPressed(js, g)
	}
	for k, b := range consumeKey {
		consumeKey[k] = b && win.Pressed(k)
	}
}