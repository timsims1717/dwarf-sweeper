package debug

import "github.com/faiface/pixel/pixelgl"

func Initialize() {
	InitializeLines()
}

func Draw(win *pixelgl.Window) {
	DrawLines(win)
}