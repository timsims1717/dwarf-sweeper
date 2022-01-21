package util

import (
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
	"math/rand"
)

func RandomPosition(orig pixel.Vec, varX, varY float64, rando *rand.Rand) pixel.Vec {
	pos := orig
	xVar := (rando.Float64() - 0.5) * varX
	yVar := (rando.Float64() - 0.5) * varY
	pos.X += xVar
	pos.Y += yVar
	return pos
}

func RandomVelocity(angle, angleVar, force, forceVar float64, rando *rand.Rand) pixel.Vec {
	v := pixel.V(1, 0).Rotated(angle + (rando.Float64()-0.5)*angleVar)
	return v.Scaled(force + (rando.Float64()-0.5)*forceVar)
}

func RandomPosAndVel(orig pixel.Vec, varX, varY, angle, angleVar, force, forceVar float64, rando *rand.Rand) (*physics.Physics, *transform.Transform) {
	tran := transform.NewTransform()
	phys := physics.New()
	tran.Pos = RandomPosition(orig, varX, varY, rando)
	phys.Velocity = RandomVelocity(angle, angleVar, force, forceVar, rando)
	phys.Velocity.X -= orig.X - tran.Pos.X
	phys.Velocity.Y -= orig.Y - tran.Pos.Y
	return phys, tran
}
