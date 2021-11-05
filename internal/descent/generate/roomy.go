package generate
import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/world"
	"fmt"
)

// requires at least 3 chunks wide
func NewRoomyCave(spriteSheet *img.SpriteSheet, biome string, level, left, right, bottom int) *cave.Cave {
	random.RandCaveSeed()
	batcher := img.NewBatcher(spriteSheet, false)
	layers := makeLayers(left, right, bottom, 5, 7, 3)
	start := random.CaveGen.Intn(3) // 0 = left, 1 = mid, 2 = right
	end := random.CaveGen.Intn(2)
	if start == 0 {
		end++
	} else if start == 1 && end == 1 {
		end = random.CaveGen.Intn(2)
	}
	startT := layers[0][start]
	exitT := layers[2][end]
	newCave := cave.NewCave(batcher, biome,true)
	newCave.SetSize(left, right, bottom)
	newCave.StartC = startT
	newCave.ExitC = exitT
	newCave.GemRate = constants.BaseGem
	newCave.ItemRate = constants.BaseItem
	newCave.BombPMin = 0.12
	newCave.BombPMax = 0.22
	for i := 1; i < level; i++ {
		newCave.BombPMin += 0.02
		newCave.BombPMax += 0.02
	}
	if newCave.BombPMin > 0.3 {
		newCave.BombPMin = 0.3
	}
	if newCave.BombPMax > 0.4 {
		newCave.BombPMax = 0.4
	}
	CreateChunks(newCave)
	// generate entrance (at y level 9, x between l + 10 and r - 10)
	Entrance(newCave, startT, 11, 5, 4, false)
	box := startT
	box.X -= 8
	box.Y -= 9
	RectRoom(newCave, box, 17, 12)
	// generate exit (between y level 4 and 10, x between l + 10 and r - 10)
	Entrance(newCave, exitT, 7, 3, 1, true)
	box = exitT
	box.X -= 5
	box.Y -= 5
	RectRoom(newCave, box, 11, 8)
	newCave.MarkAsNotChanged()
	// generate paths and/or cycles from entrance to exit
	path, deadends, marked := SemiStraightPath(newCave, startT, exitT, data.Left, false)
	var room []world.Coords
	p2, d2, m2 := SemiStraightPath(newCave, startT, exitT, data.Right, false)
	path = append(path, p2...)
	deadends = append(deadends, d2...)
	marked = append(marked, m2...)
	startT.X -= 2
	p2, d2, m2 = SemiStraightPath(newCave, startT, exitT, data.Down, false)
	path = append(path, p2...)
	marked = append(marked, m2...)
	deadends = append(deadends, d2...)
	startT.X += 4
	p2, d2, m2 = SemiStraightPath(newCave, startT, exitT, data.Down, false)
	path = append(path, p2...)
	deadends = append(deadends, d2...)
	marked = append(marked, m2...)
	// how much we should fill is based on the size of the cave
	// a 32 chunk size at 3x3 gives a value of 72 for fillVar
	// a 16 chunk size at 3x3 gives a value of 18 for fillVar
	fillVar := newCave.Width * newCave.Height / 128.
	fillVarF := float64(fillVar)
	// generate path branching from orig paths
	count := random.CaveGen.Intn(fillVar / 4) + fillVar / 3
	for i := 0; i < count && len(marked) > 0; i++ {
		sti := random.CaveGen.Intn(len(marked))
		include := marked[sti]
		marked = append(marked[:sti], marked[sti+1:]...)
		p2, d2, m2 = BranchOff(newCave, include, 8, 16)
		path = append(path, p2...)
		deadends = append(deadends, d2...)
		marked = append(marked, m2...)
	}
	// place rectangles at random marked tiles
	count = random.CaveGen.Intn(fillVar / 4) + fillVar / 3
	for i := 0; i < count && len(marked) > 0; i++ {
		sti := random.CaveGen.Intn(len(marked))
		include := marked[sti]
		marked = append(marked[:sti], marked[sti+1:]...)
		r1, m1 := RandRectRoom(newCave, 7, (constants.ChunkSize/ 4) * 3, include)
		room = append(room, r1...)
		marked = append(marked, m1...)
	}
	newCave.MarkAsNotChanged()
	/* ***********  This marks where the Base Cave Generation ends  *********** */
	/* ***********        and structures begin to be created.       *********** */
	// pockets
	count = random.CaveGen.Intn(int(fillVarF * 0.16) + int(fillVarF * 0.2))
	for i := 0; i < count && len(marked) > 0; i++ {
		sti := random.CaveGen.Intn(len(marked))
		include := marked[sti]
		marked = append(marked[:sti], marked[sti+1:]...)
		Pocket(newCave, random.CaveGen.Intn(3) + 2, world.TileSize * 2., false, include)
	}
	// rings
	count = random.CaveGen.Intn(fillVar / 6) + fillVar / 5
	for i := 0; i < count && len(marked) > 0; i++ {
		sti := random.CaveGen.Intn(len(marked))
		include := marked[sti]
		marked = append(marked[:sti], marked[sti+1:]...)
		Ring(newCave, random.CaveGen.Intn(5) + 3, world.TileSize * 3., false, include)
	}
	// noodle caves
	count = random.CaveGen.Intn(fillVar / 6) + fillVar / 5
	for i := 0; i < count && len(marked) > 0; i++ {
		sti := random.CaveGen.Intn(len(marked))
		s := marked[sti]
		marked = append(marked[:sti], marked[sti+1:]...)
		dir := RandomDirection()
		for dir == data.Down {
			dir = RandomDirection()
		}
		NoodleCave(newCave, s, dir)
	}
	newCave.MarkAsNotChanged()
	// treasure rooms
	max := 10
	count = 0
	for _, d := range deadends {
		if random.CaveGen.Intn(3) == 0 {
			// big
			TreasureRoom(newCave, 6, 8, 2, d)
		} else {
			// small
			TreasureRoom(newCave, 4, 6, 1, d)
		}
		count++
		if count > max {
			break
		}
	}
	newCave.MarkAsNotChanged()
	// bombable nodes
	if len(room) > 0 {
		count = random.CaveGen.Intn(10) + 20
		for i := 0; i < count; i++ {
			s := room[random.CaveGen.Intn(len(room))]
			BombableNode(newCave, random.CaveGen.Intn(2) + 1, world.TileSize * 2., true, s)
		}
	}
	for _, ch := range newCave.LChunks {
		FillChunk(ch)
	}
	for _, ch := range newCave.RChunks {
		FillChunk(ch)
	}
	fmt.Println("Total bombs:", descent.CaveTotalBombs)
	newCave.PrintCaveToTerminal()
	return newCave
}
