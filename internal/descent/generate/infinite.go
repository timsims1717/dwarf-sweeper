package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/world"
)

func NewInfiniteCave(spriteSheet *img.SpriteSheet, biome string) *cave.Cave {
	random.RandCaveSeed()
	batcher := img.NewBatcher(spriteSheet, false)
	newCave := cave.NewCave(batcher, biome, false)
	newCave.FillChunk = FillChunk
	newCave.StartC = world.Coords{X: 16, Y: 9}
	newCave.GemRate = constants.BaseGem
	newCave.ItemRate = constants.BaseItem
	newCave.BombPMin = 0.2
	newCave.BombPMax = 0.3
	chunk0 := cave.NewChunk(world.Coords{X: 0, Y: 0}, newCave, cave.BlockCollapse)
	FillChunk(chunk0)

	chunkr1 := cave.NewChunk(world.Coords{X: 1, Y: 0}, newCave, cave.BlockCollapse)
	chunkr2 := cave.NewChunk(world.Coords{X: 1, Y: 1}, newCave, cave.BlockCollapse)
	chunkr3 := cave.NewChunk(world.Coords{X: 0, Y: 1}, newCave, cave.BlockCollapse)
	FillChunk(chunkr1)
	FillChunk(chunkr2)
	FillChunk(chunkr3)

	chunkl1 := cave.NewChunk(world.Coords{X: -1, Y: 0}, newCave, cave.BlockCollapse)
	chunkl2 := cave.NewChunk(world.Coords{X: -1, Y: 1}, newCave, cave.BlockCollapse)
	FillChunk(chunkl1)
	FillChunk(chunkl2)

	newCave.RChunks[chunk0.Coords] = chunk0
	newCave.RChunks[chunkr1.Coords] = chunkr1
	newCave.RChunks[chunkr2.Coords] = chunkr2
	newCave.RChunks[chunkr3.Coords] = chunkr3

	newCave.LChunks[chunkl1.Coords] = chunkl1
	newCave.LChunks[chunkl2.Coords] = chunkl2
	Entrance(newCave, world.Coords{X: 16, Y: 9}, 9, 5, 3, false)
	return newCave
}
