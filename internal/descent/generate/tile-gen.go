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
		tile.Fillable = false
	}
}

func toBlockCollapse(tile *cave.Tile, perm, bomb bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Type = cave.BlockCollapse
		tile.Bomb = bomb
		tile.IsChanged = true
		tile.NeverChange = perm
		tile.Fillable = true
	}
}

func toBlockDig(tile *cave.Tile, perm, bomb bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Type = cave.BlockDig
		tile.Bomb = bomb
		tile.IsChanged = true
		tile.NeverChange = perm
		tile.Fillable = true
	}
}

func toBlockBlast(tile *cave.Tile, perm, bomb bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Type = cave.BlockBlast
		tile.Bomb = bomb
		tile.IsChanged = true
		tile.NeverChange = perm
		tile.Fillable = true
	}
}

func toWall(tile *cave.Tile, perm bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Type = cave.Wall
		tile.Bomb = false
		tile.IsChanged = true
		tile.NeverChange = perm
		tile.Fillable = false
	}
}

func wallUp(tile *cave.Tile, noBomb bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Type = cave.BlockCollapse
		tile.Fillable = true
		if noBomb {
			tile.Bomb = false
		}
		tile.IsChanged = true
		for _, n := range tile.SubCoords.Neighbors() {
			t := tile.Chunk.Get(n)
			if t != nil && !t.NeverChange && !t.IsChanged {
				t.Type = cave.Wall
			}
		}
	}
}