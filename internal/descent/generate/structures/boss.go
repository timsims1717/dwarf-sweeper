package structures

import (
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

func GnomeMineLayer(c *cave.Cave, includeL, includeR world.Coords) {
	minWidth := util.Abs(includeL.X - includeR.X)
	totalWidth := minWidth + 16 + random.CaveGen.Intn(12)
	offset := random.CaveGen.Intn(totalWidth - minWidth - 4)
	curr := includeL
	curr.X -= offset
	pillarX := curr.X + 3 + random.CaveGen.Intn(4)
	count := 0
	for count < totalWidth {
		main1 := c.GetTileInt(curr.X, curr.Y)
		if main1 != nil {
			above1 := c.GetTileInt(curr.X, curr.Y-2)
			above2 := c.GetTileInt(curr.X, curr.Y-3)
			above3 := c.GetTileInt(curr.X, curr.Y-4)
			main2 := c.GetTileInt(curr.X, curr.Y-1)
			below1 := c.GetTileInt(curr.X, curr.Y+1)
			below2 := c.GetTileInt(curr.X, curr.Y+2)
			below3 := c.GetTileInt(curr.X, curr.Y+3)
			below4 := c.GetTileInt(curr.X, curr.Y+4)
			if count < 2 || count > totalWidth - 3 {
				ToBlock(main1, cave.BlockDig, true, true)
				ToBlock(main2, cave.BlockDig, true, true)
			} else if pillarX%6 == curr.X%6 {
				ToEmpty(main1, true, true, true)
				main1.BGSpriteS = "pillar"
				main1.BGSprite = c.Batcher.Sprites["pillar"]
				main1.BGMatrix = pixel.IM
				ToEmpty(main2, true, true, true)
				main2.BGSpriteS = "pillar_top"
				main2.BGSprite = c.Batcher.Sprites["pillar_top"]
				main2.BGMatrix = pixel.IM
			} else {
				ToEmpty(main1, true, false, true)
				ToEmpty(main2, true, false, true)
			}
			ToBlock(above1, cave.BlockDig, true, true)
			ToBlock(above2, cave.BlockBlast, true, true)
			ToBlock(above3, cave.BlockBlast, true, true)
			ToBlock(below1, cave.BlockDig, true, true)
			ToBlock(below2, cave.BlockDig, true, true)
			ToBlock(below3, cave.BlockBlast, true, true)
			ToBlock(below4, cave.BlockBlast, true, true)
		}
		curr.X++
		count++
	}
}