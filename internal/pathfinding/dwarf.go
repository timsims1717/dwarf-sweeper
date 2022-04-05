package pathfinding

import "dwarf-sweeper/internal/descent/cave"

func DigNeighbors(tile *cave.Tile) []*cave.Tile {
	var neighbors []*cave.Tile
	if tile != nil {
		for _, t := range tile.Neighbors() {
			if !t.Solid() || t.Diggable() {
				// if the tile is diagonally connected, we have to check if one of tile and t's
				// shared neighbors is legal
				if tile.RCoords.X != t.RCoords.X && tile.RCoords.Y != t.RCoords.Y {
					t1 := tile.Chunk.Cave.GetTileInt(tile.RCoords.X, t.RCoords.Y)
					t2 := tile.Chunk.Cave.GetTileInt(t.RCoords.X, tile.RCoords.Y)
					if !t1.Solid() || t1.Diggable() || !t2.Solid() || t2.Diggable() {
						neighbors = append(neighbors, t)
					}
				} else {
					neighbors = append(neighbors, t)
				}
			}
		}
	}
	return neighbors
}

func DigCost(tile, to *cave.Tile) float64 {
	w := 1.
	if to.Type == cave.Collapse {
		w++
	} else if to.Type == cave.Dig {
		w += 2.
	}
	if to.RCoords.X != tile.RCoords.X && to.RCoords.Y != tile.RCoords.Y {
		w *= 1.5
	}
	return w
}