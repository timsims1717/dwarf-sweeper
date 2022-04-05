package structures

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
)

func MineComplex(c *cave.Cave, include world.Coords, dir data.Direction) {
	layer := 1
	layers := 3 + random.CaveGen.Intn(3)
	layerWidth := 12 + random.CaveGen.Intn(int(c.FillVar*0.25))
	currX := include.X
	currY := include.Y
	done := false
	for !done {
		if currY > (c.Right+1)*constants.ChunkSize-15 || random.CaveGen.Intn(layer) > layers {
			done = true
		}
		pillarX := currX + random.CaveGen.Intn(6)
		lastDown := 5
		var downX []int
		count := 0
		bg := random.CaveGen.Intn(2) == 0
		for count < layerWidth {
			if currX > c.Left*constants.ChunkSize+4 &&
				currX < (c.Right+1)*constants.ChunkSize-4 {
				mainTile := c.GetTileInt(currX, currY)
				above1 := c.GetTileInt(currX, currY-1)
				above2 := c.GetTileInt(currX, currY-2)
				above3 := c.GetTileInt(currX, currY-3)
				below1 := c.GetTileInt(currX, currY+1)
				below2 := c.GetTileInt(currX, currY+2)
				if pillarX%6 == currX%6 && (above2.Solid() || layer == 1) {
					mainTile.BG = bg
					above1.BG = bg
					bg = !bg
					ToType(mainTile, cave.Pillar, false, false)
					ToType(above1, cave.Pillar, false, false)
				} else {
					ToType(mainTile, cave.Empty, false, false)
					ToType(above1, cave.Empty, false, false)
					if !done && random.CaveGen.Intn(lastDown) > 5 {
						ToType(below1, cave.Empty, false, true)
						ToType(below2, cave.Empty, false, true)
						lastDown = 0
						downX = append(downX, currX)
					}
				}
				ToBlock(above2, false, false)
				ToBlock(above3, false, false)
				ToBlock(below1, false, false)
				ToBlock(below2, false, false)
			} else {
				break
			}
			if dir == data.Right {
				currX++
			} else {
				currX--
			}
			count++
			lastDown++
		}
		layer++
		currY += 4
		if len(downX) == 0 {
			break
		}
		x1 := downX[0]
		x2 := downX[len(downX)-1]
		minWidth := util.Max(util.Abs(x1-x2), 8)
		layerWidth = minWidth + 2 + random.CaveGen.Intn(int(c.FillVar*0.25))
		if dir == data.Right {
			currX = x1 - random.CaveGen.Intn(layerWidth-minWidth)
		} else {
			currX = x1 + random.CaveGen.Intn(layerWidth-minWidth)
		}
		downX = []int{}
	}
}

func MineLayer(c *cave.Cave, include world.Coords) {
	totalWidth := 16 + random.CaveGen.Intn(int(c.FillVar*0.25))
	offset := random.CaveGen.Intn(totalWidth)
	currX := include.X
	currX -= offset
	pillarX := currX + random.CaveGen.Intn(6)
	count := 0
	bg := random.CaveGen.Intn(2) == 0
	for count < totalWidth {
		if currX > c.Left*constants.ChunkSize+4 &&
			currX < (c.Right+1)*constants.ChunkSize-4 {
			mainTile := c.GetTileInt(currX, include.Y)
			above1 := c.GetTileInt(currX, include.Y-1)
			above2 := c.GetTileInt(currX, include.Y-2)
			above3 := c.GetTileInt(currX, include.Y-3)
			below1 := c.GetTileInt(currX, include.Y+1)
			below2 := c.GetTileInt(currX, include.Y+2)
			if count < 2 || count > totalWidth-3 {
				ToBlock(mainTile, false, false)
				ToBlock(above1, false, false)
				ToBlock(above2, false, false)
				ToBlock(below1, false, false)
			} else {
				if pillarX%6 == currX%6 {
					mainTile.BG = bg
					above1.BG = bg
					bg = !bg
					ToType(mainTile, cave.Pillar, false,  false)
					ToType(above1, cave.Pillar, false, false)
				} else {
					ToType(mainTile, cave.Empty, false, false)
					ToType(above1, cave.Empty, false, false)
				}
				ToBlock(above2, false, false)
				ToBlock(above3, false, false)
				ToBlock(below1, false, false)
				ToBlock(below2, false, false)
			}
		}
		currX++
		count++
	}
}

func Stairs(c *cave.Cave, include world.Coords, left, down bool, height, width int) {
	h := height
	if h < 1 {
		h = random.CaveGen.Intn(int(c.FillVar*0.25)) + int(c.FillVar*0.12)
	}
	w := width
	if w < 1 || w > h {
		w = random.CaveGen.Intn(int(c.FillVar*0.1)) + 5
	}
	curr := include
	count := 0
	turnC := 0
	for count < h {
		if (turnC != w-1 && down) || (turnC == 1 && !down) {
			stair := c.GetTileInt(curr.X, curr.Y+1)
			below1 := c.GetTileInt(curr.X, curr.Y+2)
			ToType(stair, cave.Wall, false, down)
			ToType(below1, cave.Wall, false, down)
		}
		above1 := c.GetTileInt(curr.X, curr.Y)
		above2 := c.GetTileInt(curr.X, curr.Y-1)
		above3 := c.GetTileInt(curr.X, curr.Y-2)
		above4 := c.GetTileInt(curr.X, curr.Y-3)
		above5 := c.GetTileInt(curr.X, curr.Y-4)
		ToType(above1, cave.Empty, true,  true)
		ToType(above2, cave.Empty, true,  true)
		ToType(above3, cave.Empty, true,  true)
		ToType(above4, cave.Wall, false, !down)
		ToType(above5, cave.Wall, false, !down)
		if turnC == w {
			turnC = 0
			left = !left
		}
		if left {
			curr.X--
		} else {
			curr.X++
		}
		if down {
			curr.Y++
		} else {
			curr.Y--
		}
		turnC++
		count++
	}
}
