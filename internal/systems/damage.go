package systems

import (
	"dwarf-sweeper/internal/character"
	"dwarf-sweeper/internal/dungeon"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"math"
)

func AreaDamageSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasAreaDamage) {
		area, ok := result.Components[myecs.AreaDmg].(*character.AreaDamage)
		if ok {
			for _, tResult := range myecs.Manager.Query(myecs.HasHealth) {
				tran, okT := tResult.Components[myecs.Transform].(*transform.Transform)
				_, okH := tResult.Components[myecs.Health].(*character.Health)
				if okT && okH {
					for _, a := range area.Area {
						if dungeon.TileInTile(tran.Pos, a) {
							kb := area.Knockback
							if area.KnockbackDecay {
								p := tran.Pos.Sub(area.Source)
								mag := math.Sqrt(p.X*p.X + p.Y*p.Y)
								kb = kb - (mag / world.TileSize)
							}
							tResult.Entity.AddComponent(myecs.Damage, &character.Damage{
								Amount:    area.Amount,
								Dazed:     area.Dazed,
								Knockback: kb,
								Source:    area.Source,
								Override:  area.Override,
							})
						}
					}
				}
			}
		}
		myecs.Manager.DisposeEntity(result.Entity)
	}
}

func DamageSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasDamage) {
		hp, okH := result.Components[myecs.Health].(*character.Health)
		phys, okP := result.Components[myecs.Physics].(*physics.Physics)
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		dmg, okD := result.Components[myecs.Damage].(*character.Damage)
		if okH && okP && okD && okT && !hp.Inv && !hp.TempInv && (!hp.Override || dmg.Override) {
			if dmg.Amount > 0 {
				hp.Curr -= dmg.Amount
				if hp.Curr <= 0 {
					hp.Curr = 0
					hp.Dead = true
				}
			}
			if dmg.Knockback > 0.0 {
				phys.CancelMovement()
				var dir pixel.Vec
				if dmg.Angle == nil {
					d := tran.Pos.Sub(dmg.Source)
					d.Y += 1.
					dir = util.Normalize(d)
				} else {
					dir = pixel.V(1., 0.).Rotated(*dmg.Angle)
					if tran.Pos.X < dmg.Source.X {
						dir.X *= -1
					}
				}
				phys.SetVelX(dir.X * dmg.Knockback * world.TileSize, 0.)
				phys.SetVelY(dir.Y * dmg.Knockback * world.TileSize, 0.)
				phys.RagDoll = true
			}
			if dmg.Dazed > 0.0 && hp.Curr > 0 {
				hp.Dazed = true
				if hp.DazeOverride {
					hp.DazedO = true
				} else {
					hp.DazedTimer = timing.New(dmg.Dazed)
				}
				if hp.DazedVFX != nil {
					hp.DazedVFX.Animation.Done = true
					hp.DazedVFX = nil
				}
			}
			if hp.TempInvSec > 0. {
				hp.TempInv = true
				hp.TempInvTimer = timing.New(hp.TempInvSec)
			}
		}
		result.Entity.RemoveComponent(myecs.Damage)
	}
}