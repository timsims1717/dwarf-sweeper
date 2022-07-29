package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/boss"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/pathfinding"
	"dwarf-sweeper/internal/profile"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/noise"
	"dwarf-sweeper/pkg/util"
	"fmt"
	"github.com/beefsack/go-astar"
)

func NewAsyncCave(build *builder.CaveBuilder, level int, signal chan bool) *cave.Cave {
	random.RandCaveSeed()
	noise.Seed(random.CaveGen)
	if build.Biome == "" {
		biome := "mine"
		if random.CaveGen.Intn(2) == 0 {
			biome = "dark"
		}
		build.Biome = biome
	}
	c := cave.NewCave(build.Biome, build.Type)
	c.Level = level
	c.Enemies = build.Enemies
	left := 0
	right := util.Max(build.Width - 1, 2)
	bottom := util.Max(build.Height - 1, 2)
	c.SetSize(left, right, bottom)
	c.BombPMin, c.BombPMax = BombLevel(level)
	structures.CreateChunks(c, cave.Blast)
	go buildCave(build, c, signal)
	return c
}

func NewCave(build *builder.CaveBuilder, level int) *cave.Cave {
	random.RandCaveSeed()
	noise.Seed(random.CaveGen)
	if build.Biome == "" {
		build.Biome = "mine"
	}
	c := cave.NewCave(build.Biome, build.Type)
	c.Level = level
	c.Enemies = build.Enemies
	left := 0
	right := util.Max(build.Width - 1, 2)
	bottom := util.Max(build.Height - 1, 2)
	c.SetSize(left, right, bottom)
	c.BombPMin, c.BombPMax = BombLevel(level)
	structures.CreateChunks(c, cave.Blast)
	buildCave(build, c, nil)
	return c
}

func buildCave(build *builder.CaveBuilder, c *cave.Cave, signal chan bool) {
	fmt.Printf("Name: %s\n", build.Name)
	fmt.Printf("Biome: %s\n", c.Biome)
	switch build.Base {
	case builder.Roomy:
		RoomyCave(c, signal)
	case builder.Blob:
		BlobCave(c, signal)
	case builder.Maze:
		MazeCave(c, signal)
	case builder.Maze2:
		TightMazeCave(c, signal)
	case builder.Custom:
		switch build.Key {
		case "gnomeBoss":
			boss.GnomeBoss(c, c.Level)
		case "minesweeper":
			MinesweeperCave(c, c.Level)
		}
	}
	if build.Base != builder.Custom {
		// generate entrance with a start inside the largest group
		tile := PickTile(c, 8, 8, 7, constants.ChunkSize * (c.Bottom+1) - 10, true)
		startC := tile.RCoords
		c.StartC = startC
		structures.Entrance(c, startC, 11, 5, 4, cave.DoorType(build.DoorType))
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
		tile = PickTile(c, 8, 8, constants.ChunkSize*(c.Bottom+1)-9, 5, true)
		exitC := tile.RCoords
		structures.Exit(c, exitC, 7, 3, 1, 0, cave.DoorType(build.DoorType))
		box = exitC
		box.X -= 5
		box.Y -= 5
		structures.RectRoom(c, box, 11, 8, 3, cave.Unknown)
		descent.Descent.Exits = []string{c.Biome}
		c.MarkAsNotChanged()
		if signal != nil {
			signal <- false
			if !<-signal {
				return
			}
		}
	}
	for _, s := range build.Structures {
		enemies := c.Enemies
		if len(s.Enemies) > 0 {
			enemies = s.Enemies
		}
		s.Defaults()
		r := s.Maximum - s.Minimum
		count := s.Minimum
		if r > 0 {
			count = s.Minimum + random.CaveGen.Intn(r+1)
		}
		for i := 0; i < count; i++ {
			if s.Chance > 0. && random.CaveGen.Float64() > s.Chance {
				continue
			}
			tile := PickTileDist(c, s.MarginL, s.MarginR, s.MarginT, s.MarginB, false, s.DigDist)
			if tile != nil {
				xDir := data.Left
				if tile.RCoords.X < c.Left*constants.ChunkSize {
					xDir = data.Right
				} else if tile.RCoords.X > c.Right*constants.ChunkSize {
					xDir = data.Left
				} else if random.CaveGen.Intn(2) == 0 {
					xDir = data.Right
				}
				switch s.Key {
				case "secretExit":
					var biome string
					tw := 0
					for _, w := range profile.CurrentProfile.BiomeExits[c.Biome] {
						tw += w
					}
					rw := random.CaveGen.Intn(tw)
					tw = 0
					for b, w := range profile.CurrentProfile.BiomeExits[c.Biome] {
						tw += w
						if rw < tw {
							biome = b
							descent.Descent.Exits = append(descent.Descent.Exits, b)
							break
						}
					}
					structures.SecretExit(c, tile.RCoords, i+1, biome)
					fmt.Printf("Secret Exit: (%d,%d)\n", tile.RCoords.X, tile.RCoords.Y)
				case "pocket":
					structures.Pocket(c, random.CaveGen.Intn(4)+2, 1.5, false, tile.RCoords, enemies)
				case "ring":
					structures.Ring(c, random.CaveGen.Intn(5)+3, 3., false, tile.RCoords, enemies)
				case "noodleCave":
					dir := structures.RandomDirection()
					for dir == data.Down {
						dir = structures.RandomDirection()
					}
					structures.NoodleCave(c, tile.RCoords, dir, enemies)
				case "treasure":
					if random.CaveGen.Intn(3) == 0 {
						// big
						structures.TreasureRoom(c, 7, 9, 2, tile.RCoords)
					} else {
						// small
						structures.TreasureRoom(c, 5, 7, 1, tile.RCoords)
					}
				case "bombable":
					structures.BombableNode(c, random.CaveGen.Intn(4)+2, 2., true, tile.RCoords)
				case "mineLayer":
					structures.MineLayer(c, tile.RCoords, enemies)
				case "mineComplex":
					structures.MineComplex(c, tile.RCoords, data.Direction(xDir), enemies)
				case "stairs":
					structures.Stairs(c, tile.RCoords, random.CaveGen.Intn(2) == 0, random.CaveGen.Intn(2) == 0, 0, 0)
				case "bigBomb":
					descent.Descent.CoordsMap["big-bomb"] = tile.RCoords
					structures.BombRoom(c, 4, 7, 7, 11, 3, c.Level, tile.RCoords)
				case "bridge":
					structures.BridgeCavern(c, tile.RCoords, data.Direction(xDir), enemies)
				case "cavern":
					structures.Cavern(c, tile.RCoords, data.Direction(xDir), enemies)
				case "smallCamp":
					structures.SmallCamp(c, tile.RCoords, data.Direction(xDir))
				}
				c.MarkAsNotChanged()
				if signal != nil {
					signal <- false
					if !<-signal {
						return
					}
				}
			} else {
				fmt.Printf("failed to generate structure %s\n", s.Key)
			}
		}
	}
	if c.Type != cave.Minesweeper {
		structures.FillCave(c)
	}
	fmt.Println("Total bombs:", c.TotalBombs)
	c.PrintCaveToTerminal()
	if signal != nil {
		signal <- true
		<-signal
	}
}

