package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
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
				coll, okC := tResult.Components[myecs.Collision].(*data.Collider)
				hp, okH1 := tResult.Components[myecs.Health].(*data.Health)
				_, okH2 := tResult.Components[myecs.Health].(*data.SimpleHealth)
				if okT && okC && (okH1 || okH2) && tran.ID != area.SourceID && tran.Load {
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
						tarHB := coll.Hitbox.Moved(tran.Pos).Moved(pixel.V(coll.Hitbox.W()*-0.5, coll.Hitbox.H()*-0.5))
						if area.Radius > 0. {
							c := pixel.C(area.Center, area.Radius)
							v := c.IntersectRect(tarHB)
							hit = v.X != 0 || v.Y != 0
						} else if area.Rect.W() > 0. && area.Rect.H() > 0. {
							dmgHB := area.Rect.Moved(area.Center).Moved(pixel.V(area.Rect.W()*-0.5, area.Rect.H()*-0.5))
							hit = dmgHB.Intersects(tarHB)
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
		myecs.Manager.DisposeEntity(result.Entity)
	}
}

func DamageSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasDamage) {
		hp, okH := result.Components[myecs.Health].(*data.Health)
		hpB, okB := result.Components[myecs.Health].(*data.SimpleHealth)
		phys, okP := result.Components[myecs.Physics].(*physics.Physics)
		trans, okT := result.Components[myecs.Transform].(*transform.Transform)
		dmg, okD := result.Components[myecs.Damage].(*data.Damage)
		if okH && okP && okD && okT && trans.ID != dmg.SourceID && trans.Load {
			if !hp.Inv {
				immune := data.Immunity{}
				for t, immunity := range hp.Immune {
					if t == dmg.Type {
						immune = immunity
						break
					}
				}
				dmgAmt := dmg.Amount
				if hp.TempInvTimer.Done() {
					if dmgAmt > 0 && !immune.DMG {
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
							d := trans.Pos.Sub(dmg.Source)
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
						dazeTime := dmg.Dazed
						if hp.DazedTime > 0. {
							dazeTime = hp.DazedTime
						}
						hp.DazedTimer = timing.New(dazeTime)
						if hp.DazedEntity != nil {
							hp.DazedEntity.AddComponent(myecs.Temp, timing.New(hp.DazedTime))
						} else {
							dazeTran := transform.New()
							dazeTran.Offset.Y += world.TileSize * 0.5 + 1.
							anim := particles.DazedAnimation()
							hp.DazedEntity = myecs.Manager.NewEntity()
							hp.DazedEntity.AddComponent(myecs.Temp, timing.New(dazeTime)).
								AddComponent(myecs.Transform, dazeTran).
								AddComponent(myecs.Parent, trans).
								AddComponent(myecs.Animation, anim).
								AddComponent(myecs.Drawable, anim).
								AddComponent(myecs.Batch, constants.ParticleKey)
						}
					}
				}
			}
		} else if okB && okD && okP && okT && trans.Load {
			immune := data.Immunity{}
			for t, immunity := range hpB.Immune {
				if t == dmg.Type {
					immune = immunity
					break
				}
			}
			dmgAmt := dmg.Amount
			if (dmgAmt > 0 || (hpB.DigMe && dmg.Type == data.Shovel)) && !immune.DMG {
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
					d := trans.Pos.Sub(dmg.Source)
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
		result.Entity.RemoveComponent(myecs.Damage)
	}
}
