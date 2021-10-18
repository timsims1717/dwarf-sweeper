package constants

import (
	"github.com/faiface/pixel"
)

const (
	Title = "DwarfSweeper"
	BaseW = 320.
	BaseH = 240.

	// Tile Constants
	TileSize = 16.0
	DrawDistance = 24.0 * TileSize

	// Batcher Keys
	EntityKey    = "entities"
	BigEntityKey = "big_entities"
	MenuSprites  = "menu_sprites"

	// Config
	LinuxDir = "/.local/share/DwarfSweeper"
	WinDir   = "/Documents/My Games/DwarfSweeper"
	MacDir   = "/Library/Application Support/DwarfSweeper"

	// Descent Constants

	ChunkSize = 32
	ChunkArea = ChunkSize * ChunkSize

	BaseFuse = 1.5
	MineFuse = 0.2
	BaseGem  = 20
	BaseItem = 50
)

var (
	// Config
	System     string
	HomeDir    string
	ConfigDir  string
	ConfigFile string

	// Graphics
	FullScreen       = false
	VSync            = true
	ChangeScreenSize = false
	Resolutions = []pixel.Vec{
		pixel.V(800, 600),
		pixel.V(1280, 960),
		pixel.V(1600, 900),
		pixel.V(1920, 1080),
	}
	ResStrings = []string{
		"800x600",
		"1280x960",
		"1600x900",
		"1920x1080",
	}
	ResIndex = 2

	// Input
	AimDedicated = true
	DigOnRelease = true
)