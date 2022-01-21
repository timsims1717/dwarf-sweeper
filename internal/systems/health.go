package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"github.com/faiface/pixel"
	"math"
)

func HealthSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasHealth) {
		hp, okH := result.Components[myecs.Health].(*data.Health)
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		if okH && okT {
			dist := camera.Cam.Pos.Sub(tran.Pos)
			if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
				if hp.Dazed {
					if hp.DazedTimer.UpdateDone() {
						hp.Dazed = false
					} else {
						if hp.DazedVFX != nil {
							hp.DazedVFX.Matrix = pixel.IM.Moved(tran.APos).Moved(pixel.V(0., 9.))
						} else if hp.DazedVFX == nil {
							hp.DazedVFX = vfx.CreateDazed(tran.APos.Add(pixel.V(0., 9.)))
						}
					}
				}
				if !hp.Dazed && hp.DazedVFX != nil {
					hp.DazedVFX.Animation.Done = true
					hp.DazedVFX = nil
				}
				hp.TempInvTimer.Update()
				if hp.Curr < 1 {
					hp.Dead = true
				}
				if hp.TempHP > 0 && hp.TempHPTimer.UpdateDone() {
					hp.TempHP = 0
				}
				if hp.Curr > hp.Max {
					hp.Curr = hp.Max
				}
			}
		}
	}
}

func HealingSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasHealing) {
		hp, ok1 := result.Components[myecs.Health].(*data.Health)
		heal, ok2 := result.Components[myecs.Healing].(*data.Heal)
		if ok1 && ok2 {
			hp.Curr += heal.Amount
			if hp.Curr > hp.Max {
				hp.Curr = hp.Max
			}
			if heal.TmpAmount > 0 {
				hp.TempHP = util.Max(heal.TmpAmount, hp.TempHP)
				hp.TempHPTimer = timing.New(8.)
			}
		}
		result.Entity.RemoveComponent(myecs.Healing)
	}
}
