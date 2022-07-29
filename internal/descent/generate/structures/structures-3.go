package structures

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/critters"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/world"
)

func Cavern(c *cave.Cave, include world.Coords, dir data.Direction, enemies []string) {
	h := random.CaveGen.Intn(2)+1
	offset := random.CaveGen.Intn(3)-1
	//b := random.CaveGen.Intn(6)-3
	w := random.CaveGen.Intn(8) + random.CaveGen.Intn(8) + 8
	currX := include.X
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
		for y := include.Y + offset; y < include.Y + offset + h; y++ {
			tile := c.GetTileInt(currX, y)
			if tile != nil {
				if stc && t < random.CaveGen.Intn(4)+1 {
					ToType(tile, cave.Growth, false, false)
					tile.BG = stcBG
				} else {
					stc = false
					if stm && b < stmI {
						ToType(tile, cave.Growth, false, false)
						tile.BG = stmBG
					} else {
						ToType(tile, cave.Empty, false, false)
						if random.CaveGen.Intn(20) == 0 {
							critters.AddRandomCritter(c, enemies, tile.Transform.Pos)
						}
					}
				}
			}
			t++
			b--
		}
		h += random.CaveGen.Intn(6) - 3
		offset += random.CaveGen.Intn(6) - 3
		if h < 2 {
			h = 2 + random.CaveGen.Intn(2)
		}
		if offset > 4 {
			offset = 4 - random.CaveGen.Intn(2)
		} else if offset < -4 {
			offset = -4 + random.CaveGen.Intn(2)
		}
		if dir == data.Right {
			currX++
		} else {
			currX--
		}
		count++
	}
}

func BridgeCavern(c *cave.Cave, include world.Coords, dir data.Direction, enemies []string) {
	h := random.CaveGen.Intn(2)+1
	offset := random.CaveGen.Intn(3)-1
	bY := random.CaveGen.Intn(6)-3
	w := random.CaveGen.Intn(12) + 8
	currX := include.X
	done := false
	count := 0
	for !done {
		if currX > (c.Right+1)*constants.ChunkSize-6 || currX < 6 || count > w {
			done = true
			h /= 2
		}
		bAdd := false
		t := 0
		b := h
		stc := random.CaveGen.Intn(4) == 0
		stcBG := random.CaveGen.Intn(2) == 0
		stm := random.CaveGen.Intn(4) == 0
		stmBG := random.CaveGen.Intn(2) == 0
		stmI := random.CaveGen.Intn(4) + 1
		for y := include.Y + offset; y < include.Y + offset + h; y++ {
			tile := c.GetTileInt(currX, y)
			if tile != nil {
				if tile.RCoords.Y == include.Y+bY {
					ToType(tile, cave.Bridge, false, false)
					bAdd = true
					above := c.GetTileInt(currX, y-1)
					if above != nil && above.Solid() {
						ToType(above, cave.Empty, false, true)
						if random.CaveGen.Intn(20) == 0 {
							critters.AddRandomCritter(c, enemies, tile.Transform.Pos)
						}
					}
				} else if tile.RCoords.Y > include.Y+bY && bAdd && tile.RCoords.X%4 == include.X%4 {
					ToType(tile, cave.Pillar, false, false)
				} else {
					if stc && t < random.CaveGen.Intn(4)+1 {
						ToType(tile, cave.Growth, false, false)
						tile.BG = stcBG
					} else {
						stc = false
						if stm && b < stmI {
							ToType(tile, cave.Growth, false, false)
							tile.BG = stmBG
						} else {
							ToType(tile, cave.Empty, false, false)
							if random.CaveGen.Intn(20) == 0 {
								critters.AddRandomCritter(c, enemies, tile.Transform.Pos)
							}
						}
					}
				}
			}
			t++
			b--
		}
		h += random.CaveGen.Intn(6) - 3
		offset += random.CaveGen.Intn(6) - 3
		if h < 2 {
			h = 2 + random.CaveGen.Intn(2)
		}
		if offset > 4 {
			offset = 4 - random.CaveGen.Intn(2)
		} else if offset < -4 {
			offset = -4 + random.CaveGen.Intn(2)
		}
		if dir == data.Right {
			currX++
		} else {
			currX--
		}
		count++
	}
}