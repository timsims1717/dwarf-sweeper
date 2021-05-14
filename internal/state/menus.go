package state

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/menu"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/transform"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"strconv"
)

const (
	Size1 = 2.
	Size2 = 2.
	Spacing = 28.
)

var (
	MainMenu    *menu.Menu
	Options     *menu.Menu
	PostGame    *menu.Menu
	Current     int
)

func InitializeMenus() {
	MainMenu = menu.NewMenu(pixel.R(0,0, cfg.BaseW, cfg.BaseH), camera.Cam)
	Options  = menu.NewMenu(pixel.R(0,0, cfg.BaseW, cfg.BaseH), camera.Cam)
	PostGame = menu.NewMenu(pixel.R(0,0, cfg.BaseW, cfg.BaseH), camera.Cam)
	Current = 0
}

func InitializeMainMenu() {
	startS := "Start Game"
	startText := menu.NewItemText(startS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
	startText.HoverColor = colornames.Darkblue
	startText.HoverSize = pixel.V(Size2, Size2)
	startR := pixel.R(0., 0., startText.Text.BoundsOf(startS).W()*Size2, startText.Text.BoundsOf(startS).H()*Size2)
	startGame := menu.NewItem(startText, startR, MainMenu.Canvas.Bounds())
	startGame.Transform.Pos = pixel.V(0., 0.)
	startGame.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	startGame.SetOnHoverFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	startGame.SetClickFn(func() {
		//camera.Cam.Effect = transform.FadeTo(camera.Cam, colornames.Black, 1.0)
		//sfx.MusicPlayer.FadeOut(1.0)
		newState = 0
	})
	MainMenu.Items["start"] = startGame

	optionsS := "Options"
	optionsText := menu.NewItemText(optionsS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
	optionsText.HoverColor = colornames.Darkblue
	optionsText.HoverSize = pixel.V(Size2, Size2)
	optionsR := pixel.R(0., 0., optionsText.Text.BoundsOf(optionsS).W()*Size2, optionsText.Text.BoundsOf(optionsS).H()*Size2)
	optionsItem := menu.NewItem(optionsText, optionsR, MainMenu.Canvas.Bounds())
	optionsItem.Transform.Pos = pixel.V(0., Spacing * -1.)
	optionsItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	optionsItem.SetOnHoverFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	optionsItem.SetClickFn(func() {
		SwitchToOptions()
	})
	MainMenu.Items["options"] = optionsItem

	creditsS := "Credits"
	creditsText := menu.NewItemText(creditsS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
	creditsText.HoverColor = colornames.Darkblue
	creditsText.HoverSize = pixel.V(Size2, Size2)
	creditsR := pixel.R(0., 0., creditsText.Text.BoundsOf(creditsS).W()*Size2, creditsText.Text.BoundsOf(creditsS).H()*Size2)
	creditsItem := menu.NewItem(creditsText, creditsR, MainMenu.Canvas.Bounds())
	creditsItem.Transform.Pos = pixel.V(0., Spacing * -2.)
	creditsItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	creditsItem.SetOnHoverFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	creditsItem.SetClickFn(func() {
		newState = 3
	})
	MainMenu.Items["credits"] = creditsItem

	exitS := "Exit"
	exitText := menu.NewItemText(exitS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
	exitText.HoverColor = colornames.Darkblue
	exitText.HoverSize = pixel.V(Size2, Size2)
	exitR := pixel.R(0., 0., exitText.Text.BoundsOf(exitS).W()*Size2, exitText.Text.BoundsOf(exitS).H()*Size2)
	exitItem := menu.NewItem(exitText, exitR, MainMenu.Canvas.Bounds())
	exitItem.Transform.Pos = pixel.V(0., Spacing * -3.)
	exitItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	exitItem.SetOnHoverFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	MainMenu.Items["exit"] = exitItem
}

func InitializeSoundOption() {
	soundVS := "Sound Volume:"
	soundVMinusS := " - "
	soundVolumeS := strconv.Itoa(sfx.GetSoundVolume())
	soundVPlusS := " + "

	soundVText := menu.NewItemText(soundVS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)

	soundVolumeText := menu.NewItemText(soundVolumeS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)

	soundVMinusText := menu.NewItemText(soundVMinusS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
	soundVMinusText.HoverColor = colornames.Darkblue
	soundVMinusText.HoverSize = pixel.V(Size2, Size2)

	soundVPlusText := menu.NewItemText(soundVPlusS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
	soundVPlusText.HoverColor = colornames.Darkblue
	soundVPlusText.HoverSize = pixel.V(Size2, Size2)

	soundVR := pixel.R(0., 0., soundVText.Text.BoundsOf(soundVS).W()*Size2, soundVText.Text.BoundsOf(soundVS).H()*Size2)
	soundVolumeR := pixel.R(0., 0., soundVolumeText.Text.BoundsOf("100").W()*Size2, soundVolumeText.Text.BoundsOf(soundVolumeS).H()*Size2)
	soundVMinusR := pixel.R(0., 0., soundVMinusText.Text.BoundsOf(soundVMinusS).W()*Size2, soundVMinusText.Text.BoundsOf(soundVMinusS).H()*Size2)
	soundVPlusR := pixel.R(0., 0., soundVPlusText.Text.BoundsOf(soundVPlusS).W()*Size2, soundVPlusText.Text.BoundsOf(soundVPlusS).H()*Size2)
	soundVText.Transform.SetParent(soundVR)
	soundVolumeText.Transform.SetParent(soundVolumeR)
	soundVMinusText.Transform.SetParent(soundVMinusR)
	soundVPlusText.Transform.SetParent(soundVPlusR)

	soundVItem := menu.NewItem(soundVText, soundVR, Options.Canvas.Bounds())
	soundVItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	soundVItem.Disabled = true
	Options.Items["sound_v"] = soundVItem

	soundVolumeItem := menu.NewItem(soundVolumeText, soundVolumeR, Options.Canvas.Bounds())
	soundVolumeItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	soundVolumeItem.Disabled = true
	Options.Items["sound_volume"] = soundVolumeItem

	soundVMinusItem := menu.NewItem(soundVMinusText, soundVMinusR, Options.Canvas.Bounds())
	soundVMinusItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	soundVMinusItem.SetOnHoverFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	soundVMinusItem.SetClickFn(func() {
		n := sfx.GetSoundVolume() - 5
		if n < 0 {
			n = 0
		}
		sfx.SetSoundVolume(n)
		soundVolumeItem.Text.SetText(strconv.Itoa(n))
		SetOptionSoundWidth()
	})
	Options.Items["sound_-"] = soundVMinusItem

	soundVPlusItem := menu.NewItem(soundVPlusText, soundVPlusR, Options.Canvas.Bounds())
	soundVPlusItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	soundVPlusItem.SetOnHoverFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	soundVPlusItem.SetClickFn(func() {
		n := sfx.GetSoundVolume() + 5
		if n > 100 {
			n = 100
		}
		sfx.SetSoundVolume(n)
		soundVolumeItem.Text.SetText(strconv.Itoa(n))
		SetOptionSoundWidth()
	})
	Options.Items["sound_+"] = soundVPlusItem
	SetOptionSoundWidth()
}

func InitializeFullscreenOption() {
	fullscreenS := "Fullscreen:"
	fullscreenOptionS := " Off"
	if cfg.FullScreen {
		fullscreenOptionS = " On"
	}

	fullscreenText := menu.NewItemText(fullscreenS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)

	fullscreenOptionText := menu.NewItemText(fullscreenOptionS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
	fullscreenOptionText.HoverColor = colornames.Darkblue
	fullscreenOptionText.HoverSize = pixel.V(Size2, Size2)

	fullscreenR := pixel.R(0., 0., fullscreenText.Text.BoundsOf(fullscreenS).W()*Size2, fullscreenText.Text.BoundsOf(fullscreenS).H()*Size2)
	fullscreenItem := menu.NewItem(fullscreenText, fullscreenR, Options.Canvas.Bounds())
	fullscreenItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	fullscreenItem.Disabled = true
	Options.Items["fullscreen"] = fullscreenItem

	fullscreenOptionR := pixel.R(0., 0., fullscreenOptionText.Text.BoundsOf(" Off").W()*Size2, fullscreenOptionText.Text.BoundsOf(" Off").H()*Size2)
	fullscreenOptionItem := menu.NewItem(fullscreenOptionText, fullscreenOptionR, Options.Canvas.Bounds())
	fullscreenOptionItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	fullscreenOptionItem.SetOnHoverFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	fullscreenOptionItem.SetClickFn(func() {
		s := " On"
		if cfg.FullScreen {
			s = " Off"
		}
		cfg.FullScreen = !cfg.FullScreen
		cfg.ChangeScreenSize = true
		fullscreenOptionItem.Text.SetText(s)
		SetOptionFullscreenWidth()
	})
	Options.Items["fullscreen_options"] = fullscreenOptionItem
	SetOptionFullscreenWidth()
}

func InitializeResolutionOption() {
	resolutionS := "Resolution:"
	resolutionOptS := cfg.ResStrings[cfg.ResIndex]

	resolutionText := menu.NewItemText(resolutionS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)

	resolutionOptText := menu.NewItemText(resolutionOptS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
	resolutionOptText.HoverColor = colornames.Darkblue
	resolutionOptText.HoverSize = pixel.V(Size2, Size2)

	resolutionR := pixel.R(0., 0., resolutionText.Text.BoundsOf(resolutionS).W()*Size2, resolutionText.Text.BoundsOf(resolutionS).H()*Size2)
	resolutionItem := menu.NewItem(resolutionText, resolutionR, Options.Canvas.Bounds())
	resolutionItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	resolutionItem.Disabled = true
	Options.Items["resolution"] = resolutionItem

	resolutionOptR := pixel.R(0., 0., resolutionOptText.Text.BoundsOf(cfg.ResStrings[len(cfg.ResStrings)-1]).W()*Size2, resolutionOptText.Text.BoundsOf(cfg.ResStrings[len(cfg.ResStrings)-1]).H()*Size2)
	resolutionOptItem := menu.NewItem(resolutionOptText, resolutionOptR, Options.Canvas.Bounds())
	resolutionOptItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	resolutionOptItem.SetOnHoverFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	resolutionOptItem.SetClickFn(func() {
		cfg.ResIndex += 1
		cfg.ResIndex %= len(cfg.Resolutions)
		cfg.ChangeScreenSize = true
		resolutionOptItem.Text.SetText(cfg.ResStrings[cfg.ResIndex])
		SetOptionResolutionWidth()
	})
	Options.Items["resolution_options"] = resolutionOptItem
	SetOptionResolutionWidth()
}

func InitializeOptionsMenu() {
	Options.Transform.Pos.X = cfg.BaseW * 1.5

	InitializeSoundOption()
	InitializeFullscreenOption()
	InitializeResolutionOption()

	backS := "Back"
	backText := menu.NewItemText(backS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
	backText.HoverColor = colornames.Darkblue
	backText.HoverSize = pixel.V(Size2, Size2)
	backR := pixel.R(0., 0., backText.Text.BoundsOf(backS).W()*Size2, backText.Text.BoundsOf(backS).H()*Size2)
	backItem := menu.NewItem(backText, backR, Options.Canvas.Bounds())
	backItem.Transform.Pos = pixel.V(0., Spacing * -3.)
	backItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	backItem.SetOnHoverFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	Options.Items["back"] = backItem
}

func SetOptionSoundWidth() {
	soundVItem := Options.Items["sound_v"]
	soundVTextW := soundVItem.Text.Text.BoundsOf(soundVItem.Text.Raw).W()
	soundVMinusItem := Options.Items["sound_-"]
	soundVMinusTextW := soundVMinusItem.Text.Text.BoundsOf(soundVMinusItem.Text.Raw).W()
	soundVolumeItem := Options.Items["sound_volume"]
	soundVVolumeTextW := soundVolumeItem.Text.Text.BoundsOf(soundVolumeItem.Text.Raw).W()
	soundVPlusItem := Options.Items["sound_+"]
	soundVPlusTextW := soundVPlusItem.Text.Text.BoundsOf(soundVPlusItem.Text.Raw).W()
	totalWidth := (soundVTextW + soundVMinusTextW + soundVVolumeTextW + soundVPlusTextW) * Size1
	soundVItem.Transform.Pos = pixel.V((totalWidth - soundVTextW * Size1) * -0.5, Spacing * 0.)
	soundVMinusItem.Transform.Pos = pixel.V(totalWidth * -0.5 + (soundVTextW + soundVMinusTextW * 0.5) * Size1, Spacing * 0.)
	soundVolumeItem.Transform.Pos = pixel.V(totalWidth * 0.5 - (soundVPlusTextW + soundVVolumeTextW * 0.5) * Size1, Spacing * 0.)
	soundVPlusItem.Transform.Pos = pixel.V((totalWidth - soundVPlusTextW * Size1) * 0.5, Spacing * 0.)
}

func SetOptionFullscreenWidth() {
	fullscreenItem := Options.Items["fullscreen"]
	fullscreenTextW := fullscreenItem.Text.Text.BoundsOf(fullscreenItem.Text.Raw).W()
	fullscreenOptionItem := Options.Items["fullscreen_options"]
	fullscreenOptionTextW := fullscreenOptionItem.Text.Text.BoundsOf(fullscreenOptionItem.Text.Raw).W()
	totalWidth := (fullscreenTextW + fullscreenOptionTextW) * Size1
	fullscreenItem.Transform.Pos = pixel.V((totalWidth - fullscreenTextW * Size1) * -0.5, Spacing * -1.)
	fullscreenOptionItem.Transform.Pos = pixel.V((totalWidth - fullscreenOptionTextW * Size1) * 0.5, Spacing * -1.)
}

func SetOptionResolutionWidth() {
	resolutionItem := Options.Items["resolution"]
	resolutionTextW := resolutionItem.Text.Text.BoundsOf(resolutionItem.Text.Raw).W()
	resolutionOptItem := Options.Items["resolution_options"]
	resolutionOptTextW := resolutionOptItem.Text.Text.BoundsOf(resolutionOptItem.Text.Raw).W()
	totalWidth := (resolutionTextW + resolutionOptTextW) * Size1
	resolutionItem.Transform.Pos = pixel.V((totalWidth - resolutionTextW * Size1) * -0.5, Spacing * -2.)
	resolutionOptItem.Transform.Pos = pixel.V((totalWidth - resolutionOptTextW * Size1) * 0.5, Spacing * -2.)
}

func InitializePostGameMenu() {
	retryS := "Retry"
	retryText := menu.NewItemText(retryS, colornames.Aliceblue, pixel.V(2., 2.), menu.Center, menu.Top)
	retryText.HoverColor = colornames.Mediumblue
	retryR := pixel.R(0., 0., retryText.Text.BoundsOf(retryS).W() * 2.5, retryText.Text.BoundsOf(retryS).H() * 2.5)
	retryItem := menu.NewItem(retryText, retryR, PostGame.Canvas.Bounds())
	retryItem.Transform.Pos = pixel.V(-75., 25.)
	retryItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Bottom,
	}
	retryItem.SetOnHoverFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	retryItem.SetClickFn(func() {
		newState = 0
	})
	PostGame.Items["retry"] = retryItem

	menuS := "Main Menu"
	menuText := menu.NewItemText(menuS, colornames.Aliceblue, pixel.V(2., 2.), menu.Center, menu.Top)
	menuText.HoverColor = colornames.Mediumblue
	menuR := pixel.R(0., 0., menuText.Text.BoundsOf(menuS).W() * 2.5, menuText.Text.BoundsOf(menuS).H() * 2.5)
	menuItem := menu.NewItem(menuText, menuR, PostGame.Canvas.Bounds())
	menuItem.Transform.Pos = pixel.V(75., 25.)
	menuItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Bottom,
	}
	menuItem.SetOnHoverFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	PostGame.Items["menu"] = menuItem
}

func SwitchToOptions() {
	Options.Transform.Pos.X = 0.
	MainMenu.Transform.Pos.X = cfg.BaseW * 1.5
	Current = 1
}

func SwitchToMain() {
	Options.Transform.Pos.X = cfg.BaseW * 1.5
	MainMenu.Transform.Pos.X = 0.
	Current = 0
}