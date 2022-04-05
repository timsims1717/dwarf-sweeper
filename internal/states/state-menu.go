package states

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/credits"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/state"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	pressAKeySec = 1.0
)

type menuState struct {
	*state.AbstractState
	First          bool
	Splash         *pixel.Sprite
	splashTran     *transform.Transform
	splashScale    float64
	Title          *pixel.Sprite
	titleTran      *transform.Transform
	titleScale     float64
	titleY         float64
	pressAKey      *typeface.Text
	pressAKeyTimer *timing.Timer
}

func (s *menuState) Unload() {
	sfx.MusicPlayer.Stop("menu")
}

func (s *menuState) Load(done chan struct{}) {
	sfx.SoundPlayer.KillAll()
	sfx.MusicPlayer.StopAllMusic()
	descent.Descent.FreeCam = false
	camera.Cam.SetZoom(1.)
	s.pressAKey = typeface.New(&camera.Cam.APos, "main", typeface.Alignment{ H: typeface.Center, V: typeface.Center }, 1.0, constants.ActualMenuSize, 0., 0.)
	s.pressAKey.SetColor(colornames.Aliceblue)
	s.pressAKey.SetText("press any key")
	s.pressAKey.Transform.Pos = pixel.V(0., -75.)
	s.pressAKey.NoShow = true
	s.pressAKeyTimer = timing.New(2.5)
	s.titleY = 70.
	s.titleTran = transform.New()
	s.titleTran.Pos = pixel.V(0., s.titleY)
	s.titleScale = 0.4
	s.splashTran = transform.New()
	s.splashScale = 0.4
	camera.Cam.SnapTo(pixel.ZV)
	if !s.First {
		s.First = true
	} else {
		OpenMenu(MainMenu)
	}
	sfx.MusicPlayer.PlayMusic("menu")
	done <- struct{}{}
}

func (s *menuState) Update(win *pixelgl.Window) {
	s.pressAKey.Update()
	if s.pressAKeyTimer.UpdateDone() {
		s.pressAKey.NoShow = !s.pressAKey.NoShow
		s.pressAKeyTimer = timing.New(pressAKeySec)
	}
	s.titleTran.Scalar = pixel.V(s.titleScale, s.titleScale)
	s.titleTran.UIPos = camera.Cam.APos
	s.titleTran.UIZoom = camera.Cam.GetZoomScale()
	s.titleTran.Update()
	w := win.Bounds().W()
	h := win.Bounds().H()
	ratio := w/h
	w = s.Splash.Frame().W()
	h = s.Splash.Frame().H()
	ratioSplash := w/h
	if ratio > ratioSplash && ratio - ratioSplash < 0.5 {
		s.splashScale = 0.4 * (1. + ratio - ratioSplash)
	} else {
		s.splashScale = 0.4
	}
	s.splashTran.Scalar = pixel.V(s.splashScale, s.splashScale)
	s.splashTran.UIPos = camera.Cam.APos
	s.splashTran.UIZoom = camera.Cam.GetZoomScale()
	s.splashTran.Update()
	if credits.Opened() {
		credits.Update()
		if pressed, _ := menuInput.AnyJustPressed(true); pressed {
			credits.Close()
		}
	} else {
		UpdateMenus(win)
		pressed, mode := menuInput.AnyJustPressed(true)
		if MenuClosed() && pressed {
			OpenMenu(MainMenu)
			data.GameInputP1.Mode = mode
		}
	}
	//debug.AddText(fmt.Sprintf("Input TLines: %d", InputMenu.TLines))
	//debug.AddText(fmt.Sprintf("Input Top: %d", InputMenu.Top))
	//debug.AddText(fmt.Sprintf("Input Curr: %d", InputMenu.Items[InputMenu.Hovered].CurrLine))
}

func (s *menuState) Draw(win *pixelgl.Window) {
	s.Splash.Draw(win, s.splashTran.Mat)
	s.Title.Draw(win, s.titleTran.Mat)
	for _, m := range menuStack {
		m.Draw(win)
	}
	if credits.Opened() {
		credits.Draw(win)
	} else if MenuClosed() {
		s.pressAKey.Draw(win)
	}
}

func (s *menuState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}