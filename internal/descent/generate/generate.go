package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/descent/generate/structures/boss"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/noise"
	"dwarf-sweeper/pkg/world"
	"fmt"
)

func NewAsyncCave(build *builder.CaveBuilder, level int, signal chan bool) *cave.Cave {
	random.RandCaveSeed()
	noise.SeedBlockType(random.CaveGen)
	if build.Biome == "" {
		biome := "mine"
		if random.CaveGen.Intn(2) == 0 {
			biome = "dark"
		}
		build.Biome = biome
	}
	c := cave.NewCave(build.Biome, build.Type)
	left := 0
	right := 2
	bottom := 2
	c.SetSize(left, right, bottom)
	c.BombPMin, c.BombPMax = BombLevel(level)
	structures.CreateChunks(c)
	go newCave(build, c, level, signal)
	return c
}

func NewCave(build *builder.CaveBuilder, level int) *cave.Cave {
	random.RandCaveSeed()
	noise.SeedBlockType(random.CaveGen)
	if build.Biome == "" {
		biome := "mine"
		if random.CaveGen.Intn(2) == 0 {
			biome = "dark"
		}
		build.Biome = biome
	}
	c := cave.NewCave(build.Biome, build.Type)
	left := 0
	right := 2
	bottom := 2
	c.SetSize(left, right, bottom)
	c.BombPMin, c.BombPMax = BombLevel(level)
	structures.CreateChunks(c)
	newCave(build, c, level, nil)
	return c
}

func newCave(build *builder.CaveBuilder, c *cave.Cave, level int, signal chan bool) {
	switch build.Base {
	case builder.Roomy:
		RoomyCave(c, level, signal)
	case builder.Blob:
		BlobCave(c, signal)
		// entrance (will be moved outside base later)
		// generate entrance with a start inside the largest group
		tile := PickTile(c, 8, 8, 7, constants.ChunkSize * (c.Bottom+1) - 10)
		startC := tile.RCoords
		c.StartC = startC
		structures.Entrance(c, startC, 11, 5, 4, false)
		box := startC
		box.X -= 8
		box.Y -= 9
		structures.RectRoom(c, box, 17, 12, 3, cave.Unknown)
		c.MarkAsNotChanged()
		if signal != nil {
			signal <- false
			if !<-signal {
				return
			}
		}
		// generate exit inside the largest group
		tile = PickTile(c, 8, 8, constants.ChunkSize * (c.Bottom+1) - 9, 5)
		exitC := tile.RCoords
		c.ExitC = exitC
		structures.Entrance(c, exitC, 7, 3, 1, true)
		box = exitC
		box.X -= 5
		box.Y -= 5
		structures.RectRoom(c, box, 11, 8,3, cave.Unknown)
		c.MarkAsNotChanged()
		if signal != nil {
			signal <- false
			if !<-signal {
				return
			}
		}
	case builder.Maze:
		RoomyCave(c, level, signal)
	case builder.Custom:
		switch build.Key {
		case "gnomeBoss":
			boss.GnomeBoss(c, level)
		case "minesweeper":
			MinesweeperCave(c, level)
		}
	}
	for _, s := range build.Structures {
		count := s.Minimum + random.CaveGen.Intn(s.Maximum-s.Minimum)
		for i := 0; i < count; i++ {
			tile := PickTile(c, s.MarginL, s.MarginR, s.MarginT, s.MarginB)
			switch s.Key {
			case "pocket":
				structures.Pocket(c, random.CaveGen.Intn(3)+2, world.TileSize*2., false, tile.RCoords)
			case "ring":
				structures.Ring(c, random.CaveGen.Intn(5)+3, world.TileSize*3., false, tile.RCoords)
			case "noodleCave":
				dir := structures.RandomDirection()
				for dir == data.Down {
					dir = structures.RandomDirection()
				}
				structures.NoodleCave(c, tile.RCoords, dir)
			case "treasure":
				if random.CaveGen.Intn(3) == 0 {
					// big
					structures.TreasureRoom(c, 6, 8, 2, tile.RCoords)
				} else {
					// small
					structures.TreasureRoom(c, 4, 6, 1, tile.RCoords)
				}
			case "bombable":
				structures.BombableNode(c, random.CaveGen.Intn(2)+1, world.TileSize*2., true, tile.RCoords)
			case "mineLayer":
				structures.MineLayer(c, tile.RCoords)
			case "mineComplex":
				dir := data.Left
				if tile.RCoords.X < c.Left*constants.ChunkSize {
					dir = data.Right
				} else if tile.RCoords.X > c.Right*constants.ChunkSize {
					dir = data.Left
				} else if random.CaveGen.Intn(2) == 0 {
					dir = data.Right
				}
				structures.MineComplex(c, tile.RCoords, data.Direction(dir))
			case "stairs":
				structures.Stairs(c, tile.RCoords, random.CaveGen.Intn(2) == 0, random.CaveGen.Intn(2) == 0, 0, 0)
			}
			if signal != nil {
				signal <- false
				if !<-signal {
					return
				}
			}
		}
	}
	if c.Type != cave.Minesweeper {
		structures.FillCave(c)
	}
	fmt.Println("Total bombs:", descent.CaveTotalBombs)
	c.PrintCaveToTerminal()
	if signal != nil {
		signal <- true
		<-signal
	}
}

func PickTile(c *cave.Cave, marginL, marginR, marginT, marginB int) *cave.Tile {
	var tX, tY int
	var tile *cave.Tile = nil
	for tile == nil || tile.Type == cave.Wall || tile.Group != c.MainGroup || tile.IsChanged {
		tX = c.Left*constants.ChunkSize + marginL + random.CaveGen.Intn((c.Right-c.Left+1) * constants.ChunkSize - (marginR + marginL))
		tY = marginT + random.CaveGen.Intn((c.Bottom+1) * constants.ChunkSize - (marginT + marginB))
		tile = c.GetTileInt(tX, tY)
	}
	return tile
}