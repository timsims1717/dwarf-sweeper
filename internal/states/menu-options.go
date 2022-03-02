package states

import (
	"dwarf-sweeper/internal/config"
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/hud"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/sfx"
	"fmt"
	"strconv"
)

func InitOptionsMenu() {
	OptionsMenu = menus.New("options", camera.Cam)
	OptionsMenu.Title = true
	optionsTitle := OptionsMenu.AddItem("title", "Options", false)
	audioOptions := OptionsMenu.AddItem("audio", "Audio", false)
	gameplayOptions := OptionsMenu.AddItem("gameplay", "Gameplay", false)
	graphicsOptions := OptionsMenu.AddItem("graphics", "Graphics", false)
	inputOptions := OptionsMenu.AddItem("input", "Input", false)
	back := OptionsMenu.AddItem("back", "Back", false)

	optionsTitle.NoHover = true
	audioOptions.SetClickFn(func() {
		OpenMenu(AudioMenu)
	})
	gameplayOptions.SetClickFn(func() {
		OpenMenu(GameplayMenu)
		sfx.SoundPlayer.PlaySound("click", 2.0)
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
	audioTitle := AudioMenu.AddItem("title", "Audio Options", false)
	soundVolume := AudioMenu.AddItem("s_volume", "Sound Volume", false)
	soundVolumeR := AudioMenu.AddItem("s_volume_r", strconv.Itoa(sfx.GetSoundVolume()), true)
	musicVolume := AudioMenu.AddItem("m_volume", "Music Volume", false)
	musicVolumeR := AudioMenu.AddItem("m_volume_r", strconv.Itoa(sfx.GetMusicVolume()), true)
	muteOnUnfocus := AudioMenu.AddItem("mute_focus", "Mute On Unfocus", false)
	muteOnUnfocusR := AudioMenu.AddItem("mute_focus_r", "No", true)
	back := AudioMenu.AddItem("back", "Back", false)

	audioTitle.NoHover = true
	soundVolume.SetRightFn(func() {
		n := sfx.GetSoundVolume() + 5
		if n > 100 {
			n = 100
		}
		sfx.SetSoundVolume(n)
		soundVolumeR.SetText(strconv.Itoa(n))
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
	})
	soundVolume.SetLeftFn(func() {
		n := sfx.GetSoundVolume() - 5
		if n < 0 {
			n = 0
		}
		sfx.SetSoundVolume(n)
		soundVolumeR.SetText(strconv.Itoa(n))
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
	})
	soundVolumeR.NoHover = true
	musicVolume.SetRightFn(func() {
		n := sfx.GetMusicVolume() + 5
		if n > 100 {
			n = 100
		}
		sfx.SetMusicVolume(n)
		musicVolumeR.SetText(strconv.Itoa(n))
	})
	musicVolume.SetLeftFn(func() {
		n := sfx.GetMusicVolume() - 5
		if n < 0 {
			n = 0
		}
		sfx.SetMusicVolume(n)
		musicVolumeR.SetText(strconv.Itoa(n))
	})
	musicVolumeR.NoHover = true
	muteFn := func() {
		constants.MuteOnUnfocused = !constants.MuteOnUnfocused
		if constants.MuteOnUnfocused {
			muteOnUnfocusR.SetText("Yes")
		} else {
			muteOnUnfocusR.SetText("No")
		}
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	muteOnUnfocus.SetRightFn(muteFn)
	muteOnUnfocus.SetLeftFn(muteFn)
	muteOnUnfocus.SetClickFn(muteFn)
	if constants.MuteOnUnfocused {
		muteOnUnfocusR.SetText("Yes")
	} else {
		muteOnUnfocusR.SetText("No")
	}
	muteOnUnfocusR.NoHover = true
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		AudioMenu.Close()
	})
}

func InitGameplayMenu() {
	GameplayMenu = menus.New("gameplay", camera.Cam)
	GameplayMenu.Title = true
	GameplayMenu.SetCloseFn(config.SaveAsync)
	gameplayTitle := GameplayMenu.AddItem("title", "Gameplay Options", false)
	showDescentTimer := GameplayMenu.AddItem("showTimer", "Show Timer", false)
	showDescentTimerR := GameplayMenu.AddItem("showTimer_r", "No", true)
	screenShake := GameplayMenu.AddItem("screenShake", "Screen Shake", false)
	screenShakeR := GameplayMenu.AddItem("screenShake_r", "No", true)
	splitScreen := GameplayMenu.AddItem("splitScreen", "Split Screen", false)
	splitScreenR := GameplayMenu.AddItem("splitScreen_r", "Horiz.", true)

	gameplayTitle.NoHover = true
	showFn := func() {
		hud.ShowTimer = !hud.ShowTimer
		if hud.ShowTimer {
			showDescentTimerR.SetText("Yes")
		} else {
			showDescentTimerR.SetText("No")
		}
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	showDescentTimer.SetClickFn(showFn)
	showDescentTimer.SetRightFn(showFn)
	showDescentTimer.SetLeftFn(showFn)
	if hud.ShowTimer {
		showDescentTimerR.SetText("Yes")
	}
	showDescentTimerR.NoHover = true

	shakeFn := func() {
		constants.ScreenShake = !constants.ScreenShake
		if constants.ScreenShake {
			screenShakeR.SetText("Yes")
		} else {
			screenShakeR.SetText("No")
		}
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	screenShake.SetClickFn(shakeFn)
	screenShake.SetRightFn(shakeFn)
	screenShake.SetLeftFn(shakeFn)
	if constants.ScreenShake {
		screenShakeR.SetText("Yes")
	}
	screenShakeR.NoHover = true

	splitFn := func() {
		constants.SplitScreenV = !constants.SplitScreenV
		if constants.SplitScreenV {
			splitScreenR.SetText("Vert.")
		} else {
			splitScreenR.SetText("Horiz.")
		}
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	splitScreen.SetClickFn(splitFn)
	splitScreen.SetRightFn(splitFn)
	splitScreen.SetLeftFn(splitFn)
	if constants.SplitScreenV {
		splitScreenR.SetText("Vert.")
	}
	splitScreenR.NoHover = true
}

func InitGraphicsMenu() {
	GraphicsMenu = menus.New("graphics", camera.Cam)
	GraphicsMenu.Title = true
	GraphicsMenu.SetCloseFn(config.SaveAsync)
	graphicsTitle := GraphicsMenu.AddItem("title", "Graphics Options", false)
	vsync := GraphicsMenu.AddItem("vsync", "VSync", false)
	vsyncR := GraphicsMenu.AddItem("vsync_r", "Off", true)
	fullscreen := GraphicsMenu.AddItem("fullscreen", "Fullscreen", false)
	fullscreenR := GraphicsMenu.AddItem("fullscreen_r", "Off", true)
	resolution := GraphicsMenu.AddItem("resolution", "Resolution", false)
	resolutionR := GraphicsMenu.AddItem("resolution_r", constants.ResStrings[constants.ResIndex], true)
	back := GraphicsMenu.AddItem("back", "Back", false)

	graphicsTitle.NoHover = true
	vFn := func() {
		constants.VSync = !constants.VSync
		if constants.VSync {
			vsyncR.SetText("On")
		} else {
			vsyncR.SetText("Off")
		}
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	vsync.SetClickFn(vFn)
	vsync.SetRightFn(vFn)
	vsync.SetLeftFn(vFn)
	if constants.VSync {
		vsyncR.SetText("On")
	}
	if constants.FullScreen {
		fullscreenR.SetText("On")
	}
	vsyncR.NoHover = true
	fsFn := func() {
		constants.FullScreen = !constants.FullScreen
		if constants.FullScreen {
			fullscreenR.SetText("On")
		} else {
			fullscreenR.SetText("Off")
		}
		constants.ChangeScreen = true
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	fullscreen.SetClickFn(fsFn)
	fullscreen.SetRightFn(fsFn)
	fullscreen.SetLeftFn(fsFn)
	fullscreenR.NoHover = true
	fn := func() {
		constants.ResIndex += 1
		constants.ResIndex %= len(constants.Resolutions)
		constants.ChangeScreen = true
		resolutionR.SetText(constants.ResStrings[constants.ResIndex])
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	resolution.SetClickFn(fn)
	resolution.SetRightFn(fn)
	resolution.SetLeftFn(func() {
		constants.ResIndex += len(constants.Resolutions) - 1
		constants.ResIndex %= len(constants.Resolutions)
		constants.ChangeScreen = true
		resolutionR.SetText(constants.ResStrings[constants.ResIndex])
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	resolutionR.NoHover = true
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		GraphicsMenu.Close()
	})
}
