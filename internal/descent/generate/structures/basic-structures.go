package structures

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
)

func RectRoom(c *cave.Cave, tl world.Coords, width, height, curve int, blockType cave.BlockType) {
	for y := tl.Y; y < tl.Y+height; y++ {
		for x := tl.X; x < tl.X+width; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				dx := util.Min(util.Abs(x-tl.X), util.Abs(x-(tl.X+width)))
				dy := util.Abs(y - tl.Y)
				if !tile.NeverChange && !tile.IsChanged && !(dx+dy+random.CaveGen.Intn(2) < curve*width/8) {
					if blockType != cave.Unknown {
						tile.Type = blockType
					} else if tile.Perlin < 0 {
						tile.Type = cave.Collapse
					} else {
						tile.Type = cave.Dig
					}
					tile.IsChanged = true
				}
			}
		}
	}
}

func RandRectRoom(c *cave.Cave, min, max int, include world.Coords) ([]world.Coords, []world.Coords) {
	s := max - min
	width, height := c.Dimensions()
	if max > width*constants.ChunkSize || max > height*constants.ChunkSize {
		fmt.Printf("WARNING: rect room not generated: max %d is greater than cave width %d or height %d\n", max, width, height)
		return nil, nil
	}
	//popperMade := false
	toMark := 20
	var roomTiles, marked []world.Coords
	w := random.CaveGen.Intn(s) + min
	h := random.CaveGen.Intn(s) + min
	sX := include.X - w + 1
	sY := include.Y - h + 1
	tlX := random.CaveGen.Intn(w-2) + sX
	tlY := random.CaveGen.Intn(h-2) + sY
	for y := tlY; y < tlY+h; y++ {
		for x := tlX; x < tlX+w; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				if !tile.NeverChange && !tile.IsChanged && (x == tlX || x == tlX+w-1 || y == tlY || y == tlY+h-1) {
					ToType(tile, cave.Wall, false, true)
				} else if !tile.NeverChange {
					ToBlock(tile, false, true)
					roomTiles = append(roomTiles, tile.RCoords)
					if random.CaveGen.Intn(toMark) == 0 {
						toMark += 12
						marked = append(marked, tile.RCoords)
					//} else {
					//	toMark--
					//	if !tile.Bomb && !popperMade && random.CaveGen.Intn(125) == 0 {
					//		p := descent.Popper{}
					//		p.Create(tile.Transform.Pos)
					//		popperMade = true
					//		tile.Special = true
					//	}
					}
				}
			}
		}
	}
	return roomTiles, marked
}
