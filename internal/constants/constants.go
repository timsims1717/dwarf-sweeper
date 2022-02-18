package constants

import (
	"github.com/faiface/pixel"
)

const (
	Title   = "DwarfSweeper"
	Release = 0
	Version = 1
	Build   = 20220114
	BaseW   = 320.
	BaseH   = 240.

	// Tile Constants
	TileSize     = 16.0
	DrawDistance = 24.0 * TileSize

	// Batcher Keys
	FogKey        = "fog"
	DwarfKey      = "dwarf"
	EntityKey     = "entities"
	BigEntityKey  = "big_entities"
	TileEntityKey = "tile_entities"
	MenuSprites   = "menu_sprites"
	ParticleKey   = "particle_sprites"
	CaveBGKey     = "cave_bg"
	CaveKey       = "cave"
	TileLayerKey  = "tile_layer"
	PuzzleKey     = "puzzle"

	// Music Keys
	GameMusic  = "gameMusic"
	PauseMusic = "pauseMusic"
	MenuMusic  = "menuMusic"

	// Config
	LinuxDir = "/.local/share/DwarfSweeper"
	WinDir   = "/Documents/My Games/DwarfSweeper"
	MacDir   = "/Library/Application Support/DwarfSweeper"

	// Descent Constants
	ChunkSize = 32
	ChunkArea = ChunkSize * ChunkSize

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
	Resolutions      = []pixel.Vec{
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
	ActualW  = 320.

	// Accessibility
	BaseMenuSize    = 1.4
	BaseHoverSize   = 1.45
	BaseHintSize    = 0.8
	TypeFaceSize    = 200.
	ActualMenuSize  = BaseMenuSize * (10 / TypeFaceSize)
	ActualHoverSize = BaseHoverSize * (10 / TypeFaceSize)
	ActualHintSize  = BaseHintSize * (10 / TypeFaceSize)
	ActualOneSize   = 10 / TypeFaceSize

	// Input
	AimDedicated = true
	DigOnRelease = true
)
