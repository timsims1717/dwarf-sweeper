package state

import (
	"dwarf-sweeper/internal/cave"
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/internal/menu"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/typeface"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"time"
)

var (
	state     = -1
	newState  = 1
	timer     time.Time
	timerKeys map[string]bool
	credits   = text.New(pixel.ZV, typeface.BasicAtlas)
	title     = text.New(pixel.ZV, typeface.BasicAtlas)
)

func Update(win *pixelgl.Window) {
	updateState()
	if state == 0 {
		if dead, ok := timerKeys["death"]; ok && dead {
			if time.Since(timer).Seconds() > 1. {
				newState = 2
			}
		}
		input.Input.Update(win)
		if input.Input.Debug {
			fmt.Println("DEBUG PAUSE")
		}
		if input.Input.Back {
			newState = 1
		}
		camera.Cam.Update(win)
		cave.CurrCave.Update(cave.Player1.Transform.Pos)
		cave.Entities.Update()
		particles.Update()
		vfx.Update()
		cave.Player1.Update()
		if dead, ok := timerKeys["death"]; (!ok || !dead) && cave.Player1.Dead {
			timer = time.Now()
			timerKeys["death"] = true
		}
	} else if state == 1 {
		input.Input.Update(win)
		camera.Cam.Update(win)
		menu.Start.Update()
		menu.Exit.Update()
		menu.Credits.Update()
		if menu.Exit.Clicked || input.Input.Back {
			win.SetClosed(true)
		} else if menu.Start.Clicked {
			newState = 0
		} else if menu.Credits.Clicked {
			newState = 3
		}
	} else if state == 2 {
		input.Input.Update(win)
		camera.Cam.Update(win)
		cave.CurrCave.Update(cave.Player1.Transform.Pos)
		cave.Entities.Update()
		particles.Update()
		vfx.Update()
		cave.Player1.Update()
		cave.BlocksDugItem.Update()
		cave.LowestLevelItem.Update()
		cave.BombsMarkedItem.Update()
		cave.BlocksMarkedItem.Update()
		cave.TotalScore.Update()
		menu.Retry.Update()
		menu.Menu.Update()
		if menu.Menu.Clicked || input.Input.Back {
			newState = 1
		} else if menu.Retry.Clicked {
			newState = 0
		}
	} else if state == 3 {
		input.Input.Update(win)
		if input.Input.Back {
			newState = 1
		}
		camera.Cam.Update(win)
	}
}

