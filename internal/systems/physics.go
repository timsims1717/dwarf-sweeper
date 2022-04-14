package systems

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
)

func PhysicsSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasPhysics) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		phys, okP := result.Components[myecs.Physics].(*physics.Physics)
		if okT && okP && tran.Load {
			tran.LastPos = tran.Pos
			if (phys.RightBound && phys.Velocity.X > 0.) || (phys.LeftBound && phys.Velocity.X < 0.) {
				phys.Velocity.X = 0.
			}
			if (phys.BottomBound && phys.Velocity.Y < 0.) || (phys.TopBound && phys.Velocity.Y > 0.) {
				phys.Velocity.Y = 0.
			}
			tran.Pos.X += timing.DT * phys.Velocity.X
			tran.Pos.Y += timing.DT * phys.Velocity.Y
			if !phys.GravityOff && !phys.YJustSet && !phys.BottomBound && !phys.Grounded {
				if phys.Velocity.Y > -phys.Terminal {
					phys.Velocity.Y -= phys.Gravity * timing.DT
				}
				if phys.Velocity.Y <= -phys.Terminal {
					phys.Velocity.Y = -phys.Terminal
				}
			}
			phys.YJustSet = false
			if !phys.FrictionOff && !phys.XJustSet {
				friction := phys.AirFriction
				if phys.Grounded {
					friction = phys.Friction
				}
				if phys.Velocity.X > 0. {
					phys.Velocity.X -= friction * timing.DT
					if phys.Velocity.X < 0. {
						phys.Velocity.X = 0
					}
				} else if phys.Velocity.X < 0. {
					phys.Velocity.X += friction * timing.DT
					if phys.Velocity.X > 0. {
						phys.Velocity.X = 0
					}
				}
			}
			phys.XJustSet = false
		}
	}
}
