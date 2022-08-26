package constants

import (
	"github.com/faiface/pixel"
	"image/color"
)

type Biome int

const (
	Title   = "DwarfSweeper"
	Release = 0
	Version = 4
	Build   = 20220729
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
	BigExpKey     = "big_explosions"
	HugeExpKey    = "huge_explosions"
	ExpKey        = "explosions"
	CaveBGFMT     = "%s_bg"
	TileLayerKey  = "tile_layer"
	PuzzleKey     = "puzzle"

	// Music Keys
	GameMusic  = "gameMusic"
	PauseMusic = "pauseMusic"
	MenuMusic  = "menuMusic"

	// Directories
	LinuxDir   = "/.local/share/DwarfSweeper"
	WinDir     = "/Documents/My Games/DwarfSweeper"
	MacDir     = "/Library/Application Support/DwarfSweeper"
	ImgBiomeBG = "assets/img/biome-bg.json"
	ImgBiomeFG = "assets/img/biome-fg.json"
	ImgCave    = "the-%s.png"

	// Descent Constants
	ChunkSize = 16
	ChunkArea = ChunkSize * ChunkSize

	// Biomes
	Mine = iota
	Moss
	Glacier
	Crystal
	Dark
)

var (
	// Config
	System     string
	HomeDir    string
	ConfigDir  string
	ConfigFile string

	// Graphics
	FullScreen   = false
	VSync        = true
	ChangeScreen = false
	Resolutions  = []pixel.Vec{
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
	BGColor  = color.RGBA{
		R: 6,
		G: 6,
		B: 8,
		A: 255,
	}

	// Audio
	MuteOnUnfocused = false

	// Accessibility
	BaseMenuSize    = 1.0
	BaseHoverSize   = 1.1
	BaseHintSize    = 0.8
	TypeFaceSize    = 200.
	ActualMenuSize  = BaseMenuSize * (10 / TypeFaceSize)
	ActualHoverSize = BaseHoverSize * (10 / TypeFaceSize)
	ActualHintSize  = BaseHintSize * (10 / TypeFaceSize)
	ActualOneSize   = 10 / TypeFaceSize

	// Gameplay
	ScreenShake = true
	SplitScreenV = true

	// Menus
	DefaultColor = color.RGBA{
		R: 74,
		G: 84,
		B: 98,
		A: 255,
	}
	HoverColor = color.RGBA{
		R: 20,
		G: 52,
		B: 100,
		A: 255,
	}
	DisabledColor = color.RGBA{
		R: 109,
		G: 117,
		B: 141,
		A: 255,
	}

	// Biomes
	Biomes = []Biome{
		Mine,
		Moss,
		Glacier,
		Crystal,
		Dark,
	}
)

func (b Biome) Key() string {
	switch b {
	case Mine:
		return "mine"
	case Moss:
		return "moss"
	case Glacier:
		return "glacier"
	case Crystal:
		return "crystal"
	case Dark:
		return "dark"
	}
	return "error"
}

func (b Biome) String() string {
	switch b {
	case Mine:
		return "Mine"
	case Moss:
		return "Moss"
	case Glacier:
		return "Glacier"
	case Crystal:
		return "Crystal"
	case Dark:
		return "Dark"
	}
	return "Unknown"
}