package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
)

func MazeCave(newCave *cave.Cave, signal chan bool) {
	if signal != nil {
		signal <- false
	}

	// using cellular automata, generate the cave
	// rule is Born if 2, Survive if 123
	PrepareCaveMaze(newCave)
	if signal != nil {
		signal <- false
		if !<-signal {
			return
		}
	}
	n := util.Max(newCave.Width, newCave.Height) / 3
	for i := 0; i < n; i++ {
		CellAutoB2S123(newCave)
		//if signal != nil {
		//	signal <- false
		//	if !<-signal {
		//		return
		//	}
		//}
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
		coords, _ := GroupTile(newCave, tile, group)
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
					tile.Type = cave.Blast
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
		if tile.Group == 0 && tile.Diggable() {
			tile.Group = 0
			tile.Type = cave.Blast
		}
	})
	if signal != nil {
		signal <- false
		if !<-signal {
			return
		}
	}
	// connect groups with >50 tiles to the largest group
	//for g, grp := range groups {
	//	if g != newCave.MainGroup {
	//		newGroup := -1
	//		for _, c := range grp.coords {
	//			t := newCave.GetTileInt(c.X, c.Y)
	//			ns := c.Neighbors()
	//			ng := -1
	//			for _, n := range ns {
	//				nt := newCave.GetTileInt(n.X, n.Y)
	//				if nt.Group != t.Group && nt.Group != 0 {
	//					ng = nt.Group
	//					break
	//				}
	//			}
	//			if ng != g && ng != -1 {
	//				newGroup = ng
	//				for _, n := range ns {
	//					nt := newCave.GetTileInt(n.X, n.Y)
	//					if nt != nil {
	//						if !nt.Diggable() {
	//							structures.ToBlock(nt, false, false)
	//						}
	//						nt.Group = newGroup
	//						groups[newGroup] = Group{
	//							count:  groups[newGroup].count,
	//							orig:   groups[newGroup].orig,
	//							coords: append(groups[newGroup].coords, nt.RCoords),
	//						}
	//					}
	//				}
	//			}
	//		}
	//		if newGroup != g && newGroup != -1 {
	//			for _, c := range grp.coords {
	//				t := newCave.GetTileInt(c.X, c.Y)
	//				t.Group = newGroup
	//			}
	//			groups[newGroup] = Group{
	//				count:  groups[newGroup].count,
	//				orig:   groups[newGroup].orig,
	//				coords: append(groups[newGroup].coords, grp.coords...),
	//			}
	//			delete(groups, g)
	//			if signal != nil {
	//				signal <- false
	//				if !<-signal {
	//					return
	//				}
	//			}
	//		}
	//	}
	//}
	newCave.MapFn(func(tile *cave.Tile) {
		if tile.Group != 0 && tile.Diggable() {
			newCave.Rooms = append(newCave.Rooms, tile.RCoords)
		}
	})
	newCave.MarkAsNotChanged()
}

func CellAutoB2S123(newCave *cave.Cave) {
	newCave.MapFn(func(tile *cave.Tile) {
		if !tile.NeverChange {
			c := 0
			for _, n := range tile.RCoords.Neighbors() {
				t := newCave.GetTileInt(n.X, n.Y)
				if t != nil && !(t.Solid() && t.Diggable()) {
					c++
				}
			}
			if tile.Diggable() {
				switch c {
				case 2:
					tile.Change = true
				}
			} else {
				switch c {
				case 0,4,5,6,7,8:
					tile.Change = true
				}
			}
		}
	})
	newCave.MapFn(func(tile *cave.Tile) {
		if !tile.NeverChange && tile.Change {
			if !tile.Diggable() {
				structures.ToBlock(tile, false, false)
				tile.IsChanged = false
			} else {
				tile.Type = cave.Blast
			}
			tile.Change = false
		}
	})
}

func PrepareCaveMaze(newCave *cave.Cave) {
	cX, cY := newCave.CenterCoords()
	// randomize all tiles, weighted away from the edges
	newCave.MapFn(func(tile *cave.Tile) {
		if !tile.NeverChange {
			w := util.Max(util.Abs(tile.RCoords.X - cX), util.Abs(tile.RCoords.Y - cY))
			if w < 6 {
				if random.CaveGen.Intn(2) == 0 {
					structures.ToBlock(tile, false, false)
				}
			} else {
				w = 0
				chCoords := tile.Chunk.Coords
				if chCoords.Y == 0 || chCoords.Y == newCave.Bottom ||
					chCoords.X == newCave.Left || chCoords.X == newCave.Right {
					if chCoords.Y == 0 {
						w = 8 - tile.SubCoords.Y
					} else if chCoords.Y == newCave.Bottom {
						w = tile.SubCoords.Y - (constants.ChunkSize - 9)
					}
					if chCoords.X == newCave.Left {
						w = util.Max(8 - tile.SubCoords.X, w)
					} else if chCoords.X == newCave.Right {
						w = util.Max(tile.SubCoords.X - (constants.ChunkSize - 9), w)
					}
					if w < 0 {
						w = 0
					}
					r := random.CaveGen.Intn(7)
					if r < w {
						tile.NeverChange = true
					}
				}
			}
		}
	})
}