package dungeon

import (
	"dwarf-sweeper/internal/character"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"github.com/google/uuid"
)

const (
	mmSpeed    = 40.
	mmAcc      = 5.
	mmAtkWait  = 2.
	entityBKey = "entities"
)

type MadMonk struct {
	ID         uuid.UUID
	Transform  *transform.Transform
	Physics    *physics.Physics
	Reanimator *reanimator.Tree
	Entity     *ecs.Entity
	created    bool
	Health     *character.Health
	faceLeft   bool
	Attack     bool
	AtkTimer   *timing.FrameTimer
}

func (m *MadMonk) Update() {
	if Dungeon.GetCave().PointLoaded(m.Transform.Pos) {
		if !m.Health.Dazed && !m.Health.Dead {
			m.AtkTimer.Update()
			if m.Physics.Grounded && !m.Attack {
				ownPos := Dungeon.GetCave().GetTile(m.Transform.Pos).RCoords
				playerPos := Dungeon.GetPlayerTile().RCoords
				if util.Abs(ownPos.X - playerPos.X) < 2 && ownPos.Y == playerPos.Y && m.AtkTimer.Done() {
					m.Attack = true
					m.faceLeft = ownPos.X > playerPos.X
				}
				if ownPos.X > playerPos.X {
					m.faceLeft = true
					m.Physics.SetVelX(-mmSpeed, mmAcc)
				} else if ownPos.X < playerPos.X {
					m.faceLeft = false
					m.Physics.SetVelX(mmSpeed, mmAcc)
				}
			}
			m.Transform.Flip = m.faceLeft
		}
	}
	if m.Health.Dead {
		m.Health.Delete()
		myecs.LazyDelete(m.Entity)
	}
}

func (m *MadMonk) Create(pos pixel.Vec) {
	m.ID = uuid.New()
	m.AtkTimer = timing.New(mmAtkWait)
	m.Transform = transform.NewTransform()
	m.Transform.Pos = pos
	m.Physics = physics.New()
	m.Physics.Terminal = 100.
	m.Health = &character.Health{
		Max:          2,
		Curr:         2,
		Dead:         false,
		Dazed:        true,
		DazedTimer:   timing.New(3.),
		TempInv:      true,
		TempInvSec:   1.,
		TempInvTimer: timing.New(1.),
	}
	m.created = true
	m.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("mm_attack", img.Batchers[entityBKey].Animations["mm_attack"].S, reanimator.Tran, map[int]func(){
				3: func() {
					m.AtkTimer = timing.New(mmAtkWait)
					ownPos := Dungeon.GetCave().GetTile(m.Transform.Pos).RCoords
					playerPos := Dungeon.GetPlayerTile().RCoords
					if util.Abs(ownPos.X - playerPos.X) < 2 && ownPos.Y == playerPos.Y {
						Dungeon.GetPlayer().Entity.AddComponent(myecs.Damage, &character.Damage{
							Amount:    1,
							Dazed:     1.,
							Knockback: 8.,
							Source:    m.Transform.Pos,
						})
					}
				},
				5: func() {
					m.Attack = false
				},
			}),
			reanimator.NewAnimFromSprites("mm_fall", img.Batchers[entityBKey].Animations["mm_fall"].S, reanimator.Loop, nil),
			reanimator.NewAnimFromSprites("mm_walk", img.Batchers[entityBKey].Animations["mm_walk"].S, reanimator.Loop, nil),
			reanimator.NewAnimFromSprites("mm_idle", img.Batchers[entityBKey].Animations["mm_idle"].S, reanimator.Hold, nil),
		),
		Check: func() int {
			if m.Attack {
				return 0
			} else if !m.Physics.Grounded {
				return 1
			} else if m.Physics.IsMovingX() {
				return 2
			} else {
				return 3
			}
		},
	}, "mm_idle")
	m.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, m).
		AddComponent(myecs.Transform, m.Transform).
		AddComponent(myecs.Animation, m.Reanimator).
		AddComponent(myecs.Physics, m.Physics).
		AddComponent(myecs.Health, m.Health).
		AddComponent(myecs.Collision, myecs.Collider{}).
		AddComponent(myecs.Batch, entityBKey)
}