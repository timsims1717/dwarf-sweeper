package particles

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

var (
	explosion     *reanimator.Anim
	bigExplosion  *reanimator.Anim
	hugeExplosion *reanimator.Anim
	dazed         *reanimator.Anim
)

func Initialize() {
	explosion = reanimator.NewAnimFromSprites("explosion", img.Batchers[constants.ExpKey].GetAnimation("explosion").S, reanimator.Done)
	bigExplosion = reanimator.NewAnimFromSprites("big-explosion", img.Batchers[constants.BigExpKey].GetAnimation("explosion").S, reanimator.Done)
	hugeExplosion = reanimator.NewAnimFromSprites("huge-explosion", img.Batchers[constants.HugeExpKey].GetAnimation("explosion").S, reanimator.Done)
	dazed = reanimator.NewAnimFromSprites("dazed", img.Batchers[constants.ParticleKey].GetAnimation("dazed").S, reanimator.Loop)
}

func CreateSmallExplosion(vec pixel.Vec) {
	e := myecs.Manager.NewEntity()
	anim := reanimator.NewSimple(explosion.Copy().SetTrigger(6, func(_ *reanimator.Anim, _ string, _ int) {
		myecs.Manager.DisposeEntity(e)
	}))
	trans := transform.New()
	trans.Pos = vec
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Temp, timing.New(0.6)).
		AddComponent(myecs.Animation, anim).
		AddComponent(myecs.Drawable, anim).
		AddComponent(myecs.Batch, constants.ExpKey)
}

func CreateBigExplosion(vec pixel.Vec) {
	e := myecs.Manager.NewEntity()
	anim := reanimator.NewSimple(bigExplosion.Copy().SetTrigger(9, func(_ *reanimator.Anim, _ string, _ int) {
		myecs.Manager.DisposeEntity(e)
	}))
	trans := transform.New()
	trans.Pos = vec
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Temp, timing.New(0.9)).
		AddComponent(myecs.Animation, anim).
		AddComponent(myecs.Drawable, anim).
		AddComponent(myecs.Batch, constants.BigExpKey)
}

func CreateHugeExplosion(vec pixel.Vec) {
	e := myecs.Manager.NewEntity()
	anim := reanimator.NewSimple(hugeExplosion.Copy().SetTrigger(22, func(_ *reanimator.Anim, _ string, _ int) {
		myecs.Manager.DisposeEntity(e)
	}))
	trans := transform.New()
	trans.Pos = vec
	trans.Pos.Y += world.TileSize * 3.
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Temp, timing.New(2.5)).
		AddComponent(myecs.Animation, anim).
		AddComponent(myecs.Drawable, anim).
		AddComponent(myecs.Batch, constants.HugeExpKey)
}

func DazedAnimation() *reanimator.Tree {
	return reanimator.NewSimple(dazed)
}