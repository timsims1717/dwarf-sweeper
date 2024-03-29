package profile

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/pkg/util"
	"fmt"
)

var (
	DefaultProfile = &data.Profile{
		Flags: data.Flags{},
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
		ItemLimits: data.ItemLimits{
			Hold: map[string]int{
				"bomb_item": 8,
				"beer": 3,
				"throw_shovel": 5,
			},
			Uses: map[string]int{
				"pickaxe": 8,
				"detector": 4,
			},
			Secs: map[string]float64{
				"xray": 16.,
			},
		},
		StartingAttr: data.StartingAttr{
			MaxHealth: 3,
		},
	}
	CurrentProfile *data.Profile
)

func AddQuest(p *data.Profile, q *data.Quest) {
	if !util.ContainsStr(q.Key, p.Quests) {
		p.Quests = append(p.Quests, q.Key)
		if !q.Hidden {
			p.QuestsShown = append(p.QuestsShown, q.Key)
			menus.NotificationHandler.AddMessage(fmt.Sprintf("New Quest: %s", q.Name))
		}
	}
}

func UpdateQuests(p *data.Profile) {
	for _, key := range p.Quests {
		q := data.Quests[key]
		if !util.ContainsStr(key, p.QuestsComplete) && q.Check(p) {
			p.QuestsComplete = append(p.QuestsComplete, key)
			menus.NotificationHandler.AddMessage(fmt.Sprintf("Quest: %s {symbol:checkmark}", q.Name))
			if q.OnFinish != nil {
				q.OnFinish(p)
			}
			// add to notifications
		}
	}
}