package quests

import (
	"dwarf-sweeper/internal/data"
)

func init() {
	data.Quests = make(map[string]*data.Quest)
	data.Quests[Flag5.Key] = Flag5
	data.Quests[Flag50.Key] = Flag50
	data.Quests[Flag150.Key] = Flag150
	data.Quests[Flag500.Key] = Flag500
	data.Quests[Flag1000.Key] = Flag1000
	data.Quests[DiscoverMoss.Key] = DiscoverMoss
	data.Quests[DiscoverCrystal.Key] = DiscoverCrystal
	data.Quests[CrystalToMine.Key] = CrystalToMine
	data.Quests[DiscoverDark.Key] = DiscoverDark
}