package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"math"
)

func CollectSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCollectible) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		coll, okC := result.Components[myecs.Collect].(*data.Collectible)
		if okT && okC {
			dist := camera.Cam.Pos.Sub(tran.Pos)
			if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
				if descent.Descent.GetPlayer() != nil &&
					!descent.Descent.GetPlayer().Health.Dazed &&
					!descent.Descent.GetPlayer().Health.Dead &&
					math.Abs(descent.Descent.GetPlayer().Transform.Pos.X-tran.Pos.X) < world.TileSize &&
					math.Abs(descent.Descent.GetPlayer().Transform.Pos.Y-tran.Pos.Y) < world.TileSize {
					if coll.OnCollect(tran.Pos) {
						coll.Collected = true
						myecs.Manager.DisposeEntity(result.Entity)
					}
				}
			}
		}
	}
}