package boss

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

func GnomeBoss(c *cave.Cave, level int) *cave.Cave {
	c.FillChunk = structures.FillChunkWall
	c.StartC = world.Coords{X: 16, Y: 9}
	c.BombPMin = 0.2
	c.BombPMax = 0.3
	chunk0 := cave.NewChunk(world.Coords{X: 0, Y: 0}, c, cave.Wall)

	chunkr1 := cave.NewChunk(world.Coords{X: 1, Y: 0}, c, cave.Wall)
	chunkr2 := cave.NewChunk(world.Coords{X: 1, Y: 1}, c, cave.Wall)
	chunkr3 := cave.NewChunk(world.Coords{X: 0, Y: 1}, c, cave.Wall)
	chunkl1 := cave.NewChunk(world.Coords{X: -1, Y: 0}, c, cave.Wall)
	chunkl2 := cave.NewChunk(world.Coords{X: -1, Y: 1}, c, cave.Wall)

	c.Chunks[chunk0.Coords] = chunk0
	c.Chunks[chunkr1.Coords] = chunkr1
	c.Chunks[chunkr2.Coords] = chunkr2
	c.Chunks[chunkr3.Coords] = chunkr3
	c.Chunks[chunkl1.Coords] = chunkl1
	c.Chunks[chunkl2.Coords] = chunkl2
	structures.Entrance(c, c.StartC, 9, 5, 3, false)
	structures.Stairs(c, world.Coords{X: 11, Y: 10}, true, true, 12, 6)
	structures.Stairs(c, world.Coords{X: 21, Y: 10}, false, true, 12, 6)
	structures.GnomeMineLayer(c, world.Coords{X: 11, Y: 23}, world.Coords{X: 21, Y: 23})

	cl := c.StartC
	cl.Y = 23
	descent.Descent.CoordsMap["current_layer"] = cl
	tt := c.StartC
	tt.Y = 22
	trans := transform.New()
	trans.Pos = c.GetTileInt(tt.X, tt.Y).Transform.Pos
	coll := data.NewCollider(pixel.R(0., 0., world.TileSize*70., world.TileSize*3.), true, true)
	coll.Debug = true

	gnome := descent.CreateGnomeBoss(5 + level/4)
	gnome.Charge = false
	gnome.SetOnDamageFn(func() {
		nextLayer(gnome)
	})
	gnome.SetOnFleeFn(func() {
		revealExit(gnome)
	})
	myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Collision, coll).
		AddComponent(myecs.Trigger, data.NewFrameFunc(func() bool {
			TriggerGnome(gnome)
			return true
		}))

	return c
}

func nextLayer(gnome *descent.GnomeBoss) {
	if gnome.State != descent.GBDig {
		cl := descent.Descent.CoordsMap["current_layer"]
		lp := gnome.Transform.Pos
		lp.X -= world.TileSize * 3.
		rp := gnome.Transform.Pos
		rp.X += world.TileSize * 3.
		l := descent.Descent.GetTile(lp)
		r := descent.Descent.GetTile(rp)
		lc := l.RCoords
		rc := r.RCoords
		lc.Y = cl.Y + 10
		rc.Y = cl.Y + 10
		descent.Descent.CoordsMap["current_layer"] = lc
		gnome.State = descent.GBDig
		updated := structures.GnomeMineLayer(descent.Descent.Cave, lc, rc)
		structures.UpdateTiles(updated)
	}
}

func revealExit(gnome *descent.GnomeBoss) {
	c := descent.Descent.Cave
	c.ExitC = c.GetTile(gnome.Transform.Pos).RCoords
	structures.Door(c, c.ExitC, true)
	descent.Descent.SetExitPopup()
	c.UpdateBatch = true
}

func TriggerGnome(gnome *descent.GnomeBoss) {
	gnome.Transform.Pos, _ = descent.EmergeCoords()
	x := gnome.Transform.Pos.X + (descent.Descent.Player.Transform.Pos.X-gnome.Transform.Pos.X)*0.5
	camera.Cam.MoveTo(pixel.V(x, camera.Cam.Pos.Y), 0.4, false)
	descent.Descent.DisableInput = true
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		startRumble(gnome)
		myecs.Manager.DisposeEntity(e)
		return false
	}, 1.25)).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
}

func startRumble(gnome *descent.GnomeBoss) {
	shake := myecs.Manager.NewEntity()
	sfx.SoundPlayer.PlaySound("rockslide", 0.)
	shake.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		camera.Cam.Shake(0.3, 8.)
		return false
	}, 0.2)).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	stopShake := myecs.Manager.NewEntity()
	stopShake.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		myecs.Manager.DisposeEntity(shake)
		myecs.Manager.DisposeEntity(stopShake)
		// todo: pan camera
		return false
	}, 1.5)).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	gnomeFn := myecs.Manager.NewEntity()
	gnomeFn.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		myecs.Manager.DisposeEntity(gnomeFn)
		gnome.Emerge(false)
		return false
	}, 6.0)).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	stop := myecs.Manager.NewEntity()
	stop.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		descent.Descent.DisableInput = false
		gnome.Charge = true
		gnome.State = descent.GBCharge
		myecs.Manager.DisposeEntity(stop)
		sfx.MusicPlayer.PlayTrack(constants.GameMusic, "hero")
		return false
	}, 9.5)).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
}
