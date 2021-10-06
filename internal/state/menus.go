package state

import (
	"dwarf-sweeper/internal/config"
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/dungeon"
	"dwarf-sweeper/internal/enchants"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/player"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/util"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"strconv"
	"strings"
)

var (
	MainMenu       *menus.DwarfMenu
	AudioMenu      *menus.DwarfMenu
	GraphicsMenu   *menus.DwarfMenu
	InputMenu      *menus.DwarfMenu
	KeybindingMenu *menus.DwarfMenu
	PauseMenu      *menus.DwarfMenu
	OptionsMenu    *menus.DwarfMenu
	EnchantMenu    *menus.DwarfMenu
	PostMenu       *menus.DwarfMenu
	KeyString      string
)

func InitializeMenus(win *pixelgl.Window) {
	InitMainMenu(win)
	InitOptionsMenu()
	// todo: accessibility
	InitAudioMenu()
	InitGraphicsMenu()
	InitInputMenu(win)
	InitKeybindingMenu()
	InitPauseMenu(win)
	InitEnchantMenu()
	InitPostGameMenu()
}

func UpdateMenus(win *pixelgl.Window) {
	if len(menuStack) > 0 {
		currMenu := menuStack[len(menuStack)-1]
		currMenu.Update(menuInput)
		if currMenu.Closed {
			if len(menuStack) > 1 {
				menuStack = menuStack[:len(menuStack)-1]
			} else {
				menuStack = []*menus.DwarfMenu{}
			}
		} else if currMenu.Key == "keybinding" && currMenu.IsOpen() {
			if menuInput.Get("inputClear").JustPressed() {
				input.ClearInput(player.GameInput, KeyString)
				menuInput.Get("inputClear").Consume()
				currMenu.Close()
			} else {
				if input.CheckAssign(win, player.GameInput, KeyString) {
					player.GameInput.Buttons[KeyString].Button.Consume()
					currMenu.Close()
				}
			}
		}
	}
}

func MenuClosed() bool {
	return len(menuStack) < 1
}

func OpenMenu(menu *menus.DwarfMenu) {
	menu.Open()
	menuStack = append(menuStack, menu)
}

func InitMainMenu(win *pixelgl.Window) {
	MainMenu = menus.New("main", camera.Cam)
	MainMenu.Title = true
	start := MainMenu.AddItem("start", "Start Game")
	options := MainMenu.AddItem("options", "Options")
	credit := MainMenu.AddItem("credits", "Credits")
	quit := MainMenu.AddItem("quit", "Quit")

	start.SetClickFn(func() {
		MainMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		newState = 4
	})
	start.Hint = "Start a new run!"
	options.SetClickFn(func() {
		OpenMenu(OptionsMenu)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	credit.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		newState = 3
	})
	quit.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		win.SetClosed(true)
	})
	quit.Hint = "You're going to leave?"
}

