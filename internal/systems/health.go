package systems

import (
	"dwarf-sweeper/internal/character"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
)

func HealthSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasHealth) {
		hp, okH := result.Components[myecs.Health].(*character.Health)
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		if okH && okT {
			if hp.Dazed {
				if hp.DazedTimer.UpdateDone() && !hp.DazedO {
					hp.Dazed = false
				}
				if !hp.Dazed {
					hp.DazedVFX.Animation.Done = true
					hp.DazedVFX = nil
				} else if hp.DazedVFX != nil {
					hp.DazedVFX.Matrix = pixel.IM.Moved(tran.APos).Moved(pixel.V(0., 9.))
				} else if hp.DazedVFX == nil {
					hp.DazedVFX = vfx.CreateDazed(tran.APos.Add(pixel.V(0., 9.)))
				}
			}
			if hp.TempInv && (hp.TempInvSec == 0. || hp.TempInvTimer.UpdateDone()) {
				hp.TempInv = false
			}
			if hp.Curr < 1 {
				hp.Dead = true
			}
		}
	}
}