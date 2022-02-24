package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/internal/util"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
	"math"
)

func CreateApple(pos pixel.Vec) {
	spr := img.Batchers[constants.EntityKey].Sprites["apple"]
	fn := func(pos pixel.Vec, d *Dwarf) bool {
		if d.Health.Curr < d.Health.Max {
			d.Entity.AddComponent(myecs.Healing, &data.Heal{
				Amount: 1,
			})
			sfx.SoundPlayer.PlaySound("bite", 1.0)
			return true
		}
		return false
	}
	CreateCollectible(pos, fn, spr)
}

func CreateGem(pos pixel.Vec) {
	spr := img.Batchers[constants.EntityKey].Sprites["gem_diamond"]
	fn := func(pos pixel.Vec, d *Dwarf) bool {
		player.OverallStats.CaveGemsFound++
		d.Player.Stats.CaveGemsFound++
		particles.CreateRandomStaticParticles(2, 4, []string{"sparkle_plus_0", "sparkle_plus_1", "sparkle_plus_2", "sparkle_x_0", "sparkle_x_1", "sparkle_x_2"}, pos, 10.0, 1.0, 0.5)
		sfx.SoundPlayer.PlaySound("clink", 1.0)
		return true
	}
	CreateCollectible(pos, fn, spr)
}

func CreateCollectible(pos pixel.Vec, fn func(pixel.Vec, *Dwarf) bool, spr *pixel.Sprite) {
	e := myecs.Manager.NewEntity()
	c := &Collectible{
		OnCollect: fn,
		Timer:     timing.New(1.),
		AutoCollect: true,
	}
	phys, trans := util.RandomPosAndVel(pos, 0., 0., math.Pi*0.5, math.Pi*0.25, 125., 10., random.Effects)
	coll := data.NewCollider(pixel.R(0., 0., spr.Frame().W(), spr.Frame().H()), true, false)
	coll.Debug = true
	hp := &data.SimpleHealth{Immune: data.ItemImmunity}
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Collision, coll).
		AddComponent(myecs.Health, hp).
		AddComponent(myecs.Temp, timing.New(10.)).
		AddComponent(myecs.Collect, c).
		AddComponent(myecs.Drawable, spr).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
			myecs.AddEffect(e, data.NewBlink(2.))
			return true
		}, 8.)).
		AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
			if hp.Dead {
				e.AddComponent(myecs.Temp, myecs.ClearFlag(true))
			}
			return false
		}))
}