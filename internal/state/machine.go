package state

import (
	"dwarf-sweeper/internal/cave"
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/debug"
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
	if debug.Debug {
		debug.AddLine(colornames.Red, imdraw.SharpEndShape, pixel.ZV, input.Input.World, 1.)
	}
	if state == 0 {
		if win.Focused() {
			reanimator.Update()
			cave.CurrCave.Update(cave.Player1.Transform.Pos)
			systems.PhysicsSystem()
			systems.TransformSystem()
			systems.CollisionSystem()
			systems.CollectSystem()
			systems.AnimationSystem()
			cave.Entities.Update()
			particles.Update()
			vfx.Update()
			cave.Player1.Update()
			if dead, ok := timerKeys["death"]; (!ok || !dead) && cave.Player1.Dead {
				timer = time.Now()
				timerKeys["death"] = true
			}
			if dead, ok := timerKeys["death"]; ok && dead {
				if time.Since(timer).Seconds() > 1. {
					newState = 2
				}
			}
			if input.Input.Back {
				newState = 1
			}
			bl, tr := cave.CurrCave.CurrentBoundaries()
			bl.X += (camera.Cam.Width / world.TileSize) + world.TileSize
			bl.Y += (camera.Cam.Height / world.TileSize) + world.TileSize
			tr.X -= (camera.Cam.Width / world.TileSize) + world.TileSize
			tr.Y -= (camera.Cam.Height / world.TileSize) + world.TileSize
			camera.Cam.Restrict(bl, tr)
		}
	} else if state == 1 {
		title.Transform.UIPos = camera.Cam.Pos
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
		cave.CurrCave.Update(cave.Player1.Transform.Pos)
		systems.PhysicsSystem()
		systems.TransformSystem()
		systems.CollisionSystem()
		systems.AnimationSystem()
		cave.Entities.Update()
		particles.Update()
		vfx.Update()
		cave.Player1.Update()
		cave.BlocksDugItem.Transform.UIPos = camera.Cam.Pos
		cave.LowestLevelItem.Transform.UIPos = camera.Cam.Pos
		cave.GemsFoundItem.Transform.UIPos = camera.Cam.Pos
		cave.BombsMarkedItem.Transform.UIPos = camera.Cam.Pos
		cave.WrongMarksItem.Transform.UIPos = camera.Cam.Pos
		cave.TotalScore.Transform.UIPos = camera.Cam.Pos
		cave.BlocksDugItem.Transform.UIZoom = camera.Cam.GetZoomScale()
		cave.LowestLevelItem.Transform.UIZoom = camera.Cam.GetZoomScale()
		cave.GemsFoundItem.Transform.UIZoom = camera.Cam.GetZoomScale()
		cave.BombsMarkedItem.Transform.UIZoom = camera.Cam.GetZoomScale()
		cave.WrongMarksItem.Transform.UIZoom = camera.Cam.GetZoomScale()
		cave.TotalScore.Transform.UIZoom = camera.Cam.GetZoomScale()
		cave.BlocksDugItem.Update(pixel.Rect{})
		cave.LowestLevelItem.Update(pixel.Rect{})
		cave.GemsFoundItem.Update(pixel.Rect{})
		cave.BombsMarkedItem.Update(pixel.Rect{})
		cave.WrongMarksItem.Update(pixel.Rect{})
		cave.TotalScore.Update(pixel.Rect{})
		PostGame.Update(input.Input.World, input.Input.Click)
		if PostGame.Items["menu"].IsClicked() || input.Input.Back {
			newState = 1
		}
	} else if state == 3 {
		credits.Transform.UIPos = camera.Cam.Pos
		credits.Transform.UIZoom = camera.Cam.GetZoomScale()
		credits.Update(pixel.Rect{})
		if input.Input.Back || input.Input.Click.JustPressed() {
			input.Input.Click.Consume()
			newState = 1
		}
	}
	camera.Cam.Update(win)
}

