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

type Mine struct {
	Transform *transform.Transform
	Timer     time.Time
	Tile      *Tile
	created   bool
	done      bool
	animation *img.Instance
}

func (m *Mine) Update() {
	if m.created && !m.done {
		m.Transform.Update(pixel.Rect{})
		m.animation.Update()
		m.animation.SetMatrix(m.Transform.Mat)
		if time.Since(m.Timer).Seconds() > 0.25 {
			for _, n := range m.Tile.Coords.Neighbors() {
				m.Tile.Chunk.Get(n).Destroy()
			}
			p := Player1.Transform.Pos.Sub(m.Tile.Transform.Pos)
			mag := math.Sqrt(p.X*p.X + p.Y*p.Y)
			if mag < world.TileSize*2. {
				Player1.Damage(4.-mag/world.TileSize, m.Tile.Transform.Pos)
			}
			vfx.CreateExplosion(m.Tile.Transform.Pos)
			m.done = true
		}
	}
}

func (m *Mine) Draw(target pixel.Target) {
	if m.created && !m.done {
		m.animation.Draw(target)
	}
}

func (m *Mine) Create(pos pixel.Vec, batcher *img.Batcher) {
	m.Transform = transform.NewTransform(true)
	m.Transform.Pos = pos
	m.created = true
	m.Timer = time.Now()
	m.animation = batcher.Animations["mine"].NewInstance()
}

func (m *Mine) Remove() bool {
	return m.done
}