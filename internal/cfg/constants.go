package cfg

import "github.com/faiface/pixel"

const (
	Title = "Dwarf Sweeper"
	BaseW = 320.
	BaseH = 240.

	// Tile Constants
	TileSize = 16.0
)

var (
	FullScreen = false
	ChangeScreenSize = false
	Resolutions = []pixel.Vec{
		pixel.V(800, 600),
		pixel.V(1280, 960),
		pixel.V(1600, 900),
		pixel.V(1920, 1080),
	}
	ResStrings = []string{
		" 800 x 600",
		" 1280 x 960",
		" 1600 x 900",
		" 1920 x 1080",
	}
	ResIndex = 2
)