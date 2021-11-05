package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"math"
)

func AnimationSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasAnimation) {
		anim, ok := result.Components[myecs.Animation].(*reanimator.Tree)
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		if ok && okT {
			dist := descent.Descent.GetPlayer().Transform.Pos.Sub(tran.Pos)
			if math.Abs(dist.X) < 24. * world.TileSize && math.Abs(dist.Y) < 24. * world.TileSize {
				anim.Update()
			}
		}
	}
}

func AnimationDraw() {
	for _, result := range myecs.Manager.Query(myecs.HasAnimDrawing) {
		anim, okA := result.Components[myecs.Animation].(*reanimator.Tree)
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		bkey, okB := result.Components[myecs.Batch].(string)
		if okA && okT && okB && !tran.Hide {
			dist := camera.Cam.Pos.Sub(tran.Pos)
			if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
				if batcher, ok := img.Batchers[bkey]; ok {
					anim.DrawColorMask(batcher.Batch(), tran.Mat, tran.Mask)
				}
			}
		}
	}
}

func SpriteDraw() {
	for _, result := range myecs.Manager.Query(myecs.HasSprDrawing) {
		spr, okS := result.Components[myecs.Sprite].(*pixel.Sprite)
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		bkey, okB := result.Components[myecs.Batch].(string)
		if okS && okT && okB && !tran.Hide {
			dist := camera.Cam.Pos.Sub(tran.Pos)
			if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
				if batcher, ok := img.Batchers[bkey]; ok {
					spr.DrawColorMask(batcher.Batch(), tran.Mat, tran.Mask)
				}
			}
		}
	}
}