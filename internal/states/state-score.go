package states

import (
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/hud"
	"dwarf-sweeper/internal/systems"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/state"
	"dwarf-sweeper/pkg/timing"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
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
	for _, d := range descent.Descent.GetPlayers() {
		d.Player.Stats.AddStats()
	}
	player.OverallStats.AddStats()
	p := descent.Descent.GetPlayers()[0].Player
	score := 0
	score += p.Stats.BlocksDug * 2
	score += p.Stats.GemsFound
	score += p.Stats.BombsFlagged * 10
	score -= p.Stats.WrongFlags * 5
	PostMenu.ItemMap["blocks_s"].SetText(fmt.Sprintf("%d x  2", p.Stats.BlocksDug))
	PostMenu.ItemMap["gem_count_s"].SetText(fmt.Sprintf("%d x  1", p.Stats.GemsFound))
	PostMenu.ItemMap["bombs_flagged_s"].SetText(fmt.Sprintf("%d x 10", p.Stats.BombsFlagged))
	PostMenu.ItemMap["wrong_flags_s"].SetText(fmt.Sprintf("%d x -5", p.Stats.WrongFlags))
	PostMenu.ItemMap["total_score_s"].SetText(fmt.Sprintf("%d", score))
	s.ScoreTimer = timing.New(5.)
	OpenMenu(PostMenu)
	sfx.MusicPlayer.PlayMusic("pause")
	done <- struct{}{}
}

func (s *scoreState) Update(win *pixelgl.Window) {
	reanimator.Update()
	systems.VFXSystem()
	descent.Update()
	hud.UpdateHUD()
	UpdateMenus(win)
	if MenuClosed() {
		SwitchState(MenuStateKey)
	}
}

func (s *scoreState) Draw(win *pixelgl.Window) {
	descent.Descent.GetCave().Draw(win)
	//descent.Descent.GetPlayer().Draw(win, data.GameInput)
	systems.DrawSystem()
	img.DrawBatches(win)
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