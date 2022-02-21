package states

import (
	"dwarf-sweeper/internal/config"
	"dwarf-sweeper/internal/constants"
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
	graphicsOptions := OptionsMenu.AddItem("graphics", "Graphics", false)
	inputOptions := OptionsMenu.AddItem("input", "Input", false)
	back := OptionsMenu.AddItem("back", "Back", false)

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
	audioTitle := AudioMenu.AddItem("title", "Audio Options", false)
	soundVolume := AudioMenu.AddItem("s_volume", "Sound Volume", false)
	soundVolumeR := AudioMenu.AddItem("s_volume_r", strconv.Itoa(sfx.GetSoundVolume()), true)
	musicVolume := AudioMenu.AddItem("m_volume", "Music Volume", false)
	musicVolumeR := AudioMenu.AddItem("m_volume_r", strconv.Itoa(sfx.GetMusicVolume()), true)
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
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		AudioMenu.Close()
	})
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
	vsync.SetClickFn(func() {
		constants.VSync = !constants.VSync
		if constants.VSync {
			vsyncR.SetText("On")
		} else {
			vsyncR.SetText("Off")
		}
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	if constants.VSync {
		vsyncR.SetText("On")
	}
	if constants.FullScreen {
		fullscreenR.SetText("On")
	}
	vsyncR.NoHover = true
	fullscreen.SetClickFn(func() {
		constants.FullScreen = !constants.FullScreen
		if constants.FullScreen {
			fullscreenR.SetText("On")
		} else {
			fullscreenR.SetText("Off")
		}
		constants.ChangeScreenSize = true
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	fullscreenR.NoHover = true
	fn := func() {
		constants.ResIndex += 1
		constants.ResIndex %= len(constants.Resolutions)
		constants.ChangeScreenSize = true
		resolutionR.SetText(constants.ResStrings[constants.ResIndex])
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	resolution.SetClickFn(fn)
	resolution.SetRightFn(fn)
	resolution.SetLeftFn(func() {
		constants.ResIndex += len(constants.Resolutions) - 1
		constants.ResIndex %= len(constants.Resolutions)
		constants.ChangeScreenSize = true
		resolutionR.SetText(constants.ResStrings[constants.ResIndex])
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	resolutionR.NoHover = true
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		GraphicsMenu.Close()
	})
}