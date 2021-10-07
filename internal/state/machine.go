package state

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/generate"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/player"
	"dwarf-sweeper/internal/systems"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/menu"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	titleString  = `DwarfSweeper`
	creditString = `DwarfSweeper


Made by Tim Sims for Ludum Dare 48
(DEEPER AND DEEPER)
using Pixel, a 2d Engine written in Go.

Sound from the PMSFX Sampler March 2021

Special Thanks:
My wife Kaylan,
Marshall and Clark,
faiface, the Ludum Dare LD48 team,
and YOU!

Thanks for playing!`
)

var (
	switchState = true
	state       = -1
	newState    = 1

	debugPause = false
	menuStack  []*menus.DwarfMenu
	timer      *timing.FrameTimer
	timerKeys  map[string]bool
	credits    = menu.NewItemText(creditString, colornames.Aliceblue, pixel.V(1., 1.), menu.Center, menu.Center)
	title      = menu.NewItemText(titleString, colornames.Aliceblue, pixel.V(3., 3.), menu.Center, menu.Center)
	debugInput = &input.Input{
		Buttons: map[string]*input.ButtonSet{
			"debugTest":   input.NewJoyless(pixelgl.KeyF7),
			"debugPause":  input.NewJoyless(pixelgl.KeyF9),
			"debugResume": input.NewJoyless(pixelgl.KeyF10),
			"debug":       input.NewJoyless(pixelgl.KeyF3),
			"debugText":   input.NewJoyless(pixelgl.KeyF4),
			"debugInv":    input.NewJoyless(pixelgl.KeyF11),
		},
		Mode: input.KeyboardMouse,
	}
	menuInput = &input.Input{
		Buttons: map[string]*input.ButtonSet{
			"menuUp": {
				Keys:    []pixelgl.Button{pixelgl.KeyW, pixelgl.KeyUp, pixelgl.KeyKP8},
				Buttons: []pixelgl.GamepadButton{pixelgl.ButtonDpadUp},
			},
			"menuDown": {
				Keys:    []pixelgl.Button{pixelgl.KeyS, pixelgl.KeyDown, pixelgl.KeyKP5},
				Buttons: []pixelgl.GamepadButton{pixelgl.ButtonDpadDown},
			},
			"menuLeft": {
				Keys:    []pixelgl.Button{pixelgl.KeyA, pixelgl.KeyLeft, pixelgl.KeyKP4},
				Buttons: []pixelgl.GamepadButton{pixelgl.ButtonDpadLeft},
			},
			"menuRight": {
				Keys:    []pixelgl.Button{pixelgl.KeyD, pixelgl.KeyRight, pixelgl.KeyKP6},
				Buttons: []pixelgl.GamepadButton{pixelgl.ButtonDpadRight},
			},
			"menuSelect": {
				Keys:    []pixelgl.Button{pixelgl.KeySpace, pixelgl.KeyEnter, pixelgl.KeyKPEnter},
				Buttons: []pixelgl.GamepadButton{pixelgl.ButtonA},
			},
			"menuBack": input.New(pixelgl.KeyEscape, pixelgl.ButtonB),
			"inputClear": input.New(pixelgl.KeyF1, pixelgl.ButtonBack),
			"click": input.NewJoyless(pixelgl.MouseButtonLeft),
			"scrollUp":  {
				Scroll: 1,
			},
			"scrollDown":  {
				Scroll: -1,
			},
			"pause": input.New(pixelgl.KeyEscape, pixelgl.ButtonStart),
		},
		Mode: input.Any,
	}
)

