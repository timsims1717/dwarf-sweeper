package structures

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/objects"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
)

func NoodleCave(c *cave.Cave, start world.Coords, iDir data.Direction) {
	t := start
	dir := iDir
	chance := random.CaveGen.Intn(40) + 40
	empty := random.CaveGen.Intn(6) > 0
	for {
		tile := c.GetTileInt(t.X, t.Y)
		if (chance < 25 && random.CaveGen.Intn(chance) == 0) || tile == nil || tile.NeverChange {
			break
		} else {
			chance--
			if empty {
				if random.CaveGen.Intn(15) == 0 {
					ToType(tile, cave.Growth, false, false)
				} else {
					ToType(tile, cave.Empty, false, false)
				}
			} else {
				ToBlock(tile, false, false)
			}
			// carve out a bit
			BlockUp(tile, cave.Unknown)
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
			empty = random.CaveGen.Intn(6) > 0
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
		fmt.Printf("WARNING: treasure room not generated: max %d is greater than cave width %d or height %d\n", max, width, height)
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
					ToType(tile, cave.Dig, true, true)
				} else if !tile.NeverChange && !tile.IsChanged {
					ToType(tile, cave.Empty, true, true)
					if y == tlY+h-2 && util.Contains(x, tC) {
						objects.AddChest(tile)
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
				if dist < fRad-random.CaveGen.Float64()*variance*world.TileSize {
					if !(tile.Type == cave.Wall && ignoreWalls) && !tile.Path {
						ToType(tile, cave.Blast, false, false)
					}
				}
			}
		}
	}
}

func Pocket(c *cave.Cave, radius int, variance float64, ignoreWalls bool, center world.Coords) {
	cPos := c.GetTileInt(center.X, center.Y).Transform.Pos
	fRad := float64(radius) * world.TileSize
	for y := center.Y - radius; y < center.Y+radius; y++ {
		for x := center.X - radius; x < center.X+radius; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				tPos := tile.Transform.Pos
				dist := util.Magnitude(cPos.Sub(tPos))
				if dist < fRad-random.CaveGen.Float64()*variance*world.TileSize {
					if !(tile.Type == cave.Wall && ignoreWalls) {
						ToType(tile, cave.Empty, false, false)
					}
				}
			}
		}
	}
	for x := center.X - radius; x < center.X+radius; x++ {
		bg := random.CaveGen.Intn(2) == 0
		for y := center.Y - radius; y < center.Y+radius; y++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				if random.CaveGen.Intn(4) == 0 {
					above := c.GetTileInt(x, y-1)
					below := c.GetTileInt(x, y+1)
					if above != nil && below != nil {
						if above.Solid() || below.Solid() || above.Type == cave.Growth || below.Type == cave.Growth {
							ToType(tile, cave.Growth, false, true)
							tile.BG = bg
							continue
						}
					}
				}
			}
			bg = random.CaveGen.Intn(2) == 0
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
				if tile.RCoords == center || (y == center.Y && dist < world.TileSize*0.5*random.CaveGen.Float64()*variance) {
					if !(tile.Type == cave.Wall && ignoreWalls) {
						ToType(tile, cave.Blast, false, false)
					}
				} else if dist < inRad-random.CaveGen.Float64()*variance*world.TileSize {
					if !(tile.Type == cave.Wall && ignoreWalls) {
						ToBlock(tile, false, false)
					}
				} else if dist < fRad-random.CaveGen.Float64()*variance*world.TileSize {
					if !(tile.Type == cave.Wall && ignoreWalls) {
						ToType(tile, cave.Empty, false, false)
					}
				}
			}
		}
	}
	for x := center.X - radius; x < center.X+radius; x++ {
		bg := random.CaveGen.Intn(2) == 0
		for y := center.Y - radius; y < center.Y+radius; y++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				if random.CaveGen.Intn(4) == 0 {
					above := c.GetTileInt(x, y-1)
					below := c.GetTileInt(x, y+1)
					if above != nil && below != nil {
						if above.Solid() || below.Solid() || above.Type == cave.Growth || below.Type == cave.Growth {
							ToType(tile, cave.Growth, false, true)
							tile.BG = bg
							continue
						}
					}
				}
			}
			bg = random.CaveGen.Intn(2) == 0
		}
	}
}

func BombRoom(c *cave.Cave, minH, maxH, minW, maxW, curve, level int, include world.Coords) {
	width, height := c.Dimensions()
	if maxW > width*constants.ChunkSize || maxH > height*constants.ChunkSize {
		fmt.Printf("WARNING: bomb room not generated: max width %d or max height %d is greater than cave width %d or height %d\n", maxW, maxH, width, height)
		return
	}
	if minH < 3 {
		minH = 3
	}
	if minW < 5 {
		minW = 5
	}
	sW := maxW - minW
	sH := maxH - minH
	w := random.CaveGen.Intn(sW) + minW
	h := random.CaveGen.Intn(sH) + minH
	sX := include.X - int(float64(w) * 0.33) + random.CaveGen.Intn(int(float64(w) * 0.167))
	sY := include.Y - h + 1
	for y := sY; y < sY+h; y++ {
		for x := sX; x < sX+w; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				dx := util.Min(util.Abs(x-sX), util.Abs(x-(sX+width)))
				dy := util.Abs(y - sY)
				if !tile.NeverChange && !tile.IsChanged && !(dx+dy+random.CaveGen.Intn(2) < curve*maxW/8) {
					ToType(tile, cave.Empty, true, true)
					if y == include.Y && x == include.X {
						objects.AddBigBomb(tile, level)
						t1 := c.GetTileInt(x, y+1)
						t2 := c.GetTileInt(x+1, y+1)
						ToType(t1, cave.Wall, true, true)
						ToType(t2, cave.Wall, true, true)
					}
				}
			}
		}
	}
}