package generate

import "dwarf-sweeper/internal/descent/cave"

func toEmpty(tile *cave.Tile, perm, blank bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Solid = false
		if blank {
			tile.Type = cave.Deco
		} else {
			tile.Type = cave.Empty
		}
		tile.Breakable = true
		tile.IsChanged = true
		tile.NeverChange = perm
		tile.UpdateSprites()
	}
}

func toBlock(tile *cave.Tile, perm, bomb bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Solid = true
		tile.Type = cave.Block
		tile.Breakable = true
		tile.Bomb = bomb
		tile.IsChanged = true
		tile.NeverChange = perm
		tile.UpdateSprites()
	}
}

func toWall(tile *cave.Tile, perm bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Solid = true
		tile.Type = cave.Wall
		tile.Breakable = false
		tile.Bomb = false
		tile.IsChanged = true
		tile.NeverChange = perm
		tile.UpdateSprites()
	}
}

func wallUp(tile *cave.Tile, noBomb bool) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Solid = true
		tile.Type = cave.Block
		tile.Breakable = true
		tile.Fillable = true
		if noBomb {
			tile.Bomb = false
		}
		tile.IsChanged = true
		tile.UpdateSprites()
		for _, n := range tile.SubCoords.Neighbors() {
			t := tile.Chunk.Get(n)
			if t != nil && !t.NeverChange && !t.IsChanged {
				t.Solid = true
				t.Type = cave.Wall
				t.Breakable = false
				t.UpdateSprites()
			}
		}
	}
}