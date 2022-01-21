package debug

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/menu"
	"dwarf-sweeper/pkg/timing"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	fpsText     *menu.ItemText
	versionText *menu.ItemText
)

func InitializeFPS() {
	col := colornames.Aliceblue
	col.A = 90
	fpsText = menu.NewItemText("", col, pixel.V(0.8, 0.8), menu.Left, menu.Bottom)
	fpsText.Transform.Pos = pixel.V(constants.BaseW*-0.5, constants.BaseH*-0.5+2)
	versionText = menu.NewItemText("", col, pixel.V(0.8, 0.8), menu.Right, menu.Bottom)
	versionText.Transform.Pos = pixel.V(constants.BaseW*0.5, constants.BaseH*-0.5+2)
}

func DrawFPS(win *pixelgl.Window) {
	fpsText.SetText(fmt.Sprintf("FPS: %s", timing.FPS))
	fpsText.Transform.UIPos = camera.Cam.APos
	fpsText.Transform.UIZoom = camera.Cam.GetZoomScale()
	fpsText.Update(pixel.Rect{})
	fpsText.Draw(win)
	versionText.SetText(fmt.Sprintf("%d.%d.%d", constants.Release, constants.Version, constants.Build))
	versionText.Transform.UIPos = camera.Cam.APos
	versionText.Transform.UIZoom = camera.Cam.GetZoomScale()
	versionText.Update(pixel.Rect{})
	versionText.Draw(win)
}
