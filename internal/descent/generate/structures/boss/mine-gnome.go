package boss

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

func GnomeBoss(c *cave.Cave, level int) *cave.Cave {
	c.FillChunk = structures.FillChunkWall
	c.StartC = world.Coords{X: 16, Y: 9}
	c.GemRate = constants.BaseGem
	c.ItemRate = constants.BaseItem
	c.BombPMin = 0.2
	c.BombPMax = 0.3
	chunk0 := cave.NewChunk(world.Coords{X: 0, Y: 0}, c, cave.Wall)

	chunkr1 := cave.NewChunk(world.Coords{X: 1, Y: 0}, c, cave.Wall)
	chunkr2 := cave.NewChunk(world.Coords{X: 1, Y: 1}, c, cave.Wall)
	chunkr3 := cave.NewChunk(world.Coords{X: 0, Y: 1}, c, cave.Wall)

	chunkl1 := cave.NewChunk(world.Coords{X: -1, Y: 0}, c, cave.Wall)
	chunkl2 := cave.NewChunk(world.Coords{X: -1, Y: 1}, c, cave.Wall)

	c.RChunks[chunk0.Coords] = chunk0
	c.RChunks[chunkr1.Coords] = chunkr1
	c.RChunks[chunkr2.Coords] = chunkr2
	c.RChunks[chunkr3.Coords] = chunkr3

	c.LChunks[chunkl1.Coords] = chunkl1
	c.LChunks[chunkl2.Coords] = chunkl2
	structures.Entrance(c, c.StartC, 9, 5, 3, false)
	structures.Stairs(c, world.Coords{X: 11, Y: 10}, true, true, 12, 6)
	structures.Stairs(c, world.Coords{X: 21, Y: 10}, false, true, 12, 6)
	structures.GnomeMineLayer(c, world.Coords{X: 11, Y: 23}, world.Coords{X: 21, Y: 23})

	tt := c.StartC
	tt.Y += 13
	trans := transform.NewTransform()
	trans.Pos = c.GetTileInt(tt.X, tt.Y).Transform.Pos
	coll := data.NewCollider(pixel.R(0., 0., world.TileSize * 70., world.TileSize * 3.), true, true)
	coll.Debug = true
	myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Collision, coll).
		AddComponent(myecs.Trigger, data.NewFrameFunc(func() bool {
			sfx.MusicPlayer.SetNextTrack(constants.GameMusic, "hero")
			return true
		}))

	return c
}