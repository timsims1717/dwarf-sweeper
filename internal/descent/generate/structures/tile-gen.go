package structures

import (
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
)

func ToType(tile *cave.Tile, tt cave.BlockType, perm, override bool) {
	if tile != nil && !tile.NeverChange && (!tile.IsChanged || override) {
		tile.Type = tt
		tile.IsChanged = true
		tile.NeverChange = perm
		//if !tile.Solid() {
		//	tile.Bomb = false
		//}
	}
}

func ToBlock(tile *cave.Tile, perm, override bool) {
	if tile != nil && !tile.NeverChange && (!tile.IsChanged || override) {
		if tile.Perlin < 0 {
			tile.Type = cave.Collapse
		} else {
			tile.Type = cave.Dig
		}
		tile.IsChanged = true
		tile.NeverChange = perm
	}
}

func BlockUp(tile *cave.Tile, tt cave.BlockType) {
	for _, n := range tile.RCoords.Neighbors() {
		t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
		if t != nil && !t.NeverChange && !t.IsChanged && !t.Path {
			if tt == cave.Unknown {
				if t.Perlin < 0 {
					t.Type = cave.Collapse
				} else {
					t.Type = cave.Dig
				}
			} else {
				t.Type = tt
			}
		}
	}
}

func WallUp(tile *cave.Tile) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Type = cave.Collapse
		tile.IsChanged = true
		for _, n := range tile.RCoords.Neighbors() {
			t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
			if t != nil && !t.NeverChange && !t.IsChanged && !t.Path {
				t.Type = cave.Wall
			}
		}
	}
}

func WallUpWidth(tile *cave.Tile, width int, left bool) {
	WallUp(tile)
	ns := tile.RCoords.Neighbors()
	if width == 3 || (width == 2 && left) {
		z := tile.Chunk.Cave.GetTileInt(ns[4].X, ns[4].Y)
		WallUp(z)
		y := tile.Chunk.Cave.GetTileInt(ns[5].X, ns[5].Y)
		WallUp(y)
		x := tile.Chunk.Cave.GetTileInt(ns[6].X, ns[6].Y)
		WallUp(x)
		w := tile.Chunk.Cave.GetTileInt(ns[7].X, ns[7].Y)
		WallUp(w)
	}
	if width == 3 || (width == 2 && !left) {
		v := tile.Chunk.Cave.GetTileInt(ns[0].X, ns[0].Y)
		WallUp(v)
		u := tile.Chunk.Cave.GetTileInt(ns[1].X, ns[1].Y)
		WallUp(u)
		t := tile.Chunk.Cave.GetTileInt(ns[2].X, ns[2].Y)
		WallUp(t)
		s := tile.Chunk.Cave.GetTileInt(ns[3].X, ns[3].Y)
		WallUp(s)
	}
}

func ToBiomeCircle(c *cave.Cave, center world.Coords, biome string, radius int, variance float64) {
	cenTile := c.GetTileInt(center.X, center.Y)
	cPos := cenTile.Transform.Pos
	fRad := float64(radius) * world.TileSize
	for y := center.Y - radius; y < center.Y+radius; y++ {
		for x := center.X - radius; x < center.X+radius; x++ {
			tile := c.GetTileInt(x, y)
			if tile != nil {
				tPos := tile.Transform.Pos
				dist := util.Magnitude(cPos.Sub(tPos))
				if dist < fRad-(random.CaveGen.Float64()*variance*world.TileSize) {
					ToBiome(tile, biome)
				}
			}
		}
	}
}

func ToBiome(tile *cave.Tile, biome string) {
	if tile != nil && !tile.NeverChange && !tile.IsChanged {
		tile.Biome = biome
	}
}