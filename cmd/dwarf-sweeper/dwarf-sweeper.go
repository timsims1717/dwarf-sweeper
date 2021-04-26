package main

import (
	"dwarf-sweeper/internal/cave"
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/state"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
)

func run() {
	//seed := int64(1619305348812219488)
	seed := time.Now().UnixNano()
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

	vfx.Initialize()
	particles.Initialize()
	cave.Entities.Initialize()

	sfx.SoundPlayer.RegisterSound("assets/sound/click.wav", "click")
	sfx.SoundPlayer.RegisterSound("assets/sound/impact1.wav", "impact1")
	sfx.SoundPlayer.RegisterSound("assets/sound/impact2.wav", "impact2")
	sfx.SoundPlayer.RegisterSound("assets/sound/impact3.wav", "impact3")
	sfx.SoundPlayer.RegisterSound("assets/sound/impact4.wav", "impact4")
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
	sfx.SetMasterVolume(25)

	timing.Reset()
	for !win.Closed() {
		timing.Update()
		debug.Clear()
		state.Update(win)

		win.Clear(colornames.Black)

		state.Draw(win)
		win.Update()
	}
}

// fire the run function (the real main function)
func main() {
	pixelgl.Run(run)
}
