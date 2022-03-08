package systems

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
)

func TemporarySystem() {
	for _, result := range myecs.Manager.Query(myecs.IsTemp) {
		temp := result.Components[myecs.Temp]
		trans, okT := result.Components[myecs.Transform].(*transform.Transform)
		if okT {
			if timer, ok := temp.(*timing.Timer); ok {
				if timer.UpdateDone() {
					trans.Hide = true
					trans.Dead = true
					myecs.Manager.DisposeEntity(result.Entity)
				}
			} else if check, ok := temp.(myecs.ClearFlag); ok {
				if check {
					trans.Hide = true
					trans.Dead = true
					myecs.Manager.DisposeEntity(result.Entity)
				}
			}
		}
	}
}

func ClearSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsTemp) {
		trans, ok := result.Components[myecs.Transform].(*transform.Transform)
		if ok {
			trans.Hide = true
			trans.Dead = true
		}
		myecs.Manager.DisposeEntity(result.Entity)
	}
}
