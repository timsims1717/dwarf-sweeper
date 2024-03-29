package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"math"
)

const (
	mmSpeed     = 40.
	mmAcc       = 5.
	mmAtkWait   = 2.
	mmKnockback = 8.
)

var (
	mmAngle = pixel.V(1., 1.).Angle()
)

type MadMonk struct {
	Transform  *transform.Transform
	Physics    *physics.Physics
	Reanimator *reanimator.Tree
	Entity     *ecs.Entity
	created    bool
	Health     *data.Health
	faceLeft   bool
	Attack     bool
	AtkTimer   *timing.Timer
}

func (m *MadMonk) Update() {
	if !m.Health.Dazed && !m.Health.Dead {
		m.AtkTimer.Update()
		if m.Physics.Grounded && !m.Attack {
			ownCoords := Descent.GetCave().GetTile(m.Transform.Pos).RCoords
			dwarf := Descent.GetClosestPlayer(m.Transform.Pos)
			playerPos := dwarf.Transform.Pos
			playerCoords := Descent.GetTile(playerPos).RCoords
			ownPos := m.Transform.Pos
			if math.Abs(ownPos.X-playerPos.X) <= world.TileSize && ownCoords.Y == playerCoords.Y && m.AtkTimer.Done() {
				m.Attack = true
				m.faceLeft = ownCoords.X > playerCoords.X
			}
			if ownCoords.X > playerCoords.X {
				m.faceLeft = true
				m.Physics.SetVelX(-mmSpeed, mmAcc)
			} else if ownCoords.X < playerCoords.X {
				m.faceLeft = false
				m.Physics.SetVelX(mmSpeed, mmAcc)
			}
		}
		m.Transform.Flip = m.faceLeft
	}
	if m.Health.Dead {
		m.Delete()
	}
}

func (m *MadMonk) Create(pos pixel.Vec) {
	m.AtkTimer = timing.New(mmAtkWait)
	m.Transform = transform.New().WithID("mad-monk")
	m.Transform.Pos = pos
	m.Physics = physics.New()
	m.Physics.Terminal = 100.
	m.Health = &data.Health{
		Max:          2,
		Curr:         2,
		Dazed:        true,
		DazedTimer:   timing.New(3.),
		TempInvTimer: timing.New(0.5),
		TempInvSec:   0.5,
		Immune:       data.EnemyImmunity,
	}
	m.created = true
	m.Reanimator = reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("mm_attack", img.Batchers[constants.EntityKey].Animations["mm_attack"].S, reanimator.Tran).
		SetTrigger(3, func() {
				m.AtkTimer = timing.New(mmAtkWait)
				atkPos := m.Transform.Pos
				if m.faceLeft {
					atkPos.X -= world.TileSize * 0.5
				} else {
					atkPos.X += world.TileSize * 0.5
				}
				myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
					SourceID:  m.Transform.ID,
					Center:    atkPos,
					Radius:    world.TileSize,
					Amount:    1,
					Dazed:     3.,
					Angle:     &mmAngle,
					Knockback: mmKnockback,
					Type:      data.Enemy,
				})
			}).
		SetTrigger(5, func() {
				m.Attack = false
			})).
		AddAnimation(reanimator.NewAnimFromSprites("mm_fall", img.Batchers[constants.EntityKey].Animations["mm_fall"].S, reanimator.Loop)).
		AddAnimation(reanimator.NewAnimFromSprites("mm_walk", img.Batchers[constants.EntityKey].Animations["mm_walk"].S, reanimator.Loop)).
		AddAnimation(reanimator.NewAnimFromSprites("mm_idle", img.Batchers[constants.EntityKey].Animations["mm_idle"].S, reanimator.Hold)).
		SetChooseFn(func() int {
			if m.Attack {
				return 0
			} else if !m.Physics.Grounded {
				return 1
			} else if m.Physics.IsMovingX() {
				return 2
			} else {
				return 3
			}
		}), "mm_idle")
	m.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, m).
		AddComponent(myecs.Transform, m.Transform).
		AddComponent(myecs.Animation, m.Reanimator).
		AddComponent(myecs.Drawable, m.Reanimator).
		AddComponent(myecs.Physics, m.Physics).
		AddComponent(myecs.Health, m.Health).
		AddComponent(myecs.Collision, &data.Collider{
			Hitbox: pixel.R(0., 0., 16., 16.),
		}).
		AddComponent(myecs.Batch, constants.EntityKey)
}

func (m *MadMonk) Delete() {
	myecs.Manager.DisposeEntity(m.Entity)
}
