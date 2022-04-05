package quests

import (
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/pkg/util"
)

var (
	DiscoverMoss = &player.Quest{
		Key:  "discoverMoss",
		Name: "The Mossy Caverns",
		Desc: "Find an entrance to the Mossy Cavern Biome.",
		Check: func(p *player.Profile) bool {
			return util.ContainsStr("moss", p.Flags.Discovered)
		},
		OnFinish: func(p *player.Profile) {
			p.BiomeExits["mine"]["moss"] = 20
			p.SecretExit["mine"] = 0.5
			p.AddQuest(DiscoverCrystal)
		},
		Hidden: false,
	}
	DiscoverCrystal = &player.Quest{
		Key:  "discoverCrystal",
		Name: "The Crystal Hoard",
		Desc: "Find an entrance to the Crystal Hoard Biome.",
		Check: func(p *player.Profile) bool {
			return util.ContainsStr("crystal", p.Flags.Discovered)
		},
		OnFinish: func(p *player.Profile) {
			p.BiomeExits["moss"]["crystal"] = 20
			p.SecretExit["moss"] = 0.5
			p.AddQuest(CrystalToMine)
			p.AddQuest(DiscoverDark)
		},
		Hidden: false,
	}
	CrystalToMine = &player.Quest{
		Key:  "crystalBackToMine",
		Name: "Crystal Shortcut",
		Desc: "Find an entrance back to the Mine Biome.",
		Check: func(p *player.Profile) bool {
			crystal := false
			for _, biome := range descent.Descent.BiomeOrder {
				if crystal && biome == "mine" {
					return true
				} else if biome == "crystal" {
					crystal = true
				} else {
					crystal = false
				}
			}
			return false
		},
		OnFinish: func(p *player.Profile) {
			p.BiomeExits["mine"]["crystal"] = 20
			p.BiomeExits["crystal"]["dark"] = 5
		},
		Hidden: false,
	}
	DiscoverDark = &player.Quest{
		Key:  "discoverDark",
		Name: "The Dark",
		Desc: "Find an entrance to the Dark Biome.",
		Check: func(p *player.Profile) bool {
			return util.ContainsStr("dark", p.Flags.Discovered)
		},
		OnFinish: func(p *player.Profile) {
			p.BiomeExits["crystal"]["dark"] = 20
			p.SecretExit["crystal"] = 0.5
		},
		Hidden: false,
	}
)