package particles

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/internal/util"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"fmt"
	"github.com/faiface/pixel"
)

func BlockParticles(pos pixel.Vec, biome string) {
	c := random.Effects.Intn(3) + 3
	for i := 0; i < c; i++ {
		CreateParticle(fmt.Sprintf("%s_%d", biome, random.Effects.Intn(8)), pos, 5., 0.75, 0.1, true)
	}
}

func BiomeParticles(pos pixel.Vec, biome string, min, max int, variance, dur, durVar float64, collide bool) {
	c := random.Effects.Intn(max - min + 1) + min
	for i := 0; i < c; i++ {
		CreateParticle(fmt.Sprintf("%s_%d", biome, random.Effects.Intn(8)), pos, variance, dur, durVar, collide)
	}
}

func CreateRandomStaticParticles(min, max int, keys []string, orig pixel.Vec, variance, dur, durVar float64) {
	c := random.Effects.Intn(max - min + 1) + min
	for i := 0; i < c; i++ {
		tran := transform.NewTransform()
		tran.Pos = util.RandomPosition(orig, variance, random.Effects)
		if random.Effects.Intn(2) == 0 {
			tran.Flip = true
		}
		if random.Effects.Intn(2) == 0 {
			tran.Flop = true
		}
		nDur := dur + (random.Effects.Float64() - 0.5) * durVar
		key := keys[random.Effects.Intn(len(keys))]
		myecs.Manager.NewEntity().
			AddComponent(myecs.Transform, tran).
			AddComponent(myecs.Sprite, img.Batchers[constants.ParticleKey].GetSprite(key)).
			AddComponent(myecs.Batch, constants.ParticleKey).
			AddComponent(myecs.Temp, timing.New(nDur))
	}
}

func CreateRandomParticles(min, max int, keys []string, orig pixel.Vec, variance, dur, durVar float64, collide bool) {
	c := random.Effects.Intn(max - min + 1) + min
	for i := 0; i < c; i++ {
		key := keys[random.Effects.Intn(len(keys))]
		CreateParticle(key, orig, variance, dur, durVar, collide)
	}
}

func CreateParticle(key string, orig pixel.Vec, variance, dur, durVar float64, collide bool) {
	phys, tran := util.RandomVelocity(orig, variance, random.Effects)
	if random.Effects.Intn(2) == 0 {
		tran.Flip = true
	}
	if random.Effects.Intn(2) == 0 {
		tran.Flop = true
	}
	nDur := dur + (random.Effects.Float64() - 0.5) * durVar
	spr := img.Batchers[constants.ParticleKey].GetSprite(key)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Transform, tran).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Sprite, spr).
		AddComponent(myecs.Batch, constants.ParticleKey).
		AddComponent(myecs.Temp, timing.New(nDur))
	if collide {
		coll := data.NewCollider(spr.Frame(), true, true)
		e.AddComponent(myecs.Collision, coll)
	}
}