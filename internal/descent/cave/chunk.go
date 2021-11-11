package cave

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/minesweeper"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

type Chunk struct {
	Coords  world.Coords
	Rows    [constants.ChunkSize][constants.ChunkSize]*Tile
	Display bool
	Reload  bool
	Cave    *Cave
}

func NewChunk(coords world.Coords, c *Cave, tt TileType) *Chunk {
	// Array of 1024 bools
	var list []bool
	// fill first BombPMin-BombPMax% with true
	if c.BombPMax > 0. && c.BombPMax > c.BombPMin {
		bCount := int(constants.ChunkArea * (random.CaveGen.Float64()*(c.BombPMax-c.BombPMin) + c.BombPMin))
		list = minesweeper.CreateBoard(constants.ChunkSize, constants.ChunkSize, bCount, random.CaveGen).AsArray()
	} else {
		list = make([]bool, constants.ChunkArea)
	}
	// create chunk, distribute bombs (trues), build tiles
	chunk := &Chunk{
		Coords:  coords,
		Rows:    [constants.ChunkSize][constants.ChunkSize]*Tile{},
		Display: true,
		Reload:  true,
		Cave:    c,
	}
	y := 0
	x := 0
	for _, b := range list {
		var tile *Tile
		if c.Type != Infinite &&
			((coords.Y == c.Bottom && y == constants.ChunkSize- 1) ||
				(coords.X == c.Left && x == 0) ||
				(coords.X == c.Right && x == constants.ChunkSize- 1)) {
			tile = NewTile(x, y, coords, false, chunk)
			tile.Type = Wall
			tile.NeverChange = true
		} else if coords.Y == 0 && y == 0 {
			tile = NewTile(x, y, coords, false, chunk)
			tile.Type = Wall
			tile.NeverChange = true
		} else {
			tile = NewTile(x, y, coords, b, chunk)
			tile.Type = tt
		}
		chunk.Rows[y][x] = tile
		x++
		if x %constants.ChunkSize == 0 {
			x = 0
			y++
		}
	}
	return chunk
}

func (ch *Chunk) Update() {
	if ch.Reload {
		for _, row := range ch.Rows {
			for _, tile := range row {
				tile.reload = true
			}
		}
		ch.Reload = false
	}
	if ch.Display {
		for _, row := range ch.Rows {
			for _, tile := range row {
				tile.Update()
			}
		}
	}
}

func (ch *Chunk) Draw(target pixel.Target) {
	if ch.Display {
		for _, row := range ch.Rows {
			for _, tile := range row {
				tile.Draw(target)
			}
		}
	}
}

func (ch *Chunk) Get(coords world.Coords) *Tile {
	if ch == nil {
		return nil
	}
	if coords.X < 0 || coords.Y < 0 || coords.X >= constants.ChunkSize || coords.Y >= constants.ChunkSize {
		ax := coords.X
		ay := coords.Y
		cx := 0
		cy := 0
		if coords.X < 0 {
			cx = -1
			ax = constants.ChunkSize - 1
		} else if coords.X >= constants.ChunkSize {
			cx = 1
			ax = 0
		}
		if coords.Y < 0 {
			cy = -1
			ay = constants.ChunkSize - 1
		} else if coords.Y >= constants.ChunkSize {
			cy = 1
			ay = 0
		}
		cc := ch.Coords
		cc.X += cx
		cc.Y += cy
		ac := world.Coords{
			X: ax,
			Y: ay,
		}
		return ch.Cave.GetChunk(cc).Get(ac)
	}
	return ch.Rows[coords.Y][coords.X]
}