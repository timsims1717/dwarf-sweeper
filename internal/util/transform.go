package util

import (
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"math/rand"
)

func RandomPosition(orig pixel.Vec, variance float64, rando *rand.Rand) pixel.Vec {
	pos := orig
	actVar := variance * world.TileSize
	xVar := (rando.Float64() - 0.5) * actVar
	yVar := (rando.Float64() - 0.5) * actVar
	pos.X += xVar
	pos.Y += yVar
	return pos
}


func RandomVelocity(orig pixel.Vec, variance float64, rando *rand.Rand) *physics.Physics {
	tran := transform.NewTransform()
	physicsT := &physics.Physics{Transform: tran}
	physicsT.Pos = orig
	actVar := variance * world.TileSize
	//if square {
	xVar := (rando.Float64() - 0.5) * actVar
	yVar := (rando.Float64() - 0.5) * actVar
	physicsT.Pos.X += xVar
	physicsT.Pos.Y += yVar
	physicsT.Velocity.X = xVar * 5.
	physicsT.Velocity.Y = 20.
	//}
	return physicsT
}
