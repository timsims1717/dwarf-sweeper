package cave

import (
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/animation"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"time"
)

type Tile struct {
	Coords     world.Coords
	Sprite     *pixel.Sprite
	bomb       bool
	destroyed  bool
	solid      bool
	Transform  *animation.Transform
	Chunk      *Chunk
	revealT    time.Time
	revealing  bool
	destroyT   time.Time
	destroying bool
}

func NewTile(x, y int, ch world.Coords, bomb bool, chunk *Chunk) *Tile {
	transform := animation.NewTransform(true)
	transform.Pos = pixel.V(float64(x + ch.X * ChunkSize) * world.TileSize, -(float64(y + ch.Y * ChunkSize) * world.TileSize))
	spr := chunk.Cave.batcher.Sprites["block"]
	if bomb {
		spr = chunk.Cave.batcher.Sprites["bomb"]
	}
	return &Tile{
		Coords:    world.Coords{ X: x, Y: y },
		Sprite:    spr,
		bomb:      bomb,
		solid:     true,
		Transform: transform,
		Chunk:     chunk,
	}
}

func (tile *Tile) Update(input *input.Input) {
 	if tile.solid && !tile.destroyed && util.PointInside(input.World, world.TileRect, tile.Transform.Mat) && input.Select.JustPressed() {
		input.Select.Consume()
		tile.Destroy()
	}
	if !tile.destroyed && tile.destroying {
		s := time.Since(tile.destroyT).Seconds()
		if s >= 0.2 {
			tile.Destroy()
		}
	}
	if tile.solid && !tile.destroyed && tile.revealing {
		s := time.Since(tile.revealT).Seconds()
		if s >= 0.2 {
			tile.Reveal()
		}
	}
}

func (tile *Tile) ToDestroy() {
	if tile != nil && !tile.destroyed && !tile.destroying {
		tile.destroyT = time.Now()
		tile.destroying = true
	}
}

func (tile *Tile) Destroy() {
	if tile != nil && !tile.destroyed {
		tile.destroying = false
		tile.solid = false
		if tile.bomb {
			tile.bomb = false
			tile.destroyed = true
			tile.Sprite = nil
			for _, n := range tile.Coords.Neighbors() {
				tile.Chunk.Get(n).ToDestroy()
			}
			vfx.CreateExplosion(tile.Transform.Pos)
		} else {
			ns := tile.Coords.Neighbors()
			c := 0
			for _, n := range ns {
				if tile.Chunk.Get(n).IsBomb() {
					c++
				}
			}
			spr := new(pixel.Sprite)
			switch c {
			case 0:
				tile.destroyed = true
				for _, n := range ns {
					tile.Chunk.Get(n).ToReveal()
				}
			case 1:
				spr = tile.Chunk.Cave.batcher.Sprites["one"]
			case 2:
				spr = tile.Chunk.Cave.batcher.Sprites["two"]
			case 3:
				spr = tile.Chunk.Cave.batcher.Sprites["three"]
			case 4:
				spr = tile.Chunk.Cave.batcher.Sprites["four"]
			case 5:
				spr = tile.Chunk.Cave.batcher.Sprites["five"]
			case 6:
				spr = tile.Chunk.Cave.batcher.Sprites["six"]
			case 7:
				spr = tile.Chunk.Cave.batcher.Sprites["seven"]
			case 8:
				spr = tile.Chunk.Cave.batcher.Sprites["eight"]
			}
			tile.Sprite = spr
			particles.BlockParticles(tile.Transform.Pos)
		}
	}
}

func (tile *Tile) ToReveal() {
	if tile != nil && !tile.revealing && tile.solid {
		tile.revealT = time.Now()
		tile.revealing = true
	}
}

func (tile *Tile) Reveal() {
	if tile != nil && !tile.bomb && tile.solid {
		tile.revealing = false
		tile.solid = false
		ns := tile.Coords.Neighbors()
		c := 0
		for _, n := range ns {
			if tile.Chunk.Get(n).IsBomb() {
				c++
			}
		}
		spr := new(pixel.Sprite)
		switch c {
		case 0:
			tile.destroyed = true
			for _, n := range ns {
				tile.Chunk.Get(n).ToReveal()
			}
		case 1:
			spr = tile.Chunk.Cave.batcher.Sprites["one"]
		case 2:
			spr = tile.Chunk.Cave.batcher.Sprites["two"]
		case 3:
			spr = tile.Chunk.Cave.batcher.Sprites["three"]
		case 4:
			spr = tile.Chunk.Cave.batcher.Sprites["four"]
		case 5:
			spr = tile.Chunk.Cave.batcher.Sprites["five"]
		case 6:
			spr = tile.Chunk.Cave.batcher.Sprites["six"]
		case 7:
			spr = tile.Chunk.Cave.batcher.Sprites["seven"]
		case 8:
			spr = tile.Chunk.Cave.batcher.Sprites["eight"]
		}
		tile.Sprite = spr
		particles.BlockParticles(tile.Transform.Pos)
	}
}

func (tile *Tile) IsBomb() bool {
	if tile == nil {
		return false
	}
	return tile.bomb
}