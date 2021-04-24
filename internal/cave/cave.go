package cave

import (
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var CurrCave *Cave

type Cave struct {
	Chunks  map[world.Coords]*Chunk
	batcher *img.Batcher
}

func NewCave(spriteSheet *img.SpriteSheet) *Cave {
	batcher := img.NewBatcher(spriteSheet)
	cave := &Cave{
		Chunks:  nil,
		batcher: batcher,
	}
	chunk := GenerateChunk(world.Coords{X: 0, Y: 0}, cave)
	cave.Chunks = make(map[world.Coords]*Chunk)
	cave.Chunks[chunk.Coords] = chunk
	return cave
}

func (cave *Cave) Update(pos pixel.Vec, input *input.Input) {
	cave.batcher.Clear()
	p := WorldToChunk(pos)
	all := append([]world.Coords{p}, p.Neighbors()...)
	for _, i := range all {
		if i.X >= 0 && i.Y >= 0 {
			if _, ok := cave.Chunks[i]; !ok {
				cave.Chunks[i] = GenerateChunk(i, cave)
			}
		}
	}
	for i, chunk := range cave.Chunks {
		chunk.display = world.CoordsIn(i, all)
		chunk.Update(input)
	}
}

func (cave *Cave) Draw(win *pixelgl.Window) {
	for _, chunk := range cave.Chunks {
		chunk.Draw(cave.batcher.Batch())
	}
	cave.batcher.Draw(win)
}

func WorldToChunk(v pixel.Vec) world.Coords {
	return world.Coords{X: int(v.X / ChunkSize / world.TileSize), Y: int(-v.Y / ChunkSize / world.TileSize)}
}