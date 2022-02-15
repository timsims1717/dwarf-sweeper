package states

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/credits"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/descend"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/player"
	"dwarf-sweeper/internal/puzzles"
	"dwarf-sweeper/internal/systems"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/menu"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	switchState = true
	state       = -1
	newState    = 1

	minePuzzle *puzzles.MinePuzzle
	puzLevel   int

	pressAKey      = menu.NewItemText("press any key", colornames.Aliceblue, pixel.V(1.4, 1.4), menu.Center, menu.Center)
	pressAKeyTimer *timing.FrameTimer
	pressAKeySec   = 1.0
	Splash         *pixel.Sprite
	splashTran     *transform.Transform
	splashScale    = 0.4
	Title          *pixel.Sprite
	titleTran      *transform.Transform
	titleScale     = 0.4
	titleY         = 70.
	debugPause     = false
	menuStack      []*menus.DwarfMenu
	timer          *timing.FrameTimer
	timerKeys      map[string]bool
	debugInput     = &input.Input{
		Buttons: map[string]*input.ButtonSet{
			"debugConsole": input.NewJoyless(pixelgl.KeyGraveAccent),
			"debug":        input.NewJoyless(pixelgl.KeyF3),
			"debugText":    input.NewJoyless(pixelgl.KeyF4),
			"debugMenu":    input.NewJoyless(pixelgl.KeyF7),
			"debugTest":    input.NewJoyless(pixelgl.KeyF8),
			"debugPause":   input.NewJoyless(pixelgl.KeyF9),
			"debugResume":  input.NewJoyless(pixelgl.KeyF10),
			"debugInv":     input.NewJoyless(pixelgl.KeyF11),
			"debugSP":      input.NewJoyless(pixelgl.KeyKPAdd),
			"debugSM":      input.NewJoyless(pixelgl.KeyKPSubtract),
			"freeCamUp":    input.NewJoyless(pixelgl.KeyP),
			"freeCamRight": input.NewJoyless(pixelgl.KeyApostrophe),
			"freeCamDown":  input.NewJoyless(pixelgl.KeySemicolon),
			"freeCamLeft":  input.NewJoyless(pixelgl.KeyL),
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
			"menuBack":   input.New(pixelgl.KeyEscape, pixelgl.ButtonB),
			"inputClear": input.New(pixelgl.KeyF1, pixelgl.ButtonBack),
			"click":      input.NewJoyless(pixelgl.MouseButtonLeft),
			"scrollUp": {
				Scroll: 1,
			},
			"scrollDown": {
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
	data.GameInput.Update(win)
	if debugInput.Get("debug").JustPressed() {
		debug.Debug = !debug.Debug
		if debug.Debug {
			fmt.Println("DEBUG ON")
		} else {
			fmt.Println("DEBUG OFF")
		}
	}
	if debugInput.Get("debugText").JustPressed() {
		debug.Text = !debug.Text
		if debug.Text {
			fmt.Println("DEBUG TEXT ON")
		} else {
			fmt.Println("DEBUG TEXT OFF")
		}
	}
	if debugInput.Get("debugInv").JustPressed() && descent.Descent.GetPlayer() != nil {
		descent.Descent.GetPlayer().Health.Inv = !descent.Descent.GetPlayer().Health.Inv
	}
	if debugInput.Get("debugMenu").JustPressed() && MenuClosed() {
		debugInput.Get("debugMenu").Consume()
		OpenMenu(DebugMenu)
	}
	if debugInput.Get("debugTest").JustPressed() {
		item := typeface.New(camera.Cam,"main", typeface.DefaultAlign, 1.2, constants.ActualHintSize,100., 0.)
		item.SetText("this is a test of {here's the tricky bit} the line width stuff man {penis} I hope this works")
		item.PrintLines()
	}
	if debugInput.Get("debugSP").JustPressed() {
		if state == 1 {
			splashScale *= 1.2
			fmt.Printf("Splash Scale: %f\n", splashScale)
		} else if descent.Descent.FreeCam {
			camera.Cam.ZoomIn(1.)
		}
	}
	if debugInput.Get("debugSM").JustPressed() {
		if state == 1 {
			splashScale /= 1.2
			fmt.Printf("Splash Scale: %f\n", splashScale)
		} else if descent.Descent.FreeCam {
			camera.Cam.ZoomIn(-1.)
		}
	}
	if debugInput.Get("freeCamUp").Pressed() && descent.Descent.FreeCam {
		camera.Cam.Up()
	} else if debugInput.Get("freeCamDown").Pressed() && descent.Descent.FreeCam {
		camera.Cam.Down()
	}
	if debugInput.Get("freeCamRight").Pressed() && descent.Descent.FreeCam {
		camera.Cam.Right()
	} else if debugInput.Get("freeCamLeft").Pressed() && descent.Descent.FreeCam {
		camera.Cam.Left()
	}
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
			reanimator.Update()
			UpdateMenus(win)
			if menuInput.Get("pause").JustPressed() || !win.Focused() {
				menuInput.Get("pause").Consume()
				if MenuClosed() && !descent.Descent.GetPlayer().Health.Dead {
					OpenMenu(PauseMenu)
					sfx.MusicPlayer.Pause(constants.GameMusic, true)
					sfx.MusicPlayer.PlayMusic("pause")
				}
			}
			if MenuClosed() {
				if !descent.UpdatePuzzle(data.GameInput) {
					descent.UpdatePlayer(data.GameInput)
					systems.TemporarySystem()
					systems.EntitySystem()
					systems.UpdateSystem()
					systems.FunctionSystem()
					systems.PhysicsSystem()
					systems.TileCollisionSystem()
					systems.CollisionSystem()
					systems.CollisionBoundSystem()
					systems.ParentSystem()
					systems.TransformSystem()
					systems.CollectSystem()
					systems.InteractSystem()
					systems.HealingSystem()
					systems.AreaDamageSystem()
					systems.DamageSystem()
					systems.HealthSystem()
					systems.PopUpSystem()
					systems.VFXSystem()
					systems.TriggerSystem()
					vfx.Update()
					descent.UpdateInventory()
					systems.AnimationSystem()
					descent.Update()
					player.UpdateHUD()
					if data.GameInput.Get("up").JustPressed() &&
						descent.Descent.GetPlayerTile().IsExit() &&
						descent.Descent.CanExit() {
						if descent.Descent.CurrDepth >= descent.Descent.Depth-1 {
							SwitchState(2)
						} else {
							SwitchState(5)
						}
					}
				}
			}
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
			bl, tr := descent.Descent.GetCave().CurrentBoundaries()
			ratio := camera.Cam.Height / constants.BaseH
			bl.X += camera.Cam.Width * 0.5 / ratio * camera.Cam.GetZoomScale()
			bl.Y += constants.BaseH * 0.5 * camera.Cam.GetZoomScale()
			tr.X -= camera.Cam.Width*0.5/ratio*camera.Cam.GetZoomScale() + world.TileSize
			tr.Y -= constants.BaseH * 0.5 * camera.Cam.GetZoomScale()
			camera.Cam.Restrict(bl, tr)
			descent.Debug(data.GameInput)
		} else if state == 1 {
			pressAKey.Transform.UIPos = camera.Cam.APos
			pressAKey.Transform.UIZoom = camera.Cam.GetZoomScale()
			pressAKey.Update(pixel.Rect{})
			if pressAKeyTimer.UpdateDone() {
				pressAKey.NoShow = !pressAKey.NoShow
				pressAKeyTimer = timing.New(pressAKeySec)
			}
			titleTran.Scalar = pixel.V(titleScale, titleScale)
			titleTran.UIPos = camera.Cam.APos
			titleTran.UIZoom = camera.Cam.GetZoomScale()
			titleTran.Update()
			splashTran.Scalar = pixel.V(splashScale, splashScale)
			splashTran.UIPos = camera.Cam.APos
			splashTran.UIZoom = camera.Cam.GetZoomScale()
			splashTran.Update()
			if credits.Opened() {
				credits.Update()
				if pressed, _ := menuInput.AnyJustPressed(true); pressed {
					credits.Close()
				}
			} else {
				UpdateMenus(win)
				pressed, mode := menuInput.AnyJustPressed(true)
				if MenuClosed() && pressed {
					OpenMenu(MainMenu)
					data.GameInput.Mode = mode
				}
			}
			//debug.AddText(fmt.Sprintf("Input TLines: %d", InputMenu.TLines))
			//debug.AddText(fmt.Sprintf("Input Top: %d", InputMenu.Top))
			//debug.AddText(fmt.Sprintf("Input Curr: %d", InputMenu.Items[InputMenu.Hovered].CurrLine))
		} else if state == 2 {
			reanimator.Update()
			systems.TemporarySystem()
			systems.PhysicsSystem()
			systems.TransformSystem()
			systems.TileCollisionSystem()
			systems.CollisionSystem()
			systems.CollisionBoundSystem()
			systems.EntitySystem()
			systems.UpdateSystem()
			systems.FunctionSystem()
			systems.VFXSystem()
			vfx.Update()
			descent.Descent.GetPlayer().Update(data.GameInput)
			systems.AnimationSystem()
			descent.Update()
			player.UpdateHUD()
			UpdateMenus(win)
			if MenuClosed() {
				SwitchState(1)
			}
		} else if state == 4 {
			SwitchState(0)
		} else if state == 5 {
			reanimator.Update()
			vfx.Update()
			descent.Update()
			player.UpdateHUD()
			UpdateMenus(win)
			if MenuClosed() {
				ClearEnchantMenu()
				SwitchState(0)
			}
		} else if state == 6 {
			if minePuzzle != nil {
				reanimator.Update()
				minePuzzle.Update(data.GameInput)
				debug.AddText(fmt.Sprintf("orig: (%d,%d)", int(minePuzzle.InfoText.Text.Orig.X), int(minePuzzle.InfoText.Text.Orig.Y)))
				if minePuzzle.Solved() {
					minePuzzle.Close()
				}
				if minePuzzle.IsClosed() && minePuzzle.Solved() {
					SwitchState(1)
				}
			}
		}
	}
	camera.Cam.Update(win)
	myecs.UpdateManager()
}

