package systems

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/transform"
	"math"
)

func CollectSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCollectible) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		coll, okC := result.Components[myecs.Collect].(*descent.Collectible)
		collider, okC1 := result.Components[myecs.Collision].(*data.Collider)
		if okT && okC && okC1 && tran.Load {
			coll.Timer.Update()
			for _, d := range descent.Descent.GetPlayers() {
				if !d.Health.Dazed && !d.Health.Dead {
					if math.Abs(d.Transform.Pos.X-tran.Pos.X) < (d.Collider.Hitbox.W()+collider.Hitbox.W())*0.5 &&
						math.Abs(d.Transform.Pos.Y-tran.Pos.Y) < (d.Collider.Hitbox.H()+collider.Hitbox.H())*0.5 {
						if coll.Timer.Done() && coll.OnCollect(tran.Pos, d) {
							coll.Collected = true
							myecs.Manager.DisposeEntity(result.Entity)
						}
					}
				}
			}
		}
	}
}
