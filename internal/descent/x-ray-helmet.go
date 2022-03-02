package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

var XRaySec = 16.

func StartXRayVision(d *Dwarf) {
	myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, transform.New()).
		AddComponent(myecs.Temp, timing.New(XRaySec)).
		AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
			n := Descent.Cave.GetTile(d.Transform.Pos).RCoords.Neighbors()
			a := world.Combine(n, n[0].Neighbors())
			a = world.Combine(a, n[2].Neighbors())
			a = world.Combine(a, n[4].Neighbors())
			a = world.Combine(a, n[6].Neighbors())
			for _, c := range a {
				tile := Descent.Cave.GetTileInt(c.X, c.Y)
				if tile != nil && tile.Breakable() && tile.Solid() && tile.XRay != "" {
					myecs.Manager.NewEntity().
						AddComponent(myecs.Drawable, img.Batchers[constants.ParticleKey].GetSprite(tile.XRay)).
						AddComponent(myecs.Transform, transform.New()).
						AddComponent(myecs.Parent, tile.Transform).
						AddComponent(myecs.Batch, constants.ParticleKey).
						AddComponent(myecs.Temp, timing.New(0.1))
				}
			}
			return false
		}))

	ring := img.Batchers[constants.ParticleKey].Sprites["x-ray-ring"]
	t1 := transform.New()
	t1.Offset = pixel.V(-world.TileSize, world.TileSize)
	myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t1).
		AddComponent(myecs.Parent, d.Transform).
		AddComponent(myecs.Temp, timing.New(XRaySec)).
		AddComponent(myecs.Drawable, ring).
		AddComponent(myecs.Batch, constants.ParticleKey)
	t2 := transform.New()
	t2.Offset = pixel.V(world.TileSize, world.TileSize)
	t2.Flip = true
	myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t2).
		AddComponent(myecs.Parent, d.Transform).
		AddComponent(myecs.Temp, timing.New(XRaySec)).
		AddComponent(myecs.Drawable, ring).
		AddComponent(myecs.Batch, constants.ParticleKey)
	t3 := transform.New()
	t3.Offset = pixel.V(world.TileSize, -world.TileSize)
	t3.Flip = true
	t3.Flop = true
	myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t3).
		AddComponent(myecs.Parent, d.Transform).
		AddComponent(myecs.Temp, timing.New(XRaySec)).
		AddComponent(myecs.Drawable, ring).
		AddComponent(myecs.Batch, constants.ParticleKey)
	t4 := transform.New()
	t4.Offset = pixel.V(-world.TileSize, -world.TileSize)
	t4.Flop = true
	myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, t4).
		AddComponent(myecs.Parent, d.Transform).
		AddComponent(myecs.Temp, timing.New(XRaySec)).
		AddComponent(myecs.Drawable, ring).
		AddComponent(myecs.Batch, constants.ParticleKey)
}