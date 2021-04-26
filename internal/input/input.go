package input

import (
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/pkg/camera"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var Input = &input{}

type input struct {
	Cursor    pixel.Vec
	World     pixel.Vec
	Click     bool
	Scroll    float64
	Debug     bool
	XDir      XDirection
	XDirC     bool
	Jumping   bool
	Jumped    bool
	IsDig     bool
	IsMark    bool
	UpDown    int
	LeftRight int
	UseCursor bool
	Back      bool
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
	i.Click = win.JustPressed(pixelgl.MouseButtonLeft)

	if win.JustPressed(pixelgl.KeyKP1) || win.JustPressed(pixelgl.KeyKP4) || win.JustPressed(pixelgl.KeyKP7) {
		i.LeftRight = -1
		i.UseCursor = false
	} else if win.JustPressed(pixelgl.KeyKP3) || win.JustPressed(pixelgl.KeyKP6) || win.JustPressed(pixelgl.KeyKP9) {
		i.LeftRight = 1
		i.UseCursor = false
	} else if win.JustPressed(pixelgl.KeyKP2) || win.JustPressed(pixelgl.KeyKP8) {
		i.LeftRight = 0
		i.UseCursor = false
	}
	if win.JustPressed(pixelgl.KeyKP1) || win.JustPressed(pixelgl.KeyKP2) || win.JustPressed(pixelgl.KeyKP3) {
		i.UpDown = -1
		i.UseCursor = false
	} else if win.JustPressed(pixelgl.KeyKP7) || win.JustPressed(pixelgl.KeyKP8) || win.JustPressed(pixelgl.KeyKP9) {
		i.UpDown = 1
		i.UseCursor = false
	} else if win.JustPressed(pixelgl.KeyKP4) || win.JustPressed(pixelgl.KeyKP6) {
		i.UpDown = 0
		i.UseCursor = false
	}

	//if win.Pressed(pixelgl.KeyLeft) {
	//	camera.Cam.Left()
	//}
	//if win.Pressed(pixelgl.KeyRight) {
	//	camera.Cam.Right()
	//}
	//if win.Pressed(pixelgl.KeyDown) {
	//	camera.Cam.Down()
	//}
	//if win.Pressed(pixelgl.KeyUp) {
	//	camera.Cam.Up()
	//}

	i.World = camera.Cam.Mat.Unproject(win.MousePosition())
	i.Debug = win.JustPressed(pixelgl.KeyP)

	debug.AddLine(colornames.Red, imdraw.SharpEndShape, pixel.ZV, i.World, 1.)

	i.Back = win.JustPressed(pixelgl.KeyEscape)

	if win.JustPressed(pixelgl.KeyD) {
		i.XDir = Right
		i.XDirC = true
	} else if win.JustPressed(pixelgl.KeyA) {
		i.XDir = Left
		i.XDirC = true
	} else if !win.Pressed(pixelgl.KeyD) && !win.Pressed(pixelgl.KeyA) {
		i.XDir = None
		i.XDirC = true
	}
	i.IsDig = win.JustPressed(pixelgl.MouseButtonLeft) || win.JustPressed(pixelgl.KeyKP0) || win.JustPressed(pixelgl.KeyKPEnter)
	i.IsMark = win.JustPressed(pixelgl.MouseButtonRight)
	i.Jumping = win.Pressed(pixelgl.KeyW)
	i.Jumped = win.JustPressed(pixelgl.KeyW)
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
