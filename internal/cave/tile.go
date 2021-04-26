package cave

import (
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/pkg/animation"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"math/rand"
	"time"
)

type Tile struct {
	Coords     world.Coords
	Sprite     *pixel.Sprite
	bomb       bool
	destroyed  bool
	Solid      bool
	Transform  *animation.Transform
	Chunk      *Chunk
	revealT    time.Time
	revealing  bool
	destroyT   time.Time
	destroying bool
	reload     bool
	marked     bool
}

func NewTile(x, y int, ch world.Coords, bomb bool, chunk *Chunk) *Tile {
	transform := animation.NewTransform(true)
	transform.Pos = pixel.V(float64(x + ch.X * ChunkSize) * world.TileSize, -(float64(y + ch.Y * ChunkSize) * world.TileSize))
	spr := chunk.Cave.batcher.Sprites["block"]
	return &Tile{
		Coords:    world.Coords{ X: x, Y: y },
		Sprite:    spr,
		bomb:      bomb,
		Solid:     true,
		Transform: transform,
		Chunk:     chunk,
	}
}

func (tile *Tile) Update() {
	if tile.reload {
		if tile.Coords.X == 0 || tile.Coords.X == ChunkSize - 1 || tile.Coords.Y == 0 || tile.Coords.Y == ChunkSize - 1 {
			for _, n := range tile.Coords.Neighbors() {
				t := tile.Chunk.Get(n)
				if t != nil && t.destroyed {
					tile.Reveal(true)
				}
			}
		}
		tile.reload = false
	}
	if !tile.destroyed && tile.destroying {
		s := time.Since(tile.destroyT).Seconds()
		if s >= 0.2 {
			tile.Destroy()
		}
	}
	if tile.Solid && !tile.destroyed && tile.revealing {
		s := time.Since(tile.revealT).Seconds()
		if s >= 0.2 {
			tile.Reveal(false)
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
		tile.Solid = false
		if tile.bomb {
			tile.bomb = false
			tile.destroyed = true
			tile.Sprite = nil
			Entities.Add(&Bomb{
				Tile:      tile,
			}, tile.Transform.Pos)
		} else {
			ns := tile.Coords.Neighbors()
			c := 0
			for _, n := range ns {
				t := tile.Chunk.Get(n)
				if t != nil && t.bomb {
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
			sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", rand.Intn(5) + 1), -1.0)
		}
	}
}

func (tile *Tile) ToReveal() {
	if tile != nil && !tile.revealing && tile.Solid {
		tile.revealT = time.Now()
		tile.revealing = true
	}
}

func (tile *Tile) Reveal(instant bool) {
	if tile != nil && !tile.bomb && tile.Solid {
		tile.revealing = false
		tile.Solid = false
		ns := tile.Coords.Neighbors()
		c := 0
		for _, n := range ns {
			t := tile.Chunk.Get(n)
			if t != nil && t.bomb {
				c++
			}
		}
		spr := new(pixel.Sprite)
		switch c {
		case 0:
			tile.destroyed = true
			for _, n := range ns {
				if instant {
					tile.Chunk.Get(n).Reveal(true)
				} else {
					tile.Chunk.Get(n).ToReveal()
				}
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
		if !instant {
			particles.BlockParticles(tile.Transform.Pos)
		}
	}
}

func (tile *Tile) Unmark() {
	tile.marked = false
}

func (tile *Tile) Mark(from pixel.Vec) {
	if tile != nil && tile.Solid && !tile.destroyed {
		correct := tile.bomb
		if !tile.marked {
			tile.marked = true
			Entities.Add(&Flag{
				Tile: tile,
			}, from)
			if correct {
				BombsMarked++
			} else {
				BlocksMarked++
			}
		} else {
			tile.marked = false
			if correct {
				BombsMarked--
			} else {
				BlocksMarked--
			}
		}
	}
}