package quests

import (
	"dwarf-sweeper/internal/data/player"
)

var (
	Flag5 = &player.Quest{
		Key:  "flag5",
		Name: "Beginner Sweeper",
		Desc: "Correctly flag 5 total bombs.",
		Check: func(p *player.Profile) bool {
			return p.Flags.CorrectFlags >= 5 || p.Stats.CorrectFlags >= 5
		},
		OnFinish: func(p *player.Profile) {
			p.AddQuest(Flag50)
		},
		Hidden: false,
	}
	Flag50 = &player.Quest{
		Key:  "flag50",
		Name: "Novice Sweeper",
		Desc: "Correctly flag 50 total bombs.",
		Check: func(p *player.Profile) bool {
			return p.Flags.CorrectFlags >= 50
		},
		OnFinish: func(p *player.Profile) {
			p.AddQuest(Flag150)
		},
		Hidden: false,
	}
	Flag150 = &player.Quest{
		Key:  "flag150",
		Name: "Apprentice Sweeper",
		Desc: "Correctly flag 150 total bombs.",
		Check: func(p *player.Profile) bool {
			return p.Flags.CorrectFlags >= 150
		},
		OnFinish: func(p *player.Profile) {
			p.AddQuest(Flag500)
		},
		Hidden: false,
	}
	Flag500 = &player.Quest{
		Key:  "flag500",
		Name: "Journeyman Sweeper",
		Desc: "Correctly flag 500 total bombs.",
		Check: func(p *player.Profile) bool {
			return p.Flags.CorrectFlags >= 500
		},
		OnFinish: func(p *player.Profile) {
			p.AddQuest(Flag1000)
		},
		Hidden: false,
	}
	Flag1000 = &player.Quest{
		Key:  "flag1000",
		Name: "Master Sweeper",
		Desc: "Correctly flag 1000 total bombs.",
		Check: func(p *player.Profile) bool {
			return p.Flags.CorrectFlags >= 1000
		},
		OnFinish: func(p *player.Profile) {

		},
		Hidden: false,
	}
)