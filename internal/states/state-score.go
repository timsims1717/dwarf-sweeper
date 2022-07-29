package states

import (
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/hud"
	"dwarf-sweeper/internal/systems"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/state"
	"dwarf-sweeper/pkg/timing"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
)

var (
	BlocksDugTimer    = 0.4
	GemsFoundTimer    = 0.8
	BombsFlaggedTimer = 1.0
	WrongFlagsTimer   = 1.2
	TotalScoreTimer   = 1.4
)

type scoreState struct {
	*state.AbstractState
	ScoreTimer *timing.Timer
}

func (s *scoreState) Unload() {
	sfx.MusicPlayer.Stop("pause")
}

func (s *scoreState) Load(done chan struct{}) {
	scoreStats := player.Stats{}
	allGems := 0
	for _, d := range descent.Descent.GetPlayers() {
		scoreStats = player.AddStats(d.Player.Stats, scoreStats)
		allGems += d.Player.Gems
	}
	score := 0
	score += scoreStats.BlocksDug
	score += allGems * 5
	score += scoreStats.CorrectFlags * 10
	score -= scoreStats.WrongFlags * 5
	PostMenu.ItemMap["blocks_s"].SetText(fmt.Sprintf("%d x  1", scoreStats.BlocksDug))
	PostMenu.ItemMap["gem_count_s"].SetText(fmt.Sprintf("%d x  5", allGems))
	PostMenu.ItemMap["bombs_flagged_s"].SetText(fmt.Sprintf("%d x 10", scoreStats.CorrectFlags))
	PostMenu.ItemMap["wrong_flags_s"].SetText(fmt.Sprintf("%d x -5", scoreStats.WrongFlags))
	PostMenu.ItemMap["total_score_s"].SetText(fmt.Sprintf("%d", score))
	s.ScoreTimer = timing.New(5.)
	OpenMenu(PostMenu)
	sfx.MusicPlayer.PlayMusic("pause")
	descent.Descent.DisableInput = true
	descent.Descent.Cave.UpdateBatch = true
	done <- struct{}{}
}

func (s *scoreState) Update(win *pixelgl.Window) {
	reanimator.Update()
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
	systems.AreaDamageSystem()
	systems.DamageSystem()
	systems.VFXSystem()
	systems.AnimationSystem()
	descent.Update()
	hud.UpdateHUD()
	descent.UpdateViews()
	UpdateMenus(win)
	if MenuClosed() {
		SwitchState(MenuStateKey)
	}
}

func (s *scoreState) Draw(win *pixelgl.Window) {
	if descent.Descent.Cave != nil {
		descent.Descent.GetCave().Draw()
	}
	systems.DrawSystem()
	for _, d := range descent.Descent.Dwarves {
		d.Player.Canvas.Clear(color.RGBA{})
		descent.Descent.GetCave().DrawBG(d.Player)
		img.Draw(d.Player.Canvas)
		debug.DrawLines(d.Player.Canvas)
		mat := pixel.IM
		mat = mat.Moved(camPos).Moved(d.Player.CanvasPos)
		d.Player.Canvas.Draw(win, mat)
	}
	for _, h := range hud.HUDs {
		h.Draw(win)
	}
	hud.DrawHUD(win)
	s.ScoreTimer.Update()
	since := s.ScoreTimer.Elapsed()
	if since > BlocksDugTimer {
		PostMenu.ItemMap["blocks"].NoDraw = false
		PostMenu.ItemMap["blocks_s"].NoDraw = false
	}
	if since > GemsFoundTimer {
		PostMenu.ItemMap["gem_count"].NoDraw = false
		PostMenu.ItemMap["gem_count_s"].NoDraw = false
	}
	if since > BombsFlaggedTimer {
		PostMenu.ItemMap["bombs_flagged"].NoDraw = false
		PostMenu.ItemMap["bombs_flagged_s"].NoDraw = false
	}
	if since > WrongFlagsTimer {
		PostMenu.ItemMap["wrong_flags"].NoDraw = false
		PostMenu.ItemMap["wrong_flags_s"].NoDraw = false
	}
	if since > TotalScoreTimer {
		PostMenu.ItemMap["total_score"].NoDraw = false
		PostMenu.ItemMap["total_score_s"].NoDraw = false
	}
	for _, m := range menuStack {
		m.Draw(win)
	}
}

func (s *scoreState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}