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
		if anim, ok := result.Components[myecs.Animation].(*reanimator.Tree); ok {
			anim.Update()
		}
	}
}

func AnimationDraw() {
	for _, result := range myecs.Manager.Query(myecs.HasAnimDrawing) {
		anim, okA := result.Components[myecs.Animation].(*reanimator.Tree)
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		bkey, okB := result.Components[myecs.Batch].(string)
		if okA && okT && okB && anim.CurrentSprite() != nil {
			if batcher, ok := img.Batchers[bkey]; ok {
				anim.CurrentSprite().Draw(batcher.Batch(), tran.Mat)
			}
		}
	}
}

func SpriteDraw() {
	for _, result := range myecs.Manager.Query(myecs.HasSprDrawing) {
		spr, okS := result.Components[myecs.Sprite].(*pixel.Sprite)
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		bkey, okB := result.Components[myecs.Batch].(string)
		if okS && okT && okB {
			if batcher, ok := img.Batchers[bkey]; ok {
				spr.Draw(batcher.Batch(), tran.Mat)
			}
		}
	}
}