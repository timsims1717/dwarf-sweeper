package dungeon

import (
	"dwarf-sweeper/pkg/timing"
)

var(
	BlocksDug        int
	BlocksDugTimer   = 0.4
	GemsFound        int
	GemsFoundTimer   = 0.8
	BombsMarked      int
	BombsMarkedTimer = 1.0
	WrongMarks       int
	WrongMarksTimer  = 1.2
	TotalScoreTimer  = 1.4
	ScoreTimer       *timing.FrameTimer
)