package structures

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/objects"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/world"
)

func SmallCamp(c *cave.Cave, include world.Coords, dir data.Direction) {
	h := random.CaveGen.Intn(2)+1
	w := random.CaveGen.Intn(6) + 8
	currX := include.X
	currY := include.Y
	done := false
	count := 0
	for !done {
		if currX > (c.Right+1)*constants.ChunkSize-6 || currX < 6 || count > w {
			done = true
			h /= 2
		}
		t := 0
		b := h
		stc := random.CaveGen.Intn(4) == 0
		stcBG := random.CaveGen.Intn(2) == 0
		stm := random.CaveGen.Intn(4) == 0
		stmBG := random.CaveGen.Intn(2) == 0
		stmI := random.CaveGen.Intn(4) + 1
		for y := currY + 2; y > currY - h; y-- {
			tile := c.GetTileInt(currX, y)
			if tile != nil {
				if y > currY {
					ToBlock(tile, true, true)
				} else {
					if stc && t < random.CaveGen.Intn(4)+1 {
						ToType(tile, cave.Growth, false, true)
						tile.BG = stcBG
					} else {
						stc = false
						if stm && b < stmI {
							ToType(tile, cave.Growth, false, true)
							tile.BG = stmBG
						} else {
							ToType(tile, cave.Empty, false, true)
							if y == currY && random.CaveGen.Intn(10) == 0 {
								objects.AddTent(tile, random.CaveGen.Intn(2) == 0)
							}
						}
					}
				}
			}
			t++
			b--
		}
		if random.CaveGen.Intn(2) == 0 {
			h += random.CaveGen.Intn(3) - 1
		}
		if h < 3 {
			h = 3
		}
		if random.CaveGen.Intn(10) == 0 {
			currY += random.CaveGen.Intn(3) - 1
		}
		if dir == data.Right {
			currX++
		} else {
			currX--
		}
		count++
	}
}