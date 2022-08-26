package load

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/img"
	"fmt"
)

func Sprites() {
	// Cave Backgrounds
	for _, b := range constants.Biomes {
		bgSheet, err := img.LoadSpriteImg(constants.ImgBiomeBG, fmt.Sprintf(constants.ImgCave, b.Key()))
		if err != nil {
			panic(err)
		}
		img.AddBatcher(fmt.Sprintf(constants.CaveBGFMT, b.Key()), bgSheet, true, false)
	}

	// Entities
	tileEntitySheet, err := img.LoadSpriteSheet("assets/img/tile_entities.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.TileEntityKey, tileEntitySheet, true, true)
	dwarfSheet, err := img.LoadSpriteSheet("assets/img/dwarf.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.DwarfKey, dwarfSheet, true, true)
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

	// Cave Foregrounds
	for _, b := range constants.Biomes {
		sheet, err := img.LoadSpriteImg(constants.ImgBiomeFG, fmt.Sprintf(constants.ImgCave, b.Key()))
		if err != nil {
			panic(err)
		}
		img.AddBatcher(b.Key(), sheet, true, false)
	}

	// Particles/VFX
	partSheet, err := img.LoadSpriteSheet("assets/img/particles.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.ParticleKey, partSheet, true, true)
	hugeExpSheet, err := img.LoadSpriteSheet("assets/img/huge-explosion.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.HugeExpKey, hugeExpSheet, true, true)
	bigExpSheet, err := img.LoadSpriteSheet("assets/img/big-explosion.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.BigExpKey, bigExpSheet, true, true)
	expSheet, err := img.LoadSpriteSheet("assets/img/explosion.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.ExpKey, expSheet, true, true)

	// Fog
	fogSheet, err := img.LoadSpriteSheet("assets/img/fog.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(constants.FogKey, fogSheet, true, false)

	// Puzzles and Menus
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
}