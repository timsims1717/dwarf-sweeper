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
			if math.Abs(dist.X) < 24.*world.TileSize && math.Abs(dist.Y) < 24.*world.TileSize {
				anim.Update()
			}
		}
	}
}

func DrawSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsDrawable) {
		draw := result.Components[myecs.Drawable]
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		bkey, okB := result.Components[myecs.Batch].(string)
		if okT && okB && !tran.Hide {
			dist := camera.Cam.Pos.Sub(tran.Pos)
			if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
				if batcher, ok := img.Batchers[bkey]; ok {
					if spr, okS := draw.(*pixel.Sprite); okS {
						spr.DrawColorMask(batcher.Batch(), tran.Mat, tran.Mask)
					} else if anim, okA := draw.(*reanimator.Tree); okA {
						anim.DrawColorMask(batcher.Batch(), tran.Mat, tran.Mask)
					}
				}
			}
		}
	}
}