func Draw(win *pixelgl.Window) {
	if state == 0 {
		cave.CurrCave.Draw(win)
		cave.Player1.Draw(win)
		cave.Entities.Draw(win)
		particles.Draw(win)
		vfx.Draw(win)
	} else if state == 1 {
		MainMenu.Draw(win)
		Options.Draw(win)
		title.Draw(win)
	} else if state == 2 {
		cave.CurrCave.Draw(win)
		cave.Player1.Draw(win)
		cave.Entities.Draw(win)
		particles.Draw(win)
		vfx.Draw(win)
		since := time.Since(cave.ScoreTimer).Seconds()
		if since > cave.BlocksDugTimer {
			cave.BlocksDugItem.Draw(win)
		}
		if since > cave.LowestLevelTimer {
			cave.LowestLevelItem.Draw(win)
		}
		if since > cave.GemsFoundTimer {
			cave.GemsFoundItem.Draw(win)
		}
		if since > cave.BombsMarkedTimer {
			cave.BombsMarkedItem.Draw(win)
		}
		if since > cave.WrongMarksTimer {
			cave.WrongMarksItem.Draw(win)
		}
		if since > cave.TotalScoreTimer {
			cave.TotalScore.Draw(win)
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
		//switch state {
		//case 0:
		//
		//}
		// initialize
		switch newState {
		case 0:
			sheet, err := img.LoadSpriteSheet("assets/img/the-dark.json")
			if err != nil {
				panic(err)
			}
			cave.CurrCave = cave.NewCave(sheet)

			cave.Player1 = cave.NewDwarf()
			camera.Cam.MoveTo(cave.Player1.Transform.Pos, 0.0, false)

			cave.BlocksDug = 0
			cave.LowestLevel = 0
			cave.BombsMarked = 0
			cave.WrongMarks = 0

			particles.Clear()
			vfx.Clear()
			cave.Entities.Clear()

			reanimator.SetFrameRate(10)
			reanimator.Reset()
		case 1:
			title.Transform.Pos = pixel.V(0., 75.)
			camera.Cam.SnapTo(pixel.ZV)
			InitializeMainMenu()
			InitializeOptionsMenu()
		case 2:
			x := cfg.BaseW * -0.5 + 15.
			yS := 16.
			score := 0
			score += cave.BlocksDug * 10
			score += cave.LowestLevel * 5
			score += cave.GemsFound * 15
			score += cave.BombsMarked * 50
			score -= cave.WrongMarks * 20
			cave.BlocksDugItem   = menu.NewItemText(fmt.Sprintf("Blocks Dug:      %d x 10", cave.BlocksDug), colornames.Aliceblue, pixel.V(1., 1.), menu.Left, menu.Top)
			cave.LowestLevelItem = menu.NewItemText(fmt.Sprintf("Lowest Level:    %d x 5", cave.LowestLevel), colornames.Aliceblue, pixel.V(1., 1.), menu.Left, menu.Top)
			cave.GemsFoundItem = menu.NewItemText(fmt.Sprintf("Gems Found:      %d x 15", cave.GemsFound), colornames.Aliceblue, pixel.V(1., 1.), menu.Left, menu.Top)
			cave.BombsMarkedItem = menu.NewItemText(fmt.Sprintf("Bombs Marked:    %d x 50", cave.BombsMarked), colornames.Aliceblue, pixel.V(1., 1.), menu.Left, menu.Top)
			cave.WrongMarksItem  = menu.NewItemText(fmt.Sprintf("Incorrect Marks: %d x -20", cave.WrongMarks), colornames.Aliceblue, pixel.V(1., 1.), menu.Left, menu.Top)
			cave.TotalScore      = menu.NewItemText(fmt.Sprintf("Total Score:     %d", score), colornames.Aliceblue, pixel.V(1., 1.), menu.Left, menu.Top)
			cave.BlocksDugItem.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 2.)
			cave.LowestLevelItem.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 3.)
			cave.GemsFoundItem.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 4.)
			cave.BombsMarkedItem.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 5.)
			cave.WrongMarksItem.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 6.)
			cave.TotalScore.Transform.Pos = pixel.V(x, cfg.BaseH * 0.5 - yS * 9.)
			cave.ScoreTimer = time.Now()
			InitializePostGameMenu()
		case 3:

		}
		state = newState
	}
}