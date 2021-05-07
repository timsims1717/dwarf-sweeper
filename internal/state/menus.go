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
	Size1 = 6.
	Size2 = 7.
	Spacing = 100.
)

var (
	MainMenu    *menu.Menu
	Options     *menu.Menu
	PostGame    *menu.Menu
	Current     int
)

func InitializeMenus() {
	MainMenu = menu.NewMenu(pixel.R(0,0, camera.Cam.Width, camera.Cam.Height), camera.Cam)
	Options  = menu.NewMenu(pixel.R(0, 0, camera.Cam.Width, camera.Cam.Height), camera.Cam)
	PostGame = menu.NewMenu(pixel.R(0,0, camera.Cam.Width, camera.Cam.Height), camera.Cam)
	Current = 0
}

func InitializeMainMenu() {
	startS := "Start Game"
	startText := menu.NewItemText(startS, colornames.Aliceblue, pixel.V(Size1, Size1), true)
	startText.Transform.Anchor.V = transform.Center
	startText.HoverColor = colornames.Darkblue
	startText.HoverSize = pixel.V(Size2, Size2)
	startR := pixel.R(0., 0., startText.Text.BoundsOf(startS).W()*Size2, startText.Text.BoundsOf(startS).H()*Size2)
	startGame := menu.NewItem(startText, startR)
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
	optionsText := menu.NewItemText(optionsS, colornames.Aliceblue, pixel.V(Size1, Size1), true)
	optionsText.Transform.Anchor.V = transform.Center
	optionsText.HoverColor = colornames.Darkblue
	optionsText.HoverSize = pixel.V(Size2, Size2)
	optionsR := pixel.R(0., 0., optionsText.Text.BoundsOf(optionsS).W()*Size2, optionsText.Text.BoundsOf(optionsS).H()*Size2)
	optionsItem := menu.NewItem(optionsText, optionsR)
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
	creditsText := menu.NewItemText(creditsS, colornames.Aliceblue, pixel.V(Size1, Size1), true)
	creditsText.Transform.Anchor.V = transform.Center
	creditsText.HoverColor = colornames.Darkblue
	creditsText.HoverSize = pixel.V(Size2, Size2)
	creditsR := pixel.R(0., 0., creditsText.Text.BoundsOf(creditsS).W()*Size2, creditsText.Text.BoundsOf(creditsS).H()*Size2)
	creditsItem := menu.NewItem(creditsText, creditsR)
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
	exitText := menu.NewItemText(exitS, colornames.Aliceblue, pixel.V(Size1, Size1), true)
	exitText.Transform.Anchor.V = transform.Center
	exitText.HoverColor = colornames.Darkblue
	exitText.HoverSize = pixel.V(Size2, Size2)
	exitR := pixel.R(0., 0., exitText.Text.BoundsOf(exitS).W()*Size2, exitText.Text.BoundsOf(exitS).H()*Size2)
	exitItem := menu.NewItem(exitText, exitR)
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

func InitializeOptionsMenu() {
	Options.Transform.Pos.X = camera.Cam.Width

	soundVS := "Sound Volume:"
	soundVMinusS := " - "
	soundVolumeS := strconv.Itoa(sfx.GetSoundVolume())
	soundVPlusS := " + "

	soundVText := menu.NewItemText(soundVS, colornames.Aliceblue, pixel.V(Size1, Size1), true)
	soundVText.Transform.Anchor.V = transform.Center

	soundVolumeText := menu.NewItemText(soundVolumeS, colornames.Aliceblue, pixel.V(Size1, Size1), true)
	soundVolumeText.Transform.Anchor.V = transform.Center

	soundVMinusText := menu.NewItemText(soundVMinusS, colornames.Aliceblue, pixel.V(Size1, Size1), true)
	soundVMinusText.Transform.Anchor.V = transform.Center
	soundVMinusText.HoverColor = colornames.Darkblue
	soundVMinusText.HoverSize = pixel.V(Size2, Size2)

	soundVPlusText := menu.NewItemText(soundVPlusS, colornames.Aliceblue, pixel.V(Size1, Size1), true)
	soundVPlusText.Transform.Anchor.V = transform.Center
	soundVPlusText.HoverColor = colornames.Darkblue
	soundVPlusText.HoverSize = pixel.V(Size2, Size2)

	soundVR := pixel.R(0., 0., soundVText.Text.BoundsOf(soundVS).W()*Size2, soundVText.Text.BoundsOf(soundVS).H()*Size2)
	soundVolumeR := pixel.R(0., 0., soundVolumeText.Text.BoundsOf("100").W()*Size2, soundVolumeText.Text.BoundsOf(soundVolumeS).H()*Size2)
	soundVMinusR := pixel.R(0., 0., soundVMinusText.Text.BoundsOf(soundVMinusS).W()*Size2, soundVMinusText.Text.BoundsOf(soundVMinusS).H()*Size2)
	soundVPlusR := pixel.R(0., 0., soundVPlusText.Text.BoundsOf(soundVPlusS).W()*Size2, soundVPlusText.Text.BoundsOf(soundVPlusS).H()*Size2)

	soundVItem := menu.NewItem(soundVText, soundVR)
	soundVItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	soundVItem.Disabled = true
	Options.Items["sound_v"] = soundVItem

	soundVolumeItem := menu.NewItem(soundVolumeText, soundVolumeR)
	soundVolumeItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	soundVolumeItem.Disabled = true
	Options.Items["sound_volume"] = soundVolumeItem

	soundVMinusItem := menu.NewItem(soundVMinusText, soundVMinusR)
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
		soundVolumeItem.Text.Raw = strconv.Itoa(n)
		SetOptionSoundWidth()
	})
	Options.Items["sound_-"] = soundVMinusItem

	soundVPlusItem := menu.NewItem(soundVPlusText, soundVPlusR)
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
		soundVolumeItem.Text.Raw = strconv.Itoa(n)
		SetOptionSoundWidth()
	})
	Options.Items["sound_+"] = soundVPlusItem
	SetOptionSoundWidth()

	fullscreenS := "Fullscreen:"
	fullscreenOptionS := " Off"
	if cfg.FullScreen {
		fullscreenOptionS = " On"
	}

	fullscreenText := menu.NewItemText(fullscreenS, colornames.Aliceblue, pixel.V(Size1, Size1), true)
	fullscreenText.Transform.Anchor.V = transform.Center

	fullscreenOptionText := menu.NewItemText(fullscreenOptionS, colornames.Aliceblue, pixel.V(Size1, Size1), true)
	fullscreenOptionText.Transform.Anchor.V = transform.Center
	fullscreenOptionText.HoverColor = colornames.Darkblue
	fullscreenOptionText.HoverSize = pixel.V(Size2, Size2)

	fullscreenR := pixel.R(0., 0., fullscreenText.Text.BoundsOf(fullscreenS).W()*Size2, fullscreenText.Text.BoundsOf(fullscreenS).H()*Size2)
	fullscreenItem := menu.NewItem(fullscreenText, fullscreenR)
	fullscreenItem.Transform.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	fullscreenItem.Disabled = true
	Options.Items["fullscreen"] = fullscreenItem

	fullscreenOptionR := pixel.R(0., 0., fullscreenOptionText.Text.BoundsOf(" Off").W()*Size2, fullscreenOptionText.Text.BoundsOf(" Off").H()*Size2)
	fullscreenOptionItem := menu.NewItem(fullscreenOptionText, fullscreenOptionR)
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
		fullscreenOptionItem.Text.Raw = s
		SetOptionSoundWidth()
	})
	Options.Items["fullscreen_options"] = fullscreenOptionItem
	SetOptionFullscreenWidth()

	backS := "Back"
	backText := menu.NewItemText(backS, colornames.Aliceblue, pixel.V(Size1, Size1), true)
	backText.Transform.Anchor.V = transform.Center
	backText.HoverColor = colornames.Darkblue
	backText.HoverSize = pixel.V(Size2, Size2)
	backR := pixel.R(0., 0., backText.Text.BoundsOf(backS).W()*Size2, backText.Text.BoundsOf(backS).H()*Size2)
	backItem := menu.NewItem(backText, backR)
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

