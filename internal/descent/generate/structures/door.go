package structures

import (
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
)

func Entrance(c *cave.Cave, door world.Coords, width, height int, roofCurve int, doorType cave.BlockType) {
	t := door.Y - height
	b := door.Y + 1
	l := door.X - width/2
	r := door.X + width/2
	for y := t; y <= b; y++ {
		for x := l; x <= r; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				dx := util.Min(util.Abs(x-l), util.Abs(x-r))
				dy := util.Abs(y - t)
				curve := roofCurve
				if (x == l || x == r || y == t || y == b) || dx+dy+random.CaveGen.Intn(2) < curve+width/8 {
					if y == b && util.Abs(x-door.X) < 2 {
						tile.Type = cave.Wall
						tile.NeverChange = true
					} else {
						tile.Type = cave.Collapse
					}
				} else {
					if (x == door.X && y == door.Y) ||
						(x == door.X-1 && y == door.Y) ||
						(x == door.X+1 && y == door.Y) ||
						(x == door.X && y == door.Y-1) ||
						(x == door.X-1 && y == door.Y-1) ||
						(x == door.X+1 && y == door.Y-1) {
						tile.Type = doorType
						tile.DoorI = c.DoorI
					} else {
						tile.Type = cave.Empty
					}
					tile.NeverChange = true
				}
				tile.Bomb = false
				tile.IsChanged = true
			}
		}
	}
	c.DoorI++
}

func Exit(c *cave.Cave, door world.Coords, width, height int, roofCurve int, exitI int, doorType cave.BlockType) {
	t := door.Y - height
	b := door.Y + 1
	l := door.X - width/2
	r := door.X + width/2
	for y := t; y <= b; y++ {
		for x := l; x <= r; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				dx := util.Min(util.Abs(x-l), util.Abs(x-r))
				dy := util.Abs(y - t)
				curve := roofCurve
				if (x == l || x == r || y == t || y == b) || dx+dy+random.CaveGen.Intn(2) < curve+width/8 {
					if y == b && util.Abs(x-door.X) < 2 {
						tile.Type = cave.Wall
						tile.NeverChange = true
					} else {
						tile.Type = cave.Collapse
					}
				} else {
					if x == door.X && y == door.Y {
						ExitTile(c, door, exitI)
						tile.Type = doorType
						tile.DoorI = c.DoorI
					} else if (x == door.X-1 && y == door.Y) ||
						(x == door.X+1 && y == door.Y) ||
						(x == door.X && y == door.Y-1) ||
						(x == door.X-1 && y == door.Y-1) ||
						(x == door.X+1 && y == door.Y-1) {
						tile.Type = doorType
						tile.DoorI = c.DoorI
					} else {
						tile.Type = cave.Empty
					}
					tile.NeverChange = true
				}
				tile.Bomb = false
				tile.IsChanged = true
			}
		}
	}
	if exitI == 0 {
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
					}
				}
			}
		}
	}
	c.DoorI++
}

func EntranceDoor(c *cave.Cave, door world.Coords, doorType cave.BlockType) {
	for y := door.Y - 1; y <= door.Y+1; y++ {
		for x := door.X - 1; x <= door.X+1; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				if x == door.X && y == door.Y {
					tile.Type = doorType
					tile.DoorI = c.DoorI
				} else if (x == door.X-1 && y == door.Y) ||
					(x == door.X+1 && y == door.Y) ||
					(x == door.X && y == door.Y-1) ||
					(x == door.X-1 && y == door.Y-1) ||
					(x == door.X+1 && y == door.Y-1) {
					tile.Type = doorType
					tile.DoorI = c.DoorI
				} else {
					ToType(tile, cave.Wall, true, true)
				}
				tile.NeverChange = true
				tile.Bomb = false
				tile.IsChanged = true
			}
		}
	}
	c.DoorI++
}

func ExitDoor(c *cave.Cave, door world.Coords, exitI int, doorType cave.BlockType) {
	for y := door.Y - 1; y <= door.Y+1; y++ {
		for x := door.X - 1; x <= door.X+1; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				if x == door.X && y == door.Y {
					ExitTile(c, door, exitI)
					tile.Type = doorType
					tile.DoorI = c.DoorI
				} else if (x == door.X-1 && y == door.Y) ||
					(x == door.X+1 && y == door.Y) ||
					(x == door.X && y == door.Y-1) ||
					(x == door.X-1 && y == door.Y-1) ||
					(x == door.X+1 && y == door.Y-1) {
					tile.Type = doorType
					tile.DoorI = c.DoorI
				} else {
					ToType(tile, cave.Wall, true, true)
				}
				tile.NeverChange = true
				tile.Bomb = false
				tile.IsChanged = true
			}
		}
	}
	c.DoorI++
}

func SecretExit(c *cave.Cave, door world.Coords, exitI int, exitBiome string) {
	ToBiomeCircle(c, door, exitBiome, 7, 2.5)
	Exit(c, door, 7, 3, 1, exitI, cave.SecretDoor)
}

func ExitTile(c *cave.Cave, door world.Coords, exitI int) {
	tile := c.GetTileInt(door.X, door.Y)
	tile.Exit = true
	tile.ExitI = exitI
	c.Exits = append(c.Exits, struct {
		Coords world.Coords
		PopUp  *menus.PopUp
		ExitI  int
		Type   cave.BlockType
	}{
		Coords: door,
		ExitI: exitI,
	})
}