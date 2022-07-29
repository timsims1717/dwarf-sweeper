package states

import (
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/sfx"
	"fmt"
)

func InitDebugMenu() {
	DebugMenu = menus.New("debug", camera.Cam)
	DebugMenu.Title = true
	debugTitle := DebugMenu.AddItem("title", "Debug Menu", false)
	back := DebugMenu.AddItem("back", "Back", false)
	freeCam := DebugMenu.AddItem("free-cam", "Free Camera", false)
	bossLevel := DebugMenu.AddItem("boss-level", "Start Boss Level", false)
	testMineSolver := DebugMenu.AddItem("test-solver", "Test Mine Solver", false)
	giveBombs := DebugMenu.AddItem("give-bombs", "Give Bombs", false)
	fogToggle := DebugMenu.AddItem("fog-toggle", "Toggle Fog", false)
	tpExit := DebugMenu.AddItem("tp-exit", "Teleport to Exit", false)
	tpSecret := DebugMenu.AddItem("tp-secret", "Teleport to Secret Exit", false)
	tpBomb := DebugMenu.AddItem("tp-bomb", "Teleport to Big Bomb", false)

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
			camera.Cam.SetZoom(1.)
		}
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	bossLevel.SetClickFn(func() {
		descent.New()
		descent.Descent.Depth = 1
		DescentState.start = true
		caveBuilders, err := builder.LoadBuilder(fmt.Sprint("assets/bosses.json"))
		if err != nil {
			panic(err)
		}
		choice := random.Effects.Intn(len(caveBuilders))
		cb := caveBuilders[choice].Copy()
		descent.Descent.Builder = &cb
		SwitchState(DescentStateKey)
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	testMineSolver.SetClickFn(func() {
		SwitchState(PuzzleStateKey)
		DebugMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	giveBombs.SetClickFn(func() {
		if len(descent.Descent.GetPlayers()) > 0 {
			descent.CreateInvItem(descent.Descent.GetPlayers()[0].Player.Inventory, "bomb_item", 5)
		}
		DebugMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	fogToggle.SetClickFn(func() {
		if descent.Descent.Cave != nil {
			descent.Descent.Cave.Fog = !descent.Descent.Cave.Fog
			descent.Descent.Cave.LoadAll = !descent.Descent.Cave.LoadAll
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
		if descent.Descent.Cave != nil && len(descent.Descent.Cave.Exits) > 0 {
			exitT := descent.Descent.Cave.GetTileInt(descent.Descent.Cave.Exits[0].Coords.X, descent.Descent.Cave.Exits[0].Coords.Y)
			if exitT != nil && exitT.Exit {
				descent.Descent.GetPlayers()[0].Transform.Pos = exitT.Transform.Pos
				DebugMenu.Close()
				sfx.SoundPlayer.PlaySound("click", 2.0)
			}
		}
	})
	tpSecret.SetClickFn(func() {
		if descent.Descent.Cave != nil && len(descent.Descent.Cave.Exits) > 1 {
			exitT := descent.Descent.Cave.GetTileInt(descent.Descent.Cave.Exits[1].Coords.X, descent.Descent.Cave.Exits[1].Coords.Y)
			if exitT != nil && exitT.Exit {
				descent.Descent.GetPlayers()[0].Transform.Pos = exitT.Transform.Pos
				DebugMenu.Close()
				sfx.SoundPlayer.PlaySound("click", 2.0)
			}
		}
	})
	tpBomb.SetClickFn(func() {
		if descent.Descent.Cave != nil {
			if bc, ok := descent.Descent.CoordsMap["big-bomb"]; ok {
				bT := descent.Descent.Cave.GetTileInt(bc.X, bc.Y)
				if bT != nil {
					descent.Descent.GetPlayers()[0].Transform.Pos = bT.Transform.Pos
					DebugMenu.Close()
					sfx.SoundPlayer.PlaySound("click", 2.0)
				}
			}
		}
	})
}
