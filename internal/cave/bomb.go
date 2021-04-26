package cave

import (
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/animation"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"time"
)

type Bomb struct {
	Transform *animation.Transform
	Timer     time.Time
	Tile      *Tile
	created   bool
	done      bool
	sprite    *pixel.Sprite
	color     color.RGBA
}

func (b *Bomb) Update() {
	if b.created && !b.done {
		b.Transform.Update(pixel.Rect{})
		if time.Since(b.Timer).Seconds() > 0.25 {
			for _, n := range b.Tile.Coords.Neighbors() {
				b.Tile.Chunk.Get(n).Destroy()
			}
			p := Player1.Transform.Pos.Sub(b.Tile.Transform.Pos)
			mag := math.Sqrt(p.X*p.X + p.Y*p.Y)
			if mag < world.TileSize*2. {
				Player1.Damage(4.-mag/world.TileSize, b.Tile.Transform.Pos)
			}
			vfx.CreateExplosion(b.Tile.Transform.Pos)
			b.done = true
		}
	}
}

func (b *Bomb) Draw(target pixel.Target) {
	if b.created && !b.done {
		b.sprite.DrawColorMask(target, b.Transform.Mat, b.color)
	}
}

func (b *Bomb) Create(pos pixel.Vec, batcher *img.Batcher) {
	b.Transform = animation.NewTransform(true)
	b.Transform.Pos = pos
	b.created = true
	b.Timer = time.Now()
	b.sprite = batcher.Sprites["bomb"]
	b.color = colornames.White
}

func (b *Bomb) Remove() bool {
	return b.done
}