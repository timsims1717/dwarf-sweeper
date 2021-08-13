package main

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/state"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
)

func run() {
	world.SetTileSize(cfg.TileSize)
	config := pixelgl.WindowConfig{
		Title:  cfg.Title,
		Bounds: pixel.R(0, 0, 1600, 900),
		VSync: true,
		Invisible: true,
	}
	win, err := pixelgl.NewWindow(config)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(false)

	camera.Cam = camera.New(true)
	camera.Cam.Opt.WindowScale = cfg.BaseH
	camera.Cam.SetZoom(4. / 3.)
	camera.Cam.SetILock(true)
	camera.Cam.SetSize(1600/900, cfg.BaseH)

	debug.Initialize()
	state.InitializeMenus()

	vfx.Initialize()
	particles.Initialize()
	sheet, err := img.LoadSpriteSheet("assets/img/entities.json")
	if err != nil {
		panic(err)
	}
	img.Batchers["entities"] = img.NewBatcher(sheet)

	sfx.SoundPlayer.RegisterSound("assets/sound/click.wav", "click")
	//sfx.SoundPlayer.RegisterSound("assets/sound/impact1.wav", "impact1")
	//sfx.SoundPlayer.RegisterSound("assets/sound/impact2.wav", "impact2")
	//sfx.SoundPlayer.RegisterSound("assets/sound/impact3.wav", "impact3")
	//sfx.SoundPlayer.RegisterSound("assets/sound/impact4.wav", "impact4")
	sfx.SoundPlayer.RegisterSound("assets/sound/rocks1.wav", "rocks1")
	sfx.SoundPlayer.RegisterSound("assets/sound/rocks2.wav", "rocks2")
	sfx.SoundPlayer.RegisterSound("assets/sound/rocks3.wav", "rocks3")
	sfx.SoundPlayer.RegisterSound("assets/sound/rocks4.wav", "rocks4")
	sfx.SoundPlayer.RegisterSound("assets/sound/rocks5.wav", "rocks5")
	sfx.SoundPlayer.RegisterSound("assets/sound/shovel.wav", "shovel")
	sfx.SoundPlayer.RegisterSound("assets/sound/step1.wav", "step1")
	sfx.SoundPlayer.RegisterSound("assets/sound/step2.wav", "step2")
	sfx.SoundPlayer.RegisterSound("assets/sound/step3.wav", "step3")
	sfx.SoundPlayer.RegisterSound("assets/sound/step4.wav", "step4")
	sfx.SoundPlayer.RegisterSound("assets/sound/clink.wav", "clink")
	sfx.SetMasterVolume(75)
	sfx.SetSoundVolume(75)

	timing.Reset()
	win.Show()
	for !win.Closed() {
		timing.Update()
		debug.Clear()
		state.Update(win)

		win.Clear(color.RGBA{
			R: 6,
			G: 6,
			B: 8,
			A: 255,
		})

		state.Draw(win)
		debug.Draw(win)
		win.Update()
		if cfg.ChangeScreenSize {
			cfg.ChangeScreenSize = false
			if (cfg.FullScreen && win.Monitor() == nil) || (!cfg.FullScreen && win.Monitor() != nil) {
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
					if picked == nil {
						picked = pixelgl.PrimaryMonitor()
					}
				}
				if cfg.FullScreen {
					win.SetMonitor(picked)
				} else {
					win.SetMonitor(nil)
				}
			}
			res := cfg.Resolutions[cfg.ResIndex]
			win.SetBounds(pixel.R(0., 0., res.X, res.Y))
			camera.Cam.SetSize(res.X / res.Y, res.Y)
		}
	}
}

// fire the run function (the real main function)
func main() {
	pixelgl.Run(run)
}
