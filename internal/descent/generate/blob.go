package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
)

func BlobCave(newCave *cave.Cave, signal chan bool) {
	if signal != nil {
		signal <- false
	}

	// using cellular automata, generate the cave
	// rule is Born if 678, Survive if 345678
	RandomizeTiles(newCave)
	if signal != nil {
		signal <- false
		if !<-signal {
			return
		}
	}
	for i := 0; i < 5; i++ {
		CellAutoB678S345678(newCave)
		if signal != nil {
			signal <- false
			if !<-signal {
				return
			}
		}
	}
	// group until the majority (90%+) of tiles are grouped
	newCave.MarkAsNotChanged()
	count := CountDigTiles(newCave)
	type Group struct{
		count  int
		orig   world.Coords
		coords []world.Coords
	}
	groups := make(map[int]Group)
	largestGroupCount := 0
	group := 1
	grouped := 0
	tile := GetRandomUngroupedTile(newCave)
	for tile != nil && float64(grouped) / float64(count) < 0.95 {
		coords := GroupTile(newCave, tile, group)
		gc := len(coords)
		//fmt.Printf("Group %d (%d,%d) has %d tiles in it\n", group, tile.RCoords.X, tile.RCoords.Y, gn)
		if largestGroupCount < gc {
			largestGroupCount = gc
			newCave.MainGroup = group
		}
		groups[group] = Group{
			count:  gc,
			orig:   tile.RCoords,
			coords: coords,
		}
		grouped += gc
		group++
		tile = GetRandomUngroupedTile(newCave)
	}
	if signal != nil {
		signal <- false
		if !<-signal {
			return
		}
	}
	// remove all tiny (<51) groups and non-grouped tiles
	for g, grp := range groups {
		if grp.count < 51 {
			newCave.MapFn(func(tile *cave.Tile) {
				if tile.Group == g {
					tile.Group = 0
					tile.Type = cave.Wall
				}
			})
			delete(groups, g)
		}
	}
	if signal != nil {
		signal <- false
		if !<-signal {
			return
		}
	}
	newCave.MapFn(func(tile *cave.Tile) {
		if tile.Group == 0 && tile.Type != cave.Wall {
			tile.Group = 0
			tile.Type = cave.Wall
		}
	})
	if signal != nil {
		signal <- false
		if !<-signal {
			return
		}
	}
	// connect groups with >50 tiles to the largest group
	for g, grp := range groups {
		if g != newCave.MainGroup {
			newGroup := -1
			for _, c := range grp.coords {
				t := newCave.GetTileInt(c.X, c.Y)
				ns := c.Neighbors()
				ng := -1
				for _, n := range ns {
					nt := newCave.GetTileInt(n.X, n.Y)
					if nt.Group != t.Group && nt.Group != 0 {
						ng = nt.Group
						break
					}
				}
				if ng != g && ng != -1 {
					newGroup = ng
					for _, n := range ns {
						nt := newCave.GetTileInt(n.X, n.Y)
						if nt != nil {
							if nt.Type == cave.Wall {
								structures.ToBlock(nt, cave.Unknown, false, false)
							}
							nt.Group = newGroup
							groups[newGroup] = Group{
								count:  groups[newGroup].count,
								orig:   groups[newGroup].orig,
								coords: append(groups[newGroup].coords, nt.RCoords),
							}
						}
					}
				}
			}
			if newGroup != g && newGroup != -1 {
				for _, c := range grp.coords {
					t := newCave.GetTileInt(c.X, c.Y)
					t.Group = newGroup
				}
				groups[newGroup] = Group{
					count:  groups[newGroup].count,
					orig:   groups[newGroup].orig,
					coords: append(groups[newGroup].coords, grp.coords...),
				}
				delete(groups, g)
				if signal != nil {
					signal <- false
					if !<-signal {
						return
					}
				}
			} else {
				// need to connect directly to large group

			}
		}
	}
	newCave.MapFn(func(tile *cave.Tile) {
		if tile.Group != 0 && tile.Type != cave.Wall {
			newCave.Rooms = append(newCave.Rooms, tile.RCoords)
		}
	})
	newCave.MarkAsNotChanged()
}

func CellAutoB678S345678(newCave *cave.Cave) {
	newCave.MapFn(func(tile *cave.Tile) {
		if !tile.NeverChange {
			c := 0
			for _, n := range tile.RCoords.Neighbors() {
				t := newCave.GetTileInt(n.X, n.Y)
				if t != nil && t.Solid() && t.Type != cave.Wall {
					c++
				}
			}
			if tile.Type == cave.Wall {
				switch c {
				case 6, 7, 8:
					tile.Change = true
				}
			} else {
				if c < 3 {
					tile.Change = true
				}
			}
		}
	})
	newCave.MapFn(func(tile *cave.Tile) {
		if !tile.NeverChange && tile.Change {
			if tile.Type == cave.Wall {
				structures.ToBlock(tile, cave.Unknown, false, false)
			} else {
				tile.Type = cave.Wall
			}
			tile.Change = false
		}
	})
}

func RandomizeTiles(newCave *cave.Cave) {
	// randomize all tiles, weighted away from the edges
	newCave.MapFn(func(tile *cave.Tile) {
		if !tile.NeverChange {
			var w int
			chCoords := tile.Chunk.Coords
			if chCoords.Y == 0 || chCoords.Y == newCave.Bottom ||
				chCoords.X == newCave.Left || chCoords.X == newCave.Right {
				if chCoords.Y == 0 {
					w = 6 - tile.SubCoords.Y
				} else if chCoords.Y == newCave.Bottom {
					w = tile.SubCoords.Y - (constants.ChunkSize - 7)
				}
				if chCoords.X == newCave.Left {
					w = util.Max(6 - tile.SubCoords.X, w)
				} else if chCoords.X == newCave.Right {
					w = util.Max(tile.SubCoords.X - (constants.ChunkSize - 7), w)
				}
			}
			if w < 0 {
				w = 0
			}
			r := random.CaveGen.Intn(9)
			if r < 5-w {
				structures.ToBlock(tile, cave.Unknown, false, false)
			}
		}
	})
}

func CountDigTiles(newCave *cave.Cave) int {
	c := 0
	newCave.MapFn(func(tile *cave.Tile) {
		if tile.Diggable() {
			c++
		}
	})
	return c
}

func GetRandomUngroupedTile(newCave *cave.Cave) *cave.Tile {
	tries := 0
	for tries < 100 {
		x := random.Effects.Intn(constants.ChunkSize * (newCave.Right - newCave.Left + 1))
		y := random.Effects.Intn(constants.ChunkSize * (newCave.Bottom + 1))
		tile := newCave.GetTileInt(x, y)
		if tile.Type != cave.Wall && tile.Group == 0 && !tile.Change {
			return tile
		}
		tries++
	}
	return nil
}

func GroupTile(newCave *cave.Cave, tile *cave.Tile, group int) ([]world.Coords) {
	tile.Change = true
	qu := []*cave.Tile{
		tile,
	}
	var list []world.Coords
	for len(qu) > 0 {
		next := qu[0]
		if next.Group != group {
			list = append(list, next.RCoords)
			next.Group = group
			for i, n := range next.RCoords.Neighbors() {
				if i % 2 == 0 {
					t := newCave.GetTileInt(n.X, n.Y)
					if t.Type != cave.Wall && !t.Change && t.Group != group {
						t.Change = true
						qu = append(qu, t)
					}
				}
			}
		}
		qu = qu[1:]
	}
	return list
}