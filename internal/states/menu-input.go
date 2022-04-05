package states

import (
	"dwarf-sweeper/internal/config"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/sfx"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
)

func InitInputMenu(win *pixelgl.Window) {
	InputMenu = menus.New("input", camera.Cam)
	InputMenu.Title = true
	InputMenu.SetCloseFn(config.SaveAsync)
	inputTitle := InputMenu.AddItem("title", "Input Options", false)
	player := InputMenu.AddItem("player", "< Player 1 >", false)
	device := InputMenu.AddItem("device", "Device", false)
	deviceR := InputMenu.AddItem("device_r", "", true)
	aimMode := InputMenu.AddItem("aim_mode", "Aim Mode", false)
	aimModeR := InputMenu.AddItem("aim_mode_r", "", true)
	digOn := InputMenu.AddItem("dig_on", "Dig On", false)
	digOnR := InputMenu.AddItem("dig_on_r", "", true)
	deadzone := InputMenu.AddItem("deadzone", "Deadzone", false)
	deadzoneR := InputMenu.AddItem("deadzone_r", fmt.Sprintf("%f", data.CurrInput.Deadzone), true)
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
	puzzLeave := InputMenu.AddItem("puzz_leave", "Leave Puzzle", false)
	puzzLeaveR := InputMenu.AddItem("puzz_leave_r", "", true)
	puzzHelp := InputMenu.AddItem("puzz_help", "Puzzle Help", false)
	puzzHelpR := InputMenu.AddItem("puzz_help_r", "", true)
	minePuzz := InputMenu.AddItem("mine_puzz", "Mine Puzzle:", false)
	minePuzzBomb := InputMenu.AddItem("mine_puzz_bomb", " Mark Bomb", false)
	minePuzzBombR := InputMenu.AddItem("mine_puzz_bomb_r", "", true)
	minePuzzSafe := InputMenu.AddItem("mine_puzz_safe", " Mark Safe", false)
	minePuzzSafeR := InputMenu.AddItem("mine_puzz_safe_r", "", true)
	back := InputMenu.AddItem("back", "Back", false)

	aimModeUpdate := func(in *input.Input) {
		if in.AimDedicated {
			aimModeR.SetText("Dedicated")
			if in.Mode == input.KeyboardMouse {
				aimMode.Hint = "Use the mouse to aim for digging, flagging, and attacking."
			} else {
				aimMode.Hint = "Use the right stick to aim for digging, flagging, and attacking."
			}
		} else {
			aimModeR.SetText("Movement")
			aimMode.Hint = "Use the movement keys to aim for digging, flagging, and attacking."
		}
	}
	aimModeUpdate(data.CurrInput)
	deviceUpdate := func(in *input.Input) {
		km := in.Mode == input.KeyboardMouse
		if km {
			device.Hint = ""
			deviceR.SetText("KB&Mouse")
		} else {
			device.Hint = win.JoystickName(in.Joystick)
			deviceR.SetText(fmt.Sprintf("Gamepad %d", in.Joystick+1))
		}
		leftStickA.Ignore = km
		leftStickB.Ignore = km
		leftStickR.Ignore = km
		deadzone.Ignore = km
		deadzoneR.Ignore = km
	}
	deviceUpdate(data.CurrInput)
	digOnUpdate := func(in *input.Input) {
		if in.DigOnRelease {
			digOnR.SetText("On Release")
			digOn.Hint = "Digging, Flagging, and Attacking happen when you release the button."
		} else {
			digOnR.SetText("On Press")
			digOn.Hint = "Digging, Flagging, and Attacking happen when you press the button."
		}
	}
	digOnUpdate(data.CurrInput)

	inputTitle.NoHover = true
	deviceSwitch := func(in *input.Input, prev bool) {
		km := in.Mode == input.KeyboardMouse
		var js int
		if prev {
			if km {
				js = input.PrevGamepad(win, -1)
			} else {
				js = input.PrevGamepad(win, int(in.Joystick))
			}
		} else {
			if km {
				js = input.NextGamepad(win, -1)
			} else {
				js = input.NextGamepad(win, int(in.Joystick))
			}
		}
		if js != -1 {
			in.Joystick = pixelgl.Joystick(js)
			in.Mode = input.Gamepad
		} else {
			in.Joystick = pixelgl.JoystickLast
			in.Mode = input.KeyboardMouse
		}
		UpdateKeybindings(in)
		sfx.SoundPlayer.PlaySound("click", 2.0)
		aimModeUpdate(in)
		deviceUpdate(in)
	}
	rfn1 := func() {
		deviceSwitch(data.CurrInput, false)
	}
	lfn1 := func() {
		deviceSwitch(data.CurrInput, true)
	}
	device.SetClickFn(rfn1)
	device.SetRightFn(rfn1)
	device.SetLeftFn(lfn1)
	deviceR.NoHover = true
	rfn2 := func() {
		data.CurrInput.AimDedicated = !data.CurrInput.AimDedicated
		sfx.SoundPlayer.PlaySound("click", 2.0)
		aimModeUpdate(data.CurrInput)
	}
	aimMode.SetClickFn(rfn2)
	aimMode.SetRightFn(rfn2)
	aimMode.SetLeftFn(rfn2)
	aimModeR.NoHover = true
	fn3 := func() {
		data.CurrInput.DigOnRelease = !data.CurrInput.DigOnRelease
		sfx.SoundPlayer.PlaySound("click", 2.0)
		digOnUpdate(data.CurrInput)
	}
	digOn.SetClickFn(fn3)
	digOn.SetRightFn(fn3)
	digOn.SetLeftFn(fn3)
	digOnR.NoHover = true
	deadzone.SetRightFn(func() {
		n := data.CurrInput.Deadzone + 0.05
		if n > 0.5 {
			n = 0.5
		}
		data.CurrInput.Deadzone = n
		deadzoneR.SetText(fmt.Sprintf("%f", n))
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	deadzone.SetLeftFn(func() {
		n := data.CurrInput.Deadzone - 0.05
		if n < 0.05 {
			n = 0.05
		}
		data.CurrInput.Deadzone = n
		deadzoneR.SetText(fmt.Sprintf("%f", n))
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	deadzoneR.NoHover = true
	rfn3 := func() {
		data.CurrInput.StickD = !data.CurrInput.StickD
		if data.CurrInput.StickD {
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
	if data.CurrInput.StickD {
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
	puzzLeave.SetClickFn(keyFn(puzzLeave))
	puzzLeaveR.NoHover = true
	puzzHelp.SetClickFn(keyFn(puzzHelp))
	puzzHelpR.NoHover = true
	minePuzz.NoHover = true
	minePuzz.Disabled = true
	minePuzzBomb.SetClickFn(keyFn(minePuzzBomb))
	minePuzzBombR.NoHover = true
	minePuzzSafe.SetClickFn(keyFn(minePuzzSafe))
	minePuzzSafeR.NoHover = true

	setProfile := func(in *input.Input) {
		data.CurrInput = in
		p := "?"
		switch data.CurrInput.Key {
		case "p1":
			p = "1"
		case "p2":
			p = "2"
		case "p3":
			p = "3"
		case "p4":
			p = "4"
		}
		player.SetText(fmt.Sprintf("< Player %s >", p))
		UpdateKeybindings(data.CurrInput)
		aimModeUpdate(data.CurrInput)
		deviceUpdate(data.CurrInput)
		digOnUpdate(data.CurrInput)
		deadzoneR.SetText(fmt.Sprintf("%f", data.CurrInput.Deadzone))
		if data.CurrInput.StickD {
			leftStickR.SetText("Yes")
		} else {
			leftStickR.SetText("No")
		}
		RegisterPlayerSymbols(data.CurrInput.Key, data.CurrInput)
	}
	player.SetLeftFn(func() {
		switch data.CurrInput.Key {
		case "p1":
			setProfile(data.GameInputP4)
		case "p2":
			setProfile(data.GameInputP1)
		case "p3":
			setProfile(data.GameInputP2)
		case "p4":
			setProfile(data.GameInputP3)
		}
	})
	player.SetRightFn(func() {
		switch data.CurrInput.Key {
		case "p1":
			setProfile(data.GameInputP2)
		case "p2":
			setProfile(data.GameInputP3)
		case "p3":
			setProfile(data.GameInputP4)
		case "p4":
			setProfile(data.GameInputP1)
		}
	})

	InputMenu.SetOpenFn(func() {
		setProfile(data.GameInputP1)
	})
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		InputMenu.Close()
	})
}