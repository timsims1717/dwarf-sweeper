package state

import (
	"dwarf-sweeper/internal/cave"
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/internal/menu"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
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
		bl := cave.Player1.Transform.Pos
		bl.X -= world.TileSize * 1.45
		bl.Y -= world.TileSize * 1.45
		tr := cave.Player1.Transform.Pos
		tr.X += world.TileSize * 1.45
		tr.Y += world.TileSize * 1.45
		input.Restrict(win, bl, tr)
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
		if input.Input.Back {
			newState = 1
		}
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
		case 2:
			cave.BlocksDugItem    = cave.NewScore(fmt.Sprintf("Blocks Dug:      %d x 10", cave.BlocksDug), pixel.V(200., camera.WindowHeightF * 0.5 + 200.), 1.)
			cave.LowestLevelItem  = cave.NewScore(fmt.Sprintf("Lowest Level:    %d x 1", cave.LowestLevel), pixel.V(200., camera.WindowHeightF * 0.5 + 150.), 1.2)
			cave.BombsMarkedItem  = cave.NewScore(fmt.Sprintf("Bombs Marked:    %d x 50", cave.BombsMarked), pixel.V(200., camera.WindowHeightF * 0.5 + 100.), 1.4)
			cave.BlocksMarkedItem = cave.NewScore(fmt.Sprintf("Incorrect Marks: %d x -20", cave.BlocksMarked), pixel.V(200., camera.WindowHeightF * 0.5 + 50.), 1.6)
			score := 0
			score += cave.BlocksDug * 10
			score += cave.LowestLevel
			score += cave.BombsMarked * 50
			score -= cave.BlocksMarked * 20
			cave.TotalScore       = cave.NewScore(fmt.Sprintf("Total Score:     %d", score), pixel.V(200., camera.WindowHeightF * 0.5 - 100.), 1.8)
			cave.ScoreTimer = time.Now()
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