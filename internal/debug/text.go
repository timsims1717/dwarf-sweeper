package debug

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/typeface"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"strings"
)

var (
	debugText *typeface.Text
	lines     = &strings.Builder{}
)

func InitializeText() {
	col := colornames.Aliceblue
	col.A = 90
	debugText = typeface.New(camera.Cam, "basic", typeface.NewAlign(typeface.Left, typeface.Top), 1.0, 0.5, 0., 0.)
}

func DrawText(win *pixelgl.Window) {
	debugText.SetText(lines.String())
	debugText.Transform.Pos = pixel.V(constants.ActualW*-0.5 + 2., constants.BaseH*0.5)
	debugText.Transform.UIPos = camera.Cam.APos
	debugText.Transform.UIZoom = camera.Cam.GetZoomScale()
	debugText.Update()
	debugText.Draw(win)
}

func AddText(s string) {
	lines.WriteString(fmt.Sprintf("%s\n", s))
}
