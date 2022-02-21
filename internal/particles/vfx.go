package particles

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
)

var explosion *reanimator.Anim

func Initialize() {
	explosion = reanimator.NewAnimFromSprites("explosion", img.Batchers[constants.ExpKey].GetAnimation("explosion").S, reanimator.Done)
}

func CreateSmallExplosion(vec pixel.Vec) {
	e := myecs.Manager.NewEntity()
	anim := reanimator.NewSimple(explosion.SetTrigger(6, func(_ *reanimator.Anim, _ string, _ int) {
		myecs.Manager.DisposeEntity(e)
	}))
	trans := transform.New()
	trans.Pos = vec
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Animation, anim).
		AddComponent(myecs.Drawable, anim).
		AddComponent(myecs.Batch, constants.ExpKey)
}