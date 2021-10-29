package state

import (
	"dwarf-sweeper/internal/config"
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/pkg/camera"
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
	inputTitle := InputMenu.AddItem("title", "Input Options")
	device := InputMenu.AddItem("device", "Device")
	deviceR := InputMenu.AddItem("device_r", "")
	aimMode := InputMenu.AddItem("aim_mode", "Aim Mode")
	aimModeR := InputMenu.AddItem("aim_mode_r", "")
	digOn := InputMenu.AddItem("dig_on", "Dig On")
	digOnR := InputMenu.AddItem("dig_on_r", "")
	deadzone := InputMenu.AddItem("deadzone", "Deadzone")
	deadzoneR := InputMenu.AddItem("deadzone_r", fmt.Sprintf("%f", input.Deadzone))
	leftStickA := InputMenu.AddItem("left_stick_a", "Move with")
	leftStickR := InputMenu.AddItem("left_stick_r", "Yes")
	leftStickB := InputMenu.AddItem("left_stick_b", " Left Stick")
	left := InputMenu.AddItem("left", "Move Left")
	leftR := InputMenu.AddItem("left_r", "")
	right := InputMenu.AddItem("right", "Move Right")
	rightR := InputMenu.AddItem("right_r", "")
	up := InputMenu.AddItem("up", "Climb Up")
	upR := InputMenu.AddItem("up_r", "")
	down := InputMenu.AddItem("down", "Climb Down")
	downR := InputMenu.AddItem("down_r", "")
	jump := InputMenu.AddItem("jump", "Jump")
	jumpR := InputMenu.AddItem("jump_r", "")
	dig := InputMenu.AddItem("dig", "Dig")
	digR := InputMenu.AddItem("dig_r", "")
	mark := InputMenu.AddItem("mark", "Mark")
	markR := InputMenu.AddItem("mark_r", "")
	interact := InputMenu.AddItem("interact", "Interact")
	interactR := InputMenu.AddItem("interact_r", "")
	use := InputMenu.AddItem("use", "Use Item")
	useR := InputMenu.AddItem("use_r", "")
	prev := InputMenu.AddItem("prev", "Prev Item")
	prevR := InputMenu.AddItem("prev_r", "")
	next := InputMenu.AddItem("next", "Next Item")
	nextR := InputMenu.AddItem("next_r", "")
	back := InputMenu.AddItem("back", "Back")

	aimModeUpdate := func() {
		if constants.AimDedicated {
			aimModeR.Raw = "Dedicated"
			if data.GameInput.Mode == input.KeyboardMouse {
				aimMode.Hint = "Use the mouse to aim for digging, marking, and attacking."
			} else {
				aimMode.Hint = "Use the right stick to aim for digging, marking, and attacking."
			}
		} else {
			aimModeR.Raw = "Movement"
			aimMode.Hint = "Use the movement keys to aim for digging, marking, and attacking."
		}
	}
	aimModeUpdate()
	deviceUpdate := func() {
		km := data.GameInput.Mode == input.KeyboardMouse
		if km {
			device.Hint = ""
			deviceR.Raw = "KB&Mouse"
		} else {
			device.Hint = win.JoystickName(data.GameInput.Joystick)
			deviceR.Raw = fmt.Sprintf("Gamepad %d", data.GameInput.Joystick+1)
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
			digOnR.Raw = "On Release"
			digOn.Hint = "Digging, Marking, and Attacking happen when you release the button."
		} else {
			digOnR.Raw = "On Press"
			digOn.Hint = "Digging, Marking, and Attacking happen when you press the button."
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
	deviceR.Right = true
	deviceR.NoHover = true
	rfn2 := func() {
		constants.AimDedicated = !constants.AimDedicated
		sfx.SoundPlayer.PlaySound("click", 2.0)
		aimModeUpdate()
	}
	aimMode.SetClickFn(rfn2)
	aimMode.SetRightFn(rfn2)
	aimMode.SetLeftFn(rfn2)
	aimModeR.Right = true
	aimModeR.NoHover = true
	fn3 := func() {
		constants.DigOnRelease = !constants.DigOnRelease
		sfx.SoundPlayer.PlaySound("click", 2.0)
		digOnUpdate()
	}
	digOn.SetClickFn(fn3)
	digOn.SetRightFn(fn3)
	digOn.SetLeftFn(fn3)
	digOnR.Right = true
	digOnR.NoHover = true
	deadzone.SetRightFn(func() {
		n := input.Deadzone + 0.05
		if n > 0.5 {
			n = 0.5
		}
		input.Deadzone = n
		deadzoneR.Raw = fmt.Sprintf("%f", n)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	deadzone.SetLeftFn(func() {
		n := input.Deadzone - 0.05
		if n < 0.05 {
			n = 0.05
		}
		input.Deadzone = n
		deadzoneR.Raw = fmt.Sprintf("%f", n)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	deadzoneR.NoHover = true
	deadzoneR.Right = true
	rfn3 := func() {
		data.GameInput.StickD = !data.GameInput.StickD
		if data.GameInput.StickD {
			leftStickR.Raw = "Yes"
		} else {
			leftStickR.Raw = "No"
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
	leftStickR.Right = true
	if data.GameInput.StickD {
		leftStickR.Raw = "Yes"
	} else {
		leftStickR.Raw = "No"
	}
	keyFn := func(item *menus.Item) func() {
		return func() {
			OpenKeybindingMenu(item.Raw, item.Key)
		}
	}
	left.SetClickFn(keyFn(left))
	leftR.Right = true
	leftR.NoHover = true
	right.SetClickFn(keyFn(right))
	rightR.Right = true
	rightR.NoHover = true
	up.SetClickFn(keyFn(up))
	upR.Right = true
	upR.NoHover = true
	down.SetClickFn(keyFn(down))
	downR.Right = true
	downR.NoHover = true
	jump.SetClickFn(keyFn(jump))
	jumpR.Right = true
	jumpR.NoHover = true
	dig.SetClickFn(keyFn(dig))
	digR.Right = true
	digR.NoHover = true
	mark.SetClickFn(keyFn(mark))
	markR.Right = true
	markR.NoHover = true
	interact.SetClickFn(keyFn(interact))
	interactR.Right = true
	interactR.NoHover = true
	use.SetClickFn(keyFn(use))
	useR.Right = true
	useR.NoHover = true
	next.SetClickFn(keyFn(next))
	nextR.Right = true
	nextR.NoHover = true
	prev.SetClickFn(keyFn(prev))
	prevR.Right = true
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
	keybindingA := KeybindingMenu.AddItem("line_a", "Set key/button ")
	keybindingA.NoHover = true
	keybindingB := KeybindingMenu.AddItem("line_b", "")
	keybindingB.NoHover = true
}

func OpenKeybindingMenu(name, key string) {
	KeybindingMenu.ItemMap["line_b"].Raw = fmt.Sprintf("for %s", name)
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
	UpdateKeybinding("mark")
	UpdateKeybinding("interact")
	UpdateKeybinding("use")
	UpdateKeybinding("prev")
	UpdateKeybinding("next")
}

func UpdateKeybinding(key string) {
	r := InputMenu.ItemMap[fmt.Sprintf("%s_r", key)]
	in := data.GameInput.Buttons[key]
	builder := strings.Builder{}
	var symKeys []string
	first := true
	if data.GameInput.Mode != input.Gamepad {
		for _, k := range in.Keys {
			if first {
				first = false
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString(typeface.SymbolItem)
			symKeys = append(symKeys, k.String())
		}
		if in.Scroll > 0 {
			if first {
				first = false
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString(typeface.SymbolItem)
			symKeys = append(symKeys, "MouseScrollUp")
		} else if in.Scroll < 0 {
			if first {
				first = false
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString(typeface.SymbolItem)
			symKeys = append(symKeys, "MouseScrollDown")
		}
	}
	if data.GameInput.Mode != input.KeyboardMouse {
		for _, b := range in.Buttons {
			if first {
				first = false
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString(typeface.SymbolItem)
			symKeys = append(symKeys, input.GamepadString(b))
		}
		if in.AxisV != 0 {
			if first {
				first = false
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString(typeface.SymbolItem)
			symKeys = append(symKeys, input.AxisDirString(in.Axis, in.AxisV > 0))
		}
	}
	r.Raw = builder.String()
	r.Symbols = symKeys
}
