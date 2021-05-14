package cave

import (
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
	reload  bool
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
		reload:  true,
		Cave:    cave,
	}
	y := 0
	x := 0
	for _, b := range list {
		var tile *Tile
		if coords.Y == 0 && y == 0 {
			tile = NewTile(x, y, coords, false, chunk)
			tile.Type = Wall
			tile.breakable = false
		} else {
			tile = NewTile(x, y, coords, b, chunk)
		}
		if b {
			if rand.Intn(2) == 0 {
				tile.AddEntity(&Bomb{
					Tile: tile,
				})
			} else {
				tile.AddEntity(&Mine{
					Tile: tile,
				})
			}
		}
		//if !b && rand.Intn(20) == 0 {
		//	tile.AddEntity(&Gem{})
		//}
		chunk.Rows[y][x] = tile
		x++
		if x % ChunkSize == 0 {
			x = 0
			y++
		}
	}
	return chunk
}

func (chunk *Chunk) Update() {
	if chunk.reload {
		for _, row := range chunk.Rows {
			for _, tile := range row {
				tile.reload = true
			}
		}
		chunk.reload = false
	}
	if chunk.display {
		for _, row := range chunk.Rows {
			for _, tile := range row {
				tile.Update()
			}
		}
	}
}

func (chunk *Chunk) Draw(target pixel.Target) {
	if chunk.display {
		for _, row := range chunk.Rows {
			for _, tile := range row {
				if !tile.destroyed {
					tile.Draw(target)
				}
			}
		}
		//ul := chunk.Rows[0][0].Transform.Pos
		//dr := chunk.Rows[ChunkSize-1][ChunkSize-1].Transform.Pos
		//half := world.TileSize*0.5
		//debug.AddLine(colornames.Green, imdraw.SharpEndShape, pixel.V(ul.X-half, ul.Y+half), pixel.V(dr.X+half, ul.Y+half), 1.0)
		//debug.AddLine(colornames.Green, imdraw.SharpEndShape, pixel.V(dr.X+half, ul.Y+half), pixel.V(dr.X+half, dr.Y-half), 1.0)
	}
}

func (chunk *Chunk) Get(coords world.Coords) *Tile {
	if chunk == nil {
		return nil
	}
	if coords.X < 0 || coords.Y < 0 || coords.X >= ChunkSize || coords.Y >= ChunkSize {
		ax := coords.X
		ay := coords.Y
		cx := 0
		cy := 0
		if coords.X < 0 {
			cx = -1
			ax = ChunkSize - 1
		} else if coords.X >= ChunkSize {
			cx = 1
			ax = 0
		}
		if coords.Y < 0 {
			cy = -1
			ay = ChunkSize - 1
		} else if coords.Y >= ChunkSize {
			cy = 1
			ay = 0
		}
		cc := chunk.Coords
		cc.X += cx
		cc.Y += cy
		ac := world.Coords{
			X: ax,
			Y: ay,
		}
		return chunk.Cave.Get(cc).Get(ac)
	}
	return chunk.Rows[coords.Y][coords.X]
}