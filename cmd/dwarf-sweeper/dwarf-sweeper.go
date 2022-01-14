package main

import (
	"dwarf-sweeper/internal/config"
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/credits"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/state"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"image/color"
)

func run() {
	world.SetTileSize(constants.TileSize)
	config.LoadConfig()
	res := constants.Resolutions[constants.ResIndex]
	conf := pixelgl.WindowConfig{
		Title:     constants.Title,
		Bounds:    pixel.R(0, 0, res.X, res.Y),
		VSync:     constants.VSync,
		Invisible: true,
	}
	if constants.FullScreen {
		constants.ChangeScreenSize = true
	}
	win, err := pixelgl.NewWindow(conf)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(false)
	mainFont, err := typeface.LoadTTF("assets/FR73PixD.ttf", 100.)
	typeface.Atlases["main"] = text.NewAtlas(mainFont, text.ASCII)

	camera.Cam = camera.New(true)
	camera.Cam.Opt.WindowScale = constants.BaseH
	camera.Cam.SetZoom(4. / 3.)
	camera.Cam.SetILock(true)
	camera.Cam.SetSize(res.X, res.Y)

	debug.Initialize()
	credits.Initialize()

	vfx.Initialize()
	splash, err := img.LoadImage("assets/img/splash.png")
	if err != nil {
		panic(err)
	}
	state.Splash = pixel.NewSprite(splash, splash.Bounds())
	title, err := img.LoadImage("assets/img/title.png")
	if err != nil {
		panic(err)
	}
	state.Title = pixel.NewSprite(title, title.Bounds())
	bgSheet, err := img.LoadSpriteSheet("assets/img/the-dark-bg.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.CaveBGKey, bgSheet, true, false)
	tileEntitySheet, err := img.LoadSpriteSheet("assets/img/tile_entities.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.TileEntityKey, tileEntitySheet, true, true)
	bigEntitySheet, err := img.LoadSpriteSheet("assets/img/big-entities.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.BigEntityKey, bigEntitySheet, true, true)
	entitySheet, err := img.LoadSpriteSheet("assets/img/entities.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.EntityKey, entitySheet, true, true)
	dwarfSheet, err := img.LoadSpriteSheet("assets/img/dwarf.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.DwarfKey, dwarfSheet, true, true)
	caveSheet, err := img.LoadSpriteSheet("assets/img/the-dark.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.CaveKey, caveSheet, true, false)
	partSheet, err := img.LoadSpriteSheet("assets/img/particles.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.ParticleKey, partSheet, true, true)
	fogSheet, err := img.LoadSpriteSheet("assets/img/fog.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.FogKey, fogSheet, true, false)
	menuSheet, err := img.LoadSpriteSheet("assets/img/menu.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.MenuSprites, menuSheet, false, true)

	menus.Initialize()
	state.InitializeMenus(win)
	descent.InitCollectibles()

	sfx.SoundPlayer.RegisterSound("assets/sound/blast1.wav", "blast1")
	sfx.SoundPlayer.RegisterSound("assets/sound/click.wav", "click")
	sfx.SoundPlayer.RegisterSound("assets/sound/emerge.wav", "emerge")
	sfx.SoundPlayer.RegisterSound("assets/sound/gnomestep.wav", "gnomestep")
	//sfx.SoundPlayer.RegisterSound("assets/sound/impact1.wav", "impact1")
	//sfx.SoundPlayer.RegisterSound("assets/sound/impact2.wav", "impact2")
	//sfx.SoundPlayer.RegisterSound("assets/sound/impact3.wav", "impact3")
	//sfx.SoundPlayer.RegisterSound("assets/sound/impact4.wav", "impact4")
	sfx.SoundPlayer.RegisterSound("assets/sound/roar.wav", "roar")
	sfx.SoundPlayer.RegisterSound("assets/sound/rocks1.wav", "rocks1")
	sfx.SoundPlayer.RegisterSound("assets/sound/rocks2.wav", "rocks2")
	sfx.SoundPlayer.RegisterSound("assets/sound/rocks3.wav", "rocks3")
	sfx.SoundPlayer.RegisterSound("assets/sound/rocks4.wav", "rocks4")
	sfx.SoundPlayer.RegisterSound("assets/sound/rocks5.wav", "rocks5")
	sfx.SoundPlayer.RegisterSound("assets/sound/rockslide.wav", "rockslide")
	sfx.SoundPlayer.RegisterSound("assets/sound/shovel.wav", "shovel")
	sfx.SoundPlayer.RegisterSound("assets/sound/step1.wav", "step1")
	sfx.SoundPlayer.RegisterSound("assets/sound/step2.wav", "step2")
	sfx.SoundPlayer.RegisterSound("assets/sound/step3.wav", "step3")
	sfx.SoundPlayer.RegisterSound("assets/sound/step4.wav", "step4")
	sfx.SoundPlayer.RegisterSound("assets/sound/clink.wav", "clink")

	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Crab Nebula.wav", "crab")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Honeybee.wav", "honey")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Prairie Oyster.wav", "prairie")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Sable.wav", "sable")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Strawberry Jam.wav", "strawberry")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/The Dawn Approaches.wav", "dawn")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/The Hero Approaches.wav", "hero")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Twin Turbo.wav", "turbo")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Voyage.wav", "voyage")

	sfx.MusicPlayer.SetTracks("menu", []string{"crab"})
	sfx.MusicPlayer.SetTracks("pause", []string{"sable"})
	sfx.MusicPlayer.NewSet(constants.GameMusic, []string{"honey", "strawberry", "dawn", "hero", "voyage", "prairie"}, sfx.Repeat, 0., 2.)

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
		sfx.MusicPlayer.Update()
		win.Update()
		win.SetVSync(constants.VSync)
		if constants.ChangeScreenSize {
			constants.ChangeScreenSize = false
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
		}
	}
}

// fire the run function (the real main function)
func main() {
	pixelgl.Run(run)
}
