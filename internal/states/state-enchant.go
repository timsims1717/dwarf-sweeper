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
	"github.com/faiface/pixel/pixelgl"
)

type enchantState struct {
	*state.AbstractState
}

func (s *enchantState) Unload() {
	sfx.MusicPlayer.Stop("pause")
}

func (s *enchantState) Load(done chan struct{}) {
	success := FillEnchantMenu()
	if !success {
		ClearEnchantMenu()
	} else {
		OpenMenu(EnchantMenu)
		sfx.MusicPlayer.PlayMusic("pause")
	}
	for _, d := range descent.Descent.GetPlayers() {
		d.Player.Stats.AddStats()
	}
	player.OverallStats.AddStats()
	done <- struct{}{}
}

func (s *enchantState) Update(win *pixelgl.Window) {
	reanimator.Update()
	descent.Update()
	hud.UpdateHUD()
	UpdateMenus(win)
	if MenuClosed() {
		ClearEnchantMenu()
		SwitchState(DescentStateKey)
	}
}

func (s *enchantState) Draw(win *pixelgl.Window) {
	descent.Descent.GetCave().Draw(win)
	systems.DrawSystem()
	img.Draw(win)
	hud.DrawHUD(win)
	for _, m := range menuStack {
		m.Draw(win)
	}
}

func (s *enchantState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}