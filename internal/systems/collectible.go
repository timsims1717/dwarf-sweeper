package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/transform"
	"math"
)

func CollectSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCollectible) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		coll, okC := result.Components[myecs.Collect].(*data.Collectible)
		collider, okC1 := result.Components[myecs.Collision].(*data.Collider)
		if okT && okC && okC1 {
			dist := camera.Cam.Pos.Sub(tran.Pos)
			if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
				if descent.Descent.GetPlayer() != nil &&
					!descent.Descent.GetPlayer().Health.Dazed &&
					!descent.Descent.GetPlayer().Health.Dead {
					if math.Abs(descent.Descent.GetPlayer().Transform.Pos.X-tran.Pos.X) < (descent.Descent.GetPlayer().Collider.Hitbox.W() + collider.Hitbox.W()) * 0.5 &&
						math.Abs(descent.Descent.GetPlayer().Transform.Pos.Y-tran.Pos.Y) < (descent.Descent.GetPlayer().Collider.Hitbox.H() + collider.Hitbox.H()) * 0.5 {
						pickUp := data.GameInput.Get("interact").JustPressed()
						if coll.AutoCollect || pickUp {
							if coll.OnCollect(tran.Pos) {
								coll.Collected = true
								myecs.Manager.DisposeEntity(result.Entity)
								if pickUp {
									data.GameInput.Get("interact").Consume()
								}
							}
						}
					}
				}
			}
		}
	}
}