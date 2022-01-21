package systems

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/timing"
)

func TemporarySystem() {
	for _, result := range myecs.Manager.Query(myecs.IsTemp) {
		temp := result.Components[myecs.Temp]
		if timer, ok := temp.(*timing.FrameTimer); ok {
			if timer.UpdateDone() {
				myecs.Manager.DisposeEntity(result.Entity)
			}
		} else if check, ok := temp.(myecs.ClearFlag); ok {
			if check {
				myecs.Manager.DisposeEntity(result.Entity)
			}
		}
	}
}

func ClearSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsTemp) {
		myecs.Manager.DisposeEntity(result.Entity)
	}
}
