package dungeon

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
	pivot   pixel.Vec
	finite  bool
	batcher *img.Batcher
	left    int
	right   int
	bottom  int
	StartC  world.Coords
}

func (cave *Cave) Update(pos pixel.Vec) {
	cave.batcher.Clear()
	cave.pivot = pos
	if !cave.finite {
		p := WorldToChunk(cave.pivot)
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
	} else {
		for _, chunk := range cave.RChunks {
			chunk.display = true
			chunk.Update()
		}
		for _, chunk := range cave.LChunks {
			chunk.display = true
			chunk.Update()
		}
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

func (cave *Cave) Dimensions() (int, int) {
	if cave.finite {
		return (cave.right - cave.left + 1) * ChunkSize, (cave.bottom + 1) * ChunkSize
	} else {
		return -1, -1
	}
}

func (cave *Cave) CurrentBoundaries() (pixel.Vec, pixel.Vec) {
	p := WorldToChunk(cave.pivot)
	var all []world.Coords
	if cave.finite {
		for _, chunk := range cave.RChunks {
			all = append(all, chunk.Coords)
		}
		for _, chunk := range cave.LChunks {
			all = append(all, chunk.Coords)
		}
	} else {
		all = append([]world.Coords{p}, p.Neighbors()...)
	}
	x1 := 10000000.
	y1 := 10000000.
	x2 := -10000000.
	y2 := -10000000.
	for _, i := range all {
		if i.X >= 0 && i.Y >= 0 {
			if chunk, ok := cave.RChunks[i]; ok {
				tr := chunk.Rows[0][ChunkSize-1].Transform.Pos
				bl := chunk.Rows[ChunkSize-1][0].Transform.Pos
				if bl.X < x1 {
					x1 = bl.X
				}
				if bl.Y < y1 {
					y1 = bl.Y
				}
				if tr.X > x2 {
					x2 = tr.X
				}
				if tr.Y > y2 {
					y2 = tr.Y
				}
			}
		} else if i.X < 0 && i.Y >= 0 {
			if chunk, ok := cave.LChunks[i]; ok {
				tr := chunk.Rows[0][ChunkSize-1].Transform.Pos
				bl := chunk.Rows[ChunkSize-1][0].Transform.Pos
				if bl.X < x1 {
					x1 = bl.X
				}
				if bl.Y < y1 {
					y1 = bl.Y
				}
				if tr.X > x2 {
					x2 = tr.X
				}
				if tr.Y > y2 {
					y2 = tr.Y
				}
			}
		}
	}
	return pixel.V(x1, y1), pixel.V(x2, y2)
}

func (cave *Cave) GetTileInt(x, y int) *Tile {
	cX := x / ChunkSize
	if x < 0 {
		cX--
	}
	tX := x % ChunkSize
	if tX < 0 {
		tX += ChunkSize
	}
	cY := y / ChunkSize
	tY := y % ChunkSize
	return cave.GetChunk(world.Coords{X: cX, Y: cY}).Get(world.Coords{X: tX, Y: tY})
}

func (cave *Cave) GetChunk(coords world.Coords) *Chunk {
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
	chunk := cave.GetChunk(ch)
	return chunk.Get(tl)
}

func (cave *Cave) GetStart() pixel.Vec {
	return cave.GetTileInt(cave.StartC.X, cave.StartC.Y).Transform.Pos
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