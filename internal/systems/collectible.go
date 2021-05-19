package systems

import (
	"dwarf-sweeper/internal/dungeon"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"math"
)

func CollectSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCollectible) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		coll, okC := result.Components[myecs.Collect].(*myecs.Collectible)
		if okT && okC {
			if dungeon.Player1 != nil && !dungeon.Player1.Hurt &&
				math.Abs(dungeon.Player1.Transform.Pos.X - tran.Pos.X) < world.TileSize &&
				math.Abs(dungeon.Player1.Transform.Pos.Y - tran.Pos.Y) < world.TileSize {
				coll.CollectedBy = true
			}
		}
	}
}