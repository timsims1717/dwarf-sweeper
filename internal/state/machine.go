package state

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/dungeon"
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/systems"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/menu"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
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
	state     = -1
	newState  = 1
	timer     time.Time
	timerKeys map[string]bool
	credits   = menu.NewItemText(creditString, colornames.Aliceblue, pixel.V(1., 1.), menu.Center, menu.Center)
	title     = menu.NewItemText(titleString, colornames.Aliceblue, pixel.V(3., 3.), menu.Center, menu.Center)
)

func Update(win *pixelgl.Window) {
	updateState()
	input.Input.Update(win)
	if input.Input.Debug {
		if debug.Debug {
			fmt.Println("DEBUG OFF")
		} else {
			fmt.Println("DEBUG ON")
		}
		debug.Debug = !debug.Debug
	}
	if input.Input.DebugInv && dungeon.Dungeon.GetPlayer() != nil {
		dungeon.Dungeon.GetPlayer().Inv = !dungeon.Dungeon.GetPlayer().Inv
	}
	if debug.Debug {
		debug.AddLine(colornames.Red, imdraw.SharpEndShape, pixel.ZV, input.Input.World, 1.)
	}
	if state == 0 {
		if win.Focused() {
			if input.Input.Debug {
				fmt.Println("DEBUG PAUSE")
			}
			reanimator.Update()
			dungeon.Dungeon.GetCave().Update(dungeon.Dungeon.GetPlayer().Transform.Pos)
			systems.PhysicsSystem()
			systems.TransformSystem()
			systems.CollisionSystem()
			systems.CollectSystem()
			systems.AnimationSystem()
			dungeon.Entities.Update()
			particles.Update()
			vfx.Update()
			dungeon.Dungeon.GetPlayer().Update()
			if dead, ok := timerKeys["death"]; (!ok || !dead) && dungeon.Dungeon.GetPlayer().Dead {
				timer = time.Now()
				timerKeys["death"] = true
			}
			if dead, ok := timerKeys["death"]; ok && dead {
				if (time.Since(timer).Seconds() > 1. && dungeon.Dungeon.GetPlayer().DeadStop) ||
					(time.Since(timer).Seconds() > 3. && dungeon.Dungeon.GetPlayer().Dead) {
					newState = 2
				}
			}
			if input.Input.Back {
				newState = 1
			}
			bl, tr := dungeon.Dungeon.GetCave().CurrentBoundaries()
			bl.X += (camera.Cam.Width / world.TileSize) + world.TileSize
			bl.Y += (camera.Cam.Height / world.TileSize) + world.TileSize
			tr.X -= (camera.Cam.Width / world.TileSize) + world.TileSize
			tr.Y -= (camera.Cam.Height / world.TileSize) + world.TileSize
			camera.Cam.Restrict(bl, tr)
			if input.Input.LookUp.JustPressed() && dungeon.Dungeon.GetPlayerTile().Exit {
				state = -1
			}
		}
	} else if state == 1 {
		title.Transform.UIPos = camera.Cam.APos
		title.Transform.UIZoom = camera.Cam.GetZoomScale()
		title.Update(pixel.Rect{})
		MainMenu.Update(input.Input.World, input.Input.Click)
		if Current == 0 && (MainMenu.Items["exit"].IsClicked() || input.Input.Back) {
			win.SetClosed(true)
		}
		Options.Update(input.Input.World, input.Input.Click)
		if Current == 1 && (Options.Items["back"].IsClicked() || input.Input.Back) {
			SwitchToMain()
		}
	} else if state == 2 {
		dungeon.Dungeon.GetCave().Update(dungeon.Dungeon.GetPlayer().Transform.Pos)
		systems.PhysicsSystem()
		systems.TransformSystem()
		systems.CollisionSystem()
		systems.AnimationSystem()
		dungeon.Entities.Update()
		particles.Update()
		vfx.Update()
		dungeon.Dungeon.GetPlayer().Update()
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
		PostGame.Update(input.Input.World, input.Input.Click)
		if PostGame.Items["menu"].IsClicked() || input.Input.Back {
			newState = 1
		}
	} else if state == 3 {
		credits.Transform.UIPos = camera.Cam.APos
		credits.Transform.UIZoom = camera.Cam.GetZoomScale()
		credits.Update(pixel.Rect{})
		if input.Input.Back || input.Input.Click.JustPressed() {
			input.Input.Click.Consume()
			newState = 1
		}
	} else if state == 4 {
		newState = 0
	}
	camera.Cam.Update(win)
}

func Draw(win *pixelgl.Window) {
	if state == 0 {
		dungeon.Dungeon.GetCave().Draw(win)
		dungeon.Dungeon.GetPlayer().Draw(win)
		dungeon.Entities.Draw(win)
		particles.Draw(win)
		vfx.Draw(win)
		debug.AddText(fmt.Sprintf("camera pos: (%f,%f)", camera.Cam.APos.X, camera.Cam.APos.Y))
		debug.AddText(fmt.Sprintf("camera zoom: %f", camera.Cam.Zoom))
	} else if state == 1 {
		MainMenu.Draw(win)
		Options.Draw(win)
		title.Draw(win)
	} else if state == 2 {
		dungeon.Dungeon.GetCave().Draw(win)
		dungeon.Dungeon.GetPlayer().Draw(win)
		dungeon.Entities.Draw(win)
		particles.Draw(win)
		vfx.Draw(win)
		since := time.Since(dungeon.ScoreTimer).Seconds()
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
			if dungeon.Dungeon.Level == 0 {
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

			dungeon.Dungeon.SetPlayer(dungeon.NewDwarf(dungeon.Dungeon.GetCave().GetStart()))
			camera.Cam.SnapTo(dungeon.Dungeon.GetPlayer().Transform.Pos)

			particles.Clear()
			vfx.Clear()
			dungeon.Entities.Clear()

			reanimator.SetFrameRate(10)
			reanimator.Reset()
		case 1:
			title.Transform.Pos = pixel.V(0., 75.)
			camera.Cam.SnapTo(pixel.ZV)
			InitializeMainMenu()
			InitializeOptionsMenu()
		case 2:
			x := cfg.BaseW * -0.5 + 8.
			yS := 16.
			score := 0
			score += dungeon.BlocksDug * 10
			score += (dungeon.LowestLevel | dungeon.LowTotal) * 5
			score += dungeon.GemsFound * 15
			score += dungeon.BombsMarked * 50
			score -= dungeon.WrongMarks * 20
			dungeon.BlocksDugItem   = menu.NewItemText(fmt.Sprintf("Blocks Dug:      %d x 10", dungeon.BlocksDug), colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Left, menu.Top)
			dungeon.LowestLevelItem = menu.NewItemText(fmt.Sprintf("Lowest Level:    %d x 5", dungeon.LowestLevel), colornames.Aliceblue, pixel.V(1.6, 1.6), menu.Left, menu.Top)
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
			dungeon.ScoreTimer = time.Now()
			InitializePostGameMenu()
		case 3:

		case 4:
			dungeon.Dungeon.Level = 0
		}
		state = newState
	}
}