func Draw(win *pixelgl.Window) {
	if state == 0 {
		cave.CurrCave.Draw(win)
		cave.Player1.Draw(win)
		cave.Entities.Draw(win)
		particles.Draw(win)
		vfx.Draw(win)
		//debug.Draw(win)
	} else if state == 1 {
		menu.Start.Draw(win)
		menu.Exit.Draw(win)
		menu.Credits.Draw(win)
		title.Draw(win, pixel.IM.Scaled(pixel.ZV, 4.3).Moved(pixel.V(0., 43.)))
		//debug.Draw(win)
	} else if state == 2 {
		cave.CurrCave.Draw(win)
		cave.Player1.Draw(win)
		cave.Entities.Draw(win)
		particles.Draw(win)
		vfx.Draw(win)
		//debug.Draw(win)
		cave.BlocksDugItem.Draw(win)
		cave.LowestLevelItem.Draw(win)
		cave.BombsMarkedItem.Draw(win)
		cave.BlocksMarkedItem.Draw(win)
		cave.TotalScore.Draw(win)
		menu.Retry.Draw(win)
		menu.Menu.Draw(win)
	} else if state == 3 {
		credits.Draw(win, camera.Cam.UITransform(pixel.V(camera.WindowWidthF * 0.5, camera.WindowHeightF - 200.), pixel.V(3., 3.), 0.))
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
			sheet, err := img.LoadSpriteSheet("assets/img/test-tiles.json")
			if err != nil {
				panic(err)
			}
			cave.CurrCave = cave.NewCave(sheet)

			cave.Player1 = cave.NewDwarf()
			camera.Cam.MoveTo(cave.Player1.Transform.Pos, 0.0, false)

			cave.BlocksDug = 0
			cave.LowestLevel = 0
			cave.BombsMarked = 0
			cave.BlocksMarked = 0

			particles.Clear()
			vfx.Clear()
			cave.Entities.Clear()
		case 1:
			camera.Cam.MoveTo(pixel.ZV, 0.0, false)
			title.Clear()
			title.Color = colornames.Aliceblue
			line := "DwarfSweeper"
			title.Dot.X -= title.BoundsOf(line).W() * 0.5
			fmt.Fprintln(title, line)
			menu.Start.Transform = &transform.Transform{
				Pos: pixel.V(camera.WindowWidthF*0.5, camera.WindowHeightF*0.5-30.),
			}
			menu.Exit.Transform = &transform.Transform{
				Pos: pixel.V(camera.WindowWidthF*0.5, camera.WindowHeightF*0.5-180.),
			}
			menu.Credits.Transform = &transform.Transform{
				Pos: pixel.V(camera.WindowWidthF*0.5, camera.WindowHeightF*0.5-330.),
			}
		case 2:
			cave.BlocksDugItem    = cave.NewScore(fmt.Sprintf("Blocks Dug:      %d x 10", cave.BlocksDug), pixel.V(200., camera.WindowHeightF * 0.5 + 200.), 0.4)
			cave.LowestLevelItem  = cave.NewScore(fmt.Sprintf("Lowest Level:    %d x 5", cave.LowestLevel), pixel.V(200., camera.WindowHeightF * 0.5 + 150.), 0.6)
			cave.BombsMarkedItem  = cave.NewScore(fmt.Sprintf("Bombs Marked:    %d x 25", cave.BombsMarked), pixel.V(200., camera.WindowHeightF * 0.5 + 100.), 0.8)
			cave.BlocksMarkedItem = cave.NewScore(fmt.Sprintf("Incorrect Marks: %d x -10", cave.BlocksMarked), pixel.V(200., camera.WindowHeightF * 0.5 + 50.), 1.0)
			score := 0
			score += cave.BlocksDug * 10
			score += cave.LowestLevel * 5
			score += cave.BombsMarked * 25
			score -= cave.BlocksMarked * 10
			cave.TotalScore       = cave.NewScore(fmt.Sprintf("Total Score:     %d", score), pixel.V(200., camera.WindowHeightF * 0.5 - 100.), 1.2)
			cave.ScoreTimer = time.Now()
			menu.Retry.Transform =  &transform.Transform{
				Pos:    pixel.V(camera.WindowWidthF*0.5-200., 200.),
			}
			menu.Menu.Transform =  &transform.Transform{
				Pos:    pixel.V(camera.WindowWidthF*0.5+200., 200.),
			}
		case 3:
			credits.Clear()
			credits.Color = colornames.Aliceblue
			c_title := "DwarfSweeper"
			credits.Dot.X -= credits.BoundsOf(c_title).W() * 0.5
			fmt.Fprintln(credits, c_title)
			fmt.Fprintln(credits)
			fmt.Fprintln(credits)
			lines := []string{
				"Made by Tim Sims for Ludum Dare 48",
				"(DEEPER AND DEEPER)",
				"using Pixel, a 2d Engine written in Go.",
				"",
				"Sound from the PMSFX Sampler March 2021",
				"",
				"Special Thanks:",
				"My wife Kaylan,",
				"Marshall and Clark,",
				"faiface, the Ludum Dare LD48 team,",
				"and YOU!",
				"",
				"Thanks for playing!",
			}
			for _, line := range lines {
				credits.Dot.X -= credits.BoundsOf(line).W() * 0.5
				fmt.Fprintln(credits, line)
			}
		}
		state = newState
	}
}