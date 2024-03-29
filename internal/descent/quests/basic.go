package quests

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/profile"
)

var (
	Flag5 = &data.Quest{
		Key:  "flag5",
		Name: "Beginner Sweeper",
		Desc: "Correctly flag 5 total bombs.",
		Check: func(p *data.Profile) bool {
			return p.Flags.CorrectFlags + p.Stats.CorrectFlags >= 5
		},
		OnFinish: func(p *data.Profile) {
			profile.AddQuest(p, Flag50)
		},
		Hidden: false,
	}
	Flag50 = &data.Quest{
		Key:  "flag50",
		Name: "Novice Sweeper",
		Desc: "Correctly flag 50 total bombs.",
		Check: func(p *data.Profile) bool {
			return p.Flags.CorrectFlags + p.Stats.CorrectFlags >= 50
		},
		OnFinish: func(p *data.Profile) {
			profile.AddQuest(p, Flag150)
		},
		Hidden: false,
	}
	Flag150 = &data.Quest{
		Key:  "flag150",
		Name: "Apprentice Sweeper",
		Desc: "Correctly flag 150 total bombs.",
		Check: func(p *data.Profile) bool {
			return p.Flags.CorrectFlags + p.Stats.CorrectFlags >= 150
		},
		OnFinish: func(p *data.Profile) {
			profile.AddQuest(p, Flag500)
		},
		Hidden: false,
	}
	Flag500 = &data.Quest{
		Key:  "flag500",
		Name: "Journeyman Sweeper",
		Desc: "Correctly flag 500 total bombs.",
		Check: func(p *data.Profile) bool {
			return p.Flags.CorrectFlags + p.Stats.CorrectFlags >= 500
		},
		OnFinish: func(p *data.Profile) {
			profile.AddQuest(p, Flag1000)
		},
		Hidden: false,
	}
	Flag1000 = &data.Quest{
		Key:  "flag1000",
		Name: "Master Sweeper",
		Desc: "Correctly flag 1000 total bombs.",
		Check: func(p *data.Profile) bool {
			return p.Flags.CorrectFlags + p.Stats.CorrectFlags >= 1000
		},
		OnFinish: func(p *data.Profile) {

		},
		Hidden: false,
	}
)