package cave

import (
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
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

func GenerateStart(cave *Cave) *Chunk {
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
		Coords:  world.Origin,
		Rows:    [32][32]*Tile{},
		display: true,
		reload:  true,
		Cave:    cave,
	}
	y := 0
	x := 0
	for _, b := range list {
		tile := NewTile(x, y, world.Origin, b, chunk)
		// starting room
		if x > 6 && x < 26 && y > 2 && y < 10 {
			tile.Solid = false
			tile.destroyed = true
			tile.bomb = false
			tile.Sprite = nil
		} else if x > 5 && x < 27 && y > 1 && y < 11 {
			tile.bomb = false
		}
		chunk.Rows[y][x] = tile
		x++
		if x % ChunkSize == 0 {
			x = 0
			y++
		}
	}
	return chunk
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
		chunk.Rows[y][x] = NewTile(x, y, coords, b, chunk)
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
				tile.Transform.Update(pixel.Rect{})
			}
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
		ul := chunk.Rows[0][0].Transform.Pos
		dr := chunk.Rows[ChunkSize-1][ChunkSize-1].Transform.Pos
		half := world.TileSize*0.5
		debug.AddLine(colornames.Green, imdraw.SharpEndShape, pixel.V(ul.X-half, ul.Y+half), pixel.V(dr.X+half, ul.Y+half), 1.0)
		debug.AddLine(colornames.Green, imdraw.SharpEndShape, pixel.V(dr.X+half, ul.Y+half), pixel.V(dr.X+half, dr.Y-half), 1.0)
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