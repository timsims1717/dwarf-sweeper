package structures

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/puzzles"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
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
	popUp := menus.NewPopUp("{symbol:player-interact}:disarm")
	popUp.Dist = world.TileSize * 1.5
	pos := blTile.Transform.Pos
	pos.X += world.TileSize * 0.5
	pos.Y += world.TileSize * 0.5
	trans := transform.New()
	trans.Pos = pos
	solved := false
	puzz := &puzzles.MinePuzzle{}
	puzz.Create(camera.Cam, level)
	puzz.OnSolveFn = func() {
		e.RemoveComponent(myecs.PopUp)
		e.RemoveComponent(myecs.Interact)
		myecs.Manager.NewEntity().
			AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
				solved = true
				return true
			}, 1.5))
		count := random.Effects.Intn(3) + 5
		for i := 0; i < count; i++ {
			descent.CreateGem(pos)
		}
	}
	interact := descent.NewInteract(
		func(pos pixel.Vec, d *descent.Dwarf) bool {
			if puzz.IsClosed() {
				d.Player.StartPuzzle(puzz)
			}
			return true
		}, world.TileSize * 1.5, false)
	anim := reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprite("big_bomb_idle", img.Batchers[constants.TileEntityKey].GetSprite("big_bomb_idle"), reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("big_bomb_defuse", img.Batchers[constants.TileEntityKey].GetAnimation("big_bomb_defuse").S, reanimator.Hold)).
		SetChooseFn(func() int {
			if solved {
				return 1
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
}