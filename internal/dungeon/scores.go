package dungeon

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
	GemsFound        int
	GemsFoundItem    *menu.ItemText
	GemsFoundTimer   = 0.8
	BombsMarked      int
	BombsMarkedItem  *menu.ItemText
	BombsMarkedTimer = 1.0
	WrongMarks       int
	WrongMarksItem   *menu.ItemText
	WrongMarksTimer  = 1.2
	TotalScore       *menu.ItemText
	TotalScoreTimer  = 1.4
	ScoreTimer       time.Time
)