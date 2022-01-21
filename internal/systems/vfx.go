package systems

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
)

func VFXSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasVFX) {
		vfx, okV := result.Components[myecs.VFX].(*data.VFX)
		trans, okT := result.Components[myecs.Transform].(*transform.Transform)
		if okV && okT {
			var remove []int
			for i, effect := range vfx.Effects {
				if blink, ok := effect.(*data.Blink); ok {
					if blink.BlinkTimer.UpdateDone() {
						blink.Blink = !blink.Blink
						if blink.Blink {
							blink.BlinkTimer = timing.New(data.BlinkSec)
						} else {
							blink.BlinkTimer = timing.New(data.ShowSec)
						}
						trans.Hide = blink.Blink
					}
					if blink.Timer.UpdateDone() {
						trans.Hide = false
						remove = append(remove, i)
					}
				} else if fade, ok := effect.(*data.FadeOut); ok {
					allFin := true
					if fade.InterR != nil {
						f, fin := fade.InterR.Update(timing.DT)
						trans.Mask.R = uint8(f)
						if fin {
							fade.InterR = nil
						} else {
							allFin = false
						}
					}
					if fade.InterG != nil {
						f, fin := fade.InterG.Update(timing.DT)
						trans.Mask.G = uint8(f)
						if fin {
							fade.InterG = nil
						} else {
							allFin = false
						}
					}
					if fade.InterB != nil {
						f, fin := fade.InterB.Update(timing.DT)
						trans.Mask.B = uint8(f)
						if fin {
							fade.InterB = nil
						} else {
							allFin = false
						}
					}
					if fade.InterA != nil {
						f, fin := fade.InterA.Update(timing.DT)
						trans.Mask.A = uint8(f)
						if fin {
							fade.InterA = nil
						} else {
							allFin = false
						}
					}
					if allFin {
						remove = append(remove, i)
					}
				}
			}
			for i := len(remove) - 1; i >= 0; i-- {
				if len(vfx.Effects) > 1 {
					vfx.Effects = append(vfx.Effects[:remove[i]], vfx.Effects[remove[i]+1:]...)
				} else {
					result.Entity.RemoveComponent(myecs.VFX)
				}
			}
		}
	}
}