func Update(win *pixelgl.Window) {
	updateState()
	debugInput.Update(win)
	menuInput.Update(win)
	player.GameInput.Update(win)
	if debugInput.Get("debug").JustPressed() {
		if debug.Debug {
			fmt.Println("DEBUG OFF")
		} else {
			fmt.Println("DEBUG ON")
		}
		debug.Debug = !debug.Debug
	}
	if debugInput.Get("debugText").JustPressed() {
		if debug.Text {
			fmt.Println("DEBUG TEXT OFF")
		} else {
			fmt.Println("DEBUG TEXT ON")
		}
		debug.Text = !debug.Text
	}
	if debugInput.Get("debugInv").JustPressed() && descent.Descent.GetPlayer() != nil {
		descent.Descent.GetPlayer().Health.Inv = !descent.Descent.GetPlayer().Health.Inv
	}
	if debugInput.Get("debugTest").JustPressed() {
		//newState = 0
		//switchState = true
		//descent.Descent.Type = descent.Minesweeper
		//descent.Descent.Level = 10
		//descent.Descent.Start = true
		descent.AddToInventory(&descent.InvItem{
			Name:   "xray",
			Sprite: img.Batchers[constants.EntityKey].Sprites["x-ray-helmet"],
			OnUse:  func() bool {
				//if Dungeon.Player.Bubble == nil {
				xray := &descent.XRayHelmet{}
				xray.Create(pixel.Vec{})
				return true
				//} else {
				//	return false
				//}
			},
			Count:  1,
			Unique: true,
		})
	}
	if win.Focused() {
		frame := false
		if debugInput.Get("debugPause").JustPressed() {
			if !debugPause {
				fmt.Println("DEBUG PAUSE")
				debugPause = true
			} else {
				frame = true
			}
		} else if debugInput.Get("debugResume").JustPressed() {
			fmt.Println("DEBUG RESUME")
			debugPause = false
		}
		if !debugPause || frame {
			if state == 0 {
				bl, tr := descent.Descent.GetCave().CurrentBoundaries()
				bl.X += (camera.Cam.Width / world.TileSize) + world.TileSize
				bl.Y += (camera.Cam.Height / world.TileSize) + world.TileSize
				tr.X -= (camera.Cam.Width / world.TileSize) + world.TileSize
				tr.Y -= (camera.Cam.Height / world.TileSize) + world.TileSize
				reanimator.Update()
				UpdateMenus(win)
				if MenuClosed() {
					descent.Update()
					systems.PhysicsSystem()
					systems.CollisionSystem()
					systems.ParentSystem()
					systems.TransformSystem()
					systems.CollectSystem()
					systems.HealingSystem()
					systems.AreaDamageSystem()
					systems.DamageSystem()
					systems.HealthSystem()
					systems.EntitySystem()
					particles.Update()
					vfx.Update()
					descent.Descent.GetPlayer().Update(player.GameInput)
					systems.AnimationSystem()
					player.UpdateHUD()
					if player.GameInput.Get("up").JustPressed() &&
						descent.Descent.GetPlayerTile().IsExit() &&
						descent.Descent.CanExit() {
						SwitchState(5)
					}
				}
				camera.Cam.Restrict(bl, tr)
				if dead, ok := timerKeys["death"]; (!ok || !dead) && descent.Descent.GetPlayer().Health.Dead {
					timer = timing.New(5.)
					timerKeys["death"] = true
				}
				if dead, ok := timerKeys["death"]; ok && dead {
					timer.Update()
					if (timer.Elapsed() > 2. && descent.Descent.GetPlayer().DeadStop) ||
						(timer.Elapsed() > 4. && descent.Descent.GetPlayer().Health.Dead) {
						SwitchState(2)
					}
				}
				if menuInput.Get("pause").JustPressed() {
					menuInput.Get("pause").Consume()
					if MenuClosed() && !descent.Descent.GetPlayer().Health.Dead {
						OpenMenu(PauseMenu)
						sfx.MusicPlayer.PauseMusic("game", true)
						sfx.MusicPlayer.UnpauseOrNext("pause")
					}
				}
			} else if state == 1 {
				title.Transform.UIPos = camera.Cam.APos
				title.Transform.UIZoom = camera.Cam.GetZoomScale()
				title.Update(pixel.Rect{})
				UpdateMenus(win)
				if MenuClosed() && menuInput.AnyJustPressed(true) {
					OpenMenu(MainMenu)
				}
				//debug.AddText(fmt.Sprintf("Input TLines: %d", InputMenu.TLines))
				//debug.AddText(fmt.Sprintf("Input Top: %d", InputMenu.Top))
				//debug.AddText(fmt.Sprintf("Input Curr: %d", InputMenu.Items[InputMenu.Hovered].CurrLine))
			} else if state == 2 {
				reanimator.Update()
				descent.Update()
				systems.PhysicsSystem()
				systems.TransformSystem()
				systems.CollisionSystem()
				systems.EntitySystem()
				particles.Update()
				vfx.Update()
				descent.Descent.GetPlayer().Update(player.GameInput)
				systems.AnimationSystem()
				player.UpdateHUD()
				UpdateMenus(win)
				if MenuClosed() {
					SwitchState(1)
				}
			} else if state == 3 {
				credits.Transform.UIPos = camera.Cam.APos
				credits.Transform.UIZoom = camera.Cam.GetZoomScale()
				credits.Update(pixel.Rect{})
				if menuInput.Get("back").JustPressed() || menuInput.Get("click").JustPressed() {
					menuInput.Get("back").Consume()
					menuInput.Get("click").Consume()
					SwitchState(1)
				}
			} else if state == 4 {
				SwitchState(0)
			} else if state == 5 {
				reanimator.Update()
				descent.Update()
				particles.Update()
				vfx.Update()
				player.UpdateHUD()
				UpdateMenus(win)
				if MenuClosed() {
					ClearEnchantMenu()
					SwitchState(0)
				}
			}
		}
	}
	camera.Cam.Update(win)
	myecs.Update()
	systems.ManagementSystem()
}

