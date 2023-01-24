package states

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/generate"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/descent/quests"
	"dwarf-sweeper/internal/hud"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/profile"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/internal/systems"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/state"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
	"image/color"
)

var (
	camPos pixel.Vec
)

type descentState struct {
	*state.AbstractState
	numPlayers  int
	gameOver    bool
	deathTimer  *timing.Timer
	start       bool
	PausePlayer *data.Player
	CurrBiome   string
}

func (s *descentState) Unload() {
	profile.CurrentProfile.Flags.CorrectFlags += profile.CurrentProfile.Stats.CorrectFlags
	profile.CurrentProfile.Flags.WrongFlags += profile.CurrentProfile.Stats.WrongFlags
	profile.CurrentProfile.Flags.BlocksDug += profile.CurrentProfile.Stats.BlocksDug
	profile.CurrentProfile.Flags.BombsBlown += profile.CurrentProfile.Stats.BombsBlown
	profile.CurrentProfile.Flags.BigBombsDisarmed += profile.CurrentProfile.Stats.BigBombsDisarmed
	sfx.MusicPlayer.Pause(constants.GameMusic, true)
	sfx.MusicPlayer.Stop("pause")
}

func (s *descentState) Load(done chan struct{}) {
	if s.start {
		profile.CurrentProfile = profile.DefaultProfile
		profile.AddQuest(profile.CurrentProfile, quests.Flag5)
		profile.AddQuest(profile.CurrentProfile, quests.DiscoverMoss)
		s.SetupPlayers()
	}
	s.deathTimer = nil
	systems.ClearSystem()
	sfx.SoundPlayer.KillAll()
	systems.DeleteAllEntities()

	s.Generate()
	s.Descend()

	reanimator.SetFrameRate(10)
	reanimator.Reset()
	done <- struct{}{}
}

func (s *descentState) Update(win *pixelgl.Window) {
	for _, d := range descent.Descent.Dwarves {
		d.Player.Input.Update(win, camera.Cam.Mat)
	}
	reanimator.Update()
	UpdateMenus(win)
	if MenuClosed() {
		systems.TemporarySystem()
		descent.UpdatePlayers()
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
		systems.AnimationSystem()
		descent.Update()
		profile.UpdateQuests(profile.CurrentProfile)
		if descent.Descent.Exited {
			if descent.Descent.CurrDepth >= descent.Descent.Depth-1 {
				SwitchState(ScoreStateKey)
			} else {
				for _, d3 := range descent.Descent.GetPlayers() {
					if d3.Health.Dead {
						d3.Health.Dead = false
						d3.Health.Curr = 1
					}
				}
				SwitchState(DescentStateKey)
			}
		}
	}
	hud.UpdateHUD()
	for _, h := range hud.HUDs {
		h.Update()
	}
	descent.UpdateViews()
	s.gameOver = true
	for _, d := range descent.Descent.GetPlayers() {
		if !d.Health.Dead {
			s.gameOver = false
		}
	}
	if s.gameOver {
		if s.deathTimer == nil {
			s.deathTimer = timing.New(5.)
		}
		s.deathTimer.Update()
		if s.deathTimer.Elapsed() > 4. {
			SwitchState(ScoreStateKey)
		}
	} else {
		for _, d := range descent.Descent.GetPlayers() {
			if d.Player.Input.Get("pause").JustPressed() || !win.Focused() {
				d.Player.Input.Get("pause").Consume()
				if MenuClosed() {
					OpenMenu(PauseMenu)
					sfx.MusicPlayer.Pause(constants.GameMusic, true)
					sfx.MusicPlayer.PlayMusic("pause")
				}
				break
			}
		}
	}
	//bl, tr := descent.Descent.GetCave().CurrentBoundaries()
	//ratio := camera.Cam.Height / constants.BaseH
	//bl.X += camera.Cam.Width * 0.5 / ratio * camera.Cam.GetZoomScale()
	//bl.Y += constants.BaseH * 0.5 * camera.Cam.GetZoomScale()
	//tr.X -= camera.Cam.Width*0.5/ratio*camera.Cam.GetZoomScale() + world.TileSize
	//tr.Y -= constants.BaseH * 0.5 * camera.Cam.GetZoomScale()
	//camera.Cam.Restrict(bl, tr)
	descent.Debug(descent.Descent.Dwarves[0].Player.Input)
}

