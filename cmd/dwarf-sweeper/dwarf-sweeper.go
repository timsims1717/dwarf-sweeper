package main

import (
	"dwarf-sweeper/internal/cave"
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/debug"
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
		win.Clear(colornames.Black)

		cave.CurrCave.Draw(win)
		//imd.Clear()
		//imd.Color = colornames.Blue
		//imd.EndShape = imdraw.SharpEndShape
		//imd.Push(pixel.V(0., 0.), in.World)
		//imd.Polygon(8.)
		//imd.Draw(win)
		//worlds.Clear()
		//fmt.Fprintf(worlds, "World (X,Y): (%d,%d)", int(in.World.X), int(in.World.Y))
		//worlds.Draw(win, camera.Cam.UITransform(pixel.V(camera.WindowWidthF * 0.5, camera.WindowHeightF * 0.5), pixel.V(1., 1.), 0.))
		//chunks.Clear()
		//ch := cave.WorldToChunk(in.World)
		//fmt.Fprintf(chunks, "Chunk (X,Y): (%d,%d)", ch.X, ch.Y)
		//chunks.Draw(win, camera.Cam.UITransform(pixel.V(camera.WindowWidthF * 0.5, camera.WindowHeightF * 0.5), pixel.V(1., 1.), 0.))
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
