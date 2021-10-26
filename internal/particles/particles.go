package particles

import (
	"dwarf-sweeper/internal/constants"
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
	c := random.Effects.Intn(3) + 4
	for i := 0; i < c; i++ {
		phys, tran := util.RandomVelocity(pos, 1.0, random.Effects)
		if random.Effects.Intn(2) == 0 {
			tran.Flip = true
		}
		if random.Effects.Intn(2) == 0 {
			tran.Flop = true
		}
		myecs.Manager.NewEntity().
			AddComponent(myecs.Transform, tran).
			AddComponent(myecs.Physics, phys).
			AddComponent(myecs.Sprite, img.Batchers[constants.ParticleKey].GetSprite(fmt.Sprintf("%s_%d", biome, random.Effects.Intn(8)))).
			AddComponent(myecs.Batch, constants.ParticleKey).
			AddComponent(myecs.Temp, timing.New(0.75))
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