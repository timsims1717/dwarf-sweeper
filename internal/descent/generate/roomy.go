package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/structures"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/world"
)

func RoomyCave(newCave *cave.Cave, signal chan bool) {
	if signal != nil {
		signal <- false
	}

	layers := makeLayers(newCave.Left, newCave.Right, newCave.Bottom, 5, 7, 3)
	start := random.CaveGen.Intn(3) // 0 = left, 1 = mid, 2 = right
	end := random.CaveGen.Intn(2)
	if start == 0 {
		end++
	} else if start == 1 && end == 1 {
		end = random.CaveGen.Intn(2)
	}
	startT := layers[0][start]
	exitT := layers[2][end]
	if signal != nil {
		signal <- false
		if !<-signal {
			return
		}
	}
	// generate paths and/or cycles from entrance to exit
	newCave.Paths, newCave.DeadEnds, newCave.Marked = structures.SemiStraightPath(newCave, startT, exitT, data.Left, false)
	var room []world.Coords
	p2, d2, m2 := structures.SemiStraightPath(newCave, startT, exitT, data.Right, false)
	newCave.Paths = append(newCave.Paths, p2...)
	newCave.DeadEnds = append(newCave.DeadEnds, d2...)
	newCave.Marked = append(newCave.Marked, m2...)
	if signal != nil {
		signal <- false
		if !<-signal {
			return
		}
	}
	startT.X -= 2
	p2, d2, m2 = structures.SemiStraightPath(newCave, startT, exitT, data.Down, false)
	newCave.Paths = append(newCave.Paths, p2...)
	newCave.Marked = append(newCave.Marked, m2...)
	newCave.DeadEnds = append(newCave.DeadEnds, d2...)
	if signal != nil {
		signal <- false
		if !<-signal {
			return
		}
	}
	startT.X += 4
	p2, d2, m2 = structures.SemiStraightPath(newCave, startT, exitT, data.Down, false)
	newCave.Paths = append(newCave.Paths, p2...)
	newCave.DeadEnds = append(newCave.DeadEnds, d2...)
	newCave.Marked = append(newCave.Marked, m2...)
	if signal != nil {
		signal <- false
		if !<-signal {
			return
		}
	}
	// generate path branching from orig paths
	count := random.CaveGen.Intn(int(newCave.FillVar*0.25)) + int(newCave.FillVar*0.3)
	for i := 0; i < count && len(newCave.Marked) > 0; i++ {
		sti := random.CaveGen.Intn(len(newCave.Marked))
		include := newCave.Marked[sti]
		newCave.Marked = append(newCave.Marked[:sti], newCave.Marked[sti+1:]...)
		p2, d2, m2 = structures.BranchOff(newCave, include, 8, 16)
		newCave.Paths = append(newCave.Paths, p2...)
		newCave.DeadEnds = append(newCave.DeadEnds, d2...)
		newCave.Marked = append(newCave.Marked, m2...)
		if signal != nil {
			signal <- false
			if !<-signal {
				return
			}
		}
	}
	// place rectangles at random marked tiles
	count = random.CaveGen.Intn(int(newCave.FillVar*0.25)) + int(newCave.FillVar*0.3)
	for i := 0; i < count && len(newCave.Marked) > 0; i++ {
		sti := random.CaveGen.Intn(len(newCave.Marked))
		include := newCave.Marked[sti]
		newCave.Marked = append(newCave.Marked[:sti], newCave.Marked[sti+1:]...)
		r1, m1 := structures.RandRectRoom(newCave, 7, int(newCave.FillVar*0.25), include)
		room = append(room, r1...)
		newCave.Marked = append(newCave.Marked, m1...)
		if signal != nil {
			signal <- false
			if !<-signal {
				return
			}
		}
	}
	newCave.MarkAsNotChanged()
}

func makeLayers(left, right, bottom, marginH, marginT, marginB int) [3][3]world.Coords {
	if marginH >= constants.ChunkSize/2 {
		marginH = constants.ChunkSize/2 - 1
	}
	layer1 := [3]world.Coords{
		{
			X: left*constants.ChunkSize + marginH + random.CaveGen.Intn(constants.ChunkSize-marginH),
			Y: marginT + random.CaveGen.Intn(3),
		},
		{
			X: (left+1)*constants.ChunkSize + random.CaveGen.Intn((right-left-1)*constants.ChunkSize),
			Y: marginT + random.CaveGen.Intn(3),
		},
		{
			X: (right+1)*constants.ChunkSize - marginH - random.CaveGen.Intn(constants.ChunkSize-marginH),
			Y: marginT + random.CaveGen.Intn(3),
		},
	}
	layer2 := [3]world.Coords{
		{
			X: left*constants.ChunkSize + marginH + random.CaveGen.Intn(constants.ChunkSize-marginH),
			Y: constants.ChunkSize + random.CaveGen.Intn((bottom-1)*constants.ChunkSize),
		},
		{
			X: (left+1)*constants.ChunkSize + random.CaveGen.Intn((right-left-1)*constants.ChunkSize),
			Y: constants.ChunkSize + random.CaveGen.Intn((bottom-1)*constants.ChunkSize),
		},
		{
			X: (right+1)*constants.ChunkSize - marginH - random.CaveGen.Intn(constants.ChunkSize-marginH),
			Y: constants.ChunkSize + random.CaveGen.Intn((bottom-1)*constants.ChunkSize),
		},
	}
	layer3 := [3]world.Coords{
		{
			X: left*constants.ChunkSize + marginH + random.CaveGen.Intn(constants.ChunkSize-marginH),
			Y: (bottom+1)*constants.ChunkSize - marginB - random.CaveGen.Intn(6),
		},
		{
			X: (left+1)*constants.ChunkSize + random.CaveGen.Intn((right-left-1)*constants.ChunkSize),
			Y: (bottom+1)*constants.ChunkSize - marginB - random.CaveGen.Intn(6),
		},
		{
			X: (right+1)*constants.ChunkSize - marginH - random.CaveGen.Intn(constants.ChunkSize-marginH),
			Y: (bottom+1)*constants.ChunkSize - marginB - random.CaveGen.Intn(6),
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