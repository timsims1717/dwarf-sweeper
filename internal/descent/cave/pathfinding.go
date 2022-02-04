package cave

import (
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/beefsack/go-astar"
	"math"
)

type PathRule struct {
	Origin     world.Coords
	Avoid      *world.Coords
	LegalTypes map[BlockType]float64
	Climb      bool
	Ceiling    bool
	Fly        bool
	Bomb       bool
}

func MakePathRule(orig world.Coords, legalTypes map[BlockType]float64, climb, ceiling, fly, bomb bool) PathRule {
	return PathRule{
		Origin:     orig,
		Avoid:      nil,
		LegalTypes: legalTypes,
		Climb:      climb,
		Ceiling:    ceiling,
		Fly:        fly,
		Bomb:       bomb,
	}
}

func MakeAvoidPathRule(orig, avoid world.Coords, legalTypes map[BlockType]float64, climb, ceiling, fly, bomb bool) PathRule {
	return PathRule{
		Origin:     orig,
		Avoid:      &avoid,
		LegalTypes: legalTypes,
		Climb:      climb,
		Ceiling:    ceiling,
		Fly:        fly,
		Bomb:       bomb,
	}
}

var (
	AllTypes = map[BlockType]float64{
		Deco:          1.,
		Empty:         1.,
		BlockCollapse: 1.,
		BlockDig:      1.,
		BlockBlast:    1.,
		Wall:          1.,
	}
	AllButWallTypes = map[BlockType]float64{
		Deco:          1.,
		Empty:         1.,
		BlockCollapse: 1.,
		BlockDig:      1.,
		BlockBlast:    1.,
	}
	EmptyTypes = map[BlockType]float64{
		Deco:  1.,
		Empty: 1.,
	}
	SolidTypes = map[BlockType]float64{
		BlockCollapse: 1.,
		BlockDig:      1.,
		BlockBlast:    1.,
		Wall:          1.,
	}
	NonWallTypes = map[BlockType]float64{
		BlockCollapse: 1.,
		BlockDig:      1.,
		BlockBlast:    1.,
	}
)

// CoordsIn returns true if Coords c are in the list.
func TypeLegal(t BlockType, list map[BlockType]float64) bool {
	if f, ok := list[t]; ok && f > 0. {
		return true
	}
	return false
}

// PathNeighbors is part of the astar implementation and returns legal
// moves to the tile's neighbors
func (tile *Tile) PathNeighbors() []astar.Pather {
	var neighbors []astar.Pather
	if tile != nil {
		currRule := tile.Chunk.Cave.PathRule
		for _, n := range tile.RCoords.Neighbors() {
			t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
			if t != nil {
				if TypeLegal(t.Type, currRule.LegalTypes) && (!t.Bomb || currRule.Bomb) {
					// if the tile is orthogonal OR
					// if both tiles + orthogonal neighbors are non-solid or the same type
					var t1, t2 *Tile
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

func (tile *Tile) PathNeighborCost(to astar.Pather) float64 {
	t := to.(*Tile)
	w := math.Max(tile.Chunk.Cave.PathRule.LegalTypes[tile.Type], tile.Chunk.Cave.PathRule.LegalTypes[t.Type])
	if t.RCoords.X != tile.RCoords.X && t.RCoords.Y != tile.RCoords.Y {
		w *= 1.5
	}
	avoid := tile.Chunk.Cave.PathRule.Avoid
	if avoid != nil {
		w += float64(util.Max(0., 4-util.Abs(avoid.X-t.RCoords.X))) * 2.
		w += float64(util.Max(0., 4-util.Abs(avoid.Y-t.RCoords.Y))) * 2.
	}
	if t.Type == tile.Type || (!t.Solid() && !tile.Solid()) {
		return w
	} else if t.Solid() && tile.Solid() {
		return w + float64(util.Abs(int(t.Type-tile.Type)))
	} else {
		return w + 3.
	}
}

func (tile *Tile) PathEstimatedCost(to astar.Pather) float64 {
	return util.Magnitude(tile.Transform.Pos.Sub(tile.Chunk.Cave.GetTileInt(tile.Chunk.Cave.PathRule.Origin.X, tile.Chunk.Cave.PathRule.Origin.Y).Transform.Pos))
}
