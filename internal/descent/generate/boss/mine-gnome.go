package boss

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/objects"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/random"
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

	chunkr0 := cave.NewChunk(world.Coords{X: 1, Y: 0}, c, cave.Wall)
	chunkr1 := cave.NewChunk(world.Coords{X: 1, Y: 1}, c, cave.Wall)
	chunk1 := cave.NewChunk(world.Coords{X: 0, Y: 1}, c, cave.Wall)
	chunkl0 := cave.NewChunk(world.Coords{X: -1, Y: 0}, c, cave.Wall)
	chunkl1 := cave.NewChunk(world.Coords{X: -1, Y: 1}, c, cave.Wall)
	chunkr2 := cave.NewChunk(world.Coords{X: 2, Y: 0}, c, cave.Wall)
	chunkr3 := cave.NewChunk(world.Coords{X: 2, Y: 1}, c, cave.Wall)

	c.Chunks[chunk0.Coords] = chunk0
	c.Chunks[chunkr0.Coords] = chunkr0
	c.Chunks[chunkr1.Coords] = chunkr1
	c.Chunks[chunk1.Coords] = chunk1
	c.Chunks[chunkl0.Coords] = chunkl0
	c.Chunks[chunkl1.Coords] = chunkl1
	c.Chunks[chunkr2.Coords] = chunkr2
	c.Chunks[chunkr3.Coords] = chunkr3
	structures.Entrance(c, c.StartC, 9, 5, 3, cave.Doorway)
	structures.Stairs(c, world.Coords{X: 11, Y: 10}, true, true, 13, 6)
	structures.Stairs(c, world.Coords{X: 21, Y: 10}, false, true, 13, 6)
	if random.CaveGen.Intn(2) == 0 {
		structures.RectRoom(c, world.Coords{X: 2, Y: 15}, 3, 2, 0, cave.Empty)
		objects.AddBombDispenser(c.GetTileInt(3, 16))
		// 3,16
	} else {
		structures.RectRoom(c, world.Coords{X: 28, Y: 15}, 3, 2, 0, cave.Empty)
		objects.AddBombDispenser(c.GetTileInt(29, 16))
		// 29,16
	}
	structures.GnomeMineLayer(c, world.Coords{X: 11, Y: 23}, world.Coords{X: 21, Y: 23})

	cl := c.StartC
	cl.Y = 23
	descent.Descent.CoordsMap["current_layer"] = cl
	tt := c.StartC
	tt.Y = 22
	triggerTrans := transform.New()
	triggerTrans.Pos = c.GetTileInt(tt.X, tt.Y).Transform.Pos
	coll := data.NewCollider(pixel.R(0., 0., world.TileSize*70., world.TileSize*3.), data.GroundOnly)
	coll.Debug = true

	gnome := descent.CreateGnomeBoss(6 + level/4)
	gnome.Charge = false
	gnome.SetOnDamageFn(func() {
		nextLayer(gnome)
	})
	gnome.SetOnFleeFn(func() {
		revealExit(gnome)
	})
	myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, triggerTrans).
		AddComponent(myecs.Collision, coll).
		AddComponent(myecs.Trigger, data.NewTriggerFunc(func(p *data.Player) bool {
			TriggerGnome(gnome, p, triggerTrans)
			return true
		}))

	return c
}

func nextLayer(gnome *descent.GnomeBoss) {
	if gnome.State != descent.GBDig && gnome.Health.Curr > 1 {
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
	} else if gnome.Health.Curr == 1 {
		gnome.State = descent.GBFlee
	}
}

func revealExit(gnome *descent.GnomeBoss) {
	c := descent.Descent.Cave
	exitC := c.GetTile(gnome.Transform.Pos).RCoords
	structures.ExitDoor(c, exitC, 0, cave.SecretOpen)
	sfx.MusicPlayer.Stop(constants.GameMusic)
	c.UpdateBatch = true
}

func TriggerGnome(gnome *descent.GnomeBoss, p *data.Player, trans *transform.Transform) {
	found := false
	for !found {
		gnome.Transform.Pos, found = gnome.EmergeCoords()
	}
	//x := gnome.Transform.Pos.X + (descent.Descent.GetPlayers()[0].Transform.Pos.X-gnome.Transform.Pos.X)*0.5
	//camera.Cam.MoveTo(pixel.V(x, camera.Cam.Pos.Y), 0.4, false)
	descent.Descent.DisableInput = true
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		startRumble(gnome)
		myecs.Manager.DisposeEntity(e)
		return false
	}, 0.5)).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
}

func startRumble(gnome *descent.GnomeBoss) {
	for _, d := range descent.Descent.GetPlayers() {
		descent.MoveCam(d, gnome.Transform.Pos, 0.5)
		d.Player.Lock = true
	}
	shake := myecs.Manager.NewEntity()
	sfx.SoundPlayer.PlaySound("rockslide", 0.)
	shake.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		for _, d := range descent.Descent.GetPlayers() {
			descent.ShakeCam(d, 0.3, 8.)
		}
		return false
	}, 0.2)).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	stopShake := myecs.Manager.NewEntity()
	stopShake.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		myecs.Manager.DisposeEntity(shake)
		myecs.Manager.DisposeEntity(stopShake)
		return false
	}, 2.25)).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	gnomeFn := myecs.Manager.NewEntity()
	gnomeFn.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		myecs.Manager.DisposeEntity(gnomeFn)
		gnome.State = descent.GBEmerge
		return false
	}, 6.0)).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	stop := myecs.Manager.NewEntity()
	stop.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		descent.Descent.DisableInput = false
		for _, d := range descent.Descent.GetPlayers() {
			descent.MoveCam(d, d.Transform.Pos, 0.5)
		}
		myecs.Manager.DisposeEntity(stop)
		sfx.MusicPlayer.PlayTrack(constants.GameMusic, "hero")
		return false
	}, 9.0)).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	startGnome := myecs.Manager.NewEntity()
	startGnome.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
		for _, d := range descent.Descent.GetPlayers() {
			d.Player.Lock = false
		}
		gnome.Charge = true
		gnome.State = descent.GBCharge
		myecs.Manager.DisposeEntity(startGnome)
		return false
	}, 9.5)).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
}
