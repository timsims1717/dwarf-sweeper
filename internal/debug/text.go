package debug

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/menu"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	debugText *menu.ItemText
)

func InitializeText() {
	debugText = menu.NewItemText("", colornames.Aliceblue, pixel.V(1., 1.), menu.Left, menu.Top)
	debugText.Transform.Pos = pixel.V(cfg.BaseW * -0.5, cfg.BaseH * 0.5)
}

func DrawText(win *pixelgl.Window) {
	debugText.Transform.UIPos = camera.Cam.Pos
	debugText.Transform.UIZoom = camera.Cam.GetZoomScale()
	debugText.Update(pixel.Rect{})
	debugText.Draw(win)
}

func AddText(s string) {
	debugText.AddText(s)
}