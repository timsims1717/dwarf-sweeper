package quests

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/profile"
	"dwarf-sweeper/pkg/util"
)

var (
	DiscoverMoss = &data.Quest{
		Key:  "discoverMoss",
		Name: "The Mossy Caverns",
		Desc: "Find an entrance to the Mossy Cavern Biome.",
		Check: func(p *data.Profile) bool {
			return util.ContainsStr("moss", p.Flags.Discovered)
		},
		OnFinish: func(p *data.Profile) {
			p.BiomeExits["mine"]["moss"] = 20
			profile.AddQuest(p, DiscoverCrystal)
		},
		Hidden: false,
	}
	DiscoverCrystal = &data.Quest{
		Key:  "discoverCrystal",
		Name: "The Crystal Hoard",
		Desc: "Find an entrance to the Crystal Hoard Biome.",
		Check: func(p *data.Profile) bool {
			return util.ContainsStr("crystal", p.Flags.Discovered)
		},
		OnFinish: func(p *data.Profile) {
			p.BiomeExits["moss"]["crystal"] = 20
			profile.AddQuest(p, CrystalToMine)
			profile.AddQuest(p, DiscoverDark)
		},
		Hidden: false,
	}
	CrystalToMine = &data.Quest{
		Key:  "crystalBackToMine",
		Name: "Crystal Shortcut",
		Desc: "Find an entrance back to the Mine Biome.",
		Check: func(p *data.Profile) bool {
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
		OnFinish: func(p *data.Profile) {
			p.BiomeExits["mine"]["crystal"] = 20
			p.BiomeExits["crystal"]["dark"] = 5
		},
		Hidden: false,
	}
	DiscoverDark = &data.Quest{
		Key:  "discoverDark",
		Name: "The Dark",
		Desc: "Find an entrance to the Dark Biome.",
		Check: func(p *data.Profile) bool {
			return util.ContainsStr("dark", p.Flags.Discovered)
		},
		OnFinish: func(p *data.Profile) {
			p.BiomeExits["crystal"]["dark"] = 20
		},
		Hidden: false,
	}
)