package systems

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/timing"
)

func ManagementSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsTemp) {
		if timer, ok := result.Components[myecs.Temp].(*timing.FrameTimer); ok {
			if timer.UpdateDone() {
				myecs.Manager.DisposeEntity(result.Entity)
			}
		} else if check, ok := result.Components[myecs.Temp].(bool); ok {
			if check {
				myecs.Manager.DisposeEntity(result.Entity)
			}
		} else {
			myecs.Manager.DisposeEntity(result.Entity)
		}
	}
}