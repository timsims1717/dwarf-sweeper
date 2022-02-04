package generate

import (
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/world"
)

func NewInfiniteCave(biome string) *cave.Cave {
	random.RandCaveSeed()
	c := cave.NewCave(biome, cave.Infinite)
	c.FillChunk = structures.FillBasic
	c.StartC = world.Coords{X: 16, Y: 9}
	c.BombPMin = 0.2
	c.BombPMax = 0.3
	chunk0 := cave.NewChunk(world.Coords{X: 0, Y: 0}, c, cave.BlockCollapse)
	structures.FillBasic(chunk0)

	chunkr1 := cave.NewChunk(world.Coords{X: 1, Y: 0}, c, cave.BlockCollapse)
	chunkr2 := cave.NewChunk(world.Coords{X: 1, Y: 1}, c, cave.BlockCollapse)
	chunkr3 := cave.NewChunk(world.Coords{X: 0, Y: 1}, c, cave.BlockCollapse)
	structures.FillBasic(chunkr1)
	structures.FillBasic(chunkr2)
	structures.FillBasic(chunkr3)

	chunkl1 := cave.NewChunk(world.Coords{X: -1, Y: 0}, c, cave.BlockCollapse)
	chunkl2 := cave.NewChunk(world.Coords{X: -1, Y: 1}, c, cave.BlockCollapse)
	structures.FillBasic(chunkl1)
	structures.FillBasic(chunkl2)

	c.Chunks[chunk0.Coords] = chunk0
	c.Chunks[chunkr1.Coords] = chunkr1
	c.Chunks[chunkr2.Coords] = chunkr2
	c.Chunks[chunkr3.Coords] = chunkr3
	c.Chunks[chunkl1.Coords] = chunkl1
	c.Chunks[chunkl2.Coords] = chunkl2
	structures.Entrance(c, world.Coords{X: 16, Y: 9}, 9, 5, 3, false)
	return c
}
