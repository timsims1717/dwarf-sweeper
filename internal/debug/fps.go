package debug

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/typeface"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	fpsText     *typeface.Text
	versionText *typeface.Text
)

func InitializeFPS() {
	col := colornames.Aliceblue
	col.A = 90
	fpsText = typeface.New(camera.Cam,  "basic", typeface.NewAlign(typeface.Left, typeface.Bottom), 1.0, 0.5, 0., 0.)
	versionText = typeface.New(camera.Cam,  "basic", typeface.NewAlign(typeface.Right, typeface.Bottom), 1.0, 0.5, 0., 0.)
}

func DrawFPS(win *pixelgl.Window) {
	fpsText.SetText(fmt.Sprintf("FPS: %s", timing.FPS))
	fpsText.Transform.Pos = pixel.V(constants.ActualW*-0.5 + 2., constants.BaseH*-0.5+2)
	fpsText.Transform.UIPos = camera.Cam.APos
	fpsText.Transform.UIZoom = camera.Cam.GetZoomScale()
	fpsText.Update()
	fpsText.Draw(win)
	versionText.SetText(fmt.Sprintf("%d.%d.%d", constants.Release, constants.Version, constants.Build))
	versionText.Transform.Pos = pixel.V(constants.ActualW*0.5 - 2., constants.BaseH*-0.5+2)
	versionText.Transform.UIPos = camera.Cam.APos
	versionText.Transform.UIZoom = camera.Cam.GetZoomScale()
	versionText.Update()
	versionText.Draw(win)
}
