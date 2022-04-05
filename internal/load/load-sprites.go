package load

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/img"
	"fmt"
)

func Sprites() {
	// Cave Backgrounds
	bgDarkSheet, err := img.LoadSpriteSheet("assets/img/the-dark-bg.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(fmt.Sprintf(constants.CaveBGFMT, "dark"), bgDarkSheet, true, false)
	bgMineSheet, err := img.LoadSpriteSheet("assets/img/the-mine-bg.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(fmt.Sprintf(constants.CaveBGFMT, "mine"), bgMineSheet, true, false)
	bgMossSheet, err := img.LoadSpriteSheet("assets/img/the-moss-bg.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(fmt.Sprintf(constants.CaveBGFMT, "moss"), bgMossSheet, true, false)
	bgCrystalSheet, err := img.LoadSpriteSheet("assets/img/the-crystal-bg.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher(fmt.Sprintf(constants.CaveBGFMT, "crystal"), bgCrystalSheet, true, false)

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
	darkSheet, err := img.LoadSpriteSheet("assets/img/the-dark.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher("dark", darkSheet, true, false)
	mineSheet, err := img.LoadSpriteSheet("assets/img/the-mine.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher("mine", mineSheet, true, false)
	mossSheet, err := img.LoadSpriteSheet("assets/img/the-moss.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher("moss", mossSheet, true, false)
	crystalSheet, err := img.LoadSpriteSheet("assets/img/the-crystal.json")
	if err != nil {
		panic(err)
	}
	img.AddBatcher("crystal", crystalSheet, true, false)

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