package generate

import (
	"dwarf-sweeper/internal/data"
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
			case data.Left:
				curr.X--
			case data.Right:
				curr.X++
			case data.Up:
				curr.Y--
			case data.Down:
				curr.Y++
			}
		}
	}
}

func RandomDirection() data.Direction {
	switch random.CaveGen.Intn(4) {
	case 0:
		return data.Left
	case 1:
		return data.Right
	case 2:
		return data.Up
	default:
		return data.Down
	}
}