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

type Bomb struct {
	Transform  *transform.Transform
	Timer      time.Time
	FuseLength float64
	Tile       *Tile
	created    bool
	done       bool
	first      bool
	explode    bool
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
}

func (b *Bomb) Update() {
	if b.created && !b.done && b.explode {
		for _, n := range b.Tile.SubCoords.Neighbors() {
			b.Tile.Chunk.Get(n).Destroy()
		}
		p := Dungeon.Player.Transform.Pos.Sub(b.Tile.Transform.Pos)
		mag := math.Sqrt(p.X*p.X + p.Y*p.Y)
		if mag < world.TileSize*2. {
			dmg := 4.-mag/world.TileSize
			Dungeon.Player.Damage(dmg, b.Tile.Transform.Pos, MineKnockback*dmg*world.TileSize)
		}
		vfx.CreateExplosion(b.Tile.Transform.Pos)
		b.done = true
	}
}

func (b *Bomb) Draw(target pixel.Target) {
	if b.created && !b.done {
		b.Reanimator.CurrentSprite().Draw(target, b.Transform.Mat)
	}
}

func (b *Bomb) Create(pos pixel.Vec, batcher *img.Batcher) {
	b.Transform = transform.NewTransform()
	b.Transform.Pos = pos
	b.created = true
	b.Timer = time.Now()
	b.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("bomb_unlit", batcher.Animations["bomb_unlit"].S, reanimator.Tran, func() {
				b.first = false
			}),
			reanimator.NewAnimFromSprites("bomb_fuse", batcher.Animations["bomb_fuse"].S, reanimator.Loop, nil),
			reanimator.NewAnimFromSprites("bomb_blow", batcher.Animations["bomb_blow"].S, reanimator.Tran, func() {
				b.explode = true
			}),
		),
		Check: func() int {
			if b.first {
				return 0
			} else if b.FuseLength - time.Since(b.Timer).Seconds() > 0.3 {
				return 1
			} else {
				return 2
			}
		},
	})
	b.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, b.Transform).
		AddComponent(myecs.Animation, b.Reanimator)
}

func (b *Bomb) Remove() bool {
	if b.done {
		myecs.Manager.DisposeEntity(b.entity)
		return true
	}
	return false
}