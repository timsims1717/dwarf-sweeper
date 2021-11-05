package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/world"
)

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