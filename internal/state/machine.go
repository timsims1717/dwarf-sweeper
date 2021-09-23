package state

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/dungeon"
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
	state      = -1
	newState   = 1
	debugPause = false
	menuStack  []*menus.DwarfMenu
	timer      *timing.FrameTimer
	timerKeys  map[string]bool
	credits    = menu.NewItemText(creditString, colornames.Aliceblue, pixel.V(1., 1.), menu.Center, menu.Center)
	title      = menu.NewItemText(titleString, colornames.Aliceblue, pixel.V(3., 3.), menu.Center, menu.Center)
	debugInput = &input.Input{
		Buttons: map[string]*input.ButtonSet{
			"debugPause":  input.NewJoyless(pixelgl.KeyF9),
			"debugResume": input.NewJoyless(pixelgl.KeyF10),
			"debug":       input.NewJoyless(pixelgl.KeyF3),
			"debugText":   input.NewJoyless(pixelgl.KeyF4),
			"debugInv":    input.NewJoyless(pixelgl.KeyF11),
		},
	}
	menuInput = &input.Input{
		Buttons: map[string]*input.ButtonSet{
			"menuUp":      input.New(pixelgl.KeyW, pixelgl.ButtonDpadUp),
			"menuDown":    input.New(pixelgl.KeyS, pixelgl.ButtonDpadDown),
			"menuLeft":    input.New(pixelgl.KeyA, pixelgl.ButtonDpadLeft),
			"menuRight":   input.New(pixelgl.KeyD, pixelgl.ButtonDpadRight),
			"menuSelect":  input.New(pixelgl.KeySpace, pixelgl.ButtonA),
			"menuBack":    input.New(pixelgl.KeyEscape, pixelgl.ButtonB),
			"pause":       input.New(pixelgl.KeyEscape, pixelgl.ButtonStart),
			"click":       input.NewJoyless(pixelgl.MouseButtonLeft),
		},
	}
	gameInput = &input.Input{
		Axes: map[string]*input.AxisSet{
			"targetX": {
				A: pixelgl.AxisRightX,
			},
			"targetY": {
				A: pixelgl.AxisRightY,
			},
		},
		Buttons: map[string]*input.ButtonSet{
			"dig": {
				Key:  pixelgl.MouseButtonLeft,
				Axis: pixelgl.AxisRightTrigger,
				GP:   1,
			},
			"mark": {
				Key:  pixelgl.MouseButtonRight,
				Axis: pixelgl.AxisLeftTrigger,
				GP:   1,
			},
			"left":      input.New(pixelgl.KeyA, pixelgl.ButtonDpadLeft),
			"right":     input.New(pixelgl.KeyD, pixelgl.ButtonDpadRight),
			"jump":      input.New(pixelgl.KeySpace, pixelgl.ButtonA),
			"lookUp":    input.New(pixelgl.KeyW, pixelgl.ButtonDpadUp),
			"lookDown":  input.New(pixelgl.KeyS, pixelgl.ButtonDpadDown),
			"climbUp":   input.New(pixelgl.KeyW, pixelgl.ButtonDpadUp),
			"climbDown": input.New(pixelgl.KeyS, pixelgl.ButtonDpadDown),
			"useItem":   input.New(pixelgl.KeyLeftShift, pixelgl.ButtonB),
			"prevItem":  {
				GPKey:  pixelgl.ButtonLeftBumper,
				Scroll: -1,
			},
			"nextItem":  {
				GPKey:  pixelgl.ButtonRightBumper,
				Scroll: 1,
			},
		},
		StickD: true,
	}
)

