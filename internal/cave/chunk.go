package cave

import (
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"math/rand"
)

const (
	ChunkSize = 32
	ChunkCnt  = ChunkSize * ChunkSize
)

type Chunk struct {
	Coords  world.Coords
	Rows    [ChunkSize][ChunkSize]*Tile
	display bool
	Cave    *Cave
}

func GenerateChunk(coords world.Coords, cave *Cave) *Chunk {
	// Array of 1024 bools
	list := [ChunkCnt]bool{}
	// fill first 10-20% with true
	bCount := rand.Intn(ChunkCnt / 10) + ChunkCnt / 10
	for i := 0; i < bCount; i++ {
		list[i] = true
	}
	// randomize list
	for i := len(list) - 1; i > 0; i-- {
		j := rand.Intn(i)
		list[i], list[j] = list[j], list[i]
	}
	// create chunk, distribute bombs (trues), build tiles
	chunk := &Chunk{
		Coords:  coords,
		Rows:    [32][32]*Tile{},
		display: true,
		Cave:    cave,
	}
	y := 0
	x := 0
	for _, b := range list {
		chunk.Rows[y][x] = NewTile(x, y, coords, b, chunk)
		x++
		if x % ChunkSize == 0 {
			x = 0
			y++
		}
	}
	return chunk
}

func (chunk *Chunk) Update(input *input.Input) {
	for _, row := range chunk.Rows {
		for _, tile := range row {
			tile.Update(input)
			tile.Transform.Update(pixel.Rect{})
		}
	}
}

func (chunk *Chunk) Draw(target pixel.Target) {
	if chunk.display {
		for _, row := range chunk.Rows {
			for _, tile := range row {
				if !tile.destroyed {
					tile.Sprite.Draw(target, tile.Transform.Mat)
				}
			}
		}
	}
}

func (chunk *Chunk) UpdatePostGen(u, l, r, d *Chunk) {
	// For each surrounding Chunk, update each edge row (corners are special case)
	// the update each tile out?
}

func (chunk *Chunk) Get(coords world.Coords) *Tile {
	if coords.X < 0 || coords.Y < 0 || coords.X >= ChunkSize || coords.Y >= ChunkSize {
		return nil
	}
	return chunk.Rows[coords.Y][coords.X]
}