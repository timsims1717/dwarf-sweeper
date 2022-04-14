package quests

import (
	"dwarf-sweeper/internal/data/player"
)

func init() {
	player.Quests = make(map[string]*player.Quest)
	player.Quests[Flag5.Key] = Flag5
	player.Quests[Flag50.Key] = Flag50
	player.Quests[Flag150.Key] = Flag150
	player.Quests[Flag500.Key] = Flag500
	player.Quests[Flag1000.Key] = Flag1000
	player.Quests[DiscoverMoss.Key] = DiscoverMoss
	player.Quests[DiscoverCrystal.Key] = DiscoverCrystal
	player.Quests[CrystalToMine.Key] = CrystalToMine
	player.Quests[DiscoverDark.Key] = DiscoverDark
}