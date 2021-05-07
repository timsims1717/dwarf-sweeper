package debug

import "github.com/faiface/pixel/pixelgl"

func Initialize() {
	InitializeLines()
	InitializeText()
}

func Draw(win *pixelgl.Window) {
	DrawLines(win)
	DrawText(win)
}

func Clear() {
	imd.Clear()
	debugText.Clear()
}