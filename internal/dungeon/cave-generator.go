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
	}
	chunk0 := GenerateChunk(world.Coords{X: 0, Y: 0}, cave)
	EntranceExit(cave, world.Coords{X: 16, Y: 9}, 9, 5, 3)

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

func NewRoomyCave(spriteSheet *img.SpriteSheet, left, right, bottom int) *Cave {
	batcher := img.NewBatcher(spriteSheet)
	startX := rand.Intn((right - left + 1) * ChunkSize - 20) + left * ChunkSize + 10
	startT := world.Coords{X: startX, Y: 9}
	cave := &Cave{
		RChunks: make(map[world.Coords]*Chunk),
		LChunks: make(map[world.Coords]*Chunk),
		finite:  true,
		batcher: batcher,
		left:    left,
		right:   right,
		bottom:  bottom,
		StartC:  startT,
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
	EntranceExit(cave, startT, 11, 5, 4)
	// generate exit (between y level 4 and 10, x between l + 10 and r - 10)
	exitX := rand.Intn((right - left + 1) * ChunkSize - 20) + left * ChunkSize + 10
	exitY := rand.Intn(6) + (bottom + 1) * ChunkSize - 10
	exitT := world.Coords{X: exitX, Y: exitY}
	EntranceExit(cave, exitT, 7, 3, 1)
	// generate paths and/or cycles from entrance to exit
	SemiStraightPath(cave, startT, exitT, Left)
	SemiStraightPath(cave, startT, exitT, Right)
	startT.X -= 2
	SemiStraightPath(cave, startT, exitT, Down)
	startT.X += 4
	SemiStraightPath(cave, startT, exitT, Down)
	// generate path branching from orig paths

	// place rectangles at random points on all paths, esp at or near dead ends

	//num := rand.Intn(8) + 12
	//for i := 0; i < num; i++ {
	//	RectRoom(cave, 5, 20)
	//}
	return cave
}
