package enchants

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/dungeon"
)

func AddEnchantment(e1 *data.Enchantment) {
	for i, e := range Enchantments {
		if e.Key == e1.Key {
			for _, in := range dungeon.Dungeon.GetPlayer().Enchants {
				if in == i {
					return
				}
			}
			e1.OnGain()
			dungeon.Dungeon.GetPlayer().Enchants = append(dungeon.Dungeon.GetPlayer().Enchants, i)
			return
		}
	}
}

func RemoveEnchantment(e1 *data.Enchantment) {
	for i, e := range Enchantments {
		if e.Key == e1.Key {
			index := -1
			for j, in := range dungeon.Dungeon.GetPlayer().Enchants {
				if in == i {
					index = j
				}
			}
			if index == -1 {
				return
			}
			e1.OnLose()
			if len(dungeon.Dungeon.GetPlayer().Enchants) > 1 {
				dungeon.Dungeon.GetPlayer().Enchants = append(dungeon.Dungeon.GetPlayer().Enchants[:index], dungeon.Dungeon.GetPlayer().Enchants[index+1:]...)
			} else {
				dungeon.Dungeon.GetPlayer().Enchants = []int{}
			}
			return
		}
	}
}

var Enchantments = []*data.Enchantment{
	{
		OnGain: func () {
			dungeon.Dungeon.GetPlayer().MaxJump++
		},
		OnLose: func () {
			dungeon.Dungeon.GetPlayer().MaxJump--
		},
		Key:    "jump1",
		Title:  "Jumping",
		Desc:   "Increases jump height.",
	},
	{
		OnGain: func () {
			dungeon.Dungeon.GetPlayer().ClimbSpeed += 15.
		},
		OnLose: func () {
			dungeon.Dungeon.GetPlayer().ClimbSpeed -= 15.
		},
		Key:    "climb1",
		Title:  "Clambering",
		Desc:   "Increases climb speed.",
	},
	{
		OnGain: func () {
			dungeon.Dungeon.GetPlayer().Speed += 25.
		},
		OnLose: func () {
			dungeon.Dungeon.GetPlayer().Speed -= 25.
		},
		Key:    "run1",
		Title:  "Running",
		Desc:   "Increases running speed.",
	},
	{
		OnGain: func () {
			dungeon.Dungeon.GetPlayer().Health.Max += 1
			dungeon.Dungeon.GetPlayer().Health.Curr += 1
		},
		OnLose: func () {
			dungeon.Dungeon.GetPlayer().Health.Max -= 1
		},
		Key:    "health",
		Title:  "Heartiness",
		Desc:   "Increases max health.",
	},
	{
		OnGain: func () {
			dungeon.Dungeon.GetPlayer().ShovelDamage += 1
		},
		OnLose: func () {
			dungeon.Dungeon.GetPlayer().ShovelDamage -= 1
		},
		Key:    "damage",
		Title:  "Sharpness",
		Desc:   "Shovel deals damage to enemies.",
	},
}