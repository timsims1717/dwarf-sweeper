package structures

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/puzzles"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/google/uuid"
)

func addChest(tile *cave.Tile) {
	popUp := menus.NewPopUp("{symbol:player-interact}:open")
	popUp.Dist = world.TileSize
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Transform, tile.Transform).
		AddComponent(myecs.Drawable, img.Batchers[constants.TileEntityKey].Sprites["chest_closed"]).
		AddComponent(myecs.Batch, constants.TileEntityKey).
		AddComponent(myecs.PopUp, popUp).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Interact, descent.NewInteract(
			func(pos pixel.Vec, _ *descent.Dwarf) bool {
				switch random.CaveGen.Intn(5) {
				case 0:
					descent.CreateBombItem(pos)
				case 1:
					descent.CreateApple(pos)
				case 2:
					descent.CreateBeerItem(pos)
				case 3:
					descent.CreateBubbleItem(pos)
				case 4:
					descent.CreateXRayItem(pos)
				}
				gemCount := 2 + random.Effects.Intn(2)
				for i := 0; i < gemCount; i++ {
					descent.CreateGem(pos)
				}
				e.AddComponent(myecs.Drawable, img.Batchers[constants.TileEntityKey].Sprites["chest_opened"])
				e.RemoveComponent(myecs.Interact)
				e.RemoveComponent(myecs.PopUp)
				return true
			}, world.TileSize, false))
}

func addBigBomb(blTile *cave.Tile, level int) {
	fmt.Printf("Bomb added here: (%d,%d)\n", blTile.RCoords.X, blTile.RCoords.Y)
	e := myecs.Manager.NewEntity()
	pe := myecs.Manager.NewEntity()
	popUp := menus.NewPopUp("{symbol:player-interact}:disarm")
	popUp.Dist = world.TileSize * 1.5
	pos := blTile.Transform.Pos
	pos.X += world.TileSize * 0.5
	pos.Y += world.TileSize * 0.5
	trans := transform.New()
	trans.Pos = pos
	solved := false
	failed := false
	puzz := &puzzles.MinePuzzle{}
	puzz.Create(camera.Cam, level)
	puzz.SetOnSolve(func() {
		e.RemoveComponent(myecs.PopUp)
		e.RemoveComponent(myecs.Interact)
		myecs.Manager.DisposeEntity(pe)
		s := myecs.Manager.NewEntity()
		s.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
			solved = true
			myecs.Manager.DisposeEntity(s)
			return true
		}, 1.5))
		count := random.Effects.Intn(3) + 5
		for i := 0; i < count; i++ {
			descent.CreateGem(pos)
		}
	})
	puzz.SetOnFail(func() {
		e.RemoveComponent(myecs.PopUp)
		e.RemoveComponent(myecs.Interact)
		myecs.Manager.DisposeEntity(pe)
		f := myecs.Manager.NewEntity()
		f.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
			failed = true
			myecs.Manager.DisposeEntity(f)
			return true
		}, 1.5))
	})
	interact := descent.NewInteract(
		func(pos pixel.Vec, d *descent.Dwarf) bool {
			if puzz.IsClosed() {
				d.Player.StartPuzzle(puzz)
				return true
			}
			return false
		}, world.TileSize * 1.5, false)
	var fuseSFX *uuid.UUID
	anim := reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprite("big_bomb_idle", img.Batchers[constants.TileEntityKey].GetSprite("big_bomb_idle"), reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("big_bomb_defuse", img.Batchers[constants.TileEntityKey].GetAnimation("big_bomb_defuse").S, reanimator.Hold).
			SetTrigger(0, func(_ *reanimator.Anim, _ string, _ int) {
				sfx.SoundPlayer.PlaySound("bigbombdefuse", 1.0)
			}).
			SetTrigger(4, func(_ *reanimator.Anim, _ string, _ int) {
				sfx.SoundPlayer.PlaySound("bigbombsmash", 1.0)
			})).
		AddAnimation(reanimator.NewAnimFromSprites("big_bomb_ignite", img.Batchers[constants.TileEntityKey].GetAnimation("big_bomb_ignite").S, reanimator.Hold).
			SetTrigger(0, func(_ *reanimator.Anim, _ string, _ int) {
				fuseSFX = sfx.SoundPlayer.PlaySound("fuselong", -0.5)
			}).
			SetTrigger(17, func(_ *reanimator.Anim, _ string, _ int) {
				bigBombDestroy(trans)
				bigBombDmg(trans)
				particles.CreateHugeExplosion(trans.Pos)
				sfx.SoundPlayer.KillSound(fuseSFX)
				sfx.SoundPlayer.PlaySound("bigblast1", -0.5)
				myecs.Manager.DisposeEntity(e)
			}),
		).
		SetChooseFn(func() int {
			if solved {
				return 1
			} else if failed {
				return 2
			} else {
				return 0
			}
		}), "big_bomb_idle")
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Animation, anim).
		AddComponent(myecs.Drawable, anim).
		AddComponent(myecs.Batch, constants.TileEntityKey).
		AddComponent(myecs.PopUp, popUp).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Interact, interact)
	pet := transform.New()
	pet.Offset.Y -= 7.
	sw := false
	pe.AddComponent(myecs.Transform, pet).
		AddComponent(myecs.Parent, trans).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
			if reanimator.FrameSwitch {
				if sw {
					pe.AddComponent(myecs.Drawable, fmt.Sprintf("big_bomb_pad_%d%d", random.Effects.Intn(5), random.Effects.Intn(5)))
				}
				sw = !sw
			}
			return false
		})).
		AddComponent(myecs.Drawable, fmt.Sprintf("big_bomb_pad_%d%d", random.Effects.Intn(5), random.Effects.Intn(5))).
		AddComponent(myecs.Batch, constants.TileEntityKey)
}

