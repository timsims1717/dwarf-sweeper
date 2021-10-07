package myecs

import (
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

var (
	Count = 0
)

var (
	Manager  = ecs.NewManager()
	Unloaded = ecs.NewManager()

	Temp = Manager.NewComponent()

	Animation = Manager.NewComponent()
	Batch     = Manager.NewComponent()
	Sprite    = Manager.NewComponent()
	Entity    = Manager.NewComponent()

	Physics   = Manager.NewComponent()
	Transform = Manager.NewComponent()
	Parent    = Manager.NewComponent()
	Collision = Manager.NewComponent()
	Collect   = Manager.NewComponent()

	Health  = Manager.NewComponent()
	Healing = Manager.NewComponent()
	Damage  = Manager.NewComponent()
	AreaDmg = Manager.NewComponent()

	IsTemp = ecs.BuildTag(Temp)

	HasAnimation   = ecs.BuildTag(Animation, Transform)
	HasAnimDrawing = ecs.BuildTag(Animation, Transform, Batch)
	HasSprDrawing  = ecs.BuildTag(Sprite, Transform, Batch)
	IsEntity       = ecs.BuildTag(Entity, Transform)

	HasTransform  = ecs.BuildTag(Transform)
	HasParent     = ecs.BuildTag(Transform, Parent)
	HasPhysics    = ecs.BuildTag(Transform, Physics)
	HasCollision  = ecs.BuildTag(Transform, Physics, Collision)
	IsCollectible = ecs.BuildTag(Transform, Collect)

	HasAreaDamage = ecs.BuildTag(AreaDmg)
	HasHealing    = ecs.BuildTag(Health, Healing)
	HasHealth     = ecs.BuildTag(Health, Transform)
	HasDamage     = ecs.BuildTag(Health, Physics, Transform, Damage)
)

func Update() {
	Count = 0
	for _, result := range Manager.Query(IsEntity) {
		if _, ok := result.Components[Entity].(AnEntity); ok {
			Count++
		}
	}
}

type AnEntity interface {
	Update()
	Create(pixel.Vec)
	Delete()
}