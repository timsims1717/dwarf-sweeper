package systems

import (
	"dwarf-sweeper/internal/myecs"
)

func EntitySystem() {
	for _, result := range myecs.Manager.Query(myecs.IsEntity) {
		if e, ok := result.Components[myecs.Entity].(myecs.AnEntity); ok {
			e.Update()
		}
	}
}