func InitOptionsMenu() {
	OptionsMenu = menus.New("options", camera.Cam)
	OptionsMenu.Title = true
	optionsTitle := OptionsMenu.AddItem("title", "Options")
	audioOptions := OptionsMenu.AddItem("audio", "Audio")
	graphicsOptions := OptionsMenu.AddItem("graphics", "Graphics")
	inputOptions := OptionsMenu.AddItem("input", "Input")
	back := OptionsMenu.AddItem("back", "Back")

	optionsTitle.NoHover = true
	audioOptions.SetClickFn(func() {
		OpenMenu(AudioMenu)
	})
	graphicsOptions.SetClickFn(func() {
		OpenMenu(GraphicsMenu)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	inputOptions.SetClickFn(func() {
		OpenMenu(InputMenu)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		OptionsMenu.Close()
	})
}

func InitAudioMenu() {
	AudioMenu = menus.New("audio", camera.Cam)
	AudioMenu.Title = true
	AudioMenu.SetCloseFn(config.SaveAsync)
	audioTitle := AudioMenu.AddItem("title", "Audio Options")
	soundVolume := AudioMenu.AddItem("s_volume", "Sound Volume")
	soundVolumeR := AudioMenu.AddItem("s_volume_r", strconv.Itoa(sfx.GetSoundVolume()))
	musicVolume := AudioMenu.AddItem("m_volume", "Music Volume")
	musicVolumeR := AudioMenu.AddItem("m_volume_r", strconv.Itoa(sfx.GetMusicVolume()))
	back := AudioMenu.AddItem("back", "Back")

	audioTitle.NoHover = true
	soundVolume.SetRightFn(func() {
		n := sfx.GetSoundVolume() + 5
		if n > 100 {
			n = 100
		}
		sfx.SetSoundVolume(n)
		soundVolumeR.Raw = strconv.Itoa(n)
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
	})
	soundVolume.SetLeftFn(func() {
		n := sfx.GetSoundVolume() - 5
		if n < 0 {
			n = 0
		}
		sfx.SetSoundVolume(n)
		soundVolumeR.Raw = strconv.Itoa(n)
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
	})
	soundVolumeR.NoHover = true
	soundVolumeR.Right = true
	musicVolume.SetRightFn(func() {
		n := sfx.GetMusicVolume() + 5
		if n > 100 {
			n = 100
		}
		sfx.SetMusicVolume(n)
		musicVolumeR.Raw = strconv.Itoa(n)
	})
	musicVolume.SetLeftFn(func() {
		n := sfx.GetMusicVolume() - 5
		if n < 0 {
			n = 0
		}
		sfx.SetMusicVolume(n)
		musicVolumeR.Raw = strconv.Itoa(n)
	})
	musicVolumeR.NoHover = true
	musicVolumeR.Right = true
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		AudioMenu.Close()
	})
}

