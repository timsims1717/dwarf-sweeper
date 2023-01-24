package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

const (
	batSpeed = 50.
	batAcc   = 10.
)

type Bat struct {
	Transform  *transform.Transform
	Physics    *physics.Physics
	Collider   *data.Collider
	Reanimator *reanimator.Tree
	Entity     *ecs.Entity
	Health     *data.Health
	faceLeft   bool
	roosting   bool
	roosted1   bool
	roost      *cave.Tile
	flight     pixel.Vec
	Evil       bool
}

func CreateBat(c *cave.Cave, pos pixel.Vec) *Bat {
	b := &Bat{}
	if random.CaveGen.Intn(100) < util.Min(c.Level - 5, 20) {
		b.Evil = true
	}
	b.Transform = transform.New().WithID("bat")
	b.Transform.Pos = pos
	tPos := pos
	tPos.Y += world.TileSize
	b.roost = c.GetTile(tPos)
	b.Physics = physics.New()
	b.Physics.GravityOff = true
	b.Health = &data.Health{
		Max:          2,
		Curr:         2,
		TempInvTimer: timing.New(0.5),
		TempInvSec:   0.5,
		Immune:       data.EnemyImmunity,
	}
	b.Reanimator = reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("bat_daze_air", img.Batchers[constants.EntityKey].GetAnimation("bat_daze_air").S, reanimator.Hold).
			SetTrigger(0, func() {
				PlaySqueak()
			})).
		AddAnimation(reanimator.NewAnimFromSprites("bat_daze_ground", img.Batchers[constants.EntityKey].GetAnimation("bat_daze_ground").S, reanimator.Hold).
			SetTrigger(0, func() {
				PlaySqueak()
			})).
		AddAnimation(reanimator.NewAnimFromSprites("bat_roost", img.Batchers[constants.EntityKey].GetAnimation("bat_roost").S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("bat_fly", img.Batchers[constants.EntityKey].GetAnimation("bat_fly").S, reanimator.Loop).
			SetTrigger(0, func() {
				if b.Transform.Load && b.roosted1 {
					sfx.SoundPlayer.PlaySound("wingflap", 0.)
				}
			})).
		AddAnimation(reanimator.NewAnimFromSprites("evil_bat_fly", img.Batchers[constants.EntityKey].GetAnimation("evil_bat_fly").S, reanimator.Loop).
			SetTrigger(0, func() {
				if b.Transform.Load && b.roosted1 {
					sfx.SoundPlayer.PlaySound("wingflap", 0.)
				}
			})).
		SetChooseFn(func() int {
			if b.Health.Dazed || b.Health.Dead {
				if b.Physics.Grounded {
					return 1
				} else {
					return 0
				}
			} else if b.roosting {
				return 2
			} else if b.Evil {
				return 4
			} else {
				return 3
			}
		}), "bat_roost")
	b.Collider = data.NewCollider(pixel.R(0., 0., 16., 16.), data.Critter)
	b.Collider.Debug = true
	b.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, b.Transform).
		AddComponent(myecs.Physics, b.Physics).
		AddComponent(myecs.Collision, b.Collider).
		AddComponent(myecs.Health, b.Health).
		AddComponent(myecs.Update, data.NewFrameFunc(b.Update)).
		AddComponent(myecs.Animation, b.Reanimator).
		AddComponent(myecs.Drawable, b.Reanimator).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	return b
}

func (b *Bat) Update() bool {
	if !b.Health.Dazed && !b.Health.Dead {
		b.Physics.GravityOff = true
		squeak := false
		if b.roosting {
			b.roosted1 = true
			b.flight = pixel.ZV
			if b.roost == nil || !b.roost.Solid() {
				b.roosting = false
			}
			if Descent.Cave.DestroyedWithin(b.roost.RCoords, 2, 1) {
				b.roosting = false
			}
			if random.Effects.Intn(100. * int(1 / timing.DT)) == 0 {
				squeak = true
			}
		}
		if !b.roosting {
			if b.flight.X == 0. && b.flight.Y == 0. {
				b.flight.X = (random.Effects.Float64() - 0.5) * 2.
				b.flight.Y = (random.Effects.Float64() - 0.5) * 2.
				b.flight = util.Normalize(b.flight)
			}
			if b.Physics.RightBound && b.Physics.LeftBound {
				b.flight.X = 0.
			} else if b.Physics.RightBound && b.flight.X > 0. {
				b.flight.X = 0.
			} else if b.Physics.LeftBound && b.flight.X < 0. {
				b.flight.X = 0.
			}
			if b.Physics.BottomBound && b.Physics.TopBound {
				b.flight.Y = 0.
			} else if b.Physics.BottomBound && b.flight.Y < 0. {
				b.flight.Y = 0.
			} else if b.Physics.TopBound && b.flight.Y > 0. {
				b.flight.Y = 0.
			}
			b.flight.Y += (random.Effects.Float64() - 0.5) * 2. * timing.DT
			b.flight.X += (random.Effects.Float64() - 0.5) * 2. * timing.DT
			b.flight = util.Normalize(b.flight)
			b.Physics.SetVelX(b.flight.X*batSpeed, batAcc)
			b.Physics.SetVelY(b.flight.Y*batSpeed, batAcc)
			b.faceLeft = b.flight.X < 0.
			// check for roost
			tPos := b.Transform.Pos
			tPos.Y += world.TileSize * 0.51
			pot := Descent.Cave.GetTile(tPos)
			if pot != nil && pot.Solid() && random.Effects.Intn(40) == 0 {
				b.roost = pot
				b.roosting = true
				b.flight = pixel.ZV
				b.Transform.Pos.Y = Descent.Cave.GetTile(b.Transform.Pos).Transform.Pos.Y
				b.Physics.CancelMovement()
			}
			if random.Effects.Intn(10. * int(1 / timing.DT)) == 0 {
				squeak = true
			}
			if b.Evil {
				myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
					SourceID:  b.Transform.ID,
					Center:    b.Transform.Pos,
					Radius:    4.,
					Amount:    1,
					Dazed:     1.,
					Knockback: 8.,
					Type:      data.Enemy,
				})
			}
		}
		b.Transform.Flip = b.faceLeft
		if squeak && b.Transform.Load {
			PlaySqueak()
		}
	} else {
		b.Physics.GravityOff = false
		b.roosting = false
		b.flight = pixel.ZV
	}
	if b.Health.Dead {
		b.Entity.RemoveComponent(myecs.Update)
		b.Entity.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
			myecs.AddEffect(b.Entity, data.NewBlink(2.))
			return true
		}, 2.))
		b.Entity.AddComponent(myecs.Temp, timing.New(4.))
	}
	return false
}