func Update(win *pixelgl.Window) {
	updateState()
	debugInput.Update(win)
	menuInput.Update(win)
	gameInput.Update(win)
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
	if debugInput.Get("debugInv").JustPressed() && dungeon.Dungeon.GetPlayer() != nil {
		dungeon.Dungeon.GetPlayer().Health.Inv = !dungeon.Dungeon.GetPlayer().Health.Inv
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
				bl, tr := dungeon.Dungeon.GetCave().CurrentBoundaries()
				bl.X += (camera.Cam.Width / world.TileSize) + world.TileSize
				bl.Y += (camera.Cam.Height / world.TileSize) + world.TileSize
				tr.X -= (camera.Cam.Width / world.TileSize) + world.TileSize
				tr.Y -= (camera.Cam.Height / world.TileSize) + world.TileSize
				camera.Cam.Restrict(bl, tr)
				reanimator.Update()
				if len(menuStack) > 0 {
					currMenu := menuStack[len(menuStack)-1]
					currMenu.Update(menuInput)
					if currMenu.Closed {
						if len(menuStack) > 1 {
							menuStack = menuStack[:len(menuStack)-1]
						} else {
							menuStack = []*menus.DwarfMenu{}
						}
					}
				} else {
					dungeon.Dungeon.GetCave().Update(dungeon.Dungeon.GetPlayer().Transform.Pos)
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
					dungeon.Dungeon.GetPlayer().Update(gameInput)
					systems.AnimationSystem()
					player.UpdateHUD()
					if gameInput.Get("lookUp").JustPressed() && dungeon.Dungeon.GetPlayerTile().IsExit() {
						newState = 5
					}
				}
				if dead, ok := timerKeys["death"]; (!ok || !dead) && dungeon.Dungeon.GetPlayer().Health.Dead {
					timer = timing.New(3.)
					timerKeys["death"] = true
				}
				if dead, ok := timerKeys["death"]; ok && dead {
					timer.Update()
					if (timer.Elapsed() > 1. && dungeon.Dungeon.GetPlayer().DeadStop) ||
						(timer.Elapsed() > 3. && dungeon.Dungeon.GetPlayer().Health.Dead) {
						newState = 2
					}
				}
				if menuInput.Get("pause").JustPressed() {
					menuInput.Get("pause").Consume()
					if len(menuStack) < 1 && !dungeon.Dungeon.GetPlayer().Health.Dead {
						PauseMenu.Open()
						menuStack = append(menuStack, PauseMenu)
					}
				}
			} else if state == 1 {
				title.Transform.UIPos = camera.Cam.APos
				title.Transform.UIZoom = camera.Cam.GetZoomScale()
				title.Update(pixel.Rect{})
				if len(menuStack) > 0 {
					currMenu := menuStack[len(menuStack)-1]
					currMenu.Update(menuInput)
					if currMenu.Closed {
						if len(menuStack) > 1 {
							menuStack = menuStack[:len(menuStack)-1]
						} else {
							menuStack = []*menus.DwarfMenu{}
						}
					}
				} else if menuInput.AnyJustPressed(true) {
					MainMenu.Open()
					menuStack = append(menuStack, MainMenu)
				}
			} else if state == 2 {
				reanimator.Update()
				dungeon.Dungeon.GetCave().Update(dungeon.Dungeon.GetPlayer().Transform.Pos)
				systems.PhysicsSystem()
				systems.TransformSystem()
				systems.CollisionSystem()
				systems.EntitySystem()
				particles.Update()
				vfx.Update()
				dungeon.Dungeon.GetPlayer().Update(gameInput)
				systems.AnimationSystem()
				player.UpdateHUD()
				dungeon.BlocksDugItem.Transform.UIPos = camera.Cam.APos
				dungeon.LowestLevelItem.Transform.UIPos = camera.Cam.APos
				dungeon.GemsFoundItem.Transform.UIPos = camera.Cam.APos
				dungeon.BombsMarkedItem.Transform.UIPos = camera.Cam.APos
				dungeon.WrongMarksItem.Transform.UIPos = camera.Cam.APos
				dungeon.TotalScore.Transform.UIPos = camera.Cam.APos
				dungeon.BlocksDugItem.Transform.UIZoom = camera.Cam.GetZoomScale()
				dungeon.LowestLevelItem.Transform.UIZoom = camera.Cam.GetZoomScale()
				dungeon.GemsFoundItem.Transform.UIZoom = camera.Cam.GetZoomScale()
				dungeon.BombsMarkedItem.Transform.UIZoom = camera.Cam.GetZoomScale()
				dungeon.WrongMarksItem.Transform.UIZoom = camera.Cam.GetZoomScale()
				dungeon.TotalScore.Transform.UIZoom = camera.Cam.GetZoomScale()
				dungeon.BlocksDugItem.Update(pixel.Rect{})
				dungeon.LowestLevelItem.Update(pixel.Rect{})
				dungeon.GemsFoundItem.Update(pixel.Rect{})
				dungeon.BombsMarkedItem.Update(pixel.Rect{})
				dungeon.WrongMarksItem.Update(pixel.Rect{})
				dungeon.TotalScore.Update(pixel.Rect{})
				PostGame.Update(menuInput.World, menuInput.Get("click"))
				if PostGame.Items["menu"].IsClicked() || menuInput.Get("back").JustPressed() {
					newState = 1
				}
			} else if state == 3 {
				credits.Transform.UIPos = camera.Cam.APos
				credits.Transform.UIZoom = camera.Cam.GetZoomScale()
				credits.Update(pixel.Rect{})
				if menuInput.Get("back").JustPressed() || menuInput.Get("click").JustPressed() {
					menuInput.Get("back").Consume()
					menuInput.Get("click").Consume()
					newState = 1
				}
			} else if state == 4 {
				newState = 0
			} else if state == 5 {
				reanimator.Update()
				dungeon.Dungeon.GetCave().Update(dungeon.Dungeon.GetPlayer().Transform.Pos)
				particles.Update()
				vfx.Update()
				player.UpdateHUD()
				if len(menuStack) > 0 {
					currMenu := menuStack[len(menuStack)-1]
					currMenu.Update(menuInput)
					if currMenu.Closed {
						if len(menuStack) > 1 {
							menuStack = menuStack[:len(menuStack)-1]
						} else {
							menuStack = []*menus.DwarfMenu{}
						}
					}
				} else {
					ClearEnchantMenu()
					newState = 0
				}
				if menuInput.Get("pause").JustPressed() {
					menuInput.Get("pause").Consume()
					PauseMenu.Open()
					menuStack = append(menuStack, PauseMenu)
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
		dungeon.Dungeon.GetCave().Draw(win)
		dungeon.Dungeon.GetPlayer().Draw(win, gameInput)
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
		dungeon.Dungeon.GetCave().Draw(win)
		dungeon.Dungeon.GetPlayer().Draw(win, gameInput)
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
		dungeon.ScoreTimer.Update()
		since := dungeon.ScoreTimer.Elapsed()
		if since > dungeon.BlocksDugTimer {
			dungeon.BlocksDugItem.Draw(win)
		}
		if since > dungeon.LowestLevelTimer {
			dungeon.LowestLevelItem.Draw(win)
		}
		if since > dungeon.GemsFoundTimer {
			dungeon.GemsFoundItem.Draw(win)
		}
		if since > dungeon.BombsMarkedTimer {
			dungeon.BombsMarkedItem.Draw(win)
		}
		if since > dungeon.WrongMarksTimer {
			dungeon.WrongMarksItem.Draw(win)
		}
		if since > dungeon.TotalScoreTimer {
			dungeon.TotalScore.Draw(win)
		}
		PostGame.Draw(win)
	} else if state == 3 {
		credits.Draw(win)
	} else if state == 5 {
		dungeon.Dungeon.GetCave().Draw(win)
		dungeon.Dungeon.GetPlayer().Draw(win, gameInput)
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
	if state != newState {
		timerKeys = make(map[string]bool)
		// uninitialize
		switch state {
		case 0:

		}
		// initialize
		switch newState {
		case 0:
			systems.DeleteAllEntities()
			if dungeon.Dungeon.Start {
				dungeon.BlocksDug = 0
				dungeon.LowestLevel = 0
				dungeon.LowTotal = 0
				dungeon.GemsFound = 0
				dungeon.BombsMarked = 0
				dungeon.WrongMarks = 0
			} else {
				dungeon.LowTotal = dungeon.LowestLevel
				dungeon.LowestLevel = 0
			}
			dungeon.Dungeon.Level++

			sheet, err := img.LoadSpriteSheet("assets/img/the-dark.json")
			if err != nil {
				panic(err)
			}
			dungeon.Dungeon.SetCave(dungeon.NewRoomyCave(sheet, dungeon.Dungeon.Level, -1, 1, 2))
			//dungeon.CurrCave = dungeon.NewInfiniteCave(sheet)

			if dungeon.Dungeon.Player != nil {
				dungeon.Dungeon.Player.Transform.Pos = dungeon.Dungeon.GetCave().GetStart()
			} else {
				dungeon.Dungeon.SetPlayer(dungeon.NewDwarf(dungeon.Dungeon.GetCave().GetStart()))
			}
			if dungeon.Dungeon.Start {
				player.InitHUD()
			}
			camera.Cam.SnapTo(dungeon.Dungeon.GetPlayer().Transform.Pos)

			particles.Clear()
			vfx.Clear()
			//dungeon.Entities.Clear()

			reanimator.SetFrameRate(10)
			reanimator.Reset()
			dungeon.Dungeon.Start = false
		case 1:
			title.Transform.Pos = pixel.V(0., 75.)
			camera.Cam.SnapTo(pixel.ZV)
			if state != -1 {
				MainMenu.Open()
				menuStack = append(menuStack, MainMenu)
			}
		case 2:
			x := cfg.BaseW * -0.5 + 8.
			yS := 16.
			score := 0
			score += dungeon.BlocksDug * 10
			score += (dungeon.LowestLevel + dungeon.LowTotal) * 5
			score += dungeon.GemsFound * 15
			score += dungeon.BombsMarked * 50
			score -= dungeon.WrongMarks * 20
			dungeon.BlocksDugItem   = menu.NewItemText(fmt.Sprintf("Blocks Dug:      %d x 10", dungeon.BlocksDug), colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Left, menu.Top)
			dungeon.LowestLevelItem = menu.NewItemText(fmt.Sprintf("Lowest Level:    %d x 5", dungeon.LowestLevel + dungeon.LowTotal), colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Left, menu.Top)
			dungeon.GemsFoundItem   = menu.NewItemText(fmt.Sprintf("Gems Found:      %d x 15", dungeon.GemsFound), colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Left, menu.Top)
			dungeon.BombsMarkedItem = menu.NewItemText(fmt.Sprintf("Bombs Marked:    %d x 50", dungeon.BombsMarked), colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Left, menu.Top)
			dungeon.WrongMarksItem  = menu.NewItemText(fmt.Sprintf("Incorrect Marks: %d x -20", dungeon.WrongMarks), colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Left, menu.Top)
			dungeon.TotalScore      = menu.NewItemText(fmt.Sprintf("Total Score:     %d", score), colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Left, menu.Top)
			dungeon.BlocksDugItem.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 2.)
			dungeon.LowestLevelItem.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 3.)
			dungeon.GemsFoundItem.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 4.)
			dungeon.BombsMarkedItem.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 5.)
			dungeon.WrongMarksItem.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 6.)
			dungeon.TotalScore.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 9.)
			dungeon.ScoreTimer = timing.New(5.)
			InitializePostGameMenu()
		case 3:

		case 4:
			dungeon.Dungeon.Level = 0
			dungeon.Dungeon.Start = true
			if dungeon.Dungeon.Player != nil {
				dungeon.Dungeon.Player.Delete()
				dungeon.Dungeon.Player = nil
			}
		case 5:
			success := FillEnchantMenu()
			if !success {
				ClearEnchantMenu()
			} else {
				EnchantMenu.Open()
				menuStack = append(menuStack, EnchantMenu)
			}
		}
		state = newState
	}
}