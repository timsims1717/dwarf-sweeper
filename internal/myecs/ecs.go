package myecs

import "github.com/bytearena/ecs"

var (
	Manager = ecs.NewManager()

	Physics   = Manager.NewComponent()
	Transform = Manager.NewComponent()
	Collision = Manager.NewComponent()

	HasTransform = ecs.BuildTag(Transform)
	HasPhysics   = ecs.BuildTag(Transform, Physics)
	HasCollision = ecs.BuildTag(Transform, Collision)
)

type Collider struct{}