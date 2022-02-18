package main

import (
	"dwarf-sweeper/internal/config"
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/credits"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/load"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/states"
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
	mainFont, err := typeface.LoadTTF("assets/FR73PixD.ttf", constants.TypeFaceSize)
	typeface.Atlases["main"] = text.NewAtlas(mainFont, text.ASCII)
	typeface.Atlases["basic"] = typeface.BasicAtlas

	camera.Cam = camera.New(true)
	camera.Cam.Opt.WindowScale = constants.BaseH
	camera.Cam.SetZoom(4. / 3.)
	camera.Cam.SetILock(true)
	camera.Cam.SetSize(res.X, res.Y)
	ratio := res.X/res.Y
	constants.ActualW = constants.BaseH * ratio

	debug.Initialize()
	credits.Initialize()

	vfx.Initialize()

	splash, err := img.LoadImage("assets/img/splash.png")
	if err != nil {
		panic(err)
	}
	states.MenuState.Splash = pixel.NewSprite(splash, splash.Bounds())

	title, err := img.LoadImage("assets/img/title.png")
	if err != nil {
		panic(err)
	}
	states.MenuState.Title = pixel.NewSprite(title, title.Bounds())

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
	puzzleSheet, err := img.LoadSpriteSheet("assets/img/puzzles.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.PuzzleKey, puzzleSheet, false, true)
	menuSheet, err := img.LoadSpriteSheet("assets/img/menu.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.MenuSprites, menuSheet, false, true)

	menus.Initialize()
	states.InitializeMenus(win)
	states.LoadingState.Load()

	load.SFX()
	load.Symbols()
	load.Music()

	timing.Reset()
	win.Show()
	for !win.Closed() {
		timing.Update()
		debug.Clear()
		states.Update(win)

		win.Clear(color.RGBA{
			R: 6,
			G: 6,
			B: 8,
			A: 255,
		})

		states.Draw(win)
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
			newRatio := res.X / res.Y
			if constants.FullScreen {
				x, y := picked.Size()
				newRatio = x / y
			}
			constants.ActualW = constants.BaseH * newRatio
		}
	}
}

// fire the run function (the real main function)
func main() {
	pixelgl.Run(run)
}
