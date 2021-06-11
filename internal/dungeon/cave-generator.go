package dungeon

import (
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/world"
	"math/rand"
)

func NewInfiniteCave(spriteSheet *img.SpriteSheet) *Cave {
	batcher := img.NewBatcher(spriteSheet)
	cave := &Cave{
		RChunks: nil,
		LChunks: nil,
		batcher: batcher,
		StartC:  world.Coords{X: 16, Y: 9},
		update:  true,
	}
	chunk0 := GenerateChunk(world.Coords{X: 0, Y: 0}, cave)
	EntranceExit(cave, world.Coords{X: 16, Y: 9}, 9, 5, 3, false)

	chunkr1 := GenerateChunk(world.Coords{X: 1, Y: 0}, cave)
	chunkr2 := GenerateChunk(world.Coords{X: 1, Y: 1}, cave)
	chunkr3 := GenerateChunk(world.Coords{X: 0, Y: 1}, cave)

	chunkl1 := GenerateChunk(world.Coords{X: -1, Y: 0}, cave)
	chunkl2 := GenerateChunk(world.Coords{X: -1, Y: 1}, cave)

	cave.RChunks = make(map[world.Coords]*Chunk)
	cave.RChunks[chunk0.Coords] = chunk0
	cave.RChunks[chunkr1.Coords] = chunkr1
	cave.RChunks[chunkr2.Coords] = chunkr2
	cave.RChunks[chunkr3.Coords] = chunkr3

	cave.LChunks = make(map[world.Coords]*Chunk)
	cave.LChunks[chunkl1.Coords] = chunkl1
	cave.LChunks[chunkl2.Coords] = chunkl2
	return cave
}

// requires at least 3 chunks wide
func NewRoomyCave(spriteSheet *img.SpriteSheet, level, left, right, bottom int) *Cave {
	batcher := img.NewBatcher(spriteSheet)
	layers := makeLayers(left, right, bottom)
	start := rand.Intn(3) // 0 = left, 1 = mid, 2 = right
	end := rand.Intn(2)
	if start == 0 {
		end++
	} else if start == 1 && end == 1 {
		end = rand.Intn(2)
	}
	startT := layers[0][start]
	exitT := layers[2][end]
	cave := &Cave{
		RChunks: make(map[world.Coords]*Chunk),
		LChunks: make(map[world.Coords]*Chunk),
		finite:  true,
		batcher: batcher,
		left:    left,
		right:   right,
		bottom:  bottom,
		StartC:  startT,
		update:  true,
	}
	cave.fuseLen = 1.0
	cave.bombPMin = 0.1
	cave.bombPMax = 0.2
	for i := 1; i < level; i += 2 {
		cave.bombPMin += 0.02
		cave.bombPMax += 0.02
	}
	if cave.bombPMin > 0.3 {
		cave.bombPMin = 0.3
	}
	if cave.bombPMax > 0.4 {
		cave.bombPMax = 0.4
	}
	for i := 2; i < level; i += 2 {
		cave.fuseLen -= 0.1
	}
	if cave.fuseLen < 0.4 {
		cave.fuseLen = 0.4
	}
	for y := 0; y <= bottom; y++ {
		for x := left; x <= right; x++ {
			chunk := GenerateChunk(world.Coords{X: x, Y: y}, cave)
			if x >= 0 {
				cave.RChunks[chunk.Coords] = chunk
			} else {
				cave.LChunks[chunk.Coords] = chunk
			}
		}
	}
	// generate entrance (at y level 9, x between l + 10 and r - 10)
	EntranceExit(cave, startT, 11, 5, 4, false)
	box := startT
	box.X -= 8
	box.Y -= 9
	RectRoom(cave, box, 17, 12)
	// generate exit (between y level 4 and 10, x between l + 10 and r - 10)
	EntranceExit(cave, exitT, 7, 3, 1, true)
	box = exitT
	box.X -= 5
	box.Y -= 5
	RectRoom(cave, box, 11, 8)
	cave.markAsNotChanged()
	// generate paths and/or cycles from entrance to exit
	path := SemiStraightPath(cave, startT, exitT, Left, false)
	path = append(path, SemiStraightPath(cave, startT, exitT, Right, false)...)
	startT.X -= 2
	path = append(path, SemiStraightPath(cave, startT, exitT, Down, false)...)
	startT.X += 4
	path = append(path, SemiStraightPath(cave, startT, exitT, Down, false)...)
	// generate path branching from orig paths

	// place rectangles at random points on all paths, esp at or near dead ends
	count := rand.Intn(15) + 5
	for i := 0; i < count; i++ {
		include := path[rand.Intn(len(path))]
		//fmt.Printf("rect room includes: (%d,%d)\n", include.X, include.Y)
		RandRectRoom(cave, 8, 24, include)
	}
	//num := rand.Intn(8) + 12
	//for i := 0; i < num; i++ {
	//	RectRoom(cave, 5, 20)
	//}
	cave.PrintCaveToTerminal()
	return cave
}

func makeLayers(left, right, bottom int) [3][3]world.Coords {
	layer1 := [3]world.Coords{
		{
			X: rand.Intn(ChunkSize - 10) + left * ChunkSize + 10,
			Y: rand.Intn(3) + 8,
		},
		{
			X: rand.Intn((right - left - 1) * ChunkSize) + left * ChunkSize,
			Y: rand.Intn(3) + 8,
		},
		{
			X: rand.Intn(ChunkSize - 10) + right * ChunkSize - 10,
			Y: rand.Intn(3) + 8,
		},
	}
	layer2 := [3]world.Coords{
		{
			X: rand.Intn(ChunkSize - 10) + left * ChunkSize + 10,
			Y: rand.Intn((bottom - 1) * ChunkSize) + ChunkSize,
		},
		{
			X: rand.Intn((right - left - 1) * ChunkSize) + left * ChunkSize,
			Y: rand.Intn((bottom - 1) * ChunkSize) + ChunkSize,
		},
		{
			X: rand.Intn(ChunkSize - 10) + right * ChunkSize - 10,
			Y: rand.Intn((bottom - 1) * ChunkSize) + ChunkSize,
		},
	}
	layer3 := [3]world.Coords{
		{
			X: rand.Intn(ChunkSize - 10) + left * ChunkSize + 10,
			Y: rand.Intn(6) + (bottom + 1) * ChunkSize - 10,
		},
		{
			X: rand.Intn((right - left - 1) * ChunkSize) + left * ChunkSize,
			Y: rand.Intn(6) + (bottom + 1) * ChunkSize - 10,
		},
		{
			X: rand.Intn(ChunkSize - 10) + right * ChunkSize - 10,
			Y: rand.Intn(6) + (bottom + 1) * ChunkSize - 10,
		},
	}
	return [3][3]world.Coords{
		layer1,
		layer2,
		layer3,
	}
}