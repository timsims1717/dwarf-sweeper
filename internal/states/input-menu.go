package states

import (
	"dwarf-sweeper/internal/config"
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/typeface"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"strings"
)

func InitInputMenu(win *pixelgl.Window) {
	InputMenu = menus.New("input", camera.Cam)
	InputMenu.Title = true
	InputMenu.SetCloseFn(config.SaveAsync)
	inputTitle := InputMenu.AddItem("title", "Input Options", false)
	device := InputMenu.AddItem("device", "Device", false)
	deviceR := InputMenu.AddItem("device_r", "", true)
	aimMode := InputMenu.AddItem("aim_mode", "Aim Mode", false)
	aimModeR := InputMenu.AddItem("aim_mode_r", "", true)
	digOn := InputMenu.AddItem("dig_on", "Dig On", false)
	digOnR := InputMenu.AddItem("dig_on_r", "", true)
	deadzone := InputMenu.AddItem("deadzone", "Deadzone", false)
	deadzoneR := InputMenu.AddItem("deadzone_r", fmt.Sprintf("%f", input.Deadzone), true)
	leftStickA := InputMenu.AddItem("left_stick_a", "Move with", false)
	leftStickR := InputMenu.AddItem("left_stick_r", "Yes", true)
	leftStickB := InputMenu.AddItem("left_stick_b", " Left Stick", false)
	left := InputMenu.AddItem("left", "Move Left", false)
	leftR := InputMenu.AddItem("left_r", "", true)
	right := InputMenu.AddItem("right", "Move Right", false)
	rightR := InputMenu.AddItem("right_r", "", true)
	up := InputMenu.AddItem("up", "Climb Up", false)
	upR := InputMenu.AddItem("up_r", "", true)
	down := InputMenu.AddItem("down", "Climb Down", false)
	downR := InputMenu.AddItem("down_r", "", true)
	jump := InputMenu.AddItem("jump", "Jump", false)
	jumpR := InputMenu.AddItem("jump_r", "", true)
	dig := InputMenu.AddItem("dig", "Dig", false)
	digR := InputMenu.AddItem("dig_r", "", true)
	flag := InputMenu.AddItem("flag", "Flag", false)
	flagR := InputMenu.AddItem("flag_r", "", true)
	interact := InputMenu.AddItem("interact", "Interact", false)
	interactR := InputMenu.AddItem("interact_r", "", true)
	use := InputMenu.AddItem("use", "Use Item", false)
	useR := InputMenu.AddItem("use_r", "", true)
	prev := InputMenu.AddItem("prev", "Prev Item", false)
	prevR := InputMenu.AddItem("prev_r", "", true)
	next := InputMenu.AddItem("next", "Next Item", false)
	nextR := InputMenu.AddItem("next_r", "", true)
	back := InputMenu.AddItem("back", "Back", false)

	aimModeUpdate := func() {
		if constants.AimDedicated {
			aimModeR.SetText("Dedicated")
			if data.GameInput.Mode == input.KeyboardMouse {
				aimMode.Hint = "Use the mouse to aim for digging, flagging, and attacking."
			} else {
				aimMode.Hint = "Use the right stick to aim for digging, flagging, and attacking."
			}
		} else {
			aimModeR.SetText("Movement")
			aimMode.Hint = "Use the movement keys to aim for digging, flagging, and attacking."
		}
	}
	aimModeUpdate()
	deviceUpdate := func() {
		km := data.GameInput.Mode == input.KeyboardMouse
		if km {
			device.Hint = ""
			deviceR.SetText("KB&Mouse")
		} else {
			device.Hint = win.JoystickName(data.GameInput.Joystick)
			deviceR.SetText(fmt.Sprintf("Gamepad %d", data.GameInput.Joystick+1))
		}
		leftStickA.Ignore = km
		leftStickB.Ignore = km
		leftStickR.Ignore = km
		deadzone.Ignore = km
		deadzoneR.Ignore = km
	}
	deviceUpdate()
	digOnUpdate := func() {
		if constants.DigOnRelease {
			digOnR.SetText("On Release")
			digOn.Hint = "Digging, Flagging, and Attacking happen when you release the button."
		} else {
			digOnR.SetText("On Press")
			digOn.Hint = "Digging, Flagging, and Attacking happen when you press the button."
		}
	}
	digOnUpdate()

	inputTitle.NoHover = true
	deviceSwitch := func(prev bool) {
		km := data.GameInput.Mode == input.KeyboardMouse
		var js int
		if prev {
			if km {
				js = input.PrevGamepad(win, -1)
			} else {
				js = input.PrevGamepad(win, int(data.GameInput.Joystick))
			}
		} else {
			if km {
				js = input.NextGamepad(win, -1)
			} else {
				js = input.NextGamepad(win, int(data.GameInput.Joystick))
			}
		}
		if js != -1 {
			data.GameInput.Joystick = pixelgl.Joystick(js)
			data.GameInput.Mode = input.Gamepad
		} else {
			data.GameInput.Joystick = pixelgl.JoystickLast
			data.GameInput.Mode = input.KeyboardMouse
		}
		UpdateKeybindings()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		aimModeUpdate()
		deviceUpdate()
	}
	rfn1 := func() {
		deviceSwitch(false)
	}
	lfn1 := func() {
		deviceSwitch(true)
	}
	device.SetClickFn(rfn1)
	device.SetRightFn(rfn1)
	device.SetLeftFn(lfn1)
	deviceR.NoHover = true
	rfn2 := func() {
		constants.AimDedicated = !constants.AimDedicated
		sfx.SoundPlayer.PlaySound("click", 2.0)
		aimModeUpdate()
	}
	aimMode.SetClickFn(rfn2)
	aimMode.SetRightFn(rfn2)
	aimMode.SetLeftFn(rfn2)
	aimModeR.NoHover = true
	fn3 := func() {
		constants.DigOnRelease = !constants.DigOnRelease
		sfx.SoundPlayer.PlaySound("click", 2.0)
		digOnUpdate()
	}
	digOn.SetClickFn(fn3)
	digOn.SetRightFn(fn3)
	digOn.SetLeftFn(fn3)
	digOnR.NoHover = true
	deadzone.SetRightFn(func() {
		n := input.Deadzone + 0.05
		if n > 0.5 {
			n = 0.5
		}
		input.Deadzone = n
		deadzoneR.SetText(fmt.Sprintf("%f", n))
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	deadzone.SetLeftFn(func() {
		n := input.Deadzone - 0.05
		if n < 0.05 {
			n = 0.05
		}
		input.Deadzone = n
		deadzoneR.SetText(fmt.Sprintf("%f", n))
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	deadzoneR.NoHover = true
	rfn3 := func() {
		data.GameInput.StickD = !data.GameInput.StickD
		if data.GameInput.StickD {
			leftStickR.SetText("Yes")
		} else {
			leftStickR.SetText("No")
		}
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	leftStickA.SetClickFn(rfn3)
	leftStickA.SetRightFn(rfn3)
	leftStickA.SetLeftFn(rfn3)
	leftStickA.SetHoverFn(func() {
		leftStickB.Hovered = true
	})
	leftStickA.SetUnhoverFn(func() {
		leftStickB.Hovered = false
	})
	leftStickB.NoHover = true
	leftStickR.NoHover = true
	if data.GameInput.StickD {
		leftStickR.SetText("Yes")
	} else {
		leftStickR.SetText("No")
	}
	keyFn := func(item *menus.Item) func() {
		return func() {
			OpenKeybindingMenu(item.Raw, item.Key)
		}
	}
	left.SetClickFn(keyFn(left))
	leftR.NoHover = true
	right.SetClickFn(keyFn(right))
	rightR.NoHover = true
	up.SetClickFn(keyFn(up))
	upR.NoHover = true
	down.SetClickFn(keyFn(down))
	downR.NoHover = true
	jump.SetClickFn(keyFn(jump))
	jumpR.NoHover = true
	dig.SetClickFn(keyFn(dig))
	digR.NoHover = true
	flag.SetClickFn(keyFn(flag))
	flagR.NoHover = true
	interact.SetClickFn(keyFn(interact))
	interactR.NoHover = true
	use.SetClickFn(keyFn(use))
	useR.NoHover = true
	next.SetClickFn(keyFn(next))
	nextR.NoHover = true
	prev.SetClickFn(keyFn(prev))
	prevR.NoHover = true

	InputMenu.SetOpenFn(func() {
		aimModeUpdate()
		deviceUpdate()
		UpdateKeybindings()
	})
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		InputMenu.Close()
	})
}

func InitKeybindingMenu() {
	KeybindingMenu = menus.New("keybinding", camera.Cam)
	KeybindingMenu.HideArrow = true
	KeybindingMenu.SetCloseFn(UpdateKeybindings)
	keybindingA := KeybindingMenu.AddItem("line_a", "Set key/button ", false)
	keybindingA.NoHover = true
	keybindingB := KeybindingMenu.AddItem("line_b", "", false)
	keybindingB.NoHover = true
}

func OpenKeybindingMenu(name, key string) {
	KeybindingMenu.ItemMap["line_b"].SetText(fmt.Sprintf("for %s", name))
	KeyString = key
	OpenMenu(KeybindingMenu)
	sfx.SoundPlayer.PlaySound("click", 2.0)
}

func UpdateKeybindings() {
	UpdateKeybinding("left")
	UpdateKeybinding("right")
	UpdateKeybinding("up")
	UpdateKeybinding("down")
	UpdateKeybinding("jump")
	UpdateKeybinding("dig")
	UpdateKeybinding("flag")
	UpdateKeybinding("interact")
	UpdateKeybinding("use")
	UpdateKeybinding("prev")
	UpdateKeybinding("next")
}

func UpdateKeybinding(key string) {
	r := InputMenu.ItemMap[fmt.Sprintf("%s_r", key)]
	in := data.GameInput.Buttons[key]
	builder := strings.Builder{}
	first := true
	if data.GameInput.Mode != input.Gamepad {
		for _, k := range in.Keys {
			if first {
				first = false
				typeface.RegisterSymbol(key, img.Batchers[constants.MenuSprites].GetSprite(k.String()), 1.)
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString(fmt.Sprintf("{symbol:%s}", k.String()))
		}
		if in.Scroll > 0 {
			if first {
				first = false
				typeface.RegisterSymbol(key, img.Batchers[constants.MenuSprites].GetSprite("MouseScrollUp"), 1.)
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString("{symbol:MouseScrollUp}")
		} else if in.Scroll < 0 {
			if first {
				first = false
				typeface.RegisterSymbol(key, img.Batchers[constants.MenuSprites].GetSprite("MouseScrollDown"), 1.)
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString("{symbol:MouseScrollDown}")
		}
	}
	if data.GameInput.Mode != input.KeyboardMouse {
		for _, b := range in.Buttons {
			if first {
				first = false
				typeface.RegisterSymbol(key, img.Batchers[constants.MenuSprites].GetSprite(input.GamepadString(b)), 1.)
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString(fmt.Sprintf("{symbol:%s}", input.GamepadString(b)))
		}
		if in.AxisV != 0 {
			if first {
				first = false
				typeface.RegisterSymbol(key, img.Batchers[constants.MenuSprites].GetSprite(input.AxisDirString(in.Axis, in.AxisV > 0)), 1.)
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString(fmt.Sprintf("{symbol:%s}", input.AxisDirString(in.Axis, in.AxisV > 0)))
		}
	}
	r.SetText(builder.String())
}
