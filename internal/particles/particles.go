package particles

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"fmt"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"math"
)

func BlockParticles(pos pixel.Vec, biome string) {
	c := random.Effects.Intn(3) + 3
	for i := 0; i < c; i++ {
		CreateParticle(fmt.Sprintf("%s_%d", biome, random.Effects.Intn(8)), pos, 5., 5., math.Pi*0.5, 1.5, 120., 5.0, 0.75, 0.2, true)
	}
}

func BiomeParticles(orig pixel.Vec, biome string, min, max int, varX, varY, angle, angleVar, force, forceVar, dur, durVar float64, collide bool) {
	c := random.Effects.Intn(max-min+1) + min
	for i := 0; i < c; i++ {
		CreateParticle(fmt.Sprintf("%s_%d", biome, random.Effects.Intn(8)), orig, varX, varY, angle, angleVar, force, forceVar, dur, durVar, collide)
	}
}

func CreateRandomStaticParticles(min, max int, keys []string, orig pixel.Vec, variance, dur, durVar float64, fade bool) {
	c := random.Effects.Intn(max-min+1) + min
	for i := 0; i < c; i++ {
		tran := transform.New()
		tran.Pos = data.RandomPosition(orig, variance, variance, random.Effects)
		if random.Effects.Intn(2) == 0 {
			tran.Flip = true
		}
		if random.Effects.Intn(2) == 0 {
			tran.Flop = true
		}
		switch random.Effects.Intn(4) {
		case 0:
			tran.Rot = math.Pi
		case 1:
			tran.Rot = math.Pi * 0.5
		case 2:
			tran.Rot = math.Pi * -0.5
		}
		nDur := dur + (random.Effects.Float64()-0.5)*durVar
		key := keys[random.Effects.Intn(len(keys))]
		e := myecs.Manager.NewEntity().
			AddComponent(myecs.Transform, tran).
			AddComponent(myecs.Drawable, img.Batchers[constants.ParticleKey].GetSprite(key)).
			AddComponent(myecs.Batch, constants.ParticleKey).
			AddComponent(myecs.Temp, timing.New(nDur))
		if fade {
			myecs.AddEffect(e, data.NewFadeOut(colornames.White, nDur))
		}
	}
}

func CreateStaticParticle(key string, orig pixel.Vec, variance, dur, durVar float64, fade bool) {
	tran := transform.New()
	tran.Pos = data.RandomPosition(orig, variance, variance, random.Effects)
	if random.Effects.Intn(2) == 0 {
		tran.Flip = true
	}
	if random.Effects.Intn(2) == 0 {
		tran.Flop = true
	}
	switch random.Effects.Intn(4) {
	case 0:
		tran.Rot = math.Pi
	case 1:
		tran.Rot = math.Pi * 0.5
	case 2:
		tran.Rot = math.Pi * -0.5
	}
	nDur := dur + (random.Effects.Float64()-0.5)*durVar
	e := myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, tran).
		AddComponent(myecs.Drawable, img.Batchers[constants.ParticleKey].GetSprite(key)).
		AddComponent(myecs.Batch, constants.ParticleKey).
		AddComponent(myecs.Temp, timing.New(nDur))
	if fade {
		myecs.AddEffect(e, data.NewFadeOut(colornames.White, nDur))
	}
}

func CreateRandomParticles(min, max int, keys []string, orig pixel.Vec, varX, varY, angle, angleVar, force, forceVar, dur, durVar float64, collide bool) {
	c := random.Effects.Intn(max-min+1) + min
	for i := 0; i < c; i++ {
		key := keys[random.Effects.Intn(len(keys))]
		CreateParticle(key, orig, varX, varY, angle, angleVar, force, forceVar, dur, durVar, collide)
	}
}

func CreateParticle(key string, orig pixel.Vec, varX, varY, angle, angleVar, force, forceVar, dur, durVar float64, collide bool) {
	phys, tran := data.RandomPosAndVel(orig, varX, varY, angle, angleVar, force, forceVar, random.Effects)
	if random.Effects.Intn(2) == 0 {
		tran.Flip = true
	}
	if random.Effects.Intn(2) == 0 {
		tran.Flop = true
	}
	switch random.Effects.Intn(4) {
	case 0:
		tran.Rot = math.Pi
	case 1:
		tran.Rot = math.Pi * 0.5
	case 2:
		tran.Rot = math.Pi * -0.5
	}
	nDur := dur + (random.Effects.Float64()-0.5)*durVar
	spr := img.Batchers[constants.ParticleKey].GetSprite(key)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Transform, tran).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Drawable, spr).
		AddComponent(myecs.Batch, constants.ParticleKey).
		AddComponent(myecs.Temp, timing.New(nDur))
	if collide {
		coll := data.NewCollider(pixel.R(0., 0., spr.Frame().W(), spr.Frame().H()), data.GroundOnly)
		e.AddComponent(myecs.Collision, coll)
	}
}
