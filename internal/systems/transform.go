package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
	"math"
)

func TransformSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasTransform) {
		if tran, ok := result.Components[myecs.Transform].(*transform.Transform); ok {
			if tran.Dispose {
				myecs.Manager.DisposeEntity(result)
			} else {
				tran.APos = tran.Pos.Add(tran.Offset)
				if tran.Shaking {
					switch tran.ShakeI {
					case 0,1,7:
						tran.APos.Y += 1.
					case 3,4,5:
						tran.APos.Y -= 1.
					}
					switch tran.ShakeI {
					case 1,2,3:
						tran.APos.X += 1.
					case 5,6,7:
						tran.APos.X -= 1.
					}
					tran.ShakeI++
					tran.ShakeI %= 8
					if tran.ShakeI == tran.ShakeE {
						tran.Shaking = false
					}
				}
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

				if descent.Descent != nil && !tran.KeepLoaded {
					p := descent.Descent.GetClosestPlayer(tran.Pos)
					if p != nil {
						dist := p.Transform.Pos.Sub(tran.Pos)
						tran.Load = math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance
					}
				}
			}
		}
	}
}

func ParentSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasParent) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		parent, okP := result.Components[myecs.Parent].(*transform.Transform)
		if okT && okP {
			if parent.Dispose {
				myecs.Manager.DisposeEntity(result)
			} else {
				tran.Pos = parent.Pos
				tran.Hide = parent.Hide
			}
		}
	}
}