func InitGraphicsMenu() {
	GraphicsMenu = menus.New("graphics", camera.Cam)
	GraphicsMenu.Title = true
	GraphicsMenu.SetCloseFn(config.SaveAsync)
	graphicsTitle := GraphicsMenu.AddItem("title", "Graphics Options")
	vsync := GraphicsMenu.AddItem("vsync", "VSync")
	vsyncR := GraphicsMenu.AddItem("vsync_r", "Off")
	fullscreen := GraphicsMenu.AddItem("fullscreen", "Fullscreen")
	fullscreenR := GraphicsMenu.AddItem("fullscreen_r", "Off")
	resolution := GraphicsMenu.AddItem("resolution", "Resolution")
	resolutionR := GraphicsMenu.AddItem("resolution_r", constants.ResStrings[constants.ResIndex])
	back := GraphicsMenu.AddItem("back", "Back")

	graphicsTitle.NoHover = true
	vsync.SetClickFn(func() {
		constants.VSync = !constants.VSync
		if constants.VSync {
			vsyncR.Raw = "On"
		} else {
			vsyncR.Raw = "Off"
		}
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	if constants.VSync {
		vsyncR.Raw = "On"
	}
	if constants.FullScreen {
		fullscreenR.Raw = "On"
	}
	vsyncR.NoHover = true
	vsyncR.Right = true
	fullscreen.SetClickFn(func() {
		constants.FullScreen = !constants.FullScreen
		if constants.FullScreen {
			fullscreenR.Raw = "On"
		} else {
			fullscreenR.Raw = "Off"
		}
		constants.ChangeScreenSize = true
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	fullscreenR.NoHover = true
	fullscreenR.Right = true
	fn := func() {
		constants.ResIndex += 1
		constants.ResIndex %= len(constants.Resolutions)
		constants.ChangeScreenSize = true
		resolutionR.Raw = constants.ResStrings[constants.ResIndex]
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	resolution.SetClickFn(fn)
	resolution.SetRightFn(fn)
	resolution.SetLeftFn(func() {
		constants.ResIndex += len(constants.Resolutions) - 1
		constants.ResIndex %= len(constants.Resolutions)
		constants.ChangeScreenSize = true
		resolutionR.Raw = constants.ResStrings[constants.ResIndex]
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	resolutionR.NoHover = true
	resolutionR.Right = true
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		GraphicsMenu.Close()
	})
}

func InitInputMenu(win *pixelgl.Window) {
	InputMenu = menus.New("input", camera.Cam)
	InputMenu.Title = true
	InputMenu.SetOpenFn(UpdateKeybindings)
	InputMenu.SetCloseFn(config.SaveAsync)
	inputTitle := InputMenu.AddItem("title", "Input Options")
	device := InputMenu.AddItem("device", "Device")
	deviceR := InputMenu.AddItem("device_r", "")
	digMode := InputMenu.AddItem("dig_mode", "Dig Mode")
	digModeR := InputMenu.AddItem("dig_mode_r", data.DigMode(constants.DigMode).String())
	deadzone := InputMenu.AddItem("deadzone", "Deadzone")
	deadzoneR := InputMenu.AddItem("deadzone_r", fmt.Sprintf("%f", input.Deadzone))
	leftStickA := InputMenu.AddItem("left_stick_a", "Move with")
	leftStickR := InputMenu.AddItem("left_stick_r", "Yes")
	leftStickB := InputMenu.AddItem("left_stick_b", " Left Stick")
	left := InputMenu.AddItem("left", "Move Left")
	leftR := InputMenu.AddItem("left_r", "A")
	right := InputMenu.AddItem("right", "Move Right")
	rightR := InputMenu.AddItem("right_r", "D")
	up := InputMenu.AddItem("up", "Climb Up")
	upR := InputMenu.AddItem("up_r", "W")
	down := InputMenu.AddItem("down", "Climb Down")
	downR := InputMenu.AddItem("down_r", "S")
	jump := InputMenu.AddItem("jump", "Jump")
	jumpR := InputMenu.AddItem("jump_r", "Space")
	dig := InputMenu.AddItem("dig", "Dig")
	digR := InputMenu.AddItem("dig_r", "LShift,LMouse")
	mark := InputMenu.AddItem("mark", "Mark")
	markR := InputMenu.AddItem("mark_r", "LCtrl,RMouse")
	use := InputMenu.AddItem("use", "Use Item")
	useR := InputMenu.AddItem("use_r", "F")
	prev := InputMenu.AddItem("prev", "Prev Item")
	prevR := InputMenu.AddItem("prev_r", "Q")
	next := InputMenu.AddItem("next", "Next Item")
	nextR := InputMenu.AddItem("next_r", "E")
	back := InputMenu.AddItem("back", "Back")

	digModeHint := func() {
		switch constants.DigMode {
		case data.Dedicated:
			if player.GameInput.Mode == input.KeyboardMouse {
				digMode.Hint = "Use the mouse to aim for digging and marking."
			} else {
				digMode.Hint = "Use the right stick to aim for digging and marking."
			}
		case data.Either:
			digMode.Hint = ""
		case data.Movement:
			digMode.Hint = "Use the movement keys to aim for digging and marking."
		}
	}
	digModeHint()
	deviceUpdate := func() {
		km := player.GameInput.Mode == input.KeyboardMouse
		if km {
			device.Hint = ""
			deviceR.Raw = "KB&Mouse"
		} else {
			device.Hint = win.JoystickName(player.GameInput.Joystick)
			deviceR.Raw = fmt.Sprintf("Gamepad %d", player.GameInput.Joystick+1)
		}
		leftStickA.NoShow = km
		leftStickB.NoShow = km
		leftStickR.NoShow = km
		deadzone.NoShow = km
		deadzoneR.NoShow = km
	}
	deviceUpdate()

	inputTitle.NoHover = true
	deviceSwitch := func(prev bool) {
		km := player.GameInput.Mode == input.KeyboardMouse
		var js int
		if prev {
			if km {
				js = input.PrevGamepad(win, -1)
			} else {
				js = input.PrevGamepad(win, int(player.GameInput.Joystick))
			}
		} else {
			if km {
				js = input.NextGamepad(win, -1)
			} else {
				js = input.NextGamepad(win, int(player.GameInput.Joystick))
			}
		}
		if js != -1 {
			player.GameInput.Joystick = pixelgl.Joystick(js)
			player.GameInput.Mode = input.Gamepad
		} else {
			player.GameInput.Joystick = pixelgl.JoystickLast
			player.GameInput.Mode = input.KeyboardMouse
		}
		UpdateKeybindings()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		digModeHint()
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
		var dm data.DigMode
		switch constants.DigMode {
		case data.Either:
			dm = data.Movement
		case data.Movement:
			dm = data.Dedicated
		case data.Dedicated:
			dm = data.Either
		}
		constants.DigMode = int(dm)
		digModeR.Raw = dm.String()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		digModeHint()
	}
	lfn2 := func() {
		var dm data.DigMode
		switch constants.DigMode {
		case data.Either:
			dm = data.Dedicated
		case data.Movement:
			dm = data.Either
		case data.Dedicated:
			dm = data.Movement
		}
		constants.DigMode = int(dm)
		digModeR.Raw = dm.String()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		digModeHint()
	}
	digMode.SetClickFn(rfn2)
	digMode.SetRightFn(rfn2)
	digMode.SetLeftFn(lfn2)
	digModeR.Right = true
	digModeR.NoHover = true
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
		player.GameInput.StickD = !player.GameInput.StickD
		if player.GameInput.StickD {
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
	if player.GameInput.StickD {
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
	use.SetClickFn(keyFn(use))
	useR.Right = true
	useR.NoHover = true
	next.SetClickFn(keyFn(next))
	nextR.Right = true
	nextR.NoHover = true
	prev.SetClickFn(keyFn(prev))
	prevR.Right = true
	prevR.NoHover = true

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
	UpdateKeybinding("use")
	UpdateKeybinding("prev")
	UpdateKeybinding("next")
}

func UpdateKeybinding(key string) {
	r := InputMenu.ItemMap[fmt.Sprintf("%s_r", key)]
	in := player.GameInput.Buttons[key]
	builder := strings.Builder{}
	var symKeys []string
	first := true
	if player.GameInput.Mode != input.Gamepad {
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
	if player.GameInput.Mode != input.KeyboardMouse {
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

func InitPauseMenu(win *pixelgl.Window) {
	PauseMenu = menus.New("pause", camera.Cam)
	PauseMenu.Title = true
	pauseTitle := PauseMenu.AddItem("title", "Paused")
	resume := PauseMenu.AddItem("resume", "Resume")
	options := PauseMenu.AddItem("options", "Options")
	mainMenu := PauseMenu.AddItem("main_menu", "Abandon Run")
	quit := PauseMenu.AddItem("quit", "Quit Game")

	pauseTitle.NoHover = true
	resume.SetClickFn(func() {
		PauseMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	options.SetClickFn(func() {
		OpenMenu(OptionsMenu)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	mainMenu.SetClickFn(func() {
		PauseMenu.CloseInstant()
		newState = 1
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	quit.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		win.SetClosed(true)
	})
}

func InitEnchantMenu() {
	EnchantMenu = menus.New("enchant", camera.Cam)
	EnchantMenu.Title = true
	chooseTitle := EnchantMenu.AddItem("title", "Enchant!")
	skip := EnchantMenu.AddItem("skip", "Skip")

	chooseTitle.NoHover = true
	skip.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		EnchantMenu.CloseInstant()
		ClearEnchantMenu()
		newState = 0
	})
}

func ClearEnchantMenu() {
	EnchantMenu.RemoveItem("option1")
	EnchantMenu.RemoveItem("option2")
	EnchantMenu.RemoveItem("option3")
}

func FillEnchantMenu() bool {
	ClearEnchantMenu()
	pe := dungeon.Dungeon.GetPlayer().Enchants
	list := enchants.Enchantments
	for _, i := range pe {
		if len(list) > 1 {
			list = append(list[:i], list[i+1:]...)
		} else {
			return false
		}
	}
	choices := util.RandomSampleRange(util.Min(len(list), 3), 0, len(list), random.CaveGen)
	var e1, e2, e3 *data.Enchantment
	for i, c := range choices {
		if i == 0 {
			e1 = list[c]
		} else if i == 1 {
			e2 = list[c]
		} else if i == 2 {
			e3 = list[c]
		}
	}
	if e1 == nil {
		return false
	}
	option1 := EnchantMenu.InsertItem("option1", e1.Title, 1)
	option1.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		enchants.AddEnchantment(e1)
		EnchantMenu.CloseInstant()
		ClearEnchantMenu()
		newState = 0
	})
	option1.Hint = e1.Desc
	if e2 != nil {
		option2 := EnchantMenu.InsertItem("option2", e2.Title, 2)
		option2.SetClickFn(func() {
			sfx.SoundPlayer.PlaySound("click", 2.0)
			enchants.AddEnchantment(e2)
			EnchantMenu.CloseInstant()
			ClearEnchantMenu()
			newState = 0
		})
		option2.Hint = e2.Desc
	}
	if e3 != nil {
		option3 := EnchantMenu.InsertItem("option3", e3.Title, 3)
		option3.SetClickFn(func() {
			sfx.SoundPlayer.PlaySound("click", 2.0)
			enchants.AddEnchantment(e3)
			EnchantMenu.CloseInstant()
			ClearEnchantMenu()
			newState = 0
		})
		option3.Hint = e3.Desc
	}
	return true
}

func InitPostGameMenu() {
	PostMenu = menus.New("post", camera.Cam)
	PostMenu.Title = true
	PostMenu.SetBackFn(func() {})
	blocksDug := PostMenu.AddItem("blocks", "Blocks Dug")
	blocksDugS := PostMenu.AddItem("blocks_s", "")
	gems := PostMenu.AddItem("gem_count", "Gems Found")
	gemsS := PostMenu.AddItem("gem_count_s", "")
	bombs := PostMenu.AddItem("bombs_marked", "Bombs Marked")
	bombsS := PostMenu.AddItem("bombs_marked_s", "")
	wrongMarks := PostMenu.AddItem("wrong_marks", "Incorrect Marks")
	wrongMarksS := PostMenu.AddItem("wrong_marks_s", "")
	totalScore := PostMenu.AddItem("total_score", "Total Score")
	totalScoreS := PostMenu.AddItem("total_score_s", "")
	retry := PostMenu.AddItem("retry", "Retry")
	backToMenu := PostMenu.AddItem("menu", "Main Menu")

	blocksDug.NoHover = true
	blocksDug.NoShow = true
	blocksDugS.Right = true
	blocksDugS.NoHover = true
	blocksDugS.NoShow = true
	gems.NoHover = true
	gems.NoShow = true
	gemsS.Right = true
	gemsS.NoHover = true
	gemsS.NoShow = true
	bombs.NoHover = true
	bombs.NoShow = true
	bombsS.Right = true
	bombsS.NoHover = true
	bombsS.NoShow = true
	wrongMarks.NoHover = true
	wrongMarks.NoShow = true
	wrongMarksS.Right = true
	wrongMarksS.NoHover = true
	wrongMarksS.NoShow = true
	totalScore.NoHover = true
	totalScore.NoShow = true
	totalScoreS.Right = true
	totalScoreS.NoHover = true
	totalScoreS.NoShow = true
	retry.SetClickFn(func() {
		PostMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		newState = 4
	})
	backToMenu.Right = true
	backToMenu.SetClickFn(func() {
		PostMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		newState = 1
	})
}