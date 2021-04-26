package cave

import (
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var CurrCave *Cave

type Cave struct {
	RChunks map[world.Coords]*Chunk
	LChunks map[world.Coords]*Chunk
	batcher *img.Batcher
}

func NewCave(spriteSheet *img.SpriteSheet) *Cave {
	batcher := img.NewBatcher(spriteSheet)
	cave := &Cave{
		RChunks: nil,
		LChunks: nil,
		batcher: batcher,
	}
	chunk0 := GenerateStart(cave)

	chunkr1 := GenerateChunk(world.Coords{X: 1, Y: 0}, cave)
	chunkr2 := GenerateChunk(world.Coords{X: 1, Y: 1}, cave)
	chunkr3 := GenerateChunk(world.Coords{X: 0, Y: 1}, cave)

	chunkl1 := GenerateChunk(world.Coords{X: -1, Y: 0}, cave)
	chunkl2 := GenerateChunk(world.Coords{X: -1, Y: 1}, cave)

	cave.RChunks = make(map[world.Coords]*Chunk)
	cave.RChunks[chunk0.Coords] = chunk0
	cave.RChunks[chunkr1.Coords] = chunkr1
	cave.RChunks[chunkr2.Coords] = chunkr2
	cave.RChunks[chunkr3.Coords] = chunkr3

	cave.LChunks = make(map[world.Coords]*Chunk)
	cave.LChunks[chunkl1.Coords] = chunkl1
	cave.LChunks[chunkl2.Coords] = chunkl2
	return cave
}

func (cave *Cave) Update(pos pixel.Vec) {
	cave.batcher.Clear()
	p := WorldToChunk(pos)
	all := append([]world.Coords{p}, p.Neighbors()...)
	for _, i := range all {
		if i.X >= 0 && i.Y >= 0 {
			if _, ok := cave.RChunks[i]; !ok {
				cave.RChunks[i] = GenerateChunk(i, cave)
			}
		} else if i.X < 0 && i.Y >= 0 {
			if _, ok := cave.LChunks[i]; !ok {
				cave.LChunks[i] = GenerateChunk(i, cave)
			}
		}
	}
	for i, chunk := range cave.RChunks {
		dis := world.CoordsIn(i, all)
		if dis && !chunk.display {
			chunk.reload = true
		}
		chunk.display = dis
		chunk.Update()
	}
	for i, chunk := range cave.LChunks {
		dis := world.CoordsIn(i, all)
		if dis && !chunk.display {
			chunk.reload = true
		}
		chunk.display = dis
		chunk.Update()
	}
}

func (cave *Cave) Draw(win *pixelgl.Window) {
	for _, chunk := range cave.RChunks {
		chunk.Draw(cave.batcher.Batch())
	}
	for _, chunk := range cave.LChunks {
		chunk.Draw(cave.batcher.Batch())
	}
	cave.batcher.Draw(win)
}

func (cave *Cave) Get(coords world.Coords) *Chunk {
	if chunkR, okR := cave.RChunks[coords]; okR {
		return chunkR
	} else if chunkL, okL := cave.LChunks[coords]; okL {
		return chunkL
	} else {
		return nil
	}
}

func (cave *Cave) GetTile(v pixel.Vec) *Tile {
	ch := WorldToChunk(v)
	tl := WorldToTile(v, ch.X < 0)
	chunk := cave.Get(ch)
	return chunk.Get(tl)
}

func WorldToChunk(v pixel.Vec) world.Coords {
	if v.X >= 0 - world.TileSize * 0.5 {
		return world.Coords{X: int((v.X+world.TileSize*0.5) / ChunkSize / world.TileSize), Y: int(-(v.Y-world.TileSize*0.5) / ChunkSize / world.TileSize)}
	} else {
		return world.Coords{X: int((v.X+world.TileSize*0.5) / ChunkSize / world.TileSize)-1, Y: int(-(v.Y-world.TileSize*0.5) / ChunkSize / world.TileSize)}
	}
}

func WorldToTile(v pixel.Vec, left bool) world.Coords {
	x, y := world.WorldToMap(v.X+world.TileSize*0.5, -(v.Y-world.TileSize*0.5))
	x = x % ChunkSize
	y = y % ChunkSize
	if left {
		x = (ChunkSize - (util.Abs(x) + 1)) % ChunkSize
	}
	return world.Coords{
		X: x % ChunkSize,
		Y: y % ChunkSize,
	}
}