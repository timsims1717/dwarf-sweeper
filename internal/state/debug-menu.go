package state

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/minesweeper"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/sfx"
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
		caveBuilder, err := builder.LoadBuilder(fmt.Sprint("assets/caves.json"))
		if err != nil {
			panic(err)
		}
		newState = 0
		switchState = true
		descent.Descent.Builder = caveBuilder[2]
		descent.Descent.Level = 1
		descent.Descent.Start = true
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	mineLevel.SetClickFn(func() {
		newState = 0
		switchState = true
		descent.Descent.Type = cave.Minesweeper
		descent.Descent.Level = 1
		descent.Descent.Start = true
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	testMineSolver.SetClickFn(func() {
		b := minesweeper.CreateBoard(5, 5, 10, random.Effects)
		b.PrintToTerminalFull()
		b.RevealTilSolvable(random.Effects)
		b.PrintToTerminal()
		fmt.Printf("Was it solvable: %t", b.Solvable())
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
}