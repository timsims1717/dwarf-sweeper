package main

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"image/color"
)

func run() {
	world.SetTileSize(constants.TileSize)
	conf := pixelgl.WindowConfig{
		Title:     constants.Title,
		Bounds:    pixel.R(0, 0, 1600, 900),
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(conf)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(false)

	camera.Cam = camera.New(true)
	camera.Cam.Opt.WindowScale = constants.BaseH
	camera.Cam.SetZoom(4. / 3.)
	camera.Cam.SetILock(true)
	camera.Cam.SetSize(1600, 900)

	mainFont, err := typeface.LoadTTF("assets/FR73PixD.ttf", constants.TypeFaceSize)
	typeface.Atlases["main"] = text.NewAtlas(mainFont, text.ASCII)

	debug.Initialize()

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

	caveSheet, err := img.LoadSpriteSheet("assets/img/the-dark.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.CaveKey, caveSheet, true, false)

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

	sfx.SoundPlayer.RegisterSound("assets/sound/click.wav", "click")

	menus.Initialize()
	initMenu(win)
	CaveMenu.Open()
	menuStack = append(menuStack, CaveMenu)

	timing.Reset()
	debug.Debug = true

	for !win.Closed() {
		timing.Update()
		debug.Clear()

		// update
		theInput.Update(win)

		if theInput.Get("camIn").JustPressed() {
			camera.Cam.ZoomIn(1.)
			camera.Cam.Opt.ScrollSpeed /= 1.2
		}
		if theInput.Get("camOut").JustPressed() {
			camera.Cam.ZoomIn(-1.)
			camera.Cam.Opt.ScrollSpeed *= 1.2
		}
		if theInput.Get("camUp").Pressed() {
			camera.Cam.Up()
		} else if theInput.Get("camDown").Pressed() {
			camera.Cam.Down()
		}
		if theInput.Get("camRight").Pressed() {
			camera.Cam.Right()
		} else if theInput.Get("camLeft").Pressed() {
			camera.Cam.Left()
		}
		if theInput.Get("center").JustPressed() && theCave != nil {
			dimX, dimY := theCave.Dimensions()
			bl, tr := theCave.CurrentBoundaries()
			camera.Cam.Pos = theCave.GetTile(pixel.V(bl.X + float64(dimX) * world.TileSize * 0.5, tr.Y - float64(dimY) * world.TileSize * 0.5)).Transform.Pos
		}
		next = theInput.Get("click").JustPressed() || theInput.Get("menuSelect").JustPressed() || theInput.Get("menuBack").JustPressed()
		updateMenu(win)
		if signal != nil && caveBuild != nil && theCave != nil && next {
			done := <-signal
			if done {
				signal = nil
				caveBuild = nil
				theCave = nil
				CaveMenu.Open()
				menuStack = append(menuStack, CaveMenu)
				camera.Cam.SetZoom(4. / 3.)
				camera.Cam.Opt.ScrollSpeed = 40.
				camera.Cam.SnapTo(pixel.ZV)
			} else {
				signal <- !theInput.Get("menuBack").JustPressed()
			}
		} else if signal == nil && caveBuild != nil {
			signal = make(chan bool)
			theCave = generate.NewAsyncCave(caveBuild, 1, signal)
			<-signal
			theCave.Fog = false
			theCave.LoadAll = true
		} else if theCave != nil {
			theCave.UpdateBatch = true
			theCave.Pivot = camera.Cam.Pos
			theCave.Update()
		}
		camera.Cam.Update(win)

		win.Clear(color.RGBA{
			R: 6,
			G: 6,
			B: 8,
			A: 255,
		})

		// draw
		img.ClearBatches()
		if theCave != nil {
			theCave.Draw(win)
		}
		img.DrawBatches(win)
		for _, m := range menuStack {
			m.Draw(win)
		}

		debug.Draw(win)
		win.Update()
	}
}

var (
	CaveMenu *menus.DwarfMenu
	focused  bool
	theInput = &input.Input{
		Buttons: map[string]*input.ButtonSet{
			"menuUp":     input.NewJoyless(pixelgl.KeyUp),
			"menuDown":   input.NewJoyless(pixelgl.KeyDown),
			"menuRight":  input.NewJoyless(pixelgl.KeyRight),
			"menuLeft":   input.NewJoyless(pixelgl.KeyLeft),
			"menuSelect": input.NewJoyless(pixelgl.KeySpace),
			"menuBack":   input.NewJoyless(pixelgl.KeyEscape),
			"click":      input.NewJoyless(pixelgl.MouseButtonLeft),
			"camUp":      input.NewJoyless(pixelgl.KeyUp),
			"camRight":   input.NewJoyless(pixelgl.KeyRight),
			"camDown":    input.NewJoyless(pixelgl.KeyDown),
			"camLeft":    input.NewJoyless(pixelgl.KeyLeft),
			"camIn":      input.NewJoyless(pixelgl.KeyKPAdd),
			"camOut":     input.NewJoyless(pixelgl.KeyKPSubtract),
			"center":     input.NewJoyless(pixelgl.KeyC),
		},
		Mode: input.KeyboardMouse,
	}
	menuStack []*menus.DwarfMenu
	signal    chan bool
	next      bool
	caveBuild *builder.CaveBuilder
	theCave   *cave.Cave
)

func initMenu(win *pixelgl.Window) {
	CaveMenu = menus.New("caves", camera.Cam)
	CaveMenu.Title = true
	CaveMenu.SetCloseFn(func() {
		if !CaveMenu.IsClosed() {
			win.SetClosed(true)
		}
	})
	gen1 := CaveMenu.AddItem("gen1", "Generate")
	quit := CaveMenu.AddItem("quit", "Quit")

	gen1.SetClickFn(func() {
		caveBuilders, err := builder.LoadBuilder(fmt.Sprint("assets/test.json"))
		if err != nil {
			panic(err)
		}
		choice := random.Effects.Intn(len(caveBuilders))
		caveBuild = caveBuilders[choice]
		CaveMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	quit.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		win.SetClosed(true)
	})
}

func updateMenu(win *pixelgl.Window) {
	if win.Focused() && !focused {
		focused = true
	} else if !win.Focused() && focused {
		focused = false
	}
	for i, me := range menuStack {
		if i == len(menuStack)-1 {
			if !win.Focused() {
				me.UnhoverAll()
			}
			me.Update(theInput)
			if me.IsClosed() {
				if len(menuStack) > 1 {
					menuStack = menuStack[:len(menuStack)-1]
				} else {
					menuStack = []*menus.DwarfMenu{}
				}
			}
		} else {
			me.Update(nil)
		}
	}
}

// fire the run function (the real main function)
func main() {
	pixelgl.Run(run)
}