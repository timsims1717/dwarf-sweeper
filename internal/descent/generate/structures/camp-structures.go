package structures

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/objects"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/world"
)

var campMat = map[string]map[string]int{
	"": {
		"tent": 5,
		"gempile": 8,
		"seat": 5,
		"barrel": 2,
		"bedroll": 6,
		"refuse": 8,
		"tools": 4,
		"woodpile": 8,
	},
	"tent": {
		"seat": 15,
		"bedroll": 10,
		"cookfire": 5,
	},
	"seat": {
		"seat": 2,
		"art": 10,
		"cookfire": 15,
		"tools": 2,
		"bedroll": 5,
	},
	"bedroll": {
		"seat": 5,
		"cookfire": 15,
	},
}

func SmallCamp(c *cave.Cave, include world.Coords, dir data.Direction) {
	h := random.CaveGen.Intn(2)+1
	w := random.CaveGen.Intn(6) + 8
	currX := include.X
	currY := include.Y
	done := false
	count := 0
	prev := ""
	next := ""
	tent := false
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
							if y == currY {
								if !tent {
									chMat, ok := campMat[prev]
									if !ok {
										chMat = campMat[""]
									}
									choice := random.CaveGen.Intn(50)
									total := 0
									next = ""
									for nxt, cnt := range chMat {
										total += cnt
										if choice < total {
											next = nxt
										}
									}
									if next != "" {
										// l means left side of the camp (so face right)
										l := (dir == data.Right && currX < include.X+w/2) || (dir == data.Left && currX < include.X-w/2)
										switch next {
										case "tent":
											objects.AddTent(tile, l)
											tent = true
										case "gempile":
											objects.AddGemPile(tile)
										case "refuse":
											if random.CaveGen.Intn(2) == 0 {
												objects.AddObject(tile, "refuse_sm", false, objects.Random)
											} else {
												objects.AddObject(tile, "refuse_lg", false, objects.Random)
											}
										case "cookfire", "tools", "woodpile":
											objects.AddObject(tile, next, false, objects.Random)
										default:
											if l {
												objects.AddObject(tile, next, false, objects.Flip)
											} else {
												objects.AddObject(tile, next, false, objects.Normal)
											}
										}
									}
									prev = next
								}
								tent = false
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