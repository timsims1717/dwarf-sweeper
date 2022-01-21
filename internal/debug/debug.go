package debug

import "github.com/faiface/pixel/pixelgl"

var (
	Debug = false
	Text  = false
)

func Initialize() {
	InitializeLines()
	InitializeText()
	InitializeFPS()
}

func Draw(win *pixelgl.Window) {
	if Debug {
		DrawLines(win)
	}
	if Text {
		DrawText(win)
	}
	DrawFPS(win)
}

func Clear() {
	imd.Clear()
	debugText.SetText("")
	fpsText.SetText("")
}
