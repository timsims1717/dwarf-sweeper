package debug

import "github.com/faiface/pixel/pixelgl"

var Debug = false

func Initialize() {
	InitializeLines()
	InitializeText()
}

func Draw(win *pixelgl.Window) {
	if Debug {
		DrawLines(win)
		DrawText(win)
	}
}

func Clear() {
	imd.Clear()
	debugText.Clear()
}