func InitializePostGameMenu() {
	retryS := "Retry"
	retryText := menu.NewItemText(retryS, colornames.Aliceblue, pixel.V(4., 4.), true)
	retryText.Transform.Anchor.V = transform.Top
	retryText.HoverColor = colornames.Mediumblue
	retryR := pixel.R(0., 0., retryText.Text.BoundsOf(retryS).W()*5., retryText.Text.BoundsOf(retryS).H()*5.)
	retryItem := menu.NewItem(retryText, retryR)
	retryItem.Transform.Pos = pixel.V(-250., 200.)
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
	menuText := menu.NewItemText(menuS, colornames.Aliceblue, pixel.V(4., 4.), true)
	menuText.Transform.Anchor.V = transform.Top
	menuText.HoverColor = colornames.Mediumblue
	menuR := pixel.R(0., 0., menuText.Text.BoundsOf(menuS).W()*5., menuText.Text.BoundsOf(menuS).H()*5.)
	menuItem := menu.NewItem(menuText, menuR)
	menuItem.Transform.Pos = pixel.V(250., 200.)
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
	MainMenu.Transform.Pos.X = camera.Cam.Width
	Current = 1
}

func SwitchToMain() {
	Options.Transform.Pos.X = camera.Cam.Width
	MainMenu.Transform.Pos.X = 0.
	Current = 0
}