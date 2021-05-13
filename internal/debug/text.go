package debug

import (
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/typeface"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

var (
	debugText *text.Text
)

func InitializeText() {
	debugText = text.New(pixel.ZV, typeface.BasicAtlas)
}

func DrawText(win *pixelgl.Window) {
	mat := camera.Cam.UITransform(pixel.V(camera.Cam.Width * -0.5 + 20., camera.Cam.Height * 0.5-40.), pixel.V(2., 2.), 0.)
	debugText.Draw(win, mat)
	debugText.Clear()
}

func AddText(s string) {
	fmt.Fprintln(debugText, s)
}