func PickTileDist(c *cave.Cave, marginL, marginR, marginT, marginB int, mainGroup bool, digDist builder.DigDist) *cave.Tile {
	var tX, tY int
	var tile *cave.Tile = nil
	tryCount := 0
	var min, max int
	mult := util.Max(c.Bottom, c.Right - c.Left) + 1
	switch digDist {
	case builder.Close:
		min = constants.ChunkSize
		max = constants.ChunkSize * mult
	case builder.Medium:
		min = constants.ChunkSize
		max = constants.ChunkSize * mult * 2
	case builder.Far:
		min = constants.ChunkSize * mult
		max = -1
	case builder.Farthest:
		min = constants.ChunkSize * mult * 2
		max = -1
	}
	for {
		tX = c.Left * constants.ChunkSize + marginL + random.CaveGen.Intn((c.Right-c.Left+1) * constants.ChunkSize - (marginR + marginL))
		tY = marginT + random.CaveGen.Intn((c.Bottom+1) * constants.ChunkSize - (marginT + marginB))
		tile = c.GetTileInt(tX, tY)
		if tile != nil && tile.Type != cave.Wall && (tile.Group == c.MainGroup || !mainGroup) && !tile.IsChanged && !tile.NeverChange {
			if digDist == builder.Any {
				return tile
			}
			cave.Origin = c.StartC
			cave.NeighborsFn = pathfinding.CaveGenNeighbors
			cave.CostFn = pathfinding.CaveGenCost
			_, dist, found := astar.Path(c.GetTileInt(c.StartC.X, c.StartC.Y), tile)
			if found && int(dist) >= min && (max == -1 || int(dist) <= max) {
				return tile
			}
		}
		tryCount++
		if tryCount > 4 {
			return nil
		}
	}
}

func PickTile(c *cave.Cave, marginL, marginR, marginT, marginB int, mainGroup bool) *cave.Tile {
	var tX, tY int
	var tile *cave.Tile = nil
	for tile == nil || tile.Type == cave.Wall ||
		(tile.Group != c.MainGroup && mainGroup) ||
		tile.IsChanged || tile.NeverChange {
		tX = c.Left*constants.ChunkSize + marginL + random.CaveGen.Intn((c.Right-c.Left+1) * constants.ChunkSize - (marginR + marginL))
		tY = marginT + random.CaveGen.Intn((c.Bottom+1) * constants.ChunkSize - (marginT + marginB))
		tile = c.GetTileInt(tX, tY)
	}
	return tile
}