package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/world"
)

var begun bool

type Challenge struct{
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

func outline(chal Challenge) []Path {
	return []Path{
		{Up, 3},
		{Right, 4},
		{Down, 1},
		{Right, 3},
		{Up, chal.Height},
		{Right, chal.Width+3},
		{Down, chal.Height},
		{Right, 3},
		{Up, 1},
		{Right, 4},
		{Down, 3},
		{Left, chal.Width+15},
	}
}

func NewMinesweeperCave(spriteSheet *img.SpriteSheet, level int) *cave.Cave {
	random.RandCaveSeed()
	begun = false
	chal := challenge(level)
	w := 1
	h := 1
	for w <= (chal.Width + 28) / constants.ChunkSize {
		w++
	}
	for h <= (chal.Height + 12) / constants.ChunkSize {
		h++
	}
	batcher := img.NewBatcher(spriteSheet, false)
	newCave := cave.NewCave(batcher, true)
	newCave.Left = 0
	newCave.Right = w
	newCave.Bottom = h-1
	newCave.StartC = world.Coords{X: 12, Y: h * constants.ChunkSize - 8}
	exitC := newCave.StartC
	exitC.X += chal.Width + 13
	newCave.ExitC = exitC
	pathS := newCave.StartC
	pathS.X -= 2
	pathS.Y += 1
	descent.CaveTotalBombs = chal.Mines
	descent.CaveBombsLeft = chal.Mines
	CreateChunks(newCave)
	Outline(newCave, pathS, outline(chal))
	Entrance(newCave, newCave.StartC, 5, 3, 0, false)
	Entrance(newCave, exitC, 5, 3, 0, true)
	for x := newCave.StartC.X+1; x < newCave.ExitC.X; x++ {
		toBlock(newCave.GetTileInt(x, newCave.StartC.Y), false, false)
	}
	newCave.MarkAsNotChanged()
	MineBlock(newCave, chal)
	newCave.PrintCaveToTerminal()
	return newCave
}

func MineBlock(c *cave.Cave, chal Challenge) {
	curr := c.StartC
	curr.X += 7
	x := curr.X
	bCount := chal.Mines
	list := make([]bool, chal.Width*chal.Height)
	for i := 0; i < bCount; i++ {
		list[i] = true
	}
	// randomize list
	for i := len(list) - 1; i > 0; i-- {
		j := random.CaveGen.Intn(i)
		list[i], list[j] = list[j], list[i]
	}
	b := 0
	for i := 0; i < chal.Height; i++ {
		for j := 0; j < chal.Width; j++ {
			tile := c.GetTileInt(curr.X, curr.Y)
			toBlock(tile, true, list[b])
			tile.Fillable = true
			tile.DigTrigger = func(tile *cave.Tile) {
				if !begun {
					descent.StartMinesweeper(c, tile)
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
	for i := 0; i < chal.Height + 1; i++ {
		for j := 0; j < chal.Width + 2; j++ {
			toEmpty(c.GetTileInt(curr.X, curr.Y), true, true)
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