func (s *descentState) Draw(win *pixelgl.Window) {
	camPos = camera.Cam.APos
	descent.Descent.GetCave().Draw()
	systems.DrawSystem()

	for _, d := range descent.Descent.Dwarves {
		d.Player.Canvas.Clear(color.RGBA{})
		descent.Descent.GetCave().DrawBG(d.Player)
		img.Draw(d.Player.Canvas)
	}

	systems.PopUpDraw()
	for _, d := range descent.Descent.Dwarves {
		if d.Hovered != nil && !d.Health.Dazed && !d.Health.Dead {
			if d.Hovered.Solid() && d.SelectLegal {
				img.Batchers[constants.ParticleKey].GetSprite("target").Draw(d.Player.Canvas, d.Hovered.Transform.Mat)
			} else {
				img.Batchers[constants.ParticleKey].GetSprite("target_blank").Draw(d.Player.Canvas, d.Hovered.Transform.Mat)
			}
		}
		debug.DrawLines(d.Player.Canvas)
		if d.Player.Puzzle != nil {
			d.Player.Puzzle.Draw(d.Player.Canvas)
		}
		d.Player.Message.Draw()
		mat := pixel.IM
		mat = mat.Moved(camPos).Moved(d.Player.CanvasPos)
		d.Player.Canvas.Draw(win, mat)
	}

	for _, h := range hud.HUDs {
		h.Draw(win)
	}
	hud.DrawHUD(win)
	for _, m := range menuStack {
		m.Draw(win)
	}
	//debug.AddText(fmt.Sprintf("camera pos: (%f,%f)", camera.Cam.APos.X, camera.Cam.APos.Y))
	//debug.AddText(fmt.Sprintf("camera zoom: %f", camera.Cam.Zoom))
	debug.AddText(fmt.Sprintf("entity count: %d", myecs.Count))
}

func (s *descentState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}

func (s *descentState) Generate() {
	descent.Descent.Exited = false
	if !s.start {
		descent.Descent.CurrDepth++
	}
	if descent.Descent.Builder == nil {
		if !s.start && descent.Descent.ExitI > -1 && descent.Descent.ExitI < len(descent.Descent.Exits) {
			s.CurrBiome = descent.Descent.NextBiome
		}
		var cb builder.CaveBuilder
		if descent.Descent.CurrDepth == descent.Descent.Depth-1 {
			caveBuilders, err := builder.LoadBuilder(fmt.Sprintf("assets/caves/mine-bosses.json"))
			if err != nil {
				panic(err)
			}
			choice := random.Effects.Intn(len(caveBuilders))
			cb = caveBuilders[choice].Copy()
		} else if (descent.Descent.CurrDepth+1)%3 == 0 {
			caveBuilders, err := builder.LoadBuilder(fmt.Sprintf("assets/caves/%s-puzzles.json", s.CurrBiome))
			if err != nil {
				panic(err)
			}
			choice := random.Effects.Intn(len(caveBuilders))
			cb = caveBuilders[choice].Copy()
		} else {
			caveBuilders, err := builder.LoadBuilder(fmt.Sprintf("assets/caves/%s-caves.json", s.CurrBiome))
			if err != nil {
				panic(err)
			}
			choice := random.Effects.Intn(len(caveBuilders))
			cb = caveBuilders[choice].Copy()
		}
		descent.Descent.Builder = &cb
	} else {
		s.CurrBiome = descent.Descent.Builder.Biome
	}
	descent.Descent.SetCave(generate.NewCave(descent.Descent.Builder, (descent.Descent.CurrDepth+1)*descent.Descent.Difficulty))
	if !util.ContainsStr(s.CurrBiome, profile.CurrentProfile.Flags.Discovered) {
		profile.CurrentProfile.Flags.Discovered = append(profile.CurrentProfile.Flags.Discovered, s.CurrBiome)
	}
}

func (s *descentState) Descend() {
	profile.CurrentProfile.Stats = data.Stats{}
	if s.start {
		hud.InitHUD()
		s.start = false
	}
	if len(descent.Descent.Builder.Tracks) > 0 {
		sfx.MusicPlayer.ChooseNextTrack(constants.GameMusic, descent.Descent.Builder.Tracks)
	} else {
		sfx.MusicPlayer.Stop(constants.GameMusic)
	}
	sfx.MusicPlayer.Resume(constants.GameMusic)

	for i, d := range descent.Descent.Dwarves {
		pos := descent.Descent.GetCave().GetStart().Transform.Pos
		if i % 2 == 0 {
			pos.X -= world.TileSize * float64(i/2)
		} else {
			pos.X += world.TileSize * float64((i+1)/2)
		}
		d.SetStart(pos)
	}
	for _, h := range hud.HUDs {
		h.Refresh = true
	}
	descent.Descent.Builder = nil
}

func (s *descentState) SetupPlayers() {
	if s.numPlayers < 1 {
		s.numPlayers = 1
	} else if s.numPlayers > 4 {
		s.numPlayers = 4
	}
	hud.HUDs = []*hud.HUD{}
	for i := 0; i < s.numPlayers; i++ {
		var in *pxginput.Input
		var code string
		if i == 1 {
			in = data.GameInputP2
			code = "p2"
		} else if i == 2 {
			in = data.GameInputP3
			code = "p3"
		} else if i == 3 {
			in = data.GameInputP4
			code = "p4"
		} else {
			in = data.GameInputP1
			code = "p1"
		}
		p := data.New(code, in)
		d := descent.NewDwarf(p)
		//descent.CreateInvItem(d.Player.Inventory, "throw_shovel", 3)
		descent.CreateInvItem(d.Player.Inventory, "bomb_item", 3)
		descent.Descent.Dwarves = append(descent.Descent.Dwarves, d)
		hud.HUDs = append(hud.HUDs, hud.New(d))
		RegisterPlayerSymbols(in.Key, in)
	}
}