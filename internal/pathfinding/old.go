package pathfinding

import (
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/pkg/util"
	"math"
)

func OldNeighbors(tile *cave.Tile) []*cave.Tile {
	var neighbors []*cave.Tile
	if tile != nil {
		currRule := tile.Chunk.Cave.PathRule
		for _, n := range tile.RCoords.Neighbors() {
			t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
			if t != nil {
				if cave.TypeLegal(t.Type, currRule.LegalTypes) && (!t.Bomb || currRule.Bomb) {
					// if the tile is orthogonal OR
					// if both tiles + orthogonal neighbors are non-solid or the same type
					var t1, t2 *cave.Tile
					if tile.RCoords.X != t.RCoords.X && tile.RCoords.Y != t.RCoords.Y {
						t1 = tile.Chunk.Cave.GetTileInt(tile.RCoords.X, t.RCoords.Y)
						t2 = tile.Chunk.Cave.GetTileInt(t.RCoords.X, tile.RCoords.Y)
					}
					if (t.RCoords.X == tile.RCoords.X || t.RCoords.Y == tile.RCoords.Y) ||
						(!tile.Solid() && !t.Solid() && currRule.Fly && t1 != nil && !t1.Solid() && t2 != nil && !t2.Solid()) ||
						(tile.Type == t.Type && t1 != nil && t1.Type == tile.Type && t2 != nil && t2.Type == tile.Type) {
						// climbing check:
						// if the tile to check is directly above the current tile
						// and flying isn't allowed, and we are in empty tiles
						// a climbing check is required
						if t.RCoords.X == tile.RCoords.X &&
							t.RCoords.Y > tile.RCoords.Y &&
							!currRule.Fly && !t.Solid() && !tile.Solid() &&
							currRule.Climb {
							rc := tile.RCoords
							rc.X += 1
							lc := tile.RCoords
							lc.X -= 1
							rt := tile.Chunk.Cave.GetTileInt(rc.X, rc.Y)
							lt := tile.Chunk.Cave.GetTileInt(lc.X, lc.Y)
							if !rt.Solid() && !lt.Solid() {
								// no climb
								continue
							}
							// walkable check:
							// if the tile below is empty, and the tile below
							// the tile we are checking is empty,
							// and we can't fly, don't allow it
						} else if !currRule.Fly {
							bc := tile.RCoords
							bc2 := t.RCoords
							bc.Y -= 1
							bc2.Y -= 1
							bt := tile.Chunk.Cave.GetTileInt(bc.X, bc.Y)
							bt2 := tile.Chunk.Cave.GetTileInt(bc2.X, bc2.Y)
							if !bt.Solid() && !bt2.Solid() {
								// no walk
								continue
							}
						}
						neighbors = append(neighbors, t)
					}
				}
			}
		}
	}
	return neighbors
}

func OldCost(tile, to *cave.Tile) float64 {
	w := math.Max(tile.Chunk.Cave.PathRule.LegalTypes[tile.Type], tile.Chunk.Cave.PathRule.LegalTypes[to.Type])
	if to.RCoords.X != tile.RCoords.X && to.RCoords.Y != tile.RCoords.Y {
		w *= 1.5
	}
	avoid := tile.Chunk.Cave.PathRule.Avoid
	if avoid != nil {
		w += float64(util.Max(0., 4-util.Abs(avoid.X-to.RCoords.X))) * 2.
		w += float64(util.Max(0., 4-util.Abs(avoid.Y-to.RCoords.Y))) * 2.
	}
	if to.Type == tile.Type || (!to.Solid() && !tile.Solid()) {
		return w
	} else if to.Solid() && tile.Solid() {
		return w + float64(util.Abs(int(to.Type-tile.Type)))
	} else {
		return w + 3.
	}
}
