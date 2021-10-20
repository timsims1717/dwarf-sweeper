package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
)

func NoodleCave(c *cave.Cave, start world.Coords, iDir data.Direction) {
	t := start
	dir := iDir
	currType := cave.BlockCollapse
	chance := random.CaveGen.Intn(40) + 40
	switch random.CaveGen.Intn(4) {
	case 0,1,2:
		currType = cave.Empty
	case 3:
		currType = cave.BlockCollapse
	//case 2:
	//	currType = cave.Block1
	}
	for {
		tile := c.GetTileInt(t.X, t.Y)
		if (chance < 25 && random.CaveGen.Intn(chance) == 0) || tile == nil || tile.NeverChange {
			break
		} else {
			chance--
			if tile.Type == cave.Wall {
				tile.Special = true
			}
			tile.Type = cave.TileType(currType)
			tile.UpdateSprites()
			tile.IsChanged = true
			// carve out a bit
			ns := tile.RCoords.Neighbors()
			for _, n := range ns {
				tmp := c.GetTileInt(n.X, n.Y)
				toBlock(tmp, false, tmp.Bomb)
			}
			// change to next tile
			switch dir {
			case data.Left:
				t.X--
			case data.Right:
				t.X++
			case data.Up:
				t.Y--
			case data.Down:
				t.Y++
			}
			// change type
			switch random.CaveGen.Intn(2) {
			case 0:
				currType = cave.Empty
			case 1:
				currType = cave.BlockCollapse
			}
			// maybe change direction
			change := random.CaveGen.Intn(4)
			switch iDir {
			case data.Left:
				switch change {
				case 0:
					dir = data.Up
				case 1:
					dir = data.Down
				default:
					dir = data.Left
				}
			case data.Right:
				switch change {
				case 0:
					dir = data.Up
				case 1:
					dir = data.Down
				default:
					dir = data.Right
				}
			case data.Up:
				switch change {
				case 0:
					dir = data.Left
				case 1:
					dir = data.Right
				default:
					dir = data.Up
				}
			case data.Down:
				switch change {
				case 0:
					dir = data.Left
				case 1:
					dir = data.Right
				default:
					dir = data.Down
				}
			}
		}
	}
}

func TreasureRoom(c *cave.Cave, min, max, tTotal int, include world.Coords) {
	s := max - min
	width, height := c.Dimensions()
	if max > width * constants.ChunkSize || max > height * constants.ChunkSize {
		fmt.Printf("WARNING: rect room not generated: max %d is greater than cave width %d or height %d\n", max, width, height)
		return
	}
	w := random.CaveGen.Intn(s) + min
	h := random.CaveGen.Intn(s) + min
	sX := include.X - w + 1
	sY := include.Y - h + 1
	tlX := random.CaveGen.Intn(w - 2) + sX
	tlY := random.CaveGen.Intn(h - 2) + sY
	tW := w-2
	tC := util.RandomSampleRange(tTotal, tlX+1, tlX+tW, random.CaveGen)
	for y := tlY; y < tlY + h; y++ {
		for x := tlX; x < tlX+w; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				if !tile.NeverChange && !tile.IsChanged && (x == tlX || x == tlX+w-1 || y == tlY || y == tlY+h-1) && !tile.Breakable() {
					tile.Type = cave.Wall
					tile.UpdateSprites()
				} else if !tile.NeverChange && !tile.IsChanged {
					if y == tlY+h-2 && util.Contains(x, tC) {
						addChest(tile)
					}
					tile.Type = cave.Empty
					tile.IsChanged = true
					tile.Fillable = false
					tile.UpdateSprites()
				}
			}
		}
	}
}