package structures

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

func MineLayer(c *cave.Cave, include world.Coords) {
	totalWidth := 16 + random.CaveGen.Intn(int(c.FillVar * 0.25))
	offset := random.CaveGen.Intn(totalWidth)
	currX := include.X
	currX -= offset
	pillarX := currX + random.CaveGen.Intn(6)
	count := 0
	for count < totalWidth {
		if currX > c.Left*constants.ChunkSize+4 &&
			currX < (c.Right+1)*constants.ChunkSize-4 {
			mainTile := c.GetTileInt(currX, include.Y)
			above1 := c.GetTileInt(currX, include.Y-1)
			above2 := c.GetTileInt(currX, include.Y-2)
			above3 := c.GetTileInt(currX, include.Y-3)
			below1 := c.GetTileInt(currX, include.Y+1)
			below2 := c.GetTileInt(currX, include.Y+2)
			if count < 2 || count > totalWidth - 3 {
				ToBlock(mainTile, cave.BlockDig, false, false)
				ToBlock(above1, cave.BlockDig, false, false)
				ToBlock(above2, cave.BlockDig, false, false)
				ToBlock(below1, cave.BlockDig, false, false)
			} else {
				if pillarX % 6 == currX % 6 {
					ToEmpty(mainTile, false, true, false)
					mainTile.BGSpriteS = "pillar"
					mainTile.BGSprite = c.Batcher.Sprites["pillar"]
					mainTile.BGMatrix = pixel.IM
					ToEmpty(above1, false, true, false)
					above1.BGSpriteS = "pillar_top"
					above1.BGSprite = c.Batcher.Sprites["pillar_top"]
					above1.BGMatrix = pixel.IM
				} else {
					ToEmpty(mainTile, false, false, false)
					ToEmpty(above1, false, false, false)
				}
				ToBlock(above2, cave.BlockDig, false, false)
				ToBlock(above3, cave.BlockDig, false, false)
				ToBlock(below1, cave.BlockDig, false, false)
				ToBlock(below2, cave.BlockDig, false, false)
			}
		}
		currX++
		count++
	}
}

func Stairs(c *cave.Cave, include world.Coords, left, down bool, height, width int) {
	h := height
	if h < 1 {
		h = random.CaveGen.Intn(int(c.FillVar * 0.25)) + int(c.FillVar * 0.12)
	}
	w := width
	if w < 1 || w > h {
		w = random.CaveGen.Intn(int(c.FillVar * 0.1)) + 5
	}
	curr := include
	count := 0
	turnC := 0
	done := false
	for count < h && !done {
		if (turnC != w - 1 && down) || (turnC == 1 && !down) {
			stair := c.GetTileInt(curr.X, curr.Y+1)
			below1 := c.GetTileInt(curr.X, curr.Y+2)
			ToBlock(stair, cave.Wall, false, down)
			ToBlock(below1, cave.Wall, false, down)
		}
		above1 := c.GetTileInt(curr.X, curr.Y)
		above2 := c.GetTileInt(curr.X, curr.Y-1)
		above3 := c.GetTileInt(curr.X, curr.Y-2)
		above4 := c.GetTileInt(curr.X, curr.Y-3)
		above5 := c.GetTileInt(curr.X, curr.Y-4)
		ToEmpty(above1, true, false, true)
		ToEmpty(above2, true, false, true)
		ToEmpty(above3, true, false, true)
		ToBlock(above4, cave.Wall, false, !down)
		ToBlock(above5, cave.Wall, false, !down)
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