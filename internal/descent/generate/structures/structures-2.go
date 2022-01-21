package structures

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
	case 0, 1, 2:
		currType = cave.Empty
	case 3:
		currType = cave.BlockCollapse
	}
	for {
		tile := c.GetTileInt(t.X, t.Y)
		if (chance < 25 && random.CaveGen.Intn(chance) == 0) || tile == nil || tile.NeverChange {
			break
		} else {
			chance--
			if currType == cave.Empty {
				ToEmpty(tile, false, false, false)
			} else {
				ToBlock(tile, cave.TileType(currType), false, true)
			}
			// carve out a bit
			ns := tile.RCoords.Neighbors()
			for _, n := range ns {
				tmp := c.GetTileInt(n.X, n.Y)
				if random.CaveGen.Intn(4) == 0 {
					ToBlock(tmp, cave.BlockDig, false, true)
				} else {
					ToBlock(tmp, cave.BlockCollapse, false, true)
				}
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
	if max > width*constants.ChunkSize || max > height*constants.ChunkSize {
		fmt.Printf("WARNING: rect room not generated: max %d is greater than cave width %d or height %d\n", max, width, height)
		return
	}
	w := random.CaveGen.Intn(s) + min
	h := random.CaveGen.Intn(s) + min
	sX := include.X - w + 1
	sY := include.Y - h + 1
	tlX := random.CaveGen.Intn(w-2) + sX
	tlY := random.CaveGen.Intn(h-2) + sY
	tW := w - 2
	tC := util.RandomSampleRange(tTotal, tlX+1, tlX+tW, random.CaveGen)
	for y := tlY; y < tlY+h; y++ {
		for x := tlX; x < tlX+w; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				if !tile.NeverChange && !tile.IsChanged && (x == tlX || x == tlX+w-1 || y == tlY || y == tlY+h-1) {
					ToBlock(tile, cave.BlockDig, true, true)
				} else if !tile.NeverChange && !tile.IsChanged {
					ToEmpty(tile, true, false, true)
					if y == tlY+h-2 && util.Contains(x, tC) {
						addChest(tile)
					}
				}
			}
		}
	}
}

func BombableNode(c *cave.Cave, radius int, variance float64, ignoreWalls bool, center world.Coords) {
	cPos := c.GetTileInt(center.X, center.Y).Transform.Pos
	fRad := float64(radius) * world.TileSize
	for y := center.Y - radius; y < center.Y+radius; y++ {
		for x := center.X - radius; x < center.X+radius; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				tPos := tile.Transform.Pos
				dist := util.Magnitude(cPos.Sub(tPos))
				if dist < fRad+random.CaveGen.Float64()*variance {
					if !(tile.Type == cave.Wall && ignoreWalls) && !tile.Path {
						ToBlock(tile, cave.BlockBlast, false, false)
					}
				}
			}
		}
	}
}

func Pocket(c *cave.Cave, radius int, variance float64, ignoreWalls bool, center world.Coords) {
	cPos := c.GetTileInt(center.X, center.Y).Transform.Pos
	fRad := float64(radius) * world.TileSize
	inRad := fRad * 0.35
	for y := center.Y - radius; y < center.Y+radius; y++ {
		for x := center.X - radius; x < center.X+radius; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				tPos := tile.Transform.Pos
				dist := util.Magnitude(cPos.Sub(tPos))
				if dist < inRad+random.CaveGen.Float64()*variance && !tile.Bomb {
					if !(tile.Type == cave.Wall && ignoreWalls) {
						ToEmpty(tile, false, false, false)
					}
				} else if dist < fRad+random.CaveGen.Float64()*variance {
					if !(tile.Type == cave.Wall && ignoreWalls) {
						ToBlock(tile, cave.BlockCollapse, false, true)
					}
				}
			}
		}
	}
}

func Ring(c *cave.Cave, radius int, variance float64, ignoreWalls bool, center world.Coords) {
	cPos := c.GetTileInt(center.X, center.Y).Transform.Pos
	fRad := float64(radius) * world.TileSize
	inRad := fRad * 0.5
	for y := center.Y - radius; y < center.Y+radius; y++ {
		for x := center.X - radius; x < center.X+radius; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				tPos := tile.Transform.Pos
				dist := util.Magnitude(cPos.Sub(tPos))
				if tile.RCoords == center || (y == center.Y && dist < world.TileSize*0.5+random.CaveGen.Float64()*variance && !tile.Bomb) {
					if !(tile.Type == cave.Wall && ignoreWalls) {
						if tile.Path {
							ToBlock(tile, cave.BlockDig, true, true)
						} else {
							ToBlock(tile, cave.Wall, true, true)
						}
					}
				} else if dist < inRad+random.CaveGen.Float64()*variance && !tile.Bomb {
					if !(tile.Type == cave.Wall && ignoreWalls) {
						ToEmpty(tile, false, false, false)
					}
				} else if dist < fRad+random.CaveGen.Float64()*variance {
					if !(tile.Type == cave.Wall && ignoreWalls) {
						ToBlock(tile, cave.BlockCollapse, false, false)
					}
				}
			}
		}
	}
}
