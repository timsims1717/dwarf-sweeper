package myecs

import (
	"dwarf-sweeper/internal/data"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

var (
	Count = 0
)

var (
	Manager  = ecs.NewManager()

	Temp = Manager.NewComponent()
	Func = Manager.NewComponent()

	Animation = Manager.NewComponent()
	Batch     = Manager.NewComponent()
	Sprite    = Manager.NewComponent()
	Entity    = Manager.NewComponent()
	PopUp     = Manager.NewComponent()
	VFX       = Manager.NewComponent()

	Physics   = Manager.NewComponent()
	Transform = Manager.NewComponent()
	Parent    = Manager.NewComponent()
	Collision = Manager.NewComponent()
	Collect   = Manager.NewComponent()
	Interact  = Manager.NewComponent()

	Health  = Manager.NewComponent()
	Healing = Manager.NewComponent()
	Damage  = Manager.NewComponent()
	AreaDmg = Manager.NewComponent()

	IsTemp  = ecs.BuildTag(Temp, Transform)
	HasFunc = ecs.BuildTag(Func)

	HasAnimation   = ecs.BuildTag(Animation, Transform)
	HasAnimDrawing = ecs.BuildTag(Animation, Transform, Batch)
	HasSprDrawing  = ecs.BuildTag(Sprite, Transform, Batch)
	IsEntity       = ecs.BuildTag(Entity, Transform)
	HasPopUp       = ecs.BuildTag(PopUp, Transform)
	HasVFX         = ecs.BuildTag(VFX, Transform)

	HasTransform  = ecs.BuildTag(Transform)
	HasParent     = ecs.BuildTag(Transform, Parent)
	HasPhysics    = ecs.BuildTag(Transform, Physics)
	HasCollision  = ecs.BuildTag(Transform, Physics, Collision)
	IsCollectible = ecs.BuildTag(Transform, Collision, Collect)
	CanInteract   = ecs.BuildTag(Transform, Interact)

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

type ClearFlag bool

func AddEffect(entity *ecs.Entity, effect interface{}) {
	if entity.HasComponent(VFX) {
		if vfxC, ok := entity.GetComponentData(VFX); ok {
			if vfx, ok := vfxC.(*data.VFX); ok {
				vfx.Effects = append(vfx.Effects, effect)
			}
		}
	} else {
		entity.AddComponent(VFX, &data.VFX{Effects: []interface{}{effect}})
	}
}