func Draw(win *pixelgl.Window) {
	for _, batcher := range img.Batchers {
		batcher.Clear()
	}
	if state == 0 {
		descent.Descent.GetCave().Draw(win)
		descent.Descent.GetPlayer().Draw(win, player.GameInput)
		//dungeon.Entities.Draw(win)
		systems.AnimationDraw()
		systems.SpriteDraw()
		for _, batcher := range img.Batchers {
			if batcher.AutoDraw {
				batcher.Draw(win)
			}
		}
		particles.Draw(win)
		vfx.Draw(win)
		player.DrawHUD(win)
		for _, m := range menuStack {
			m.Draw(win)
		}
		debug.AddText(fmt.Sprintf("camera pos: (%f,%f)", camera.Cam.APos.X, camera.Cam.APos.Y))
		debug.AddText(fmt.Sprintf("camera zoom: %f", camera.Cam.Zoom))
		debug.AddText(fmt.Sprintf("entity count: %d", myecs.Count))
	} else if state == 1 {
		title.Draw(win)
		for _, m := range menuStack {
			m.Draw(win)
		}
	} else if state == 2 {
		descent.Descent.GetCave().Draw(win)
		descent.Descent.GetPlayer().Draw(win, player.GameInput)
		systems.AnimationDraw()
		systems.SpriteDraw()
		for _, batcher := range img.Batchers {
			if batcher.AutoDraw {
				batcher.Draw(win)
			}
		}
		particles.Draw(win)
		vfx.Draw(win)
		player.DrawHUD(win)
		descent.ScoreTimer.Update()
		since := descent.ScoreTimer.Elapsed()
		if since > descent.BlocksDugTimer {
			PostMenu.ItemMap["blocks"].NoDraw = false
			PostMenu.ItemMap["blocks_s"].NoDraw = false
		}
		if since > descent.GemsFoundTimer {
			PostMenu.ItemMap["gem_count"].NoDraw = false
			PostMenu.ItemMap["gem_count_s"].NoDraw = false
		}
		if since > descent.BombsMarkedTimer {
			PostMenu.ItemMap["bombs_marked"].NoDraw = false
			PostMenu.ItemMap["bombs_marked_s"].NoDraw = false
		}
		if since > descent.WrongMarksTimer {
			PostMenu.ItemMap["wrong_marks"].NoDraw = false
			PostMenu.ItemMap["wrong_marks_s"].NoDraw = false
		}
		if since > descent.TotalScoreTimer {
			PostMenu.ItemMap["total_score"].NoDraw = false
			PostMenu.ItemMap["total_score_s"].NoDraw = false
		}
		for _, m := range menuStack {
			m.Draw(win)
		}
	} else if state == 3 {
		credits.Draw(win)
	} else if state == 5 {
		descent.Descent.GetCave().Draw(win)
		descent.Descent.GetPlayer().Draw(win, player.GameInput)
		systems.AnimationDraw()
		systems.SpriteDraw()
		for _, batcher := range img.Batchers {
			if batcher.AutoDraw {
				batcher.Draw(win)
			}
		}
		particles.Draw(win)
		vfx.Draw(win)
		player.DrawHUD(win)
		for _, m := range menuStack {
			m.Draw(win)
		}
	}
}

