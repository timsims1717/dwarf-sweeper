package dungeon

import (
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"math/rand"
)

func EntranceExit(cave *Cave, door world.Coords, width, height int, roofCurve int) {
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
					if (x == l || x == r || y == t || y == b) || dx + dy + rand.Intn(2) < curve + width / 8 {
						if y == b && util.Abs(x - door.X) < 2 {
							tile.Type = Wall
							tile.breakable = false
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
						tile.BGSpriteS = s
						if s != "" {
							tile.BGSprite = cave.batcher.Sprites[s]
						} else {
							tile.BGSprite = nil
						}
						tile.breakable = false
					}
					tile.bomb = false
					tile.neverChange = true
					tile.isChanged = true
					tile.Entities = []Entity{}
					tile.UpdateSprites()
				}
			}
		}
	}
}

func RectRoom(cave *Cave, min, max int) {
	s := max - min
	width, height := cave.Dimensions()
	if max > width * ChunkSize || max > height * ChunkSize {
		fmt.Printf("WARNING: rect room not generated: max %d is greater than cave width %d or height %d\n", max, width, height)
		return
	}
	w := rand.Intn(s) + min
	h := rand.Intn(s) + min
	sX := rand.Intn(width - w)
	sY := rand.Intn(height - h)
	for y := sY; y < sY + h; y++ {
		for x := sX; x < sX+w; x++ {
			tile := cave.GetTileInt(x, y)
			if tile != nil {
				if !tile.neverChange && !tile.isChanged && (x == sX || x == sX+w-1 || y == sY || y == sY+h-1) {
					tile.Type = Wall
					tile.breakable = false
					tile.Solid = true
					tile.bomb = false
					tile.Entities = []Entity{}
					tile.UpdateSprites()
				} else {
					tile.isChanged = true
				}
			}
		}
	}
}