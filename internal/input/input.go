package input

import (
	"dwarf-sweeper/pkg/camera"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type AxisSet struct {

}

type BoolSet struct {
	B bool
	S pixelgl.Button
}

func NewBool(button pixelgl.Button) *BoolSet {
	return &BoolSet{
		S: button,
	}
}

type ButtonSet struct {
	B Button
	S pixelgl.Button
}

func NewButton(button pixelgl.Button) *ButtonSet {
	return &ButtonSet{
		S: button,
	}
}

type Input struct {
	Cursor  pixel.Vec
	World   pixel.Vec
	Axes    map[string]float64
	Bools   map[string]*BoolSet
	Buttons map[string]*ButtonSet
}

func (i *Input) Update(win *pixelgl.Window) {
	i.Cursor = win.MousePosition()
	i.World = camera.Cam.Mat.Unproject(win.MousePosition())

	for _, set := range i.Bools {
		set.B = win.JustPressed(set.S)
	}

	for _, set := range i.Buttons {
		set.B.Set(win, set.S)
	}
}

func (i *Input) GetBool(s string) bool {
	if b, ok := i.Bools[s]; ok {
		return b.B
	}
	return false
}

func (i *Input) GetButton(s string) *Button {
	if b, ok := i.Buttons[s]; ok {
		return &b.B
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

func (t *Button) Set(win *pixelgl.Window, button pixelgl.Button) {
	t.justPressed = win.JustPressed(button)
	t.pressed = win.Pressed(button)
	t.justReleased = win.JustReleased(button)
	t.consumed = t.consumed && (t.justPressed || t.pressed || t.justReleased)
}

func (t *Button) SetBool(pressed bool) {
	if pressed {
		if !t.pressed {
			t.justPressed = true
		} else {
			t.justPressed = false
		}
		t.pressed = true
		t.justReleased = false
	} else {
		if t.pressed {
			t.justReleased = true
		}
		t.pressed = false
		t.justPressed = false
	}
	t.consumed = t.consumed && !t.justPressed && !t.pressed && !t.justReleased
}
