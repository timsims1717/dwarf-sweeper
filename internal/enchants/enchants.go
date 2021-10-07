package enchants

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
)

func AddEnchantment(e1 *data.Enchantment) {
	for i, e := range Enchantments {
		if e.Key == e1.Key {
			for _, in := range descent.Descent.GetPlayer().Enchants {
				if in == i {
					return
				}
			}
			e1.OnGain()
			descent.Descent.GetPlayer().Enchants = append(descent.Descent.GetPlayer().Enchants, i)
			return
		}
	}
}

func RemoveEnchantment(e1 *data.Enchantment) {
	for i, e := range Enchantments {
		if e.Key == e1.Key {
			index := -1
			for j, in := range descent.Descent.GetPlayer().Enchants {
				if in == i {
					index = j
				}
			}
			if index == -1 {
				return
			}
			e1.OnLose()
			if len(descent.Descent.GetPlayer().Enchants) > 1 {
				descent.Descent.GetPlayer().Enchants = append(descent.Descent.GetPlayer().Enchants[:index], descent.Descent.GetPlayer().Enchants[index+1:]...)
			} else {
				descent.Descent.GetPlayer().Enchants = []int{}
			}
			return
		}
	}
}

var Enchantments = []*data.Enchantment{
	{
		OnGain: func () {
			descent.Descent.GetPlayer().MaxJump++
		},
		OnLose: func () {
			descent.Descent.GetPlayer().MaxJump--
		},
		Key:    "jump1",
		Title:  "Jumping",
		Desc:   "Increases your jump height.",
	},
	{
		OnGain: func () {
			descent.Descent.GetPlayer().ClimbSpeed += 15.
		},
		OnLose: func () {
			descent.Descent.GetPlayer().ClimbSpeed -= 15.
		},
		Key:    "climb1",
		Title:  "Clambering",
		Desc:   "Increases your climb speed.",
	},
	{
		OnGain: func () {
			descent.Descent.GetPlayer().Speed += 25.
		},
		OnLose: func () {
			descent.Descent.GetPlayer().Speed -= 25.
		},
		Key:    "run1",
		Title:  "Running",
		Desc:   "Increases your running speed.",
	},
	{
		OnGain: func () {
			descent.Descent.GetPlayer().Health.Max += 1
			descent.Descent.GetPlayer().Health.Curr += 1
		},
		OnLose: func () {
			descent.Descent.GetPlayer().Health.Max -= 1
		},
		Key:    "health",
		Title:  "Heartiness",
		Desc:   "Increases your max health.",
	},
	{
		OnGain: func () {
			descent.Descent.GetPlayer().ShovelDamage += 1
		},
		OnLose: func () {
			descent.Descent.GetPlayer().ShovelDamage -= 1
		},
		Key:    "damage",
		Title:  "Sharpness",
		Desc:   "Shovel deals damage to enemies.",
	},
	{
		OnGain: func () {
			descent.Descent.GetPlayer().ShovelKnockback += 0.3
		},
		OnLose: func () {
			descent.Descent.GetPlayer().ShovelKnockback -= 0.3
		},
		Key:    "knockback",
		Title:  "Batting",
		Desc:   "Increases your shovel's knockback.",
	},
}