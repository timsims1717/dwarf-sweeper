package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

type Mine struct {
	Transform  *transform.Transform
	Timer      *timing.FrameTimer
	Tile       *cave.Tile
	created    bool
	explode    bool
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
}

func (m *Mine) Update() {
	if m.created {
		if m.Timer.UpdateDone() || m.explode {
			CaveBombsLeft--
			CaveBlownUpBombs++
			area := []pixel.Vec{ m.Transform.Pos }
			for _, n := range m.Tile.SubCoords.Neighbors() {
				t := m.Tile.Chunk.Get(n)
				t.Destroy(false)
				area = append(area, t.Transform.Pos)
			}
			myecs.Manager.NewEntity().
				AddComponent(myecs.AreaDmg, &data.AreaDamage{
					Center:         m.Transform.Pos,
					Radius:         MineBaseRadius * world.TileSize,
					Amount:         1,
					Dazed:          3.,
					Knockback:      MineBaseKnockback,
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
	m.Timer = timing.New(constants.MineFuse)
	m.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("mine_1", img.Batchers[constants.EntityKey].Animations["mine_1"].S, reanimator.Hold, nil),
			reanimator.NewAnimFromSprites("mine_2", img.Batchers[constants.EntityKey].Animations["mine_2"].S, reanimator.Tran, map[int]func() {
				2: func() {
					m.explode = true
				},
			}),
		),
		Check: func() int {
			if constants.MineFuse * 0.5 > m.Timer.Elapsed() {
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
		AddComponent(myecs.Batch, constants.EntityKey)
}

func (m *Mine) Delete() {
	myecs.Manager.DisposeEntity(m.entity)
}