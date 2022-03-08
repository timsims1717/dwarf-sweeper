package systems

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
)

func AnimationSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasAnimation) {
		anim, ok := result.Components[myecs.Animation].(*reanimator.Tree)
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		if ok && okT && tran.Load {
			anim.Update()
		}
	}
}

func DrawSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsDrawable) {
		draw := result.Components[myecs.Drawable]
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		bkey, okB := result.Components[myecs.Batch].(string)
		if okT && okB && tran.Load {
			if batcher, ok := img.Batchers[bkey]; ok {
				if sprStr, okSS := draw.(string); okSS {
					batcher.DrawSpriteColor(sprStr, tran.Mat, tran.Mask)
				} else if spr, okS := draw.(*pixel.Sprite); okS {
					spr.DrawColorMask(batcher.Batch(), tran.Mat, tran.Mask)
				} else if anim, okA := draw.(*reanimator.Tree); okA {
					anim.DrawColorMask(batcher.Batch(), tran.Mat, tran.Mask)
				}
			}
		}
	}
}