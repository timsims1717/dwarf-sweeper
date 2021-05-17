package myecs

import (
	"github.com/bytearena/ecs"
)

var (
	Manager = ecs.NewManager()

	Physics   = Manager.NewComponent()
	Transform = Manager.NewComponent()
	Collision = Manager.NewComponent()
	Animation = Manager.NewComponent()
	Collect   = Manager.NewComponent()

	HasTransform  = ecs.BuildTag(Transform)
	HasPhysics    = ecs.BuildTag(Transform, Physics)
	HasCollision  = ecs.BuildTag(Transform, Physics, Collision)
	HasAnimation  = ecs.BuildTag(Animation)
	IsCollectible = ecs.BuildTag(Transform, Collect)
)

type Collider struct{}
type Collectible struct{
	CollectedBy bool
}