func updateState() {
	if state != newState || switchState {
		timerKeys = make(map[string]bool)
		// uninitialize
		clearMenus()
		switch state {
		case 0:
			sfx.MusicPlayer.PauseMusic("game", true)
		case 1:
			sfx.MusicPlayer.PauseMusic("menu", true)
		case 2:
			sfx.MusicPlayer.PauseMusic("pause", true)
		case 3:
		case 4:
		case 5:
			sfx.MusicPlayer.PauseMusic("pause", true)
		}
		// initialize
		switch newState {
		case 0:
			systems.DeleteAllEntities()
			if descent.Descent.Start {
				if descent.Descent.Player != nil {
					descent.Descent.Player.Delete()
					descent.Descent.Player = nil
				}
				descent.ResetStats()
				sfx.MusicPlayer.PlayNext("game")
			} else {
				descent.ResetCaveStats()
				sfx.MusicPlayer.PauseMusic("game", false)
			}
			descent.Descent.Level++

			sheet, err := img.LoadSpriteSheet("assets/img/the-dark.json")
			if err != nil {
				panic(err)
			}
			switch descent.Descent.Type {
			case descent.Normal:
				descent.Descent.SetCave(generate.NewRoomyCave(sheet, descent.Descent.Level, -1, 1, 2))
			case descent.Infinite:
				descent.Descent.SetCave(generate.NewInfiniteCave(sheet))
			case descent.Minesweeper:
				descent.Descent.SetCave(generate.NewMinesweeperCave(sheet, descent.Descent.Level))
			}

			if descent.Descent.Player != nil {
				descent.Descent.Player.Transform.Pos = descent.Descent.GetCave().GetStart()
			} else {
				descent.Descent.SetPlayer(descent.NewDwarf(descent.Descent.GetCave().GetStart()))
			}
			if descent.Descent.Start {
				player.InitHUD()
				descent.Inventory = []*descent.InvItem{}
			}
			camera.Cam.SnapTo(descent.Descent.GetPlayer().Transform.Pos)

			particles.Clear()
			vfx.Clear()
			//dungeon.Entities.Clear()

			reanimator.SetFrameRate(10)
			reanimator.Reset()
			descent.Descent.Start = false
		case 1:
			title.Transform.Pos = pixel.V(0., 75.)
			camera.Cam.SnapTo(pixel.ZV)
			if state != -1 {
				OpenMenu(MainMenu)
			}
			sfx.MusicPlayer.PlayNext("menu")
		case 2:
			descent.AddStats()
			score := 0
			score += descent.BlocksDug * 2
			score += descent.GemsFound
			score += descent.BombsMarked * 10
			score -= descent.WrongMarks * 5
			PostMenu.ItemMap["blocks_s"].Raw = fmt.Sprintf("%d x  2", descent.BlocksDug)
			PostMenu.ItemMap["gem_count_s"].Raw = fmt.Sprintf("%d x  1", descent.GemsFound)
			PostMenu.ItemMap["bombs_marked_s"].Raw = fmt.Sprintf("%d x 10", descent.BombsMarked)
			PostMenu.ItemMap["wrong_marks_s"].Raw = fmt.Sprintf("%d x -5", descent.WrongMarks)
			PostMenu.ItemMap["total_score_s"].Raw = fmt.Sprintf("%d", score)
			descent.ScoreTimer = timing.New(5.)
			OpenMenu(PostMenu)
			sfx.MusicPlayer.UnpauseOrNext("pause")
		case 3:

		case 4:
			descent.Descent.Level = 0
			descent.Descent.Start = true
		case 5:
			success := FillEnchantMenu()
			if !success {
				ClearEnchantMenu()
			} else {
				OpenMenu(EnchantMenu)
				sfx.MusicPlayer.UnpauseOrNext("pause")
			}
			descent.AddStats()
			if descent.Descent.Type == descent.Normal {
				descent.Descent.Type = descent.Minesweeper
			} else {
				descent.Descent.Type = descent.Normal
			}
		}
		state = newState
		switchState = false
	}
}

func SwitchState(s int) {
	if !switchState {
		switchState = true
		newState = s
	}
}