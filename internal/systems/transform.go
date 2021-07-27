package systems

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
	"math"
)

func TransformSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasTransform) {
		if tran, ok := result.Components[myecs.Transform].(*transform.Transform); ok {
			tran.APos = tran.Pos
			tran.APos.X = math.Floor(tran.APos.X)
			tran.APos.Y = math.Floor(tran.APos.Y)
			tran.Mat = pixel.IM
			if tran.Flip && tran.Flop {
				tran.Mat = tran.Mat.Scaled(pixel.ZV, -1.)
			} else if tran.Flip {
				tran.Mat = tran.Mat.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
			} else if tran.Flop {
				tran.Mat = tran.Mat.ScaledXY(pixel.ZV, pixel.V(1., -1.))
			}
			tran.Mat = tran.Mat.ScaledXY(pixel.ZV, tran.Scalar.Scaled(tran.UIZoom))
			tran.Mat = tran.Mat.Rotated(pixel.ZV, tran.Rot)
			tran.Mat = tran.Mat.Moved(tran.APos.Scaled(tran.UIZoom))
			tran.Mat = tran.Mat.Moved(tran.UIPos)
		}
	}
}