package systems

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/transform"
	"math"
)

func InteractSystem() {
	for _, result := range myecs.Manager.Query(myecs.CanInteract) {
		trans, okT := result.Components[myecs.Transform].(*transform.Transform)
		inter, okI := result.Components[myecs.Interact].(*data.Interact)
		if okT && okI {
			if descent.Descent.GetPlayer() != nil &&
				!descent.Descent.GetPlayer().Health.Dazed &&
				!descent.Descent.GetPlayer().Health.Dead {
				if math.Abs(descent.Descent.GetPlayer().Transform.Pos.X-trans.Pos.X) < inter.Distance &&
					math.Abs(descent.Descent.GetPlayer().Transform.Pos.Y-trans.Pos.Y) < inter.Distance {
					if data.GameInput.Get("interact").JustPressed() {
						data.GameInput.Get("interact").Consume()
						if inter.OnInteract(trans.APos) {
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