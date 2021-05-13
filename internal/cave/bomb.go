package cave

import (
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"math"
	"time"
)

type Bomb struct {
	Transform *transform.Transform
	Timer     time.Time
	Tile      *Tile
	created   bool
	done      bool
	animation *img.Instance
}

func (b *Bomb) Update() {
	if b.created && !b.done {
		b.Transform.Update()
		b.animation.Update()
		b.animation.SetMatrix(b.Transform.Mat)
		if time.Since(b.Timer).Seconds() > 0.75 {
			for _, n := range b.Tile.Coords.Neighbors() {
				b.Tile.Chunk.Get(n).Destroy()
			}
			p := Player1.Transform.Pos.Sub(b.Tile.Transform.Pos)
			mag := math.Sqrt(p.X*p.X + p.Y*p.Y)
			if mag < world.TileSize*2. {
				dmg := 4.-mag/world.TileSize
				Player1.Damage(dmg, b.Tile.Transform.Pos, MineKnockback*dmg*world.TileSize)
			}
			vfx.CreateExplosion(b.Tile.Transform.Pos)
			b.done = true
		}
	}
}

func (b *Bomb) Draw(target pixel.Target) {
	if b.created && !b.done {
		b.animation.Draw(target)
	}
}

func (b *Bomb) Create(pos pixel.Vec, batcher *img.Batcher) {
	b.Transform = transform.NewTransform()
	b.Transform.Pos = pos
	b.created = true
	b.Timer = time.Now()
	b.animation = batcher.Animations["bomb"].NewInstance()
}

func (b *Bomb) Remove() bool {
	return b.done
}