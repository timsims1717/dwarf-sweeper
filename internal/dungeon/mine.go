package dungeon

import (
	"dwarf-sweeper/internal/character"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

const (
	MineKnockback = 25.
)

type Mine struct {
	EID        int
	Transform  *transform.Transform
	Timer      *timing.FrameTimer
	FuseLength float64
	Tile       *Tile
	created    bool
	explode    bool
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
}

func (m *Mine) Update() {
	if m.created && m.explode {
		if m.Timer.UpdateDone() {
			area := []pixel.Vec{ m.Transform.Pos }
			for _, n := range m.Tile.SubCoords.Neighbors() {
				t := m.Tile.Chunk.Get(n)
				t.Destroy(false)
				area = append(area, t.Transform.Pos)
			}
			myecs.Manager.NewEntity().
				AddComponent(myecs.AreaDmg, &character.AreaDamage{
					Area:           area,
					Amount:         1,
					Dazed:          3.,
					Knockback:      MineKnockback,
					KnockbackDecay: true,
					Source:         m.Transform.Pos,
					Override:       true,
				})
			vfx.CreateExplosion(m.Tile.Transform.Pos)
			sfx.SoundPlayer.PlaySound("blast1", -1.0)
			m.Delete()
		}
	}
}

func (m *Mine) Create(pos pixel.Vec) {
	m.Transform = transform.NewTransform()
	m.Transform.Pos = pos
	m.created = true
	m.Timer = timing.New(0.25)
	m.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("mine_1", img.Batchers[entityBKey].Animations["mine_1"].S, reanimator.Hold, nil),
			reanimator.NewAnimFromSprites("mine_2", img.Batchers[entityBKey].Animations["mine_2"].S, reanimator.Tran, map[int]func() {
				2: func() {
					m.explode = true
				},
			}),
		),
		Check: func() int {
			if m.FuseLength * 0.5 > m.Timer.Elapsed() {
				return 0
			} else {
				return 1
			}
		},
	}, "mine_1")
	m.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, m).
		AddComponent(myecs.Transform, m.Transform).
		AddComponent(myecs.Animation, m.Reanimator).
		AddComponent(myecs.Batch, entityBKey)
	Dungeon.AddEntity(m)
}

func (m *Mine) Delete() {
	myecs.Manager.DisposeEntity(m.entity)
	Dungeon.RemoveEntity(m.EID)
}

func (m *Mine) SetId(i int) {
	m.EID = i
}