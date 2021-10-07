package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/transform"
	"math"
)

func EntitySystem() {
	for _, result := range myecs.Manager.Query(myecs.IsEntity) {
		e, ok := result.Components[myecs.Entity].(myecs.AnEntity)
		t, okT := result.Components[myecs.Transform].(*transform.Transform)
		if ok && okT {
			dist := descent.Descent.GetPlayer().Transform.Pos.Sub(t.Pos)
			if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
				e.Update()
			}
		}
	}
}

func DeleteAllEntities() {
	for _, result := range myecs.Manager.Query(myecs.IsEntity) {
		if e, ok := result.Components[myecs.Entity].(myecs.AnEntity); ok {
			e.Delete()
		}
	}
}