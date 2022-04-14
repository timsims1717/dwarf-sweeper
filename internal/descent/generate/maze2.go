package generate

import (
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
)

func TightMazeCave(newCave *cave.Cave, signal chan bool) {
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
	m := util.Max(newCave.Width, newCave.Height) * 4 / 3
	for i := 0; i < m; i++ {
		CellAutoB3S1234(newCave)
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
	newCave.MarkAsNotChanged()
	if signal != nil {
		signal <- false
		if !<-signal {
			return
		}
	}
	// remove all tiny (<5) groups and non-grouped tiles
	for g, grp := range groups {
		if grp.count < 6 {
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
	newCave.MarkAsNotChanged()
	// expand groups to connect
	count = CountDigTiles(newCave)
	gm := random.CaveGen.Intn(3) + 1
	for len(groups) > gm && float64(len(groups[newCave.MainGroup].coords)) / float64(count) < 0.95 {
		index := random.CaveGen.Intn(len(groups))
		key := 0
		ip := 0
		for k := range groups {
			if ip == index {
				key = k
			}
			ip++
		}
		if key != newCave.MainGroup {
			grp := groups[key]
			grpTile := newCave.GetTileInt(grp.orig.X, grp.orig.Y)
			for _, tc := range grp.coords {
				ns := tc.Neighbors()
				c := 0
				jn := -1
				for j, n := range ns {
					if j % 2 == 0 {
						nt := newCave.GetTileInt(n.X, n.Y)
						if nt != nil && nt.Diggable() {
							c++
							jn = j
						}
					}
				}
				if c == 1 {
					for j, n := range ns {
						if j != jn && (j - jn + 8) % 8 == 4 {
							jt := newCave.GetTileInt(n.X, n.Y)
							structures.ToBlock(jt, false, false)
						}
					}
				}
			}
			coords, packed := GroupTile(newCave, grpTile, key)
			for _, p := range packed {
				delete(groups, p)
			}
			groups[key] = Group{
				count:  len(coords),
				orig:   grp.orig,
				coords: coords,
			}
			newCave.MainGroup = key
			newCave.MarkAsNotChanged()
			if signal != nil {
				signal <- false
				if !<-signal {
					return
				}
			}
		}
	}
	newCave.MapFn(func(tile *cave.Tile) {
		if tile.Group != 0 && tile.Diggable() {
			newCave.Rooms = append(newCave.Rooms, tile.RCoords)
		}
	})
	newCave.MarkAsNotChanged()
}

func CellAutoB3S1234(newCave *cave.Cave) {
	newCave.MapFn(func(tile *cave.Tile) {
		if !tile.NeverChange {
			c := 0
			for _, n := range tile.RCoords.Neighbors() {
				t := newCave.GetTileInt(n.X, n.Y)
				if t != nil && t.Solid() && t.Diggable() {
					c++
				}
			}
			if !tile.Diggable() {
				switch c {
				case 3:
					tile.Change = true
				}
			} else {
				switch c {
				case 0,5,6,7,8:
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