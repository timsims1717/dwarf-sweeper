package dungeon

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"math"
	"time"
)

const (
	MineKnockback = 10.
)

type Mine struct {
	Transform  *transform.Transform
	Timer      time.Time
	FuseLength float64
	Tile       *Tile
	created    bool
	done       bool
	explode    bool
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
}

func (m *Mine) Update() {
	if m.created && !m.done && m.explode {
		if time.Since(m.Timer).Seconds() > 0.25 {
			for _, n := range m.Tile.SubCoords.Neighbors() {
				m.Tile.Chunk.Get(n).Destroy()
			}
			p := Dungeon.Player.Transform.Pos.Sub(m.Tile.Transform.Pos)
			mag := math.Sqrt(p.X*p.X + p.Y*p.Y)
			if mag < world.TileSize*2. {
				dmg := 4.-mag/world.TileSize
				Dungeon.Player.Damage(dmg, m.Tile.Transform.Pos, MineKnockback*dmg*world.TileSize)
			}
			vfx.CreateExplosion(m.Tile.Transform.Pos)
			m.done = true
		}
	}
}

func (m *Mine) Draw(target pixel.Target) {
	if m.created && !m.done {
		m.Reanimator.CurrentSprite().Draw(target, m.Transform.Mat)
	}
}

func (m *Mine) Create(pos pixel.Vec, batcher *img.Batcher) {
	m.Transform = transform.NewTransform()
	m.Transform.Pos = pos
	m.created = true
	m.Timer = time.Now()
	m.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("mine_1", batcher.Animations["mine_1"].S, reanimator.Hold, nil),
			reanimator.NewAnimFromSprites("mine_2", batcher.Animations["mine_2"].S, reanimator.Tran, map[int]func() {
				1: func() {
					m.explode = true
				},
			}),
		),
		Check: func() int {
			if m.FuseLength - time.Since(m.Timer).Seconds() > 0.3 {
				return 0
			} else {
				return 1
			}
		},
	})
	m.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, m.Transform).
		AddComponent(myecs.Animation, m.Reanimator)
}

func (m *Mine) Done() bool {
	return m.done
}

func (m *Mine) Delete() {
	myecs.Manager.DisposeEntity(m.entity)
}