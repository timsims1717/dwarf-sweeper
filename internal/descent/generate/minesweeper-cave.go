package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/minesweeper"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/world"
)

var begun bool

type Challenge struct {
	Width, Height, Mines int
}

var (
	Easy = Challenge{
		Width:  10,
		Height: 10,
		Mines:  10,
	}
	Medium = Challenge{
		Width:  12,
		Height: 12,
		Mines:  20,
	}
	Hard = Challenge{
		Width:  16,
		Height: 16,
		Mines:  40,
	}
	Expert = Challenge{
		Width:  30,
		Height: 16,
		Mines:  99,
	}
)

func outline(chal Challenge) []structures.Path {
	return []structures.Path{
		{data.Up, 3},
		{data.Right, 4},
		{data.Down, 1},
		{data.Right, 3},
		{data.Up, chal.Height},
		{data.Right, chal.Width + 3},
		{data.Down, chal.Height},
		{data.Right, 3},
		{data.Up, 1},
		{data.Right, 4},
		{data.Down, 3},
		{data.Left, chal.Width + 15},
	}
}

func MinesweeperCave(c *cave.Cave, level int) *cave.Cave {
	random.RandCaveSeed()
	begun = false
	chal := challenge(level)
	w := 1
	h := 1
	for w <= (chal.Width+28)/constants.ChunkSize {
		w++
	}
	for h <= (chal.Height+12)/constants.ChunkSize {
		h++
	}
	c.SetSize(0, w, h-1)
	c.StartC = world.Coords{X: 12, Y: h*constants.ChunkSize - 8}
	exitC := c.StartC
	exitC.X += chal.Width + 13
	c.ExitC = exitC
	pathS := c.StartC
	pathS.X -= 2
	pathS.Y += 1
	player.CaveTotalBombs = chal.Mines
	player.CaveBombsLeft = chal.Mines
	structures.CreateChunks(c, cave.Wall)
	structures.Outline(c, pathS, outline(chal))
	structures.Entrance(c, c.StartC, 5, 2, 0, false)
	structures.Entrance(c, exitC, 5, 2, 0, true)
	for x := c.StartC.X + 1; x < c.ExitC.X; x++ {
		tile := c.GetTileInt(x, c.StartC.Y)
		structures.ToBlock(tile, cave.BlockCollapse, false, true)
		tile.Bomb = false
	}
	c.MarkAsNotChanged()
	MineBlock(c, chal)
	c.UpdateAllTileSprites()
	return c
}

func MineBlock(c *cave.Cave, chal Challenge) {
	curr := c.StartC
	curr.X += 7
	x := curr.X
	list := minesweeper.CreateBoard(chal.Width, chal.Height, chal.Mines, random.CaveGen).AsArray()
	b := 0
	for i := 0; i < chal.Height; i++ {
		for j := 0; j < chal.Width; j++ {
			tile := c.GetTileInt(curr.X, curr.Y)
			structures.ToBlock(tile, cave.BlockCollapse, true, true)
			tile.Bomb = list[b]
			tile.DestroyTrigger = func(p *player.Player, tile *cave.Tile) {
				if !begun {
					structures.StartMinesweeper(c, tile)
					begun = true
				}
			}
			curr.X++
			b++
		}
		curr.Y--
		curr.X = x
	}
	curr = c.StartC
	curr.X += 6
	x = curr.X
	for i := 0; i < chal.Height+1; i++ {
		for j := 0; j < chal.Width+2; j++ {
			structures.ToEmpty(c.GetTileInt(curr.X, curr.Y), true, true, true)
			curr.X++
		}
		curr.Y--
		curr.X = x
	}
}

func challenge(level int) Challenge {
	if level < 5 {
		return Easy
	} else if level < 9 {
		return Medium
	} else if level < 13 {
		return Hard
	} else {
		return Expert
	}
}
