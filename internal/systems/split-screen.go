package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"math"
)

func SplitScreenDrawSystem(win *pixelgl.Window, canvas1 *pixelgl.Canvas, pos pixel.Vec) {
	canvas1.Clear(color.RGBA{})
	cPos := camera.Cam.Pos
	canvas1.SetBounds(pixel.R(cPos.X - constants.ActualW*0.5, cPos.Y - constants.BaseH*0.5, cPos.X + constants.ActualW*0.5, cPos.Y + constants.BaseH*0.5))
	for _, result := range myecs.Manager.Query(myecs.IsDrawable) {
		draw := result.Components[myecs.Drawable]
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		//bkey, okB := result.Components[myecs.Batch].(string)
		if okT && !tran.Hide {
			dist := cPos.Sub(tran.Pos)
			if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
				//if batcher, ok := img.Batchers[bkey]; ok {
					if spr, okS := draw.(*pixel.Sprite); okS {
						spr.DrawColorMask(canvas1, tran.Mat, tran.Mask)
					} else if anim, okA := draw.(*reanimator.Tree); okA {
						anim.DrawColorMask(canvas1, tran.Mat, tran.Mask)
					}
				//}
			}
		}
	}
	mat := pixel.IM
	//mat := pixel.IM.Scaled(pixel.ZV, camera.Cam.GetZoomScale())
	//mat = mat.Moved(pixel.ZV.Scaled(camera.Cam.GetZoomScale()))
	mat = mat.Moved(pos)
	canvas1.Draw(win, mat)
}