package cave

import (
	"dwarf-sweeper/pkg/menu"
	"time"
)

var(
	BlocksDug        int
	BlocksDugItem    *menu.ItemText
	BlocksDugTimer   = 0.4
	LowestLevel      int
	LowestLevelItem  *menu.ItemText
	LowestLevelTimer = 0.6
	BombsMarked      int
	BombsMarkedItem  *menu.ItemText
	BombsMarkedTimer = 0.8
	WrongMarks       int
	WrongMarksItem   *menu.ItemText
	WrongMarksTimer  = 1.0
	TotalScore       *menu.ItemText
	TotalScoreTimer  = 1.2
	ScoreTimer       time.Time
)