package descent

import (
	"dwarf-sweeper/pkg/timing"
)

var(
	CaveBlocksDug    int
	CaveGemsFound    int
	CaveTotalBombs   int
	CaveBlownUpBombs int
	CaveBombsLeft    int
	CaveBombsMarked  int
	CaveCorrectMarks int
	CaveWrongMarks   int

	BlocksDug        int
	GemsFound        int
	BombsMarked      int
	CorrectMarks     int
	WrongMarks       int
	BlocksDugTimer   = 0.4
	GemsFoundTimer   = 0.8
	BombsMarkedTimer = 1.0
	WrongMarksTimer  = 1.2
	TotalScoreTimer  = 1.4
	ScoreTimer       *timing.FrameTimer
)

func ResetStats() {
	BlocksDug = 0
	GemsFound = 0
	BombsMarked = 0
	WrongMarks = 0
	ResetCaveStats()
}

func ResetCaveStats() {
	CaveBlocksDug = 0
	CaveGemsFound = 0
	CaveTotalBombs = 0
	CaveBlownUpBombs = 0
	CaveBombsLeft = 0
	CaveBombsMarked = 0
	CaveWrongMarks = 0
}

func AddStats() {
	BlocksDug += CaveBlocksDug
	GemsFound += CaveGemsFound
	BombsMarked += CaveBombsMarked
	WrongMarks += CaveWrongMarks
}