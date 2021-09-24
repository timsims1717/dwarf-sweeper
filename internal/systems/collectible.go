package systems

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/dungeon"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"math"
)

func CollectSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCollectible) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		coll, okC := result.Components[myecs.Collect].(*data.Collectible)
		if okT && okC {
			if dungeon.Dungeon.GetPlayer() != nil &&
				!dungeon.Dungeon.GetPlayer().Health.Dazed &&
				!dungeon.Dungeon.GetPlayer().Health.Dead &&
				math.Abs(dungeon.Dungeon.GetPlayer().Transform.Pos.X - tran.Pos.X) < world.TileSize &&
				math.Abs(dungeon.Dungeon.GetPlayer().Transform.Pos.Y - tran.Pos.Y) < world.TileSize {
				if coll.OnCollect(tran.Pos) {
					coll.Collected = true
					myecs.Manager.DisposeEntity(result.Entity)
				}
			}
		}
	}
}