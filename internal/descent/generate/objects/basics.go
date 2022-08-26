package objects

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/profile"
	"dwarf-sweeper/internal/puzzles"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/google/uuid"
	"math"
)

type FlipOpt int

const (
	Normal = iota
	Flip
	Random
)

func AddObject(tile *cave.Tile, key string, digMe bool, flipOpt FlipOpt) {
	spr := img.Batchers[constants.TileEntityKey].Sprites[key]
	hp := &data.SimpleHealth{
		Immune: data.EnemyImmunity,
		DigMe:  digMe,
	}
	if !digMe {
		hp.Immune = data.ItemImmunity2
	}
	coll := data.NewCollider(pixel.R(0., 0., spr.Frame().W(), spr.Frame().H()), data.Item)
	phys := physics.New()
	trans := transform.New().WithID(key)
	trans.Pos = tile.Transform.Pos
	switch flipOpt {
	case Flip:
		trans.Flip = true
	case Random:
		trans.Flip = random.CaveGen.Intn(2) == 0
	}
	xDiff := (world.TileSize - spr.Frame().W()) * 0.5
	if xDiff > 2. {
		trans.Pos.X += float64(random.Effects.Intn(int(xDiff))) - xDiff
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
			if hp.Dead {
				myecs.Manager.DisposeEntity(e)
			}
			return false
		}))
}

func AddChest(tile *cave.Tile) {
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
				gemCount := 2 + random.Effects.Intn(2)
				for i := 0; i < gemCount; i++ {
					descent.CreateGem(pos)
				}
				descent.CreateItemPickUp(pos, "throw_shovel", 1)
				switch random.CaveGen.Intn(5) {
				case 0:
					descent.CreateItemPickUp(pos, "bomb_item", 3)
				case 1:
					descent.CreateApple(pos)
				case 2:
					descent.CreateItemPickUp(pos, "beer", 1)
				case 3:
					descent.CreateItemPickUp(pos, "xray", 1)
				case 4:
					descent.CreateItemPickUp(pos, "throw_shovel", 1)
				}
				e.AddComponent(myecs.Drawable, img.Batchers[constants.TileEntityKey].Sprites["chest_opened"])
				e.RemoveComponent(myecs.Interact)
				e.RemoveComponent(myecs.PopUp)
				return true
			}, world.TileSize, false))
}

func AddBigBomb(blTile *cave.Tile, level int) {
	e := myecs.Manager.NewEntity()
	pe := myecs.Manager.NewEntity()
	popUp := menus.NewPopUp("{symbol:player-interact}:disarm")
	popUp.Dist = world.TileSize * 1.5
	pos := blTile.Transform.Pos
	pos.X += world.TileSize * 0.5
	pos.Y += world.TileSize * 0.5
	trans := transform.New().WithID("big-bomb")
	trans.Pos = pos
	solved := false
	failed := false
	var timer *timing.Timer
	puzz := &puzzles.MinePuzzle{}
	puzz.Create(nil, level)
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
	interact := descent.NewInteract(nil, world.TileSize * 1.5, false)
	interact.OnInteract = func(pos pixel.Vec, d *descent.Dwarf) bool {
		if !interact.Interacted || timer == nil {
			timer = timing.New(60.)
			puzz.SetTimer(timer)
		}
		if puzz.IsClosed() {
			d.Player.StartPuzzle(puzz)
			return true
		}
		return false
	}
	var fuseSFX *uuid.UUID
	anim := reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprite("big_bomb_idle", img.Batchers[constants.TileEntityKey].GetSprite("big_bomb_idle"), reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("big_bomb_defuse", img.Batchers[constants.TileEntityKey].GetAnimation("big_bomb_defuse").S, reanimator.Hold).
		SetTrigger(0, func() {
				sfx.SoundPlayer.PlaySound("bigbombdefuse", 1.0)
			}).
		SetTrigger(4, func() {
				sfx.SoundPlayer.PlaySound("bigbombsmash", 1.0)
			})).
		AddAnimation(reanimator.NewAnimFromSprites("big_bomb_ignite", img.Batchers[constants.TileEntityKey].GetAnimation("big_bomb_ignite").S, reanimator.Hold).
			SetTrigger(0, func() {
				fuseSFX = sfx.SoundPlayer.PlaySound("fuselong", -0.5)
			}).
			SetTrigger(17, func() {
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
		AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
			if timer != nil && interact.Interacted {
				timer.Update()
			}
			if timer != nil && timer.Done() {
				e.RemoveComponent(myecs.PopUp)
				e.RemoveComponent(myecs.Interact)
				myecs.Manager.DisposeEntity(pe)
				f := myecs.Manager.NewEntity()
				f.AddComponent(myecs.Func, data.NewFrameFunc(func() bool {
					if puzz.IsOpen() && puzz.Player != nil && !profile.CurrentProfile.Flags.BigBombFail {
						puzz.Player.GiveMessage("Uh oh! Better run!", func() {
							failed = true
						})
						profile.CurrentProfile.Flags.BigBombFail = true
					} else {
						failed = true
					}
					myecs.Manager.DisposeEntity(f)
					return true
				}))
			} else {
				if puzz.IsOpen() {
					timeLeft := timer.Sec() - timer.Elapsed()
					if timeLeft < 0. {
						timeLeft = 0.
					}
					secs := int(math.Round(timeLeft))
					min := secs / 60
					sec := secs % 60
					popUp.SetText(fmt.Sprintf("%02d:%02d", min, sec))
				} else if interact.Interacted {
					timeLeft := timer.Sec() - timer.Elapsed()
					if timeLeft < 0. {
						timeLeft = 0.
					}
					secs := int(math.Round(timeLeft))
					min := secs / 60
					sec := secs % 60
					popUp.SetText(fmt.Sprintf("{symbol:player-interact}:disarm\n%02d:%02d", min, sec))
				} else {
					popUp.SetText("{symbol:player-interact}:disarm")
				}
			}
			return false
		})).
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