func Draw(win *pixelgl.Window) {
	img.ClearBatches()
	if state == 0 {
		descent.Descent.GetCave().Draw(win)
		//descent.Descent.GetPlayer().Draw(win, data.GameInput)
		//dungeon.Entities.Draw(win)
		systems.AnimationDraw()
		systems.SpriteDraw()
		img.DrawBatches(win)
		vfx.Draw(win)
		systems.PopUpDraw(win)
		player.DrawHUD(win)
		if descent.Descent.Puzzle != nil {
			descent.Descent.Puzzle.Draw(win)
		}
		for _, m := range menuStack {
			m.Draw(win)
		}
		debug.AddText(fmt.Sprintf("camera pos: (%f,%f)", camera.Cam.APos.X, camera.Cam.APos.Y))
		debug.AddText(fmt.Sprintf("camera zoom: %f", camera.Cam.Zoom))
		debug.AddText(fmt.Sprintf("entity count: %d", myecs.Count))
	} else if state == 1 {
		Splash.Draw(win, splashTran.Mat)
		Title.Draw(win, titleTran.Mat)
		for _, m := range menuStack {
			m.Draw(win)
		}
		if credits.Opened() {
			credits.Draw(win)
		} else if MenuClosed() {
			pressAKey.Draw(win)
		}
	} else if state == 2 {
		descent.Descent.GetCave().Draw(win)
		//descent.Descent.GetPlayer().Draw(win, data.GameInput)
		systems.AnimationDraw()
		systems.SpriteDraw()
		img.DrawBatches(win)
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
		if since > descent.BombsFlaggedTimer {
			PostMenu.ItemMap["bombs_flagged"].NoDraw = false
			PostMenu.ItemMap["bombs_flagged_s"].NoDraw = false
		}
		if since > descent.WrongFlagsTimer {
			PostMenu.ItemMap["wrong_flags"].NoDraw = false
			PostMenu.ItemMap["wrong_flags_s"].NoDraw = false
		}
		if since > descent.TotalScoreTimer {
			PostMenu.ItemMap["total_score"].NoDraw = false
			PostMenu.ItemMap["total_score_s"].NoDraw = false
		}
		for _, m := range menuStack {
			m.Draw(win)
		}
	} else if state == 5 {
		descent.Descent.GetCave().Draw(win)
		//descent.Descent.GetPlayer().Draw(win, data.GameInput)
		systems.AnimationDraw()
		systems.SpriteDraw()
		img.DrawBatches(win)
		vfx.Draw(win)
		player.DrawHUD(win)
		for _, m := range menuStack {
			m.Draw(win)
		}
	} else if state == 6 {
		if minePuzzle != nil {
			minePuzzle.Draw(win)
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
			sfx.MusicPlayer.Pause(constants.GameMusic, true)
			sfx.MusicPlayer.Stop("pause")
		case 1:
			sfx.MusicPlayer.Stop("menu")
		case 2:
			sfx.MusicPlayer.Stop("pause")
		case 4:
		case 5:
			sfx.MusicPlayer.Stop("pause")
		}
		// initialize
		switch newState {
		case 0:
			systems.ClearSystem()
			systems.DeleteAllEntities()

			descend.Descend()

			reanimator.SetFrameRate(10)
			reanimator.Reset()
		case 1:
			descent.Descent.FreeCam = false
			camera.Cam.SetZoom(4. / 3.)
			pressAKey.Transform.Pos = pixel.V(0., -75.)
			pressAKey.NoShow = true
			pressAKeyTimer = timing.New(2.5)
			titleTran = transform.New()
			titleTran.Pos = pixel.V(0., titleY)
			splashTran = transform.New()
			camera.Cam.SnapTo(pixel.ZV)
			if state != -1 {
				OpenMenu(MainMenu)
			}
			sfx.MusicPlayer.PlayMusic("menu")
		case 2:
			descent.AddStats()
			score := 0
			score += descent.BlocksDug * 2
			score += descent.GemsFound
			score += descent.BombsFlagged * 10
			score -= descent.WrongFlags * 5
			PostMenu.ItemMap["blocks_s"].Raw = fmt.Sprintf("%d x  2", descent.BlocksDug)
			PostMenu.ItemMap["gem_count_s"].Raw = fmt.Sprintf("%d x  1", descent.GemsFound)
			PostMenu.ItemMap["bombs_flagged_s"].Raw = fmt.Sprintf("%d x 10", descent.BombsFlagged)
			PostMenu.ItemMap["wrong_flags_s"].Raw = fmt.Sprintf("%d x -5", descent.WrongFlags)
			PostMenu.ItemMap["total_score_s"].Raw = fmt.Sprintf("%d", score)
			descent.ScoreTimer = timing.New(5.)
			OpenMenu(PostMenu)
			sfx.MusicPlayer.PlayMusic("pause")
		case 4:
			descend.Generate()
		case 5:
			success := FillEnchantMenu()
			if !success {
				ClearEnchantMenu()
			} else {
				OpenMenu(EnchantMenu)
				sfx.MusicPlayer.PlayMusic("pause")
			}
			descent.AddStats()
		case 6:
			puzLevel++
			reanimator.SetFrameRate(10)
			reanimator.Reset()
			minePuzzle = &puzzles.MinePuzzle{}
			minePuzzle.Create(camera.Cam, puzLevel)
			minePuzzle.Open()
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
