package structures

import (
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
)

func GnomeMineLayer(c *cave.Cave, includeL, includeR world.Coords) []*cave.Tile {
	var updated []*cave.Tile
	minWidth := util.Abs(includeL.X - includeR.X)
	totalWidth := minWidth + 16 + random.CaveGen.Intn(6)
	offset := random.CaveGen.Intn(totalWidth - minWidth - 4)
	curr := includeL
	curr.X -= offset
	pillarX := curr.X + 3 + random.CaveGen.Intn(4)
	count := 0
	bg := random.CaveGen.Intn(2) == 0
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
			if count < 2 || count > totalWidth-3 {
				ToType(main1, cave.Dig, true, true)
				ToType(main2, cave.Dig, true, true)
			} else if pillarX%6 == curr.X%6 {
				main1.BG = bg
				main2.BG = bg
				bg = !bg
				ToType(main1, cave.Pillar, true, true)
				ToType(main2, cave.Pillar, true, true)
			} else {
				ToType(main1, cave.Empty, true, true)
				ToType(main2, cave.Empty, true, true)
			}
			ToType(above1, cave.Dig, true, true)
			ToType(above2, cave.Blast, true, true)
			ToType(above3, cave.Blast, true, true)
			ToType(below1, cave.Dig, true, true)
			ToType(below2, cave.Dig, true, true)
			ToType(below3, cave.Blast, true, true)
			ToType(below4, cave.Blast, true, true)
			updated = append(updated, []*cave.Tile{
				main1, main2, above1, above2, above3, below1, below2, below3, below4,
			}...)
		}
		curr.X++
		count++
	}
	return updated
}
