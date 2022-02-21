package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/player"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

const (
	MineFuse = 0.4
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
			player.CaveBombsLeft--
			player.CaveBlownUpBombs++
			area := []pixel.Vec{m.Transform.Pos}
			for _, n := range m.Tile.SubCoords.Neighbors() {
				t := m.Tile.Chunk.Get(n)
				t.Destroy(false)
				area = append(area, t.Transform.Pos)
			}
			myecs.Manager.NewEntity().
				AddComponent(myecs.AreaDmg, &data.AreaDamage{
					SourceID:       m.Transform.ID,
					Center:         m.Transform.Pos,
					Radius:         MineBaseRadius * world.TileSize,
					Amount:         1,
					Dazed:          3.,
					Knockback:      MineBaseKnockback,
					KnockbackDecay: true,
				})
			particles.CreateSmallExplosion(m.Tile.Transform.Pos)
			sfx.SoundPlayer.PlaySound("blast1", -1.0)
			m.Delete()
		}
	}
}

func (m *Mine) Create(pos pixel.Vec) {
	m.Transform = transform.New()
	m.Transform.Pos = pos
	m.created = true
	m.Timer = timing.New(MineFuse)
	m.Reanimator = reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("mine_1", img.Batchers[constants.EntityKey].Animations["mine_1"].S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("mine_2", img.Batchers[constants.EntityKey].Animations["mine_2"].S, reanimator.Tran).
			SetTrigger(2, func(_ *reanimator.Anim, _ string, _ int) {
				m.explode = true
			}),
		).
		SetChooseFn(func() int {
			if MineFuse*0.5 > m.Timer.Elapsed() {
				return 0
			} else {
				return 1
			}
		}), "mine_1")
	m.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, m).
		AddComponent(myecs.Transform, m.Transform).
		AddComponent(myecs.Animation, m.Reanimator).
		AddComponent(myecs.Drawable, m.Reanimator).
		AddComponent(myecs.Batch, constants.EntityKey)
}

func (m *Mine) Delete() {
	myecs.Manager.DisposeEntity(m.entity)
}
