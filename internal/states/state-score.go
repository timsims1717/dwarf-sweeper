package states

import (
	"dwarf-sweeper/internal/descent"
	player2 "dwarf-sweeper/internal/descent/player"
	"dwarf-sweeper/internal/player"
	"dwarf-sweeper/internal/systems"
	"dwarf-sweeper/internal/vfx"
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
	ScoreTimer *timing.FrameTimer
}

func (s *scoreState) Unload() {
	sfx.MusicPlayer.Stop("pause")
}

func (s *scoreState) Load(done chan struct{}) {
	player2.AddStats()
	score := 0
	score += player2.BlocksDug * 2
	score += player2.GemsFound
	score += player2.BombsFlagged * 10
	score -= player2.WrongFlags * 5
	PostMenu.ItemMap["blocks_s"].SetText(fmt.Sprintf("%d x  2", player2.BlocksDug))
	PostMenu.ItemMap["gem_count_s"].SetText(fmt.Sprintf("%d x  1", player2.GemsFound))
	PostMenu.ItemMap["bombs_flagged_s"].SetText(fmt.Sprintf("%d x 10", player2.BombsFlagged))
	PostMenu.ItemMap["wrong_flags_s"].SetText(fmt.Sprintf("%d x -5", player2.WrongFlags))
	PostMenu.ItemMap["total_score_s"].SetText(fmt.Sprintf("%d", score))
	s.ScoreTimer = timing.New(5.)
	OpenMenu(PostMenu)
	sfx.MusicPlayer.PlayMusic("pause")
	done <- struct{}{}
}

func (s *scoreState) Update(win *pixelgl.Window) {
	reanimator.Update()
	systems.VFXSystem()
	vfx.Update()
	descent.Update()
	player.UpdateHUD()
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
	vfx.Draw(win)
	player.DrawHUD(win)
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