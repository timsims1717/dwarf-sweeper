package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/transform"
	"math"
)

func EntitySystem() {
	for _, result := range myecs.Manager.Query(myecs.IsEntity) {
		e, ok := result.Components[myecs.Entity].(myecs.AnEntity)
		t, okT := result.Components[myecs.Transform].(*transform.Transform)
		if ok && okT {
			dist := camera.Cam.Pos.Sub(t.Pos)
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

func FunctionSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasFunc) {
		fnA := result.Components[myecs.Func]
		if fnT, ok := fnA.(*data.TimerFunc); ok {
			if fnT.Timer.UpdateDone() {
				if fnT.Func() {
					result.Entity.RemoveComponent(myecs.Func)
				} else {
					fnT.Timer.Reset()
				}
			}
		} else if fnF, ok := fnA.(*data.FrameFunc); ok {
			if fnF.Func() {
				result.Entity.RemoveComponent(myecs.Func)
			}
		}
	}
}

func UpdateSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasUpdate) {
		fnA := result.Components[myecs.Update]
		if fnT, ok := fnA.(*data.TimerFunc); ok {
			if fnT.Timer.UpdateDone() {
				fnT.Func()
				fnT.Timer.Reset()
			}
		} else if fnF, ok := fnA.(*data.FrameFunc); ok {
			fnF.Func()
		}
	}
}