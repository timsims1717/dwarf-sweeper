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
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
)

func addChest(tile *cave.Tile) {
	popUp := menus.NewPopUp(fmt.Sprintf("%s to open", typeface.SymbolItem), nil)
	popUp.Symbols = []string{data.GameInput.FirstKey("interact")}
	popUp.Dist = world.TileSize
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Transform, tile.Transform).
		AddComponent(myecs.Sprite, img.Batchers[constants.TileEntityKey].Sprites["chest_closed"]).
		AddComponent(myecs.Batch, constants.TileEntityKey).
		AddComponent(myecs.PopUp, popUp).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Interact, &data.Interact{
			OnInteract: func(pos pixel.Vec) bool {
				collectible := ""
				switch random.CaveGen.Intn(5) {
				case 0:
					item := &descent.BombItem{}
					item.Create(pos)
				case 1:
					collectible = descent.Apple
				case 2:
					collectible = descent.Beer
				case 3:
					collectible = descent.BubbleItem
				case 4:
					collectible = descent.XRayItem
				}
				if collectible != "" {
					descent.CreateCollectible(pos, collectible)
				}
				e.AddComponent(myecs.Sprite, img.Batchers[constants.TileEntityKey].Sprites["chest_opened"])
				e.RemoveComponent(myecs.Interact)
				e.RemoveComponent(myecs.PopUp)
				return true
			},
			Distance:   world.TileSize,
			Interacted: false,
			Remove:     false,
		})
}

func addBigBomb(blTile *cave.Tile, level int) {
	fmt.Printf("Bomb added here: (%d,%d)\n", blTile.RCoords.X, blTile.RCoords.Y)
	e := myecs.Manager.NewEntity()
	popUp := menus.NewPopUp(fmt.Sprintf("%s to defuse", typeface.SymbolItem), nil)
	popUp.Symbols = []string{data.GameInput.FirstKey("interact")}
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
			descent.CreateCollectible(pos, descent.GemDiamond)
		}
	}
	interact := &data.Interact{
		OnInteract: func(pos pixel.Vec) bool {
			if puzz.IsClosed() {
				descent.StartPuzzle(puzz)
			}
			return true
		},
		Distance:   world.TileSize * 1.5,
		Interacted: false,
		Remove:     false,
	}
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
		AddComponent(myecs.Batch, constants.TileEntityKey).
		AddComponent(myecs.PopUp, popUp).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Interact, interact)
}