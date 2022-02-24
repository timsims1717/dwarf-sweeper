package descent

var Enchantments = []*Enchantment{
	{
		OnGain: func(d *Dwarf) {
			d.MaxJump++
		},
		OnLose: func(d *Dwarf) {
			d.MaxJump--
		},
		Key:   "jump1",
		Title: "Jumping",
		Desc:  "Increases your jump height.",
	}, {
		OnGain: func(d *Dwarf) {
			d.MaxJump++
		},
		OnLose: func(d *Dwarf) {
			d.MaxJump--
		},
		Key:     "jump2",
		Title:   "Jumping II",
		Desc:    "Increases your jump height even more.",
		Require: "jump1",
	},
	{
		OnGain: func(d *Dwarf) {
			d.ClimbSpeed += 15.
		},
		OnLose: func(d *Dwarf) {
			d.ClimbSpeed -= 15.
		},
		Key:   "climb1",
		Title: "Clambering",
		Desc:  "Increases your climb speed.",
	},
	{
		OnGain: func(d *Dwarf) {
			d.Speed += 25.
		},
		OnLose: func(d *Dwarf) {
			d.Speed -= 25.
		},
		Key:   "run1",
		Title: "Running",
		Desc:  "Increases your running speed.",
	},
	{
		OnGain: func(d *Dwarf) {
			d.Health.Max += 1
			d.Health.Curr += 1
		},
		OnLose: func(d *Dwarf) {
			d.Health.Max -= 1
		},
		Key:   "health1",
		Title: "Heartiness",
		Desc:  "Gives you an extra heart.",
	},
	{
		OnGain: func(d *Dwarf) {
			d.Health.Max += 1
			d.Health.Curr += 1
		},
		OnLose: func(d *Dwarf) {
			d.Health.Max -= 1
		},
		Key:     "health2",
		Title:   "Heartiness II",
		Desc:    "Gives you an extra heart.",
		Require: "health1",
	},
	{
		OnGain: func(d *Dwarf) {
			d.ShovelKnockback += 2.
		},
		OnLose: func(d *Dwarf) {
			d.ShovelKnockback -= 2.
		},
		Key:   "knockback1",
		Title: "Batting",
		Desc:  "Increases your shovel's knockback.",
	},
	{
		OnGain: func(d *Dwarf) {
			d.ShovelKnockback += 2.
		},
		OnLose: func(d *Dwarf) {
			d.ShovelKnockback -= 2.
		},
		Key:     "knockback2",
		Title:   "Batting II",
		Desc:    "Increases your shovel's knockback even more.",
		Require: "knockback1",
	},
	{
		OnGain: func(d *Dwarf) {
			d.Player.Attr.GemRate += 0.5
		},
		OnLose: func(d *Dwarf) {
			d.Player.Attr.GemRate -= 0.5
		},
		Key:   "gemmagnet1",
		Title: "Gem Magnet",
		Desc:  "Increases your chance to find Gems.",
	},
	{
		OnGain: func(d *Dwarf) {
			d.Player.Attr.GemRate += 1.0
		},
		OnLose: func(d *Dwarf) {
			d.Player.Attr.GemRate -= 1.0
		},
		Key:     "gemmagnet2",
		Title:   "Gem Magnet II",
		Desc:    "Increases your chance to find Gems even more.",
		Require: "gemmagnet1",
	},
}
