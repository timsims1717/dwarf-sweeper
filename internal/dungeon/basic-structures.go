package dungeon

import (
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
)

func Entrance(cave *Cave, door world.Coords, width, height int, roofCurve int, exit bool) {
	t := door.Y - height
	b := door.Y + 1
	l := door.X - width / 2
	r := door.X + width / 2
	for y := t; y <= b; y++ {
		for x := l; x <= r; x++ {
			tile := cave.GetTileInt(x, y)
			if tile != nil {
				if !tile.neverChange && !tile.isChanged {
					dx := util.Min(util.Abs(x - l), util.Abs(x - r))
					dy := util.Abs(y - t)
					curve := roofCurve
					if (x == l || x == r || y == t || y == b) || dx + dy + random.CaveGen.Intn(2) < curve + width / 8 {
						if y == b && util.Abs(x - door.X) < 2 {
							tile.Type = Wall
							tile.breakable = false
							tile.neverChange = true
						} else {
							tile.Type = Block
							tile.breakable = true
						}
						tile.Solid = true
					} else {
						ty := Empty
						s := ""
						if x == door.X && y == door.Y {
							ty = Deco
							s = "door"
							if exit {
								tile.Exit = true
							}
						} else if x == door.X - 1 && y == door.Y {
							ty = Deco
							s = "door_l"
						} else if x == door.X + 1 && y == door.Y {
							ty = Deco
							s = "door_r"
						} else if x == door.X && y == door.Y - 1 {
							ty = Deco
							s = "door_t"
						} else if x == door.X - 1 && y == door.Y - 1 {
							ty = Deco
							s = "door_tl"
						} else if x == door.X + 1 && y == door.Y - 1 {
							ty = Deco
							s = "door_tr"
						}
						tile.Solid = false
						tile.Type = TileType(ty)
						tile.neverChange = true
						tile.BGSpriteS = s
						if s != "" {
							tile.BGSprite = cave.batcher.Sprites[s]
						} else {
							tile.BGSprite = nil
						}
						tile.breakable = false
					}
					tile.bomb = false
					tile.isChanged = true
					tile.Entity = nil
					tile.UpdateSprites()
				}
			}
		}
	}
}

func Exit(cave *Cave, door world.Coords, width, height int, roofCurve int, exit bool) {
	t := door.Y - height
	b := door.Y + 1
	l := door.X - width / 2
	r := door.X + width / 2
	for y := t; y <= b; y++ {
		for x := l; x <= r; x++ {
			tile := cave.GetTileInt(x, y)
			if tile != nil {
				if !tile.neverChange && !tile.isChanged {
					dx := util.Min(util.Abs(x - l), util.Abs(x - r))
					dy := util.Abs(y - t)
					curve := roofCurve
					if (x == l || x == r || y == t || y == b) || dx + dy + random.CaveGen.Intn(2) < curve + width / 8 {
						if y == b && util.Abs(x - door.X) < 2 {
							tile.Type = Wall
							tile.breakable = false
							tile.neverChange = true
						} else {
							tile.Type = Block
							tile.breakable = true
						}
						tile.Solid = true
					} else {
						ty := Empty
						s := ""
						if x == door.X && y == door.Y {
							ty = Deco
							s = "door"
							if exit {
								tile.Exit = true
							}
						} else if x == door.X - 1 && y == door.Y {
							ty = Deco
							s = "door_l"
						} else if x == door.X + 1 && y == door.Y {
							ty = Deco
							s = "door_r"
						} else if x == door.X && y == door.Y - 1 {
							ty = Deco
							s = "door_t"
						} else if x == door.X - 1 && y == door.Y - 1 {
							ty = Deco
							s = "door_tl"
						} else if x == door.X + 1 && y == door.Y - 1 {
							ty = Deco
							s = "door_tr"
						}
						tile.Solid = false
						tile.Type = TileType(ty)
						tile.neverChange = true
						tile.BGSpriteS = s
						if s != "" {
							tile.BGSprite = cave.batcher.Sprites[s]
						} else {
							tile.BGSprite = nil
						}
						tile.breakable = false
					}
					tile.bomb = false
					tile.isChanged = true
					tile.Entity = nil
					tile.UpdateSprites()
				}
			}
		}
	}
	_, h := cave.Dimensions()
	for y := door.Y + 2; y <= h; y++ {
		for x := door.X - 1; x <= door.X + 1; x++ {
			tile := cave.GetTileInt(x, y)
			if tile != nil {
				if !tile.neverChange && !tile.isChanged {
					tile.Type = Wall
					tile.breakable = false
					tile.neverChange = true
					tile.bomb = false
					tile.isChanged = true
					tile.Entity = nil
					tile.UpdateSprites()
				}
			}
		}
	}
}

func RectRoom(cave *Cave, tl world.Coords, width, height int) {
	for y := tl.Y; y < tl.Y + height; y++ {
		for x := tl.X; x < tl.X + width; x++ {
			tile := cave.GetTileInt(x, y)
			if tile != nil {
				if !tile.neverChange && !tile.isChanged && (x == tl.X || x == tl.X+width-1 || y == tl.Y || y == tl.Y+height-1) {
					tile.Type = Wall
					tile.breakable = false
					tile.Solid = true
					tile.UpdateSprites()
				} else if !tile.neverChange && !tile.isChanged {
					tile.Solid = true
					tile.Type = Block
					tile.breakable = true
					tile.isChanged = true
					tile.UpdateSprites()
				}
			}
		}
	}
}

func RandRectRoom(cave *Cave, min, max int, include world.Coords) {
	s := max - min
	width, height := cave.Dimensions()
	if max > width * ChunkSize || max > height * ChunkSize {
		fmt.Printf("WARNING: rect room not generated: max %d is greater than cave width %d or height %d\n", max, width, height)
		return
	}
	w := random.CaveGen.Intn(s) + min
	h := random.CaveGen.Intn(s) + min
	sX := include.X - w + 1
	sY := include.Y - h + 1
	tlX := random.CaveGen.Intn(w - 2) + sX
	tlY := random.CaveGen.Intn(h - 2) + sY
	for y := tlY; y < tlY + h; y++ {
		for x := tlX; x < tlX+w; x++ {
			tile := cave.GetTileInt(x, y)
			if tile != nil {
				if !tile.neverChange && !tile.isChanged && (x == tlX || x == tlX+w-1 || y == tlY || y == tlY+h-1) {
					tile.Type = Wall
					tile.breakable = false
					tile.Solid = true
					tile.UpdateSprites()
				} else if !tile.neverChange && !tile.isChanged {
					tile.Solid = true
					tile.Type = Block
					tile.breakable = true
					tile.isChanged = true
					tile.UpdateSprites()
				}
			}
		}
	}
}