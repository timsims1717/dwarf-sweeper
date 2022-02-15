package cave

import (
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/beefsack/go-astar"
)

var (
	Origin      world.Coords
	NeighborsFn func(*Tile) []*Tile
	CostFn      func(*Tile, *Tile) float64
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
	ns := NeighborsFn(tile)
	for _, n := range ns {
		neighbors = append(neighbors, n)
	}
	return neighbors
}

func (tile *Tile) PathNeighborCost(to astar.Pather) float64 {
	return CostFn(tile, to.(*Tile))
}

func (tile *Tile) PathEstimatedCost(to astar.Pather) float64 {
	return util.Magnitude(tile.Transform.Pos.Sub(tile.Chunk.Cave.GetTileInt(tile.Chunk.Cave.PathRule.Origin.X, tile.Chunk.Cave.PathRule.Origin.Y).Transform.Pos))
}

