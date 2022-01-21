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
			tran.APos = tran.Pos.Add(tran.Offset)
			tran.APos.X = math.Round(tran.APos.X)
			tran.APos.Y = math.Round(tran.APos.Y)
			tran.Mat = pixel.IM
			if tran.Flip && tran.Flop {
				tran.Mat = tran.Mat.Scaled(pixel.ZV, -1.)
			} else if tran.Flip {
				tran.Mat = tran.Mat.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
			} else if tran.Flop {
				tran.Mat = tran.Mat.ScaledXY(pixel.ZV, pixel.V(1., -1.))
			}
			tran.Mat = tran.Mat.ScaledXY(pixel.ZV, tran.Scalar.Scaled(tran.UIZoom))
			tran.Mat = tran.Mat.Rotated(pixel.ZV, math.Pi*tran.Rot)
			tran.Mat = tran.Mat.Moved(tran.APos.Scaled(tran.UIZoom))
			tran.Mat = tran.Mat.Moved(tran.UIPos)
		}
	}
}

func ParentSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasParent) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		parent, okP := result.Components[myecs.Parent].(*transform.Transform)
		if okT && okP {
			tran.Pos = parent.Pos
		}
	}
}
