package generate

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/descent/generate/structures/boss"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
)

func NewCave(build *builder.CaveBuilder, level int) *cave.Cave {
	random.RandCaveSeed()
	if build.Biome == "" {
		biome := "mine"
		if random.CaveGen.Intn(2) == 0 {
			biome = "dark"
		}
		build.Biome = biome
	}
	sheet, err := img.LoadSpriteSheet(fmt.Sprintf("assets/img/the-%s.json", build.Biome))
	if err != nil {
		panic(err)
	}
	batcher := img.NewBatcher(sheet, false)
	newCave := cave.NewCave(batcher, build.Biome, build.Type)
	switch build.Base {
	case builder.Roomy:
		RoomyCave(newCave, level)
	case builder.Maze:
		RoomyCave(newCave, level)
	case builder.Custom:
		switch build.Key {
		case "gnomeBoss":
			boss.GnomeBoss(newCave, level)
		}
	}
	for _, s := range build.Structures {
		count := s.Minimum
		if s.MinMult > 0. {
			count = int(newCave.FillVar * s.MinMult)
		}
		if s.RandMult > 0. {
			count += random.CaveGen.Intn(int(newCave.FillVar * s.RandMult))
		}
		if s.Maximum > 0 {
			count = util.Min(s.Maximum, count)
		}
		l := 0
		switch s.Seed {
		case builder.Path:
			l = len(newCave.Paths)
		case builder.Marked:
			l = len(newCave.Marked)
		case builder.DeadEnd:
			l = len(newCave.DeadEnds)
		case builder.Room:
			l = len(newCave.Rooms)
		}
		for i := 0; i < count && l > 0; i++ {
			sti := random.CaveGen.Intn(l)
			var include world.Coords
			switch s.Seed {
			case builder.Path:
				include = newCave.Paths[sti]
			case builder.Marked:
				include = newCave.Marked[sti]
			case builder.DeadEnd:
				include = newCave.DeadEnds[sti]
			case builder.Room:
				include = newCave.Rooms[sti]
			}
			switch s.Key {
			case "pocket":
				structures.Pocket(newCave, random.CaveGen.Intn(3) + 2, world.TileSize * 2., false, include)
			case "ring":
				structures.Ring(newCave, random.CaveGen.Intn(5) + 3, world.TileSize * 3., false, include)
			case "noodleCave":
				dir := structures.RandomDirection()
				for dir == data.Down {
					dir = structures.RandomDirection()
				}
				structures.NoodleCave(newCave, include, dir)
			case "treasure":
				if random.CaveGen.Intn(3) == 0 {
					// big
					structures.TreasureRoom(newCave, 6, 8, 2, include)
				} else {
					// small
					structures.TreasureRoom(newCave, 4, 6, 1, include)
				}
			case "bombable":
				structures.BombableNode(newCave, random.CaveGen.Intn(2) + 1, world.TileSize * 2., true, include)
			case "mineLayer":
				structures.MineLayer(newCave, include)
			case "stairs":
				structures.Stairs(newCave, include, random.CaveGen.Intn(2) == 0, random.CaveGen.Intn(2) == 0, 0, 0)
			}
			switch s.Seed {
			case builder.Path:
				newCave.Paths = append(newCave.Paths[:sti], newCave.Paths[sti+1:]...)
				l = len(newCave.Paths)
			case builder.Marked:
				newCave.Marked = append(newCave.Marked[:sti], newCave.Marked[sti+1:]...)
				l = len(newCave.Marked)
			case builder.DeadEnd:
				newCave.DeadEnds = append(newCave.DeadEnds[:sti], newCave.DeadEnds[sti+1:]...)
				l = len(newCave.DeadEnds)
			case builder.Room:
				newCave.Rooms = append(newCave.Rooms[:sti], newCave.Rooms[sti+1:]...)
				l = len(newCave.Rooms)
			}
		}
	}
	for _, ch := range newCave.LChunks {
		structures.FillBasic(ch)
	}
	for _, ch := range newCave.RChunks {
		structures.FillBasic(ch)
	}
	fmt.Println("Total bombs:", descent.CaveTotalBombs)
	newCave.PrintCaveToTerminal()
	return newCave
}