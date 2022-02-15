package pathfinding

import "dwarf-sweeper/internal/descent/cave"

func init() {
	cave.NeighborsFn = DigNeighbors
	cave.CostFn = DigCost
}