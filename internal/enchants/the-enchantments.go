package enchants

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
)

var Enchantments = []*data.Enchantment{
	{
		OnGain: func() {
			descent.Descent.GetPlayer().MaxJump++
		},
		OnLose: func() {
			descent.Descent.GetPlayer().MaxJump--
		},
		Key:   "jump1",
		Title: "Jumping",
		Desc:  "Increases your jump height.",
	}, {
		OnGain: func() {
			descent.Descent.GetPlayer().MaxJump++
		},
		OnLose: func() {
			descent.Descent.GetPlayer().MaxJump--
		},
		Key:     "jump2",
		Title:   "Jumping II",
		Desc:    "Increases your jump height even more.",
		Require: "jump1",
	},
	{
		OnGain: func() {
			descent.Descent.GetPlayer().ClimbSpeed += 15.
		},
		OnLose: func() {
			descent.Descent.GetPlayer().ClimbSpeed -= 15.
		},
		Key:   "climb1",
		Title: "Clambering",
		Desc:  "Increases your climb speed.",
	},
	{
		OnGain: func() {
			descent.Descent.GetPlayer().Speed += 25.
		},
		OnLose: func() {
			descent.Descent.GetPlayer().Speed -= 25.
		},
		Key:   "run1",
		Title: "Running",
		Desc:  "Increases your running speed.",
	},
	{
		OnGain: func() {
			descent.Descent.GetPlayer().Health.Max += 1
			descent.Descent.GetPlayer().Health.Curr += 1
		},
		OnLose: func() {
			descent.Descent.GetPlayer().Health.Max -= 1
		},
		Key:   "health1",
		Title: "Heartiness",
		Desc:  "Gives you an extra heart.",
	},
	{
		OnGain: func() {
			descent.Descent.GetPlayer().Health.Max += 1
			descent.Descent.GetPlayer().Health.Curr += 1
		},
		OnLose: func() {
			descent.Descent.GetPlayer().Health.Max -= 1
		},
		Key:     "health2",
		Title:   "Heartiness II",
		Desc:    "Gives you an extra heart.",
		Require: "health1",
	},
	{
		OnGain: func() {
			descent.Descent.GetPlayer().ShovelDamage += 1
		},
		OnLose: func() {
			descent.Descent.GetPlayer().ShovelDamage -= 1
		},
		Key:   "damage1",
		Title: "Sharpness",
		Desc:  "Shovel deals damage to enemies.",
	},
	{
		OnGain: func() {
			descent.Descent.GetPlayer().ShovelDamage += 1
		},
		OnLose: func() {
			descent.Descent.GetPlayer().ShovelDamage -= 1
		},
		Key:     "damage2",
		Title:   "Sharpness II",
		Desc:    "Shovel deals increased damage to enemies.",
		Require: "damage1",
	},
	{
		OnGain: func() {
			descent.Descent.GetPlayer().ShovelKnockback += 2.
		},
		OnLose: func() {
			descent.Descent.GetPlayer().ShovelKnockback -= 2.
		},
		Key:   "knockback1",
		Title: "Batting",
		Desc:  "Increases your shovel's knockback.",
	},
	{
		OnGain: func() {
			descent.Descent.GetPlayer().ShovelKnockback += 2.
		},
		OnLose: func() {
			descent.Descent.GetPlayer().ShovelKnockback -= 2.
		},
		Key:     "knockback2",
		Title:   "Batting II",
		Desc:    "Increases your shovel's knockback even more.",
		Require: "knockback1",
	},
	{
		OnGain: func() {
			descent.Descent.GetPlayer().GemRate += 0.5
		},
		OnLose: func() {
			descent.Descent.GetPlayer().GemRate -= 0.5
		},
		Key:   "gemmagnet1",
		Title: "Gem Magnet",
		Desc:  "Increases your chance to find Gems.",
	},
	{
		OnGain: func() {
			descent.Descent.GetPlayer().GemRate += 1.0
		},
		OnLose: func() {
			descent.Descent.GetPlayer().GemRate -= 1.0
		},
		Key:     "gemmagnet2",
		Title:   "Gem Magnet II",
		Desc:    "Increases your chance to find Gems even more.",
		Require: "gemmagnet1",
	},
}
