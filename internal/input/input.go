package input

import (
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Input struct {
	Cursor  pixel.Vec
	World   pixel.Vec
	Coords  world.Coords
	Select  toggle
	Cancel  toggle
	Scroll  float64
	HotKeys map[pixelgl.Button]toggle
	HotFunc map[pixelgl.Button]func()
	Debug   bool
}

func NewInput() *Input {
	return &Input{
		HotKeys: make(map[pixelgl.Button]toggle),
		HotFunc: make(map[pixelgl.Button]func()),
	}
}

func (i *Input) Update(win *pixelgl.Window) {
	i.Cursor = win.MousePosition()
	i.World = camera.Cam.Mat.Unproject(win.MousePosition())
	i.Coords.X, i.Coords.Y = world.WorldToMap(i.World.X, i.World.Y)
	i.Select.Set(win, pixelgl.MouseButtonLeft)
	i.Cancel.Set(win, pixelgl.MouseButtonRight)
	i.Debug = win.JustPressed(pixelgl.KeyP)

	debug.AddLine(colornames.Red, imdraw.SharpEndShape, pixel.ZV, i.World, 1.)

	if win.JustPressed(pixelgl.KeyE) {
		particles.PartBatcher.Sprites["b_8"].Draw(win, pixel.IM.Moved(i.World))
	}
	if win.Pressed(pixelgl.KeyEscape) {
		win.SetClosed(true)
	}
	if win.Pressed(pixelgl.KeyLeft) {
		camera.Cam.Left()
	}
	if win.Pressed(pixelgl.KeyRight) {
		camera.Cam.Right()
	}
	if win.Pressed(pixelgl.KeyDown) {
		camera.Cam.Down()
	}
	if win.Pressed(pixelgl.KeyUp) {
		camera.Cam.Up()
	}
	camera.Cam.ZoomIn(win.MouseScroll().Y)

	for key, tog := range i.HotKeys {
		tog.Set(win, key)
		if tog.JustPressed() {
			tog.Consume()
			if fn, ok := i.HotFunc[key]; ok && fn != nil {
				fn()
			}
		}
	}
}

func (i *Input) SetHotKey(btn pixelgl.Button, fn func()) {
	i.HotKeys[btn] = toggle{}
	i.HotFunc[btn] = fn
}

func (i *Input) RemoveHotKey(btn pixelgl.Button) {
	delete(i.HotKeys, btn)
	delete(i.HotFunc, btn)
}

func (i *Input) RemoveHotKeys() {
	for key := range i.HotKeys {
		delete(i.HotKeys, key)
	}
	for key := range i.HotFunc {
		delete(i.HotFunc, key)
	}
}

type toggle struct {
	justPressed  bool
	pressed      bool
	justReleased bool
	consumed     bool
}

func (t *toggle) JustPressed() bool {
	return t.justPressed && !t.consumed
}

func (t *toggle) Pressed() bool {
	return t.pressed && !t.consumed
}

func (t *toggle) JustReleased() bool {
	return t.justReleased && !t.consumed
}

func (t *toggle) Consume() {
	t.consumed = true
}

func (t *toggle) Set(win *pixelgl.Window, button pixelgl.Button) {
	t.justPressed = win.JustPressed(button)
	t.pressed = win.Pressed(button)
	t.justReleased = win.JustReleased(button)
	t.consumed = t.consumed && !t.justPressed && !t.pressed && !t.justReleased
}

func (t *toggle) SetBool(pressed bool) {
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
