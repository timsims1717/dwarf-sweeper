package state

import (
	"dwarf-sweeper/internal/cave"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/menu"
	"dwarf-sweeper/pkg/typeface"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"time"
)

const (
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
	credits   = menu.NewItemText(creditString, colornames.Aliceblue, pixel.V(3., 3.), menu.Center, menu.Center)
	title     = text.New(pixel.ZV, typeface.BasicAtlas)
)

func Update(win *pixelgl.Window) {
	updateState()
	input.Input.Update(win)
	if input.Input.Debug {
		if debug.Debug {
			fmt.Println("DEBUG ON")
		} else {
			fmt.Println("DEBUG OFF")
		}
		debug.Debug = !debug.Debug
	}
	camera.Cam.Update(win)
	if state == 0 {
		if dead, ok := timerKeys["death"]; ok && dead {
			if time.Since(timer).Seconds() > 1. {
				newState = 2
			}
		}
		if input.Input.Back {
			newState = 1
		}
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
		cave.Entities.Update()
		particles.Update()
		vfx.Update()
		cave.Player1.Update()
		cave.BlocksDugItem.Update()
		cave.LowestLevelItem.Update()
		cave.BombsMarkedItem.Update()
		cave.BlocksMarkedItem.Update()
		cave.TotalScore.Update()
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
		MainMenu.Draw(win)
		Options.Draw(win)
		title.Draw(win, camera.Cam.UITransform(pixel.V(0., 200.), pixel.V(13., 13.), 0.))
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
			InitializeMainMenu()
			InitializeOptionsMenu()
		case 2:
			x := camera.Cam.Width * -0.5 + 200.
			cave.BlocksDugItem    = cave.NewScore(fmt.Sprintf("Blocks Dug:      %d x 10", cave.BlocksDug), pixel.V(x, 200.), 0.4)
			cave.LowestLevelItem  = cave.NewScore(fmt.Sprintf("Lowest Level:    %d x 5", cave.LowestLevel), pixel.V(x, 150.), 0.6)
			cave.BombsMarkedItem  = cave.NewScore(fmt.Sprintf("Bombs Marked:    %d x 25", cave.BombsMarked), pixel.V(x, 100.), 0.8)
			cave.BlocksMarkedItem = cave.NewScore(fmt.Sprintf("Incorrect Marks: %d x -10", cave.BlocksMarked), pixel.V(x, 50.), 1.0)
			score := 0
			score += cave.BlocksDug * 10
			score += cave.LowestLevel * 5
			score += cave.BombsMarked * 25
			score -= cave.BlocksMarked * 10
			cave.TotalScore       = cave.NewScore(fmt.Sprintf("Total Score:     %d", score), pixel.V(x, camera.Cam.Height * 0.5 - 100.), 1.2)
			cave.ScoreTimer = time.Now()
			InitializePostGameMenu()
		case 3:

		}
		state = newState
	}
}