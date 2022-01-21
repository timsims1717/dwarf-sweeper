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
	created    bool
	Health     *data.Health
	faceLeft   bool
	roosting   bool
	roost      *cave.Tile
	flight     pixel.Vec
}

func (b *Bat) Update() {
	if !b.Health.Dazed && !b.Health.Dead {
		b.Physics.GravityOff = true
		if b.roosting {
			b.flight = pixel.ZV
			if b.roost == nil || !b.roost.Solid() {
				b.roosting = false
			}
			// todo: disturbances, bombs, etc
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
				b.flight.X -= batAcc * timing.DT
			} else if b.Physics.LeftBound && b.flight.X < 0. {
				b.flight.X += batAcc * timing.DT
			}
			if b.Physics.BottomBound && b.Physics.TopBound {
				b.flight.Y = 0.
			} else if b.Physics.BottomBound && b.flight.Y < 0. {
				b.flight.Y += batAcc * timing.DT
			} else if b.Physics.TopBound && b.flight.Y > 0. {
				b.flight.Y -= batAcc * timing.DT
			}
			b.flight = util.Normalize(b.flight)
			b.Physics.SetVelX(b.flight.X*batSpeed, batAcc)
			b.Physics.SetVelY(b.flight.Y*batSpeed, batAcc)
			b.faceLeft = b.flight.X < 0.
			// randomize a bit
			b.flight.X += (random.Effects.Float64() - 0.5) * timing.DT
			b.flight.Y += (random.Effects.Float64() - 0.5) * timing.DT
			b.flight = util.Normalize(b.flight)
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
		}
		b.Transform.Flip = b.faceLeft
	} else {
		b.Physics.GravityOff = false
		b.roosting = false
		b.flight = pixel.ZV
	}
	if b.Health.Dead {
		b.Delete()
	}
}

func (b *Bat) Create(pos pixel.Vec) {
	b.Transform = transform.NewTransform()
	b.Transform.Pos = pos
	tPos := pos
	tPos.Y += world.TileSize
	b.roost = Descent.Cave.GetTile(tPos)
	b.Physics = physics.New()
	b.Physics.GravityOff = true
	b.Health = &data.Health{
		Max:          2,
		Curr:         2,
		TempInvTimer: timing.New(0.5),
		TempInvSec:   0.5,
		Immune:       data.EnemyImmunity,
	}
	b.created = true
	b.Reanimator = reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("bat_fly", []*pixel.Sprite{img.Batchers[constants.EntityKey].GetFrame("bat_fly", 0)}, reanimator.Loop)).
		AddAnimation(reanimator.NewAnimFromSprites("bat_roost", img.Batchers[constants.EntityKey].GetAnimation("bat_roost").S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("bat_fly", img.Batchers[constants.EntityKey].GetAnimation("bat_fly").S, reanimator.Loop)).
		SetChooseFn(func() int {
			if b.Health.Dazed {
				return 0
			} else if b.roosting {
				return 1
			} else {
				return 2
			}
		}), "bat_roost")
	b.Collider = data.NewCollider(pixel.R(0., 0., 16., 16.), true, false)
	b.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, b).
		AddComponent(myecs.Transform, b.Transform).
		AddComponent(myecs.Animation, b.Reanimator).
		AddComponent(myecs.Physics, b.Physics).
		AddComponent(myecs.Health, b.Health).
		AddComponent(myecs.Collision, b.Collider).
		AddComponent(myecs.Batch, constants.EntityKey)
}

func (b *Bat) Delete() {
	b.Health.Delete()
	myecs.Manager.DisposeEntity(b.Entity)
}
