package state

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/puzzles"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/world"
	"fmt"
)

func InitDebugMenu() {
	DebugMenu = menus.New("debug", camera.Cam)
	DebugMenu.Title = true
	debugTitle := DebugMenu.AddItem("title", "Debug Menu")
	back := DebugMenu.AddItem("back", "Back")
	freeCam := DebugMenu.AddItem("free-cam", "Free Camera")
	mineLevel := DebugMenu.AddItem("mine-level", "Start Mine Level")
	bossLevel := DebugMenu.AddItem("boss-level", "Start Boss Level")
	testMineSolver := DebugMenu.AddItem("test-solver", "Test Mine Solver")
	giveBombs := DebugMenu.AddItem("give-bombs", "Give Bombs")
	fogToggle := DebugMenu.AddItem("fog-toggle", "Toggle Fog")
	tpExit := DebugMenu.AddItem("tp-exit", "Teleport to Exit")

	debugTitle.NoHover = true
	back.SetClickFn(func() {
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	freeCam.SetClickFn(func() {
		descent.Descent.FreeCam = !descent.Descent.FreeCam
		if descent.Descent.FreeCam {
			fmt.Println("DEBUG FREE CAM ON")
		} else {
			fmt.Println("DEBUG FREE CAM OFF")
			camera.Cam.SetZoom(4. / 3.)
		}
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	bossLevel.SetClickFn(func() {
		descent.New()
		caveBuilders, err := builder.LoadBuilder(fmt.Sprint("assets/bosses.json"))
		if err != nil {
			panic(err)
		}
		choice := random.Effects.Intn(len(caveBuilders))
		descent.Descent.Builders = append(descent.Descent.Builders, []builder.CaveBuilder{caveBuilders[choice].Copy()})
		SwitchState(0)
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	mineLevel.SetClickFn(func() {
		descent.New()
		caveBuilders, err := builder.LoadBuilder(fmt.Sprint("assets/puzzles.json"))
		if err != nil {
			panic(err)
		}
		choice := random.Effects.Intn(len(caveBuilders))
		descent.Descent.Builders = append(descent.Descent.Builders, []builder.CaveBuilder{caveBuilders[choice].Copy()})
		SwitchState(0)
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	testMineSolver.SetClickFn(func() {
		minePuzzle = puzzles.NewMinePuzzle(camera.Cam, 1)
		minePuzzle.Open()
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	giveBombs.SetClickFn(func() {
		if descent.Descent.Player != nil {
			descent.AddToInventory(&descent.InvItem{
				Name:   "bomb",
				Sprite: img.Batchers[constants.EntityKey].Sprites["bomb_item"],
				OnUse: func() {
					tile := descent.Descent.GetPlayerTile()
					descent.CreateBomb(tile.Transform.Pos)
				},
				Count: 3,
				Limit: 3,
			})
		}
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	fogToggle.SetClickFn(func() {
		if descent.Descent.Cave != nil {
			descent.Descent.Cave.Fog = !descent.Descent.Cave.Fog
			if descent.Descent.Cave.Fog {
				fmt.Println("DEBUG FOG ON")
				descent.Descent.GetCave().UpdateBatch = true
			} else {
				fmt.Println("DEBUG FOG OFF")
				descent.Descent.GetCave().UpdateBatch = true
			}
		}
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	tpExit.SetClickFn(func() {
		if descent.Descent.Cave != nil && descent.Descent.Cave.ExitC != world.Origin {
			exitT := descent.Descent.Cave.GetTileInt(descent.Descent.Cave.ExitC.X, descent.Descent.Cave.ExitC.Y)
			if exitT != nil && exitT.Exit {
				descent.Descent.Player.Transform.Pos = exitT.Transform.Pos
				DebugMenu.Close()
				sfx.SoundPlayer.PlaySound("click", 2.0)
			}
		}
	})
}
