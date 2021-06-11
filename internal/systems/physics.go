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
		if okT && okP {
			tran.LastPos = tran.Pos
			tran.Pos.X += timing.DT * phys.Velocity.X
			tran.Pos.Y += timing.DT * phys.Velocity.Y
			if !phys.GravityOff && !phys.YJustSet {
				if phys.Velocity.Y > -500. {
					phys.Velocity.Y -= 750. * timing.DT
				}
			}
			phys.YJustSet = false
			if !phys.FrictionOff && !phys.XJustSet {
				friction := 25.
				if phys.Grounded {
					friction = 400.
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