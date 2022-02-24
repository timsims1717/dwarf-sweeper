package systems

import (
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/transform"
	"math"
)

func InteractSystem() {
	for _, result := range myecs.Manager.Query(myecs.CanInteract) {
		trans, okT := result.Components[myecs.Transform].(*transform.Transform)
		inter, okI := result.Components[myecs.Interact].(*descent.Interact)
		if okT && okI && trans.Load {
			for _, d := range descent.Descent.GetPlayers() {
				if !d.Health.Dazed &&
					!d.Health.Dead {
					if math.Abs(d.Transform.Pos.X-trans.Pos.X) < inter.Distance &&
						math.Abs(d.Transform.Pos.Y-trans.Pos.Y) < inter.Distance {
						if d.Player.Input.Get("interact").JustPressed() {
							d.Player.Input.Get("interact").Consume()
							if inter.OnInteract == nil || inter.OnInteract(trans.APos, d) {
								inter.Interacted = true
								if inter.Remove {
									myecs.Manager.DisposeEntity(result.Entity)
								}
							}
						}
					}
				}
			}
		}
	}
}