func bigBombDestroy(trans *transform.Transform) {
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Func, data.NewFrameFunc(func() bool {
		t := trans.Pos
		t.Y -= world.TileSize * 0.5
		t.X -= world.TileSize * 0.5
		tile := descent.Descent.Cave.GetTile(t)
		x0 := 0
		x := tile.RCoords.X - 7
		y := tile.RCoords.Y
		for x0 < 16 {
			descent.Descent.Cave.GetTileInt(x, y).Destroy(nil, false)
			if x0 > 2 && x0 < 13 {
				descent.Descent.Cave.GetTileInt(x, y+1).Destroy(nil, false)
				descent.Descent.Cave.GetTileInt(x, y-1).Destroy(nil, false)
				descent.Descent.Cave.GetTileInt(x, y-2).Destroy(nil, false)
				descent.Descent.Cave.GetTileInt(x, y-3).Destroy(nil, false)
				if x0 > 4 && x0 < 11 {
					descent.Descent.Cave.GetTileInt(x, y+2).Destroy(nil, false)
					descent.Descent.Cave.GetTileInt(x, y-4).Destroy(nil, false)
					if x0 > 6 && x0 < 9 {
						descent.Descent.Cave.GetTileInt(x, y-5).Destroy(nil, false)
						descent.Descent.Cave.GetTileInt(x, y-6).Destroy(nil, false)
					}
				}
			}
			x++
			x0++
		}
		myecs.Manager.DisposeEntity(e)
		return false
	}))
}

func bigBombDmg(trans *transform.Transform) {
	pos1 := trans.Pos
	pos1.Y += world.TileSize * 0.5
	myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
		SourceID:       trans.ID,
		Center:         pos1,
		Rect:           pixel.R(0., 0., world.TileSize * 10., world.TileSize * 5.),
		Amount:         3,
		Dazed:          3.,
		Knockback:      25.,
		KnockbackDecay: false,
	})
	pos2 := trans.Pos
	pos2.Y -= world.TileSize * 0.5
	myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
		SourceID:       trans.ID,
		Center:         pos2,
		Rect:           pixel.R(0., 0., world.TileSize * 16., world.TileSize),
		Amount:         3,
		Dazed:          3.,
		Knockback:      25.,
		KnockbackDecay: false,
	})
	pos3 := trans.Pos
	pos3.Y += world.TileSize * 0.5
	myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
		SourceID:       trans.ID,
		Center:         pos3,
		Rect:           pixel.R(0., 0., world.TileSize * 6., world.TileSize * 7.),
		Amount:         3,
		Dazed:          3.,
		Knockback:      25.,
		KnockbackDecay: false,
	})
	pos4 := trans.Pos
	pos4.Y += world.TileSize * 1.5
	myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
		SourceID:       trans.ID,
		Center:         pos4,
		Rect:           pixel.R(0., 0., world.TileSize * 2., world.TileSize * 9.),
		Amount:         3,
		Dazed:          3.,
		Knockback:      25.,
		KnockbackDecay: false,
	})
}