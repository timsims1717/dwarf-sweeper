package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
	"math"
)

func CollisionSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasCollision) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		coll, okC := result.Components[myecs.Collision].(*data.Collider)
		phys, okP := result.Components[myecs.Physics].(*physics.Physics)
		if okT && okC && okP {
			if !coll.GroundOnly || coll.Damage != nil {
				dist := camera.Cam.Pos.Sub(tran.Pos)
				if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
					var hb pixel.Rect
					if math.Abs(tran.Rot) == 0.5 {
						hb = pixel.R(0., 0., coll.Hitbox.H()*collisionDistance, coll.Hitbox.W()*collisionDistance)
					} else {
						hb = pixel.R(0., 0., coll.Hitbox.W()*collisionDistance, coll.Hitbox.H()*collisionDistance)
					}
					hb = hb.Moved(tran.Pos).Moved(pixel.V(hb.W()*-0.5, hb.H()*-0.5))
					// check for overlap with other collision boxes
					for _, result1 := range myecs.Manager.Query(myecs.HasCollision) {
						tran1, okT1 := result1.Components[myecs.Transform].(*transform.Transform)
						coll1, okC1 := result1.Components[myecs.Collision].(*data.Collider)
						phys1, okP1 := result1.Components[myecs.Physics].(*physics.Physics)
						if okT1 && okC1 && okP1 && tran1.ID != tran.ID {
							dist = camera.Cam.Pos.Sub(tran1.Pos)
							if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
								var hb1 pixel.Rect
								if math.Abs(tran1.Rot) == 0.5 {
									hb1 = pixel.R(0., 0., coll1.Hitbox.H()*collisionDistance, coll1.Hitbox.W()*collisionDistance)
								} else {
									hb1 = pixel.R(0., 0., coll1.Hitbox.W()*collisionDistance, coll1.Hitbox.H()*collisionDistance)
								}
								hb1 = hb1.Moved(tran1.Pos).Moved(pixel.V(hb1.W()*-0.5, hb1.H()*-0.5))
								overlap := hb.Intersects(hb1)
								// check if the two hitboxes overlap
								if overlap {
									// if neither hitbox is 'ground only', coll1 will push coll out of its space
									// if coll can't pass, then it will be stopped/bounced by the otherwise it just is 'pushed'
									if coll.Damage != nil && coll.Damage.SourceID != tran1.ID {
										coll.Collided = true
										coll1.Collided = true
										result1.Entity.AddComponent(myecs.Damage, coll.Damage)
									}
									if coll1.Damage != nil && coll1.Damage.SourceID != tran.ID {
										coll.Collided = true
										coll1.Collided = true
										result.Entity.AddComponent(myecs.Damage, coll1.Damage)
									}
									if !coll1.GroundOnly {
										coll.Collided = true
										coll1.Collided = true
										pX := tran.Pos.X - tran1.Pos.X
										pY := tran.Pos.Y - tran1.Pos.Y
										overlapX := -(math.Abs(pX) - (hb.W()+hb1.W())*0.5)
										overlapY := -(math.Abs(pY) - (hb.H()+hb1.H())*0.5)
										if (overlapX >= 0. && overlapX < overlapY) || (coll.CanPass || coll1.CanPass) {
											if coll.CanPass || coll1.CanPass {
												if pX < 0. {
													phys.Velocity.X -= collisionPush * timing.DT
													phys1.Velocity.X += collisionPush * timing.DT
												} else {
													phys.Velocity.X += collisionPush * timing.DT
													phys1.Velocity.X -= collisionPush * timing.DT
												}
											} else {
												if pX < 0. {
													// if tran is left of tran1
													if phys.LeftBound && !phys1.RightBound {
														tran1.Pos.X += overlapX * 0.5
														BounceX(phys1, true)
													} else if !phys.LeftBound && phys1.RightBound {
														tran.Pos.X -= overlapX * 0.5
														BounceX(phys, false)
													} else if !phys.LeftBound && !phys1.RightBound {
														tran.Pos.X -= overlapX * 0.25
														tran1.Pos.X += overlapX * 0.25
														BounceX(phys1, true)
														BounceX(phys, false)
													}
													coll1.LeftBound = true
													coll.RightBound = true
												} else {
													if phys1.LeftBound && !phys.RightBound {
														tran.Pos.X += overlapX * 0.5
														BounceX(phys, true)
													} else if !phys1.LeftBound && phys.RightBound {
														tran1.Pos.X -= overlapX * 0.5
														BounceX(phys1, false)
													} else if !phys1.LeftBound && !phys.RightBound {
														tran1.Pos.X -= overlapX * 0.25
														tran.Pos.X += overlapX * 0.25
														BounceX(phys, true)
														BounceX(phys1, false)
													}
													coll.LeftBound = true
													coll1.RightBound = true
												}
											}
										} else if overlapY >= 0. && overlapY < overlapX {
											if pY < 0. {
												// if tran is below tran1
												if phys.BottomBound && !phys1.TopBound {
													tran1.Pos.Y += overlapY * 0.5
													BounceY(phys1, true)
												} else if !phys.BottomBound && phys1.TopBound {
													tran.Pos.Y -= overlapY * 0.5
													BounceY(phys, false)
												} else if !phys.BottomBound && !phys1.TopBound {
													tran.Pos.Y -= overlapY * 0.25
													tran1.Pos.Y += overlapY * 0.25
													BounceY(phys1, true)
													BounceY(phys, false)
												}
												coll1.BottomBound = true
												coll.TopBound = true
											} else {
												if phys1.BottomBound && !phys.TopBound {
													tran.Pos.Y += overlapY * 0.49
													BounceY(phys, true)
												} else if !phys1.BottomBound && phys.TopBound {
													tran1.Pos.Y -= overlapY * 0.49
													BounceY(phys1, false)
												} else if !phys1.BottomBound && !phys.TopBound {
													tran1.Pos.Y -= overlapY * 0.25
													tran.Pos.Y += overlapY * 0.25
													BounceY(phys, true)
													BounceY(phys1, false)
												}
												coll.BottomBound = true
												coll1.TopBound = true
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func BounceX(phys *physics.Physics, left bool) bool {
	if (phys.Velocity.X < 0. && left) || (phys.Velocity.X > 0. && !left) {
		if phys.RagDollX && math.Abs(phys.Velocity.X) > BounceThreshold {
			phys.Velocity.X *= -phys.Bounciness
			return true
		} else {
			phys.Velocity.X = 0.
		}
	}
	return false
}

func BounceY(phys *physics.Physics, down bool) bool {
	if (phys.Velocity.Y < 0. && down) || (phys.Velocity.Y > 0. && !down) {
		if phys.RagDollY && math.Abs(phys.Velocity.Y) > BounceThreshold {
			phys.Velocity.Y *= -phys.Bounciness
			return true
		} else {
			phys.Velocity.Y = 0.
		}
	}
	return false
}