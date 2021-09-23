package dungeon

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

//type XRayItem struct {
//	Physics   *physics.Physics
//	Transform *transform.Transform
//	created   bool
//	collect   *data.Collectible
//	sprite    *pixel.Sprite
//	entity    *ecs.Entity
//	health    *data.Health
//}
//
//func (x *XRayItem) Update() {
//	if x.created {
//		if x.collect.CollectedBy {
//			if AddToInventory(&InvItem{
//				Name:   "x-ray",
//				Sprite: x.sprite,
//				OnUse:  func() bool {
//					//if Dungeon.Player.Bubble == nil {
//						xray := &XRayHelmet{}
//						xray.Create(pixel.Vec{})
//						return true
//					//} else {
//					//	return false
//					//}
//				},
//				Count:  1,
//				Unique: true,
//			}) {
//				// todo: effects
//				x.Delete()
//			} else {
//				x.collect.CollectedBy = false
//			}
//		} else if x.health.Dead {
//			x.Delete()
//		}
//	}
//}
//
//func (x *XRayItem) Create(pos pixel.Vec) {
//	x.Physics, x.Transform = util.RandomVelocity(pos, 1.0, random.Effects)
//	x.Transform.Pos = pos
//	x.created = true
//	x.sprite = img.Batchers[entityKey].Sprites["x-ray-helmet"]
//	x.collect = &data.Collectible{}
//	x.health = &data.Health{
//		Max:        1,
//		Curr:       1,
//		Override:   true,
//	}
//	x.entity = myecs.Manager.NewEntity().
//		AddComponent(myecs.Entity, x).
//		AddComponent(myecs.Transform, x.Transform).
//		AddComponent(myecs.Physics, x.Physics).
//		AddComponent(myecs.Collision, data.Collider{ GroundOnly: true }).
//		AddComponent(myecs.Collect, x.collect).
//		AddComponent(myecs.Health, x.health).
//		AddComponent(myecs.Sprite, x.sprite).
//		AddComponent(myecs.Batch, entityKey)
//}
//
//func (x *XRayItem) Delete() {
//	x.health.Delete()
//	myecs.Manager.DisposeEntity(x.entity)
//}

type XRayHelmet struct {
	sprite   pixel.Sprite
	entity   *ecs.Entity
	entities [4]*ecs.Entity
	Timer    *timing.FrameTimer
	created  bool
}

func (x *XRayHelmet) Update() {
	if x.created {
		if x.Timer.UpdateDone() {
			x.Delete()
		} else {
			n := Dungeon.GetPlayerTile().SubCoords.Neighbors()
			a := world.Combine(n, n[0].Neighbors())
			a = world.Combine(a, n[2].Neighbors())
			a = world.Combine(a, n[4].Neighbors())
			a = world.Combine(a, n[6].Neighbors())
			for _, c := range a {
				tile := Dungeon.GetPlayerTile().Chunk.Get(c)
				if tile.breakable && tile.Solid && tile.XRay != nil && !util.IsNil(tile.Entity) {
					myecs.Manager.NewEntity().
						AddComponent(myecs.Sprite, tile.XRay).
						AddComponent(myecs.Transform, transform.NewTransform()).
						AddComponent(myecs.Parent, tile.Transform).
						AddComponent(myecs.Batch, cfg.EntityKey).
						AddComponent(myecs.Temp, timing.New(0.25))
				}
			}
		}
	}
}

func (x *XRayHelmet) Create(_ pixel.Vec) {
	x.Timer = timing.New(16.)
	x.created = true
	t1 := transform.NewTransform()
	t1.Offset = pixel.V(-world.TileSize, world.TileSize)
	e1 := myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t1).
		AddComponent(myecs.Parent, Dungeon.Player.Transform).
		AddComponent(myecs.Sprite, img.Batchers[cfg.BigEntityKey].Sprites["x-ray-ring"]).
		AddComponent(myecs.Batch, cfg.BigEntityKey)
	t2 := transform.NewTransform()
	t2.Offset = pixel.V(world.TileSize, world.TileSize)
	t2.Flip = true
	e2 := myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t2).
		AddComponent(myecs.Parent, Dungeon.Player.Transform).
		AddComponent(myecs.Sprite, img.Batchers[cfg.BigEntityKey].Sprites["x-ray-ring"]).
		AddComponent(myecs.Batch, cfg.BigEntityKey)
	t3 := transform.NewTransform()
	t3.Offset = pixel.V(world.TileSize, -world.TileSize)
	t3.Flip = true
	t3.Flop = true
	e3 := myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t3).
		AddComponent(myecs.Parent, Dungeon.Player.Transform).
		AddComponent(myecs.Sprite, img.Batchers[cfg.BigEntityKey].Sprites["x-ray-ring"]).
		AddComponent(myecs.Batch, cfg.BigEntityKey)
	t4 := transform.NewTransform()
	t4.Offset = pixel.V(-world.TileSize, -world.TileSize)
	t4.Flop = true
	e4 := myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t4).
		AddComponent(myecs.Parent, Dungeon.Player.Transform).
		AddComponent(myecs.Sprite, img.Batchers[cfg.BigEntityKey].Sprites["x-ray-ring"]).
		AddComponent(myecs.Batch, cfg.BigEntityKey)
	x.entities = [4]*ecs.Entity{e1, e2, e3, e4}
	x.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, x)
}

func (x *XRayHelmet) Delete() {
	for _, e := range x.entities {
		myecs.Manager.DisposeEntity(e)
	}
	myecs.Manager.DisposeEntity(x.entity)
}