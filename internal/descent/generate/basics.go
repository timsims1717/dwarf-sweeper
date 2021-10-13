package generate

import (
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/world"
)

func CreateChunks(newCave *cave.Cave) {
	for y := 0; y <= newCave.Bottom; y++ {
		for x := newCave.Left; x <= newCave.Right; x++ {
			chunk := cave.NewChunk(world.Coords{X: x, Y: y}, newCave)
			if x >= 0 {
				newCave.RChunks[chunk.Coords] = chunk
			} else {
				newCave.LChunks[chunk.Coords] = chunk
			}
		}
	}
}

func Outline(c *cave.Cave, s world.Coords, fullPath []Path) {
	curr := s
	for _, path := range fullPath {
		for i := 0; i < path.Count; i++ {
			toWall(c.GetTileInt(curr.X, curr.Y), true)
			switch path.Dir {
			case Left:
				curr.X--
			case Right:
				curr.X++
			case Up:
				curr.Y--
			case Down:
				curr.Y++
			}
		}
	}
}

func RandomDirection() Direction {
	switch random.CaveGen.Intn(4) {
	case 0:
		return Left
	case 1:
		return Right
	case 2:
		return Up
	default:
		return Down
	}
}