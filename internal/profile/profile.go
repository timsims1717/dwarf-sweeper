package profile

import (
	"dwarf-sweeper/internal/data/player"
)

var (
	DefaultProfile = &player.Profile{
		Flags: &player.Flags{},
		BiomeExits: map[string]map[string]int{
			"mine": {
				"moss": 5,
			},
			"moss": {
				"crystal": 5,
			},
			"crystal": {
				"mine": 5,
				"dark": 1,
			},
			"dark": {},
		},
		SecretExit: map[string]float64{
			"mine": 0.2,
			"moss": 0.2,
			"crystal": 0.2,
			"dark": 0.2,
		},
	}
	CurrentProfile *player.Profile
)