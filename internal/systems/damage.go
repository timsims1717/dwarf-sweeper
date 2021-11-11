package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"math"
)

func AreaDamageSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasAreaDamage) {
		area, ok := result.Components[myecs.AreaDmg].(*data.AreaDamage)
		if ok {
			dist := camera.Cam.Pos.Sub(area.Center)
			if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
				if debug.Debug {
					if area.Radius > 0. {
						col := colornames.White
						debug.AddCircle(col, area.Center, area.Radius, 0.5)
					} else if area.Rect.W() > 0. && area.Rect.H() > 0. {
						col := colornames.White
						debug.AddRect(col, area.Center, area.Rect, 0.5)
					}
				}
				for _, tResult := range myecs.Manager.Query(myecs.HasHealth) {
					tran, okT := tResult.Components[myecs.Transform].(*transform.Transform)
					hp, okH1 := tResult.Components[myecs.Health].(*data.Health)
					_, okH2 := tResult.Components[myecs.Health].(*data.SimpleHealth)
					if okT && (okH1 || okH2) && tran.ID != area.SourceID {
						immune := false
						if okH1 {
							for t, i := range hp.Immune {
								if t == area.Type {
									immune = i.KB && i.Dazed && i.DMG
									break
								}
							}
						}
						if !immune {
							hit := false
							if area.Radius > 0. {
								xt := area.Center.X - tran.Pos.X
								yt := area.Center.Y - tran.Pos.Y
								d2 := xt*xt + yt*yt
								hit = d2 < area.Radius*area.Radius
							} else if area.Rect.W() > 0. && area.Rect.H() > 0. {
								hit = area.Rect.Moved(area.Center).Moved(pixel.V(area.Rect.W()*-0.5, area.Rect.H()*-0.5)).Contains(tran.Pos)
							}
							if hit {
								kb := area.Knockback
								if area.KnockbackDecay {
									p := tran.Pos.Sub(area.Center)
									mag := math.Sqrt(p.X*p.X + p.Y*p.Y)
									kb = kb - (mag / world.TileSize)
								}
								tResult.Entity.AddComponent(myecs.Damage, &data.Damage{
									Amount:    area.Amount,
									Dazed:     area.Dazed,
									Knockback: kb,
									Angle:     area.Angle,
									Source:    area.Center,
									Type:      area.Type,
								})
							}
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
		hp, okH := result.Components[myecs.Health].(*data.Health)
		hpB, okB := result.Components[myecs.Health].(*data.SimpleHealth)
		phys, okP := result.Components[myecs.Physics].(*physics.Physics)
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		dmg, okD := result.Components[myecs.Damage].(*data.Damage)
		if okH && okP && okD && okT && tran.ID != dmg.SourceID {
			if !hp.Inv {
				dist := camera.Cam.Pos.Sub(tran.Pos)
				if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
					immune := data.Immunity{}
					for t, immunity := range hp.Immune {
						if t == dmg.Type {
							immune = immunity
							break
						}
					}
					dmgAmt := dmg.Amount
					if dmgAmt > 0 && hp.TempInvTimer.Done() && !immune.DMG {
						tmp := hp.TempHP
						hp.TempHP -= dmgAmt
						if hp.TempHP < 0 {
							hp.TempHP = 0
						}
						dmgAmt -= tmp
						if dmgAmt > 0 {
							hp.Curr -= dmgAmt
							if hp.Curr <= 0 {
								hp.Curr = 0
								hp.Dead = true
							}
						}
						if hp.TempInvSec > 0. {
							hp.TempInvTimer = timing.New(hp.TempInvSec)
							if hp.TempInvSec > 1.0 {
								myecs.AddEffect(result.Entity, data.NewBlink(hp.TempInvSec))
							}
						}
					}
					if dmg.Knockback > 0.0 && !immune.KB {
						if phys.GravityOff {
							phys.GravityOff = false
						}
						if phys.FrictionOff {
							phys.FrictionOff = false
						}
						phys.CancelMovement()
						var dir pixel.Vec
						if dmg.Angle == nil {
							d := tran.Pos.Sub(dmg.Source)
							d.Y += 1.
							dir = util.Normalize(d)
						} else {
							dir = pixel.V(1., 0.).Rotated(*dmg.Angle)
						}
						phys.SetVelX(dir.X*dmg.Knockback*world.TileSize, 0.)
						phys.SetVelY(dir.Y*dmg.Knockback*world.TileSize, 0.)
						phys.RagDollX = true
						if dir.Y < 0 && math.Abs(dir.X) < 4. && dmg.Knockback > 20. {
							phys.RagDollY = true
						}
					}
					if dmg.Dazed > 0.0 && hp.Curr > 0 && !immune.Dazed {
						hp.Dazed = true
						if hp.DazedTime > 0. {
							hp.DazedTimer = timing.New(hp.DazedTime)
						} else {
							hp.DazedTimer = timing.New(dmg.Dazed)
						}
						if hp.DazedVFX != nil {
							hp.DazedVFX.Animation.Done = true
							hp.DazedVFX = nil
						}
					}
				}
			}
		} else if okB && okD && okP && okT {
			dist := camera.Cam.Pos.Sub(tran.Pos)
			if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
				immune := data.Immunity{}
				for t, immunity := range hpB.Immune {
					if t == dmg.Type {
						immune = immunity
						break
					}
				}
				dmgAmt := dmg.Amount
				if dmgAmt > 0 && !immune.DMG {
					hpB.Dead = true
				}
				if dmg.Knockback > 0.0 && !immune.KB {
					if phys.GravityOff {
						phys.GravityOff = false
					}
					if phys.FrictionOff {
						phys.FrictionOff = false
					}
					phys.CancelMovement()
					var dir pixel.Vec
					if dmg.Angle == nil {
						d := tran.Pos.Sub(dmg.Source)
						d.Y += 1.
						dir = util.Normalize(d)
					} else {
						dir = pixel.V(1., 0.).Rotated(*dmg.Angle)
					}
					phys.SetVelX(dir.X*dmg.Knockback*world.TileSize, 0.)
					phys.SetVelY(dir.Y*dmg.Knockback*world.TileSize, 0.)
					phys.RagDollX = true
					phys.RagDollY = true
				}
			}
		}
		result.Entity.RemoveComponent(myecs.Damage)
	}
}