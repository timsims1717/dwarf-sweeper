package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
)

func Entrance(c *cave.Cave, door world.Coords, width, height int, roofCurve int, exit bool) {
	t := door.Y - height
	b := door.Y + 1
	l := door.X - width / 2
	r := door.X + width / 2
	for y := t; y <= b; y++ {
		for x := l; x <= r; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				if !tile.NeverChange && !tile.IsChanged {
					dx := util.Min(util.Abs(x - l), util.Abs(x - r))
					dy := util.Abs(y - t)
					curve := roofCurve
					if (x == l || x == r || y == t || y == b) || dx + dy + random.CaveGen.Intn(2) < curve + width / 8 {
						if y == b && util.Abs(x - door.X) < 2 {
							tile.Type = cave.Wall
							tile.NeverChange = true
						} else {
							tile.Type = cave.BlockCollapse
						}
					} else {
						ty := cave.Empty
						s := ""
						if x == door.X && y == door.Y {
							ty = cave.Deco
							s = "door"
							if exit {
								tile.Exit = true
							}
						} else if x == door.X - 1 && y == door.Y {
							ty = cave.Deco
							s = "door_l"
						} else if x == door.X + 1 && y == door.Y {
							ty = cave.Deco
							s = "door_r"
						} else if x == door.X && y == door.Y - 1 {
							ty = cave.Deco
							s = "door_t"
						} else if x == door.X - 1 && y == door.Y - 1 {
							ty = cave.Deco
							s = "door_tl"
						} else if x == door.X + 1 && y == door.Y - 1 {
							ty = cave.Deco
							s = "door_tr"
						}
						tile.Type = cave.TileType(ty)
						tile.NeverChange = true
						tile.BGSpriteS = s
						if s != "" {
							tile.BGSpriteS = s
							tile.BGSprite = c.Batcher.Sprites[s]
							tile.BGMatrix = pixel.IM
						} else {
							tile.BGSprite = nil
						}
					}
					tile.Bomb = false
					tile.IsChanged = true
					tile.Entity = nil
				}
			}
		}
	}
	if exit {
		_, h := c.Dimensions()
		for y := door.Y + 2; y <= h; y++ {
			for x := door.X - 1; x <= door.X+1; x++ {
				tile := c.GetTileInt(x, y)
				if tile != nil {
					if !tile.NeverChange && !tile.IsChanged {
						tile.Type = cave.Wall
						tile.NeverChange = true
						tile.Bomb = false
						tile.IsChanged = true
						tile.Entity = nil
					}
				}
			}
		}
	}
}

func RectRoom(c *cave.Cave, tl world.Coords, width, height int) {
	for y := tl.Y; y < tl.Y + height; y++ {
		for x := tl.X; x < tl.X + width; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				if !tile.NeverChange && !tile.IsChanged && (x == tl.X || x == tl.X+width-1 || y == tl.Y || y == tl.Y+height-1) {
					tile.Type = cave.Wall
				} else if !tile.NeverChange && !tile.IsChanged {
					tile.Type = cave.BlockCollapse
					tile.IsChanged = true
				}
			}
		}
	}
}

func RandRectRoom(c *cave.Cave, min, max int, include world.Coords) ([]world.Coords, []world.Coords) {
	s := max - min
	width, height := c.Dimensions()
	if max > width * constants.ChunkSize || max > height * constants.ChunkSize {
		fmt.Printf("WARNING: rect room not generated: max %d is greater than cave width %d or height %d\n", max, width, height)
		return nil, nil
	}
	popperMade := false
	toMark := 20
	var roomTiles, marked []world.Coords
	w := random.CaveGen.Intn(s) + min
	h := random.CaveGen.Intn(s) + min
	sX := include.X - w + 1
	sY := include.Y - h + 1
	tlX := random.CaveGen.Intn(w - 2) + sX
	tlY := random.CaveGen.Intn(h - 2) + sY
	for y := tlY; y < tlY + h; y++ {
		for x := tlX; x < tlX+w; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				if !tile.NeverChange && !tile.IsChanged && (x == tlX || x == tlX+w-1 || y == tlY || y == tlY+h-1) {
					toWall(tile, false)
				} else if !tile.NeverChange && !tile.IsChanged {
					if random.CaveGen.Intn(10) == 0 {
						toBlockCollapse(tile, false)
					} else {
						toBlockDig(tile, false)
					}
					roomTiles = append(roomTiles, tile.RCoords)
					tile.Room = true
					if random.CaveGen.Intn(toMark) == 0 {
						toMark += 12
						marked = append(marked, tile.RCoords)
					} else {
						toMark--
						if !tile.Bomb && !popperMade && random.CaveGen.Intn(20) == 0 {
							p := descent.Popper{}
							p.Create(tile.Transform.Pos)
							popperMade = true
							tile.Special = true
						}
					}
				}
			}
		}
	}
	return roomTiles, marked
}