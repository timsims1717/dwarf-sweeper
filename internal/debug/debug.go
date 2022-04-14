package debug

import (
	"github.com/faiface/pixel/pixelgl"
	"strings"
)

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
	if Text {
		DrawText(win)
	}
	DrawFPS(win)
}

func Clear() {
	imd.Clear()
	lines = &strings.Builder{}
}
