package states

import (
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/state"
	"dwarf-sweeper/pkg/timing"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
)

const (
	MenuStateKey    = "menu_state"
	DescentStateKey = "descent_state"
	ScoreStateKey   = "score_state"
	PuzzleStateKey  = "puzzle_state"
	LoadingStateKey = "loading_state"
)

var (
	MenuState    = &menuState{}
	DescentState = &descentState{
		//CurrBiome: "mine",
	}
	ScoreState   = &scoreState{}
	PuzzleState  = &puzzleState{}
	LoadingState = &loadingState{}
	States       = map[string]*state.AbstractState{
		MenuStateKey:    state.New(MenuState, true),
		DescentStateKey: state.New(DescentState, true),
		ScoreStateKey:   state.New(ScoreState, false),
		PuzzleStateKey:  state.New(PuzzleState, true),
	}
)

var (
	switchState = true
	currState   = "unknown"
	nextState   = "menu_state"
	loading     = false
	loadingDone = false
	done        = make(chan struct{})

	debugPause     = false
	menuStack      []*menus.DwarfMenu
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
	debugInput.Update(win)
	menuInput.Update(win)
	updateState()
	if loading {
		select{
		case <-done:
			loading = false
			loadingDone = true
			currState = nextState
		default:
			LoadingState.Update(win)
		}
	} else {
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
		if debugInput.Get("debugInv").JustPressed() {
			for _, d := range descent.Descent.Dwarves {
				d.Health.Inv = !d.Health.Inv
			}
		}
		if debugInput.Get("debugMenu").JustPressed() && MenuClosed() {
			debugInput.Get("debugMenu").Consume()
			OpenMenu(DebugMenu)
		}
		if debugInput.Get("debugTest").JustPressed() && descent.Descent.Cave != nil {
			//if len(descent.Descent.GetPlayers()) > 0 {
			//	descent.CreateXRayItem(descent.Descent.GetPlayers()[0].Transform.Pos)
			//}
			//if len(descent.Descent.GetPlayers()) > 0 {
			//	particles.CreateRandomStaticParticles(1, 1, []string{"sparkle_plus_0", "sparkle_plus_1", "sparkle_plus_2", "sparkle_x_0", "sparkle_x_1", "sparkle_x_2"}, descent.Descent.GetPlayers()[0].Transform.Pos, 10.0, 15.0, 0.5)
			//}
			//menus.NotificationHandler.AddMessage("It's a message!")
			player := descent.Descent.GetPlayers()[0].Player
			descent.CreateSlug(descent.Descent.Cave, descent.Descent.Cave.GetTile(player.CamPos.Sub(player.CanvasPos.Sub(debugInput.World))).Transform.Pos)
		}
		if debugInput.Get("debugSP").JustPressed() {
			if currState == MenuStateKey {
				MenuState.splashScale *= 1.2
				fmt.Printf("Splash Scale: %f\n", MenuState.splashScale)
			} else if descent.Descent.FreeCam {
				camera.Cam.ZoomIn(1.)
			}
		}
		if debugInput.Get("debugSM").JustPressed() {
			if currState == MenuStateKey {
				MenuState.splashScale /= 1.2
				fmt.Printf("Splash Scale: %f\n", MenuState.splashScale)
			} else if descent.Descent.FreeCam {
				camera.Cam.ZoomIn(-1.)
			}
		}
		if descent.Descent.FreeCam && len(descent.Descent.Dwarves) > 0 {
			if debugInput.Get("freeCamUp").Pressed() {
				//camera.Cam.Up()
				descent.Descent.Dwarves[0].Player.CamPos.Y += 100. * timing.DT
			} else if debugInput.Get("freeCamDown").Pressed() && descent.Descent.FreeCam {
				//camera.Cam.Down()
				descent.Descent.Dwarves[0].Player.CamPos.Y -= 100. * timing.DT
			}
			if debugInput.Get("freeCamRight").Pressed() && descent.Descent.FreeCam {
				//camera.Cam.Right()
				descent.Descent.Dwarves[0].Player.CamPos.X += 100. * timing.DT
			} else if debugInput.Get("freeCamLeft").Pressed() && descent.Descent.FreeCam {
				//camera.Cam.Left()
				descent.Descent.Dwarves[0].Player.CamPos.X -= 100. * timing.DT
			}
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
			if cState, ok := States[currState]; ok {
				cState.Update(win)
			}
		}
	}
	camera.Cam.Update(win)
	myecs.UpdateManager()
}

func Draw(win *pixelgl.Window) {
	img.Clear()
	cState, ok1 := States[currState]
	nState, ok2 := States[nextState]
	if !ok2 {
		panic(fmt.Sprintf("state %s doesn't exist", nextState))
	}
	if loading && nState.ShowLoad || !ok1 {
		LoadingState.Draw(win)
	} else {
		cState.Draw(win)
	}
}

func updateState() {
	if !loading && (currState != nextState || switchState) {
		// uninitialize
		clearMenus()
		img.FullClear()
		if cState, ok := States[currState]; ok {
			go cState.Unload()
		}
		// initialize
		if nState, ok := States[nextState]; ok {
			go nState.Load(done)
			loading = true
			loadingDone = false
		}
		switchState = false
	}
}

func SwitchState(s string) {
	if !switchState {
		switchState = true
		nextState = s
	}
}
