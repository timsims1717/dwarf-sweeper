package myecs

import (
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

var (
	Manager  = ecs.NewManager()
	toDelete = []ecs.EntityID{}

	Animation = Manager.NewComponent()
	Batch     = Manager.NewComponent()
	Sprite    = Manager.NewComponent()
	Entity    = Manager.NewComponent()

	Physics   = Manager.NewComponent()
	Transform = Manager.NewComponent()
	Collision = Manager.NewComponent()
	Collect   = Manager.NewComponent()

	Health  = Manager.NewComponent()
	Damage  = Manager.NewComponent()
	AreaDmg = Manager.NewComponent()

	HasAnimation   = ecs.BuildTag(Animation)
	HasAnimDrawing = ecs.BuildTag(Animation, Transform, Batch)
	HasSprDrawing  = ecs.BuildTag(Sprite, Transform, Batch)
	IsEntity       = ecs.BuildTag(Entity)

	HasTransform  = ecs.BuildTag(Transform)
	HasPhysics    = ecs.BuildTag(Transform, Physics)
	HasCollision  = ecs.BuildTag(Transform, Physics, Collision)
	IsCollectible = ecs.BuildTag(Transform, Collect)

	HasAreaDamage = ecs.BuildTag(AreaDmg)
	HasHealth     = ecs.BuildTag(Health, Transform)
	HasDamage     = ecs.BuildTag(Health, Physics, Transform, Damage)
)

type Collider struct{
	GroundOnly bool
	CanPass    bool
}
type Collectible struct{
	CollectedBy bool
}
type AnEntity interface {
	Update()
	Create(pixel.Vec)
	Delete()
	SetId(int)
}

//func LazyDelete(e *ecs.Entity) {
//	toDelete = append(toDelete, e.ID)
//}
//
//func Flush() {
//	for _, i := range toDelete {
//		Manager.DisposeEntity(Manager.GetEntityByID(i))
//	}
//	toDelete = []ecs.EntityID{}
//}
//
//func Clear() {
//	Manager = ecs.NewManager()
//}