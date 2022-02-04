package structures

import "dwarf-sweeper/internal/descent/cave"

func ToEmpty(tile *cave.Tile, perm, asDeco, override bool) {
	if tile != nil && !tile.NeverChange && (!tile.IsChanged || override) {
		if asDeco {
			tile.Type = cave.Deco
		} else {
			tile.Type = cave.Empty
		}
		tile.IsChanged = true
		tile.NeverChange = perm
	}
}

func ToBlock(tile *cave.Tile, tt cave.BlockType, perm, override bool) {
	if tile != nil && !tile.NeverChange && (!tile.IsChanged || override) && (tt <= cave.BlockDig || !tile.Path) {
		if tt == cave.Unknown {
			if tile.Perlin < 0 {
				tile.Type = cave.BlockCollapse
			} else {
				tile.Type = cave.BlockDig
			}
		} else {
			tile.Type = tt
		}
		tile.IsChanged = true
		tile.NeverChange = perm
	}
}

func BlockUp(tile *cave.Tile, tt cave.BlockType) {
	for _, n := range tile.SubCoords.Neighbors() {
		t := tile.Chunk.Get(n)
		if t != nil && !t.NeverChange && !t.IsChanged && !t.Path {
			if tt == cave.Unknown {
				if tile.Perlin < 0 {
					tile.Type = cave.BlockCollapse
				} else {
					tile.Type = cave.BlockDig
				}
			} else {
				tile.Type = tt
			}
		}
	}
}

func WallUp(tile *cave.Tile) {
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

func WallUpWidth(tile *cave.Tile, width int, left bool) {
	WallUp(tile)
	ns := tile.SubCoords.Neighbors()
	if width == 3 || (width == 2 && left) {
		z := tile.Chunk.Get(ns[4])
		WallUp(z)
		y := tile.Chunk.Get(ns[5])
		WallUp(y)
		x := tile.Chunk.Get(ns[6])
		WallUp(x)
		w := tile.Chunk.Get(ns[7])
		WallUp(w)
	}
	if width == 3 || (width == 2 && !left) {
		v := tile.Chunk.Get(ns[0])
		WallUp(v)
		u := tile.Chunk.Get(ns[1])
		WallUp(u)
		t := tile.Chunk.Get(ns[2])
		WallUp(t)
		s := tile.Chunk.Get(ns[3])
		WallUp(s)
	}
}
