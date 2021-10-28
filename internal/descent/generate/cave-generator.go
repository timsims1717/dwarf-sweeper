package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
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
	chunk0 := cave.NewChunk(world.Coords{X: 0, Y: 0}, newCave)
	FillChunk(chunk0)

	chunkr1 := cave.NewChunk(world.Coords{X: 1, Y: 0}, newCave)
	chunkr2 := cave.NewChunk(world.Coords{X: 1, Y: 1}, newCave)
	chunkr3 := cave.NewChunk(world.Coords{X: 0, Y: 1}, newCave)
	FillChunk(chunkr1)
	FillChunk(chunkr2)
	FillChunk(chunkr3)

	chunkl1 := cave.NewChunk(world.Coords{X: -1, Y: 0}, newCave)
	chunkl2 := cave.NewChunk(world.Coords{X: -1, Y: 1}, newCave)
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
	newCave.BombPMin = 0.1
	newCave.BombPMax = 0.2
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
	Exit(newCave, exitT, 7, 3, 1, true)
	box = exitT
	box.X -= 5
	box.Y -= 5
	RectRoom(newCave, box, 11, 8)
	newCave.MarkAsNotChanged()
	// generate paths and/or cycles from entrance to exit
	path, deadends := SemiStraightPath(newCave, startT, exitT, data.Left, false)
	var room []world.Coords
	p2, d2 := SemiStraightPath(newCave, startT, exitT, data.Right, false)
	path = append(path, p2...)
	deadends = append(deadends, d2...)
	startT.X -= 2
	p2, d2 = SemiStraightPath(newCave, startT, exitT, data.Down, false)
	path = append(path, p2...)
	deadends = append(deadends, d2...)
	startT.X += 4
	p2, d2 = SemiStraightPath(newCave, startT, exitT, data.Down, false)
	path = append(path, p2...)
	deadends = append(deadends, d2...)
	// generate path branching from orig paths

	// place rectangles at random points on all paths, esp at or near dead ends
	count := random.CaveGen.Intn(constants.ChunkSize / 3) + constants.ChunkSize / 3
	for i := 0; i < count; i++ {
		include := path[random.CaveGen.Intn(len(path))]
		//fmt.Printf("rect room includes: (%d,%d)\n", include.X, include.Y)
		r1 := RandRectRoom(newCave, 7, (constants.ChunkSize/ 4) * 3, include)
		room = append(room, r1...)
	}
	newCave.MarkAsNotChanged()
	count = random.CaveGen.Intn(4) + 4
	for i := 0; i < count; i++ {
		s := path[random.CaveGen.Intn(len(path))]
		if world.Distance(newCave.StartC, s) > 8 && world.Distance(newCave.ExitC, s) > 8 {
			dir := RandomDirection()
			for dir == data.Down {
				dir = RandomDirection()
			}
			NoodleCave(newCave, s, dir)
		}
	}
	newCave.MarkAsNotChanged()
	for _, d := range deadends {
		if world.Distance(newCave.StartC, d) > 8 && world.Distance(newCave.ExitC, d) > 8 {
			if random.CaveGen.Intn(3) == 0 {
				// big
				TreasureRoom(newCave, 6, 8, 2, d)
			} else {
				// small
				TreasureRoom(newCave, 4, 6, 1, d)
			}
		}
	}
	newCave.MarkAsNotChanged()
	if len(room) > 0 {
		count = random.CaveGen.Intn(10) + 20
		for i := 0; i < count; i++ {
			s := room[random.CaveGen.Intn(len(room))]
			if world.Distance(newCave.StartC, s) > 8 && world.Distance(newCave.ExitC, s) > 8 {
				BombableNode(newCave, random.CaveGen.Intn(2) + 2, world.TileSize * 2., true, s)
			}
		}
	}
	for _, ch := range newCave.LChunks {
		FillChunk(ch)
	}
	for _, ch := range newCave.RChunks {
		FillChunk(ch)
	}
	newCave.PrintCaveToTerminal()
	newCave.UpdateAllTileSprites()
	return newCave
}

func makeLayers(left, right, bottom, marginH, marginT, marginB int) [3][3]world.Coords {
	if marginH >= constants.ChunkSize/ 2 {
		marginH = constants.ChunkSize/ 2 - 1
	}
	layer1 := [3]world.Coords{
		{
			X: left *constants.ChunkSize + marginH + random.CaveGen.Intn(constants.ChunkSize- marginH),
			Y: marginT + random.CaveGen.Intn(3),
		},
		{
			X: (left + 1) *constants.ChunkSize + random.CaveGen.Intn((right - left - 1) *constants.ChunkSize),
			Y: marginT + random.CaveGen.Intn(3),
		},
		{
			X: (right + 1) *constants.ChunkSize - marginH - random.CaveGen.Intn(constants.ChunkSize- marginH),
			Y: marginT + random.CaveGen.Intn(3),
		},
	}
	layer2 := [3]world.Coords{
		{
			X: left *constants.ChunkSize + marginH + random.CaveGen.Intn(constants.ChunkSize- marginH),
			Y: constants.ChunkSize + random.CaveGen.Intn((bottom - 1) *constants.ChunkSize),
		},
		{
			X: (left + 1) *constants.ChunkSize + random.CaveGen.Intn((right - left - 1) *constants.ChunkSize),
			Y: constants.ChunkSize + random.CaveGen.Intn((bottom - 1) *constants.ChunkSize),
		},
		{
			X: (right + 1) *constants.ChunkSize - marginH - random.CaveGen.Intn(constants.ChunkSize- marginH),
			Y: constants.ChunkSize + random.CaveGen.Intn((bottom - 1) *constants.ChunkSize),
		},
	}
	layer3 := [3]world.Coords{
		{
			X: left *constants.ChunkSize + marginH + random.CaveGen.Intn(constants.ChunkSize- marginH),
			Y: (bottom + 1) *constants.ChunkSize - marginB - random.CaveGen.Intn(6),
		},
		{
			X: (left + 1) *constants.ChunkSize + random.CaveGen.Intn((right - left - 1) *constants.ChunkSize),
			Y: (bottom + 1) *constants.ChunkSize - marginB - random.CaveGen.Intn(6),
		},
		{
			X: (right + 1) *constants.ChunkSize - marginH - random.CaveGen.Intn(constants.ChunkSize- marginH),
			Y: (bottom + 1) *constants.ChunkSize - marginB - random.CaveGen.Intn(6),
		},
	}
	//fmt.Println("Layers:")
	//fmt.Println("TOP - Left:", layer1[0], "Mid:", layer1[1], "Right:", layer1[2])
	//fmt.Println("MID - Left:", layer2[0], "Mid:", layer2[1], "Right:", layer2[2])
	//fmt.Println("BOT - Left:", layer3[0], "Mid:", layer3[1], "Right:", layer3[2])
	return [3][3]world.Coords{
		layer1,
		layer2,
		layer3,
	}
}