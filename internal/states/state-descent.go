package states

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/generate"
	"dwarf-sweeper/internal/descent/generate/builder"
	player2 "dwarf-sweeper/internal/descent/player"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/player"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/internal/systems"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/state"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type descentState struct {
	*state.AbstractState
	deathTimer *timing.FrameTimer
	start      bool
}

func (s *descentState) Unload() {
	sfx.MusicPlayer.Pause(constants.GameMusic, true)
	sfx.MusicPlayer.Stop("pause")
}

func (s *descentState) Load(done chan struct{}) {
	if s.start {
		s.start = false
		s.Generate()
	}
	s.deathTimer = nil
	systems.ClearSystem()
	systems.DeleteAllEntities()

	s.Descend()

	reanimator.SetFrameRate(10)
	reanimator.Reset()
	done <- struct{}{}
}

func (s *descentState) Update(win *pixelgl.Window) {
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
			player2.UpdateInventory()
			systems.AnimationSystem()
			descent.Update()
			if data.GameInput.Get("up").JustPressed() &&
				descent.Descent.GetPlayerTile().IsExit() &&
				descent.Descent.CanExit() {
				if descent.Descent.CurrDepth >= descent.Descent.Depth-1 {
					SwitchState(ScoreStateKey)
				} else {
					SwitchState(EnchantStateKey)
				}
			}
		}
	}
	player.UpdateHUD()
	if descent.Descent.GetPlayer().Health.Dead {
		if s.deathTimer == nil {
			s.deathTimer = timing.New(5.)
		}
		s.deathTimer.Update()
		if (s.deathTimer.Elapsed() > 2. && descent.Descent.GetPlayer().DeadStop) ||
			(s.deathTimer.Elapsed() > 4. && descent.Descent.GetPlayer().Health.Dead) {
			SwitchState(ScoreStateKey)
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
}

func (s *descentState) Draw(win *pixelgl.Window) {
	descent.Descent.GetCave().Draw(win)
	systems.DrawSystem()
	img.DrawBatches(win)
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
}

func (s *descentState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}

func (s *descentState) Generate() {
	descent.New()
	for i := 0; i < descent.Descent.Depth; i++ {
		var cb builder.CaveBuilder
		if i == descent.Descent.Depth-1 {
			caveBuilders, err := builder.LoadBuilder(fmt.Sprint("assets/bosses.json"))
			if err != nil {
				panic(err)
			}
			choice := random.Effects.Intn(len(caveBuilders))
			cb = caveBuilders[choice].Copy()
		} else if i%2 == 0 {
			caveBuilders, err := builder.LoadBuilder(fmt.Sprint("assets/caves.json"))
			if err != nil {
				panic(err)
			}
			choice := random.Effects.Intn(len(caveBuilders))
			cb = caveBuilders[choice].Copy()
		} else {
			caveBuilders, err := builder.LoadBuilder(fmt.Sprint("assets/puzzles.json"))
			if err != nil {
				panic(err)
			}
			choice := random.Effects.Intn(len(caveBuilders))
			cb = caveBuilders[choice].Copy()
		}
		descent.Descent.Builders = append(descent.Descent.Builders, []builder.CaveBuilder{cb})
	}
}

func (s *descentState) Descend() {
	if descent.Descent.Start {
		if descent.Descent.Player != nil {
			descent.Descent.Player.Delete()
			descent.Descent.Player = nil
		}
		descent.Descent.SetPlayer(descent.NewDwarf(pixel.Vec{}))
		player.InitHUD()
		player2.Inventory = []*player2.InvItem{}
		player2.ResetStats()
		descent.Descent.Start = false
	} else {
		player2.ResetCaveStats()
		descent.Descent.CurrDepth++
	}
	descent.Descent.Builder = &descent.Descent.Builders[descent.Descent.CurrDepth][0]
	descent.Descent.SetCave(generate.NewCave(descent.Descent.Builder, descent.Descent.CurrDepth*descent.Descent.Difficulty))
	if len(descent.Descent.Builder.Tracks) > 0 {
		sfx.MusicPlayer.ChooseNextTrack(constants.GameMusic, descent.Descent.Builder.Tracks)
	} else {
		sfx.MusicPlayer.Stop(constants.GameMusic)
	}
	sfx.MusicPlayer.Resume(constants.GameMusic)

	descent.Descent.Player.Transform.Pos = descent.Descent.GetCave().GetStart().Transform.Pos
	camera.Cam.SnapTo(descent.Descent.GetPlayer().Transform.Pos)
}
