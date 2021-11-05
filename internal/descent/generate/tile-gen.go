package generate

import "dwarf-sweeper/internal/descent/cave"

func toEmpty(tile *cave.Tile, perm, asDeco bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		if asDeco {
			tile.Type = cave.Deco
		} else {
			tile.Type = cave.Empty
		}
		tile.IsChanged = true
		tile.NeverChange = perm
	}
}

func toBlockCollapse(tile *cave.Tile, perm bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Type = cave.BlockCollapse
		tile.IsChanged = true
		tile.NeverChange = perm
	}
}

func toBlockDig(tile *cave.Tile, perm bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Type = cave.BlockDig
		tile.IsChanged = true
		tile.NeverChange = perm
	}
}

func toBlockBlast(tile *cave.Tile, perm bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged && !tile.Path {
		tile.Type = cave.BlockBlast
		tile.IsChanged = true
		tile.NeverChange = perm
	}
}

func toWall(tile *cave.Tile, perm bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged && !tile.Path {
		tile.Type = cave.Wall
		tile.IsChanged = true
		tile.NeverChange = perm
	}
}

func wallUp(tile *cave.Tile) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Type = cave.BlockCollapse
		tile.IsChanged = true
		for _, n := range tile.SubCoords.Neighbors() {
			t := tile.Chunk.Get(n)
			if t != nil && !t.NeverChange && !t.IsChanged && !t.Path {
				t.Type = cave.Wall
			}
		}
	}
}

func wallUpWidth(tile *cave.Tile, width int, left bool) {
	wallUp(tile)
	ns := tile.SubCoords.Neighbors()
	if width == 3 || (width == 2 && left) {
		z := tile.Chunk.Get(ns[4])
		wallUp(z)
		y := tile.Chunk.Get(ns[5])
		wallUp(y)
		x := tile.Chunk.Get(ns[6])
		wallUp(x)
		w := tile.Chunk.Get(ns[7])
		wallUp(w)
	}
	if width == 3 || (width == 2 && !left) {
		v := tile.Chunk.Get(ns[0])
		wallUp(v)
		u := tile.Chunk.Get(ns[1])
		wallUp(u)
		t := tile.Chunk.Get(ns[2])
		wallUp(t)
		s := tile.Chunk.Get(ns[3])
		wallUp(s)
	}
}