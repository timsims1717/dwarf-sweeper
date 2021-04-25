package main

import (
	"dwarf-sweeper/internal/cave"
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/dwarf"
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
)

func run() {
	seed := int64(1619305348812219488)
	//seed := time.Now().UnixNano()
	rand.Seed(seed)
	fmt.Println("Seed:", seed)
	world.SetTileSize(cfg.TileSize)
	camera.SetWindowSize(1600, 900)
	config := pixelgl.WindowConfig{
		Title:  cfg.Title,
		Bounds: pixel.R(0, 0, camera.WindowWidthF, camera.WindowHeightF),
		//VSync: true,
	}
	win, err := pixelgl.NewWindow(config)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(false)

	camera.Cam = camera.New()
	camera.Cam.SetZoom(4.0)

	debug.Initialize()

	in := input.NewInput()
	//worlds := text.New(pixel.ZV, typeface.BasicAtlas)
	//chunks := text.New(pixel.ZV, typeface.BasicAtlas)

	vfx.Initialize()
	particles.Initialize()

	sheet, err := img.LoadSpriteSheet("assets/img/test-tiles.json")
	if err != nil {
		panic(err)
	}
	cave.CurrCave = cave.NewCave(sheet)
	//imd := imdraw.New(nil)

	player := dwarf.NewDwarf()

	timing.Reset()
	for !win.Closed() {
		timing.Update()

		in.Update(win)
		if in.Debug {
			fmt.Println("DEBUG PAUSE")
		}
		camera.Cam.Update(win)
		cave.CurrCave.Update(in.World, in)
		particles.Update()
		vfx.Update()
		player.Update(in)
		
		win.Clear(colornames.Black)

		cave.CurrCave.Draw(win)
		player.Draw(win)
		particles.Draw(win)
		vfx.Draw(win)
		debug.Draw(win)
		win.Update()
	}
}

// fire the run function (the real main function)
func main() {
	pixelgl.Run(run)
}
