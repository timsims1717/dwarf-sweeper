package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

var XRaySec = 16.

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
			n := Descent.GetPlayerTile().SubCoords.Neighbors()
			a := world.Combine(n, n[0].Neighbors())
			a = world.Combine(a, n[2].Neighbors())
			a = world.Combine(a, n[4].Neighbors())
			a = world.Combine(a, n[6].Neighbors())
			for _, c := range a {
				tile := Descent.GetPlayerTile().Chunk.Get(c)
				if tile != nil && tile.Breakable() && tile.Solid() && tile.XRay != nil && !util.IsNil(tile.Entity) {
					myecs.Manager.NewEntity().
						AddComponent(myecs.Sprite, tile.XRay).
						AddComponent(myecs.Transform, transform.NewTransform()).
						AddComponent(myecs.Parent, tile.Transform).
						AddComponent(myecs.Batch, constants.EntityKey).
						AddComponent(myecs.Temp, timing.New(0.25))
				}
			}
		}
	}
}

func (x *XRayHelmet) Create(_ pixel.Vec) {
	x.Timer = timing.New(XRaySec)
	x.created = true
	ring := img.Batchers[constants.ParticleKey].Sprites["x-ray-ring"]
	t1 := transform.NewTransform()
	t1.Offset = pixel.V(-world.TileSize, world.TileSize)
	e1 := myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t1).
		AddComponent(myecs.Parent, Descent.Player.Transform).
		AddComponent(myecs.Sprite, ring).
		AddComponent(myecs.Batch, constants.ParticleKey)
	t2 := transform.NewTransform()
	t2.Offset = pixel.V(world.TileSize, world.TileSize)
	t2.Flip = true
	e2 := myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t2).
		AddComponent(myecs.Parent, Descent.Player.Transform).
		AddComponent(myecs.Sprite, ring).
		AddComponent(myecs.Batch, constants.ParticleKey)
	t3 := transform.NewTransform()
	t3.Offset = pixel.V(world.TileSize, -world.TileSize)
	t3.Flip = true
	t3.Flop = true
	e3 := myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t3).
		AddComponent(myecs.Parent, Descent.Player.Transform).
		AddComponent(myecs.Sprite, ring).
		AddComponent(myecs.Batch, constants.ParticleKey)
	t4 := transform.NewTransform()
	t4.Offset = pixel.V(-world.TileSize, -world.TileSize)
	t4.Flop = true
	e4 := myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t4).
		AddComponent(myecs.Parent, Descent.Player.Transform).
		AddComponent(myecs.Sprite, ring).
		AddComponent(myecs.Batch, constants.ParticleKey)
	x.entities = [4]*ecs.Entity{e1, e2, e3, e4}
	x.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, Descent.Player.Transform).
		AddComponent(myecs.Entity, x)
}

func (x *XRayHelmet) Delete() {
	for _, e := range x.entities {
		myecs.Manager.DisposeEntity(e)
	}
	myecs.Manager.DisposeEntity(x.entity)
}
