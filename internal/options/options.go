package options

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/sfx"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func Update(win *pixelgl.Window) {
	sfx.MuteMaster(!win.Focused() && constants.MuteOnUnfocused)
	win.SetVSync(constants.VSync)
	if constants.ChangeScreen {
		constants.ChangeScreen = false
		pos := win.GetPos()
		pos.X += win.Bounds().W() * 0.5
		pos.Y += win.Bounds().H() * 0.5
		var picked *pixelgl.Monitor
		if len(pixelgl.Monitors()) > 1 {
			for _, m := range pixelgl.Monitors() {
				x, y := m.Position()
				w, h := m.Size()
				if pos.X >= x && pos.X <= x+w && pos.Y >= y && pos.Y <= y+h {
					picked = m
					break
				}
			}
			if picked == nil {
				pos = win.GetPos()
				for _, m := range pixelgl.Monitors() {
					x, y := m.Position()
					w, h := m.Size()
					if pos.X >= x && pos.X <= x+w && pos.Y >= y && pos.Y <= y+h {
						picked = m
						break
					}
				}
			}
		}
		if picked == nil {
			picked = pixelgl.PrimaryMonitor()
		}
		if constants.FullScreen {
			win.SetMonitor(picked)
		} else {
			win.SetMonitor(nil)
		}
		res := constants.Resolutions[constants.ResIndex]
		win.SetBounds(pixel.R(0., 0., res.X, res.Y))
		camera.Cam.SetSize(res.X, res.Y)
		newRatio := res.X / res.Y
		if constants.FullScreen {
			x, y := picked.Size()
			newRatio = x / y
		}
		constants.ActualW = constants.BaseH * newRatio
	}
}