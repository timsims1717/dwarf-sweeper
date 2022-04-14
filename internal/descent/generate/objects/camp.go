package objects

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

func AddTent(tile *cave.Tile, right bool) {
	spr := img.Batchers[constants.TileEntityKey].Sprites["tent"]
	hp := &data.SimpleHealth{
		Immune: map[data.DamageType]data.Immunity{
			data.Shovel: {
				KB:    true,
				DMG:   true,
				Dazed: true,
			},
			data.Enemy: {
				KB:    true,
				DMG:   true,
				Dazed: true,
			},
		},
	}
	coll := data.NewCollider(pixel.R(0., 0., spr.Frame().W(), spr.Frame().H()), true, false)
	phys := physics.New()
	trans := transform.New()
	trans.Pos = tile.Transform.Pos
	if right {
		trans.Pos.X -= world.TileSize * 0.5
		trans.Flip = true
	} else {
		trans.Pos.X += world.TileSize * 0.5
	}
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Health, hp).
		AddComponent(myecs.Collision, coll).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Drawable, spr).
		AddComponent(myecs.Batch, constants.TileEntityKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
			x2 := tile.RCoords.X+1
			if right {
				x2 -= 2
			}
			below1 := descent.Descent.Cave.GetTileInt(tile.RCoords.X, tile.RCoords.Y+1)
			below2 := descent.Descent.Cave.GetTileInt(x2, tile.RCoords.Y+1)
			if below1 == nil || !below1.Solid() || below2 == nil || !below2.Solid() || hp.Dead {
				// death
				myecs.Manager.DisposeEntity(e)
			}
			return false
		}))
}

func AddGemPile(tile *cave.Tile) {
	spr := img.Batchers[constants.TileEntityKey].Sprites["gempile"]
	hp := &data.SimpleHealth{
		Immune: data.EnemyImmunity,
		DigMe:  true,
	}
	coll := data.NewCollider(pixel.R(0., 0., spr.Frame().W(), spr.Frame().H()), true, false)
	phys := physics.New()
	trans := transform.New()
	trans.Pos = tile.Transform.Pos
	xDiff := (world.TileSize - spr.Frame().H()) * 0.5
	trans.Pos.X += float64(random.Effects.Intn(int(xDiff))) - xDiff
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Health, hp).
		AddComponent(myecs.Collision, coll).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Drawable, spr).
		AddComponent(myecs.Batch, constants.TileEntityKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
			below := descent.Descent.Cave.GetTileInt(tile.RCoords.X, tile.RCoords.Y+1)
			if below == nil || !below.Solid() || hp.Dead {
				myecs.Manager.DisposeEntity(e)
				gemCount := 2 + random.Effects.Intn(2)
				for i := 0; i < gemCount; i++ {
					descent.CreateGem(trans.Pos)
				}
			}
			return false
		}))
}

func AddSeat(tile *cave.Tile) {
	AddObject(tile, "seat", false)
}

func AddArt(tile *cave.Tile) {
	AddObject(tile, "art", false)
}

func AddBarrel(tile *cave.Tile) {
	AddObject(tile, "barrel", false)
}

func AddBedroll(tile *cave.Tile) {
	AddObject(tile, "bedroll", false)
}

func AddRefuse(tile *cave.Tile) {
	if random.CaveGen.Intn(2) == 0 {
		AddObject(tile, "refuse_sm", false)
	} else {
		AddObject(tile, "refuse_lg", false)
	}
}

func AddCookfire(tile *cave.Tile) {
	AddObject(tile, "cookfire", false)
}

func AddTools(tile *cave.Tile) {
	AddObject(tile, "tools", false)
}

func AddWoodpile(tile *cave.Tile) {
	AddObject(tile, "woodpile", false)
}