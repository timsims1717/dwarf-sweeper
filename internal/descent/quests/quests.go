package quests

import (
	"dwarf-sweeper/internal/data/player"
)

var Quests map[string]*player.Quest

func init() {
	Quests = make(map[string]*player.Quest)
	Quests[Flag5.Key] = Flag5
	Quests[Flag50.Key] = Flag50
	Quests[Flag150.Key] = Flag150
	Quests[Flag500.Key] = Flag500
	Quests[Flag1000.Key] = Flag1000
	Quests[DiscoverMoss.Key] = DiscoverMoss
	Quests[DiscoverCrystal.Key] = DiscoverCrystal
	Quests[CrystalToMine.Key] = CrystalToMine
	Quests[DiscoverDark.Key] = DiscoverDark
}