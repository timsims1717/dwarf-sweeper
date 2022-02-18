package states

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/typeface"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type loadingState struct {
	CurrPerc float64
	LoadText *typeface.Text
	Mine     *pixel.Sprite
	Trans    *transform.Transform
}

func (s *loadingState) Load() {
	mine, err := img.LoadImage("assets/img/mine.png")
	if err != nil {
		panic(err)
	}
	s.Mine = pixel.NewSprite(mine, mine.Bounds())
	s.LoadText = typeface.New(camera.Cam, "main", typeface.Alignment{ H: typeface.Center, V: typeface.Center }, 1., constants.ActualMenuSize, 0., 0.)
	s.LoadText.SetColor(menus.DefaultColor)
	s.LoadText.SetText("loading...")
	s.LoadText.Transform.Pos = pixel.V(s.LoadText.Text.BoundsOf("...").W()*s.LoadText.RelativeSize, -15.)
	s.Trans = transform.New()
}

func (s *loadingState) Update(_ *pixelgl.Window) {
	if currState == LoadingStateKey {
		SwitchState(MenuStateKey)
	}
	s.CurrPerc = States[currState].LoadPrc
	s.Trans.UIZoom = camera.Cam.GetZoomScale()
	s.Trans.UIPos = camera.Cam.APos
	s.Trans.Update()
	s.LoadText.Update()
}

func (s *loadingState) Draw(win *pixelgl.Window) {
	s.Mine.Draw(win, s.Trans.Mat)
	s.LoadText.Draw(win)
}