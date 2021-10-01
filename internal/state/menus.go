package state

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/dungeon"
	"dwarf-sweeper/internal/enchants"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/util"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"strconv"
)

var (
	MainMenu       *menus.DwarfMenu
	AudioMenu      *menus.DwarfMenu
	GraphicsMenu   *menus.DwarfMenu
	InputMenu      *menus.DwarfMenu
	//KeybindingMenu *menus.DwarfMenu
	PauseMenu      *menus.DwarfMenu
	OptionsMenu    *menus.DwarfMenu
	EnchantMenu    *menus.DwarfMenu
	PostMenu       *menus.DwarfMenu
)

func InitializeMenus(win *pixelgl.Window) {
	InitMainMenu(win)
	InitOptionsMenu()
	// todo: accessibility
	InitAudioMenu()
	InitGraphicsMenu()
	InitInputMenu(win)
	//InitKeybindingMenu()
	InitPauseMenu(win)
	InitEnchantMenu()
	InitPostGameMenu()
}

func InitMainMenu(win *pixelgl.Window) {
	MainMenu = menus.New("main", camera.Cam)
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
		OptionsMenu.Open()
		menuStack = append(menuStack, OptionsMenu)
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
	optionsTitle := OptionsMenu.AddItem("title", "Options")
	audioOptions := OptionsMenu.AddItem("audio", "Audio")
	graphicsOptions := OptionsMenu.AddItem("graphics", "Graphics")
	inputOptions := OptionsMenu.AddItem("input", "Input")
	back := OptionsMenu.AddItem("back", "Back")

	optionsTitle.NoHover = true
	audioOptions.SetClickFn(func() {
		AudioMenu.Open()
		menuStack = append(menuStack, AudioMenu)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	graphicsOptions.SetClickFn(func() {
		GraphicsMenu.Open()
		menuStack = append(menuStack, GraphicsMenu)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	inputOptions.SetClickFn(func() {
		InputMenu.Open()
		menuStack = append(menuStack, InputMenu)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		OptionsMenu.Close()
	})
}

func InitAudioMenu() {
	AudioMenu = menus.New("audio", camera.Cam)
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
	graphicsTitle := GraphicsMenu.AddItem("title", "Graphics Options")
	vsync := GraphicsMenu.AddItem("vsync", "VSync")
	vsyncR := GraphicsMenu.AddItem("vsync_r", "On")
	fullscreen := GraphicsMenu.AddItem("fullscreen", "Fullscreen")
	fullscreenR := GraphicsMenu.AddItem("fullscreen_r", "Off")
	resolution := GraphicsMenu.AddItem("resolution", "Resolution")
	resolutionR := GraphicsMenu.AddItem("resolution_r", cfg.ResStrings[cfg.ResIndex])
	back := GraphicsMenu.AddItem("back", "Back")

	graphicsTitle.NoHover = true
	vsync.SetClickFn(func() {
		s := "On"
		if cfg.VSync {
			s = "Off"
		}
		cfg.VSync = !cfg.VSync
		vsyncR.Raw = s
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	vsyncR.NoHover = true
	vsyncR.Right = true
	fullscreen.SetClickFn(func() {
		s := "On"
		if cfg.FullScreen {
			s = "Off"
		}
		cfg.FullScreen = !cfg.FullScreen
		cfg.ChangeScreenSize = true
		fullscreenR.Raw = s
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	fullscreenR.NoHover = true
	fullscreenR.Right = true
	fn := func() {
		cfg.ResIndex += 1
		cfg.ResIndex %= len(cfg.Resolutions)
		cfg.ChangeScreenSize = true
		resolutionR.Raw = cfg.ResStrings[cfg.ResIndex]
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	resolution.SetClickFn(fn)
	resolution.SetRightFn(fn)
	resolution.SetLeftFn(func() {
		cfg.ResIndex += len(cfg.Resolutions) - 1
		cfg.ResIndex %= len(cfg.Resolutions)
		cfg.ChangeScreenSize = true
		resolutionR.Raw = cfg.ResStrings[cfg.ResIndex]
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
	inputTitle := InputMenu.AddItem("title", "Input Options")
	device := InputMenu.AddItem("device", "Device")
	deviceR := InputMenu.AddItem("device_r", "KB&Mouse")
	digMode := InputMenu.AddItem("dig_mode", "Dig Mode")
	digModeR := InputMenu.AddItem("dig_mode_r", "Either")
	//keybindings := InputMenu.AddItem("keybindings", "Keybindings")
	back := InputMenu.AddItem("back", "Back")

	inputTitle.NoHover = true
	rfn1 := func() {
		km := gameInput.Mode == input.KeyboardMouse
		var js int
		if km {
			js = input.NextGamepad(win, -1)
		} else {
			js = input.NextGamepad(win, int(gameInput.Joystick))
		}
		if js != -1 {
			deviceR.Raw = fmt.Sprintf("Gamepad %d", js)
			gameInput.Joystick = pixelgl.Joystick(js)
			gameInput.Mode = input.Gamepad
			device.Hint = win.JoystickName(pixelgl.Joystick(js))
			if cfg.DigMode == data.Dedicated {
				digMode.Hint = "Use the right stick to aim for digging and marking."
			}
		} else {
			deviceR.Raw = "KB&Mouse"
			gameInput.Joystick = pixelgl.JoystickLast
			gameInput.Mode = input.KeyboardMouse
			device.Hint = ""
			if cfg.DigMode == data.Dedicated {
				digMode.Hint = "Use the mouse to aim for digging and marking."
			}
		}
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	lfn1 := func() {
		km := gameInput.Mode == input.KeyboardMouse
		var js int
		if km {
			js = input.PrevGamepad(win, -1)
		} else {
			js = input.PrevGamepad(win, int(gameInput.Joystick))
		}
		if js != -1 {
			deviceR.Raw = fmt.Sprintf("Gamepad %d (%s)", js, win.JoystickName(pixelgl.Joystick(js)))
			gameInput.Joystick = pixelgl.Joystick(js)
			gameInput.Mode = input.Gamepad
			if cfg.DigMode == data.Dedicated {
				digMode.Hint = "Use the right stick to aim for digging and marking."
			}
		} else {
			deviceR.Raw = "KB&Mouse"
			gameInput.Joystick = pixelgl.JoystickLast
			gameInput.Mode = input.KeyboardMouse
			if cfg.DigMode == data.Dedicated {
				digMode.Hint = "Use the mouse to aim for digging and marking."
			}
		}
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	device.SetClickFn(rfn1)
	device.SetRightFn(rfn1)
	device.SetLeftFn(lfn1)
	deviceR.Right = true
	deviceR.NoHover = true
	rfn2 := func() {
		var dm data.DigMode
		switch cfg.DigMode {
		case data.Either:
			dm = data.Movement
			digMode.Hint = "Use the movement keys to aim for digging and marking."
		case data.Movement:
			dm = data.Dedicated
			if gameInput.Mode == input.KeyboardMouse {
				digMode.Hint = "Use the mouse to aim for digging and marking."
			} else {
				digMode.Hint = "Use the right stick to aim for digging and marking."
			}
		case data.Dedicated:
			dm = data.Either
			digMode.Hint = ""
		}
		cfg.DigMode = int(dm)
		digModeR.Raw = dm.String()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	lfn2 := func() {
		var dm data.DigMode
		switch cfg.DigMode {
		case data.Either:
			dm = data.Dedicated
			if gameInput.Mode == input.KeyboardMouse {
				digMode.Hint = "Use the mouse to aim for digging and marking."
			} else {
				digMode.Hint = "Use the right stick to aim for digging and marking."
			}
		case data.Movement:
			dm = data.Either
			digMode.Hint = ""
		case data.Dedicated:
			dm = data.Movement
			digMode.Hint = "Use the movement keys to aim for digging and marking."
		}
		cfg.DigMode = int(dm)
		digModeR.Raw = dm.String()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	digMode.SetClickFn(rfn2)
	digMode.SetRightFn(rfn2)
	digMode.SetLeftFn(lfn2)
	digModeR.Right = true
	digModeR.NoHover = true
	//keybindings.SetClickFn(func() {
	//	KeybindingMenu.Open()
	//	menuStack = append(menuStack, KeybindingMenu)
	//	sfx.SoundPlayer.PlaySound("click", 2.0)
	//})
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		InputMenu.Close()
	})
}

//func InitKeybindingMenu() {
//	KeybindingMenu = menus.New("keybindings", pixel.R(0., 0., 64., 64.), camera.Cam)
//	keybindingTitle := KeybindingMenu.AddItem("title", "Keybindings")
//	left := KeybindingMenu.AddItem("left", "Move Left")
//	leftR := KeybindingMenu.AddItem("left_r", "A,DPad-L")
//	right := KeybindingMenu.AddItem("right", "Move Right")
//	up := KeybindingMenu.AddItem("up", "Climb Up")
//	down := KeybindingMenu.AddItem("down", "Climb Down")
//	jump := KeybindingMenu.AddItem("jump", "Jump")
//	dig := KeybindingMenu.AddItem("dig", "Dig")
//	mark := KeybindingMenu.AddItem("mark", "Mark")
//	use := KeybindingMenu.AddItem("use", "Use Item")
//	prev := KeybindingMenu.AddItem("prev", "Prev Item")
//	next := KeybindingMenu.AddItem("next", "Next Item")
//	back := KeybindingMenu.AddItem("back", "Back")
//
//	keybindingTitle.NoHover = true
//	back.SetClickFn(func() {
//		KeybindingMenu.Close()
//	})
//}

func InitPauseMenu(win *pixelgl.Window) {
	PauseMenu = menus.New("pause", camera.Cam)
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
		OptionsMenu.Open()
		menuStack = append(menuStack, OptionsMenu)
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