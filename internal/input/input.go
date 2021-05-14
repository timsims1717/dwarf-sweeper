package input

import (
	"dwarf-sweeper/pkg/camera"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var Input = &input{}

type input struct {
	Cursor     pixel.Vec
	World      pixel.Vec
	Click      Button
	Scroll     float64
	Debug      bool
	DebugPause bool
	XDir       XDirection
	XDirC      bool
	Jumping    Button
	IsDig      bool
	IsMark     bool
	LookUp     Button
	LookDown   Button
	UseCursor  bool
	Back       bool
	Fullscreen bool
}

type XDirection int

const (
	None = iota
	Left
	Right
)

func Restrict(win *pixelgl.Window, bl, tr pixel.Vec) {
	world := camera.Cam.Mat.Unproject(win.MousePosition())
	if bl.X <= tr.X {
		if bl.X > world.X {
			world.X = bl.X
		} else if tr.X < world.X {
			world.X = tr.X
		}
	}
	if bl.Y <= tr.Y {
		if bl.Y > world.Y {
			world.Y = bl.Y
		} else if tr.Y < world.Y {
			world.Y = tr.Y
		}
	}
	win.SetMousePosition(camera.Cam.Mat.Project(world))
}

func (i *input) Update(win *pixelgl.Window) {
	i.Cursor = win.MousePosition()
	i.Click.Set(win, pixelgl.MouseButtonLeft)

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

	i.World = camera.Cam.Mat.Unproject(win.MousePosition())
	i.DebugPause = win.JustPressed(pixelgl.KeyF9)
	i.Debug = win.JustPressed(pixelgl.KeyF3)

	i.Back = win.JustPressed(pixelgl.KeyEscape)

	if win.JustPressed(pixelgl.KeyD) {
		i.XDir = Right
		i.XDirC = true
	} else if win.JustPressed(pixelgl.KeyA) {
		i.XDir = Left
		i.XDirC = true
	} else if (i.XDir == Right && !win.Pressed(pixelgl.KeyD)) || (i.XDir == Left && !win.Pressed(pixelgl.KeyA)) {
		i.XDir = None
		i.XDirC = true
	}
	i.IsDig = win.JustPressed(pixelgl.MouseButtonLeft) || win.JustPressed(pixelgl.KeyKP0) || win.JustPressed(pixelgl.KeyKPEnter)
	i.IsMark = win.JustPressed(pixelgl.MouseButtonRight)
	i.Jumping.Set(win, pixelgl.KeySpace)
	i.Fullscreen = win.JustPressed(pixelgl.KeyF)
	i.LookUp.Set(win, pixelgl.KeyW)
	i.LookDown.Set(win, pixelgl.KeyS)
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
	t.consumed = t.consumed && !t.justPressed && !t.pressed && !t.justReleased
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
