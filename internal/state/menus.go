package state

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/dungeon"
	"dwarf-sweeper/internal/enchants"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/menu"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"strconv"
)

const (
	Size1 = 2.
	Size2 = 2.
	Spacing = 28.
)

var (
	//Main        *menu.Menu
	//Options     *menu.Menu
	PostGame    *menu.Menu
	//EnchantShop *menu.Menu
	MainMenu    *menus.DwarfMenu
	PauseMenu   *menus.DwarfMenu
	OptionsMenu *menus.DwarfMenu
	EnchantMenu *menus.DwarfMenu
	//Current     int
)

func InitializeMenus(win *pixelgl.Window) {
	//Main = menu.NewMenu(pixel.R(0,0, cfg.BaseW, cfg.BaseH), camera.Cam)
	//Options     = menu.NewMenu(pixel.R(0,0, cfg.BaseW, cfg.BaseH), camera.Cam)
	PostGame    = menu.NewMenu(pixel.R(0,0, cfg.BaseW, cfg.BaseH), camera.Cam)
	//EnchantShop = menu.NewMenu(pixel.R(0, 0, cfg.BaseW, cfg.BaseH), camera.Cam)
	InitMainMenu(win)
	InitOptionsMenu()
	InitPauseMenu(win)
	InitEnchantMenu()
	//Current = 0
}

//func SwitchToOptions() {
//	Options.Transform.Pos.X = 0.
//	Main.Transform.Pos.X = cfg.BaseW * 1.5
//	Current = 1
//}
//
//func SwitchToMain() {
//	Options.Transform.Pos.X = cfg.BaseW * 1.5
//	Main.Transform.Pos.X = 0.
//	Current = 0
//}

//func InitializeMainMenu() {
//	startS := "Start Game"
//	startText := menu.NewItemText(startS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//	startText.HoverColor = colornames.Darkblue
//	startText.HoverSize = pixel.V(Size2, Size2)
//	startR := pixel.R(0., 0., startText.Text.BoundsOf(startS).W()*Size2, startText.Text.BoundsOf(startS).H()*Size2)
//	startGame := menu.NewItem(startText, startR, Main.Canvas.Bounds())
//	startGame.Transform.Pos = pixel.V(0., 0.)
//	startGame.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	startGame.SetOnHoverFn(func() {
//		sfx.SoundPlayer.PlaySound("click", 2.0)
//	})
//	startGame.SetClickFn(func() {
//		//camera.Cam.Effect = transform.FadeTo(camera.Cam, colornames.Black, 1.0)
//		//sfx.MusicPlayer.FadeOut(1.0)
//		newState = 4
//	})
//	Main.Items["start"] = startGame
//
//	optionsS := "Options"
//	optionsText := menu.NewItemText(optionsS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//	optionsText.HoverColor = colornames.Darkblue
//	optionsText.HoverSize = pixel.V(Size2, Size2)
//	optionsR := pixel.R(0., 0., optionsText.Text.BoundsOf(optionsS).W()*Size2, optionsText.Text.BoundsOf(optionsS).H()*Size2)
//	optionsItem := menu.NewItem(optionsText, optionsR, Main.Canvas.Bounds())
//	optionsItem.Transform.Pos = pixel.V(0., Spacing * -1.)
//	optionsItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	optionsItem.SetOnHoverFn(func() {
//		sfx.SoundPlayer.PlaySound("click", 2.0)
//	})
//	optionsItem.SetClickFn(func() {
//		SwitchToOptions()
//	})
//	Main.Items["options"] = optionsItem
//
//	creditsS := "Credits"
//	creditsText := menu.NewItemText(creditsS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//	creditsText.HoverColor = colornames.Darkblue
//	creditsText.HoverSize = pixel.V(Size2, Size2)
//	creditsR := pixel.R(0., 0., creditsText.Text.BoundsOf(creditsS).W()*Size2, creditsText.Text.BoundsOf(creditsS).H()*Size2)
//	creditsItem := menu.NewItem(creditsText, creditsR, Main.Canvas.Bounds())
//	creditsItem.Transform.Pos = pixel.V(0., Spacing * -2.)
//	creditsItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	creditsItem.SetOnHoverFn(func() {
//		sfx.SoundPlayer.PlaySound("click", 2.0)
//	})
//	creditsItem.SetClickFn(func() {
//		newState = 3
//	})
//	Main.Items["credits"] = creditsItem
//
//	exitS := "Exit"
//	exitText := menu.NewItemText(exitS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//	exitText.HoverColor = colornames.Darkblue
//	exitText.HoverSize = pixel.V(Size2, Size2)
//	exitR := pixel.R(0., 0., exitText.Text.BoundsOf(exitS).W()*Size2, exitText.Text.BoundsOf(exitS).H()*Size2)
//	exitItem := menu.NewItem(exitText, exitR, Main.Canvas.Bounds())
//	exitItem.Transform.Pos = pixel.V(0., Spacing * -3.)
//	exitItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	exitItem.SetOnHoverFn(func() {
//		sfx.SoundPlayer.PlaySound("click", 2.0)
//	})
//	Main.Items["exit"] = exitItem
//}

//func InitializeSoundOption() {
//	soundVS := "Sound Volume:"
//	soundVMinusS := " - "
//	soundVolumeS := strconv.Itoa(sfx.GetSoundVolume())
//	soundVPlusS := " + "
//
//	soundVText := menu.NewItemText(soundVS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//
//	soundVolumeText := menu.NewItemText(soundVolumeS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//
//	soundVMinusText := menu.NewItemText(soundVMinusS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//	soundVMinusText.HoverColor = colornames.Darkblue
//	soundVMinusText.HoverSize = pixel.V(Size2, Size2)
//
//	soundVPlusText := menu.NewItemText(soundVPlusS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//	soundVPlusText.HoverColor = colornames.Darkblue
//	soundVPlusText.HoverSize = pixel.V(Size2, Size2)
//
//	soundVR := pixel.R(0., 0., soundVText.Text.BoundsOf(soundVS).W()*Size2, soundVText.Text.BoundsOf(soundVS).H()*Size2)
//	soundVolumeR := pixel.R(0., 0., soundVolumeText.Text.BoundsOf("100").W()*Size2, soundVolumeText.Text.BoundsOf(soundVolumeS).H()*Size2)
//	soundVMinusR := pixel.R(0., 0., soundVMinusText.Text.BoundsOf(soundVMinusS).W()*Size2, soundVMinusText.Text.BoundsOf(soundVMinusS).H()*Size2)
//	soundVPlusR := pixel.R(0., 0., soundVPlusText.Text.BoundsOf(soundVPlusS).W()*Size2, soundVPlusText.Text.BoundsOf(soundVPlusS).H()*Size2)
//	soundVText.Transform.SetParent(soundVR)
//	soundVolumeText.Transform.SetParent(soundVolumeR)
//	soundVMinusText.Transform.SetParent(soundVMinusR)
//	soundVPlusText.Transform.SetParent(soundVPlusR)
//
//	soundVItem := menu.NewItem(soundVText, soundVR, Options.Canvas.Bounds())
//	soundVItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	soundVItem.Disabled = true
//	Options.Items["sound_v"] = soundVItem
//
//	soundVolumeItem := menu.NewItem(soundVolumeText, soundVolumeR, Options.Canvas.Bounds())
//	soundVolumeItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	soundVolumeItem.Disabled = true
//	Options.Items["sound_volume"] = soundVolumeItem
//
//	soundVMinusItem := menu.NewItem(soundVMinusText, soundVMinusR, Options.Canvas.Bounds())
//	soundVMinusItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	soundVMinusItem.SetOnHoverFn(func() {
//		sfx.SoundPlayer.PlaySound("click", 2.0)
//	})
//	soundVMinusItem.SetClickFn(func() {
//		n := sfx.GetSoundVolume() - 5
//		if n < 0 {
//			n = 0
//		}
//		sfx.SetSoundVolume(n)
//		soundVolumeItem.Text.SetText(strconv.Itoa(n))
//		SetOptionSoundWidth()
//	})
//	Options.Items["sound_-"] = soundVMinusItem
//
//	soundVPlusItem := menu.NewItem(soundVPlusText, soundVPlusR, Options.Canvas.Bounds())
//	soundVPlusItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	soundVPlusItem.SetOnHoverFn(func() {
//		sfx.SoundPlayer.PlaySound("click", 2.0)
//	})
//	soundVPlusItem.SetClickFn(func() {
//		n := sfx.GetSoundVolume() + 5
//		if n > 100 {
//			n = 100
//		}
//		sfx.SetSoundVolume(n)
//		soundVolumeItem.Text.SetText(strconv.Itoa(n))
//		SetOptionSoundWidth()
//	})
//	Options.Items["sound_+"] = soundVPlusItem
//	SetOptionSoundWidth()
//}
//
//func InitializeVSyncOption() {
//	vsyncS := "VSync:"
//	vsyncOptionS := " Off"
//	if cfg.VSync {
//		vsyncOptionS = " On"
//	}
//
//	vsyncText := menu.NewItemText(vsyncS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//
//	vsyncOptionText := menu.NewItemText(vsyncOptionS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//	vsyncOptionText.HoverColor = colornames.Darkblue
//	vsyncOptionText.HoverSize = pixel.V(Size2, Size2)
//
//	vsyncR := pixel.R(0., 0., vsyncText.Text.BoundsOf(vsyncS).W()*Size2, vsyncText.Text.BoundsOf(vsyncS).H()*Size2)
//	vsyncItem := menu.NewItem(vsyncText, vsyncR, Options.Canvas.Bounds())
//	vsyncItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	vsyncItem.Disabled = true
//	Options.Items["vsync"] = vsyncItem
//
//	vsyncOptionR := pixel.R(0., 0., vsyncOptionText.Text.BoundsOf(" Off").W()*Size2, vsyncOptionText.Text.BoundsOf(" Off").H()*Size2)
//	vsyncOptionItem := menu.NewItem(vsyncOptionText, vsyncOptionR, Options.Canvas.Bounds())
//	vsyncOptionItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	vsyncOptionItem.SetOnHoverFn(func() {
//		sfx.SoundPlayer.PlaySound("click", 2.0)
//	})
//	vsyncOptionItem.SetClickFn(func() {
//		s := " On"
//		if cfg.VSync {
//			s = " Off"
//		}
//		cfg.VSync = !cfg.VSync
//		vsyncOptionItem.Text.SetText(s)
//		SetOptionVSyncWidth()
//	})
//	Options.Items["vsync_options"] = vsyncOptionItem
//	SetOptionVSyncWidth()
//}
//
//func InitializeFullscreenOption() {
//	fullscreenS := "Fullscreen:"
//	fullscreenOptionS := " Off"
//	if cfg.FullScreen {
//		fullscreenOptionS = " On"
//	}
//
//	fullscreenText := menu.NewItemText(fullscreenS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//
//	fullscreenOptionText := menu.NewItemText(fullscreenOptionS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//	fullscreenOptionText.HoverColor = colornames.Darkblue
//	fullscreenOptionText.HoverSize = pixel.V(Size2, Size2)
//
//	fullscreenR := pixel.R(0., 0., fullscreenText.Text.BoundsOf(fullscreenS).W()*Size2, fullscreenText.Text.BoundsOf(fullscreenS).H()*Size2)
//	fullscreenItem := menu.NewItem(fullscreenText, fullscreenR, Options.Canvas.Bounds())
//	fullscreenItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	fullscreenItem.Disabled = true
//	Options.Items["fullscreen"] = fullscreenItem
//
//	fullscreenOptionR := pixel.R(0., 0., fullscreenOptionText.Text.BoundsOf(" Off").W()*Size2, fullscreenOptionText.Text.BoundsOf(" Off").H()*Size2)
//	fullscreenOptionItem := menu.NewItem(fullscreenOptionText, fullscreenOptionR, Options.Canvas.Bounds())
//	fullscreenOptionItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	fullscreenOptionItem.SetOnHoverFn(func() {
//		sfx.SoundPlayer.PlaySound("click", 2.0)
//	})
//	fullscreenOptionItem.SetClickFn(func() {
//		s := " On"
//		if cfg.FullScreen {
//			s = " Off"
//		}
//		cfg.FullScreen = !cfg.FullScreen
//		cfg.ChangeScreenSize = true
//		fullscreenOptionItem.Text.SetText(s)
//		SetOptionFullscreenWidth()
//	})
//	Options.Items["fullscreen_options"] = fullscreenOptionItem
//	SetOptionFullscreenWidth()
//}
//
//func InitializeResolutionOption() {
//	resolutionS := "Resolution:"
//	resolutionOptS := cfg.ResStrings[cfg.ResIndex]
//
//	resolutionText := menu.NewItemText(resolutionS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//
//	resolutionOptText := menu.NewItemText(resolutionOptS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//	resolutionOptText.HoverColor = colornames.Darkblue
//	resolutionOptText.HoverSize = pixel.V(Size2, Size2)
//
//	resolutionR := pixel.R(0., 0., resolutionText.Text.BoundsOf(resolutionS).W()*Size2, resolutionText.Text.BoundsOf(resolutionS).H()*Size2)
//	resolutionItem := menu.NewItem(resolutionText, resolutionR, Options.Canvas.Bounds())
//	resolutionItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	resolutionItem.Disabled = true
//	Options.Items["resolution"] = resolutionItem
//
//	resolutionOptR := pixel.R(0., 0., resolutionOptText.Text.BoundsOf(cfg.ResStrings[len(cfg.ResStrings)-1]).W()*Size2, resolutionOptText.Text.BoundsOf(cfg.ResStrings[len(cfg.ResStrings)-1]).H()*Size2)
//	resolutionOptItem := menu.NewItem(resolutionOptText, resolutionOptR, Options.Canvas.Bounds())
//	resolutionOptItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	resolutionOptItem.SetOnHoverFn(func() {
//		sfx.SoundPlayer.PlaySound("click", 2.0)
//	})
//	resolutionOptItem.SetClickFn(func() {
//		cfg.ResIndex += 1
//		cfg.ResIndex %= len(cfg.Resolutions)
//		cfg.ChangeScreenSize = true
//		resolutionOptItem.Text.SetText(cfg.ResStrings[cfg.ResIndex])
//		SetOptionResolutionWidth()
//	})
//	Options.Items["resolution_options"] = resolutionOptItem
//	SetOptionResolutionWidth()
//}
//
//func InitializeOptionsMenu() {
//	Options.Transform.Pos.X = cfg.BaseW * 1.5
//
//	InitializeSoundOption()
//	InitializeVSyncOption()
//	InitializeFullscreenOption()
//	InitializeResolutionOption()
//
//	backS := "Back"
//	backText := menu.NewItemText(backS, colornames.Aliceblue, pixel.V(Size1, Size1), menu.Center, menu.Center)
//	backText.HoverColor = colornames.Darkblue
//	backText.HoverSize = pixel.V(Size2, Size2)
//	backR := pixel.R(0., 0., backText.Text.BoundsOf(backS).W()*Size2, backText.Text.BoundsOf(backS).H()*Size2)
//	backItem := menu.NewItem(backText, backR, Options.Canvas.Bounds())
//	backItem.Transform.Pos = pixel.V(0., Spacing * -3.)
//	backItem.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Center,
//	}
//	backItem.SetOnHoverFn(func() {
//		sfx.SoundPlayer.PlaySound("click", 2.0)
//	})
//	Options.Items["back"] = backItem
//}
//
//func SetOptionSoundWidth() {
//	soundVItem := Options.Items["sound_v"]
//	soundVTextW := soundVItem.Text.Text.BoundsOf(soundVItem.Text.Raw).W()
//	soundVMinusItem := Options.Items["sound_-"]
//	soundVMinusTextW := soundVMinusItem.Text.Text.BoundsOf(soundVMinusItem.Text.Raw).W()
//	soundVolumeItem := Options.Items["sound_volume"]
//	soundVVolumeTextW := soundVolumeItem.Text.Text.BoundsOf(soundVolumeItem.Text.Raw).W()
//	soundVPlusItem := Options.Items["sound_+"]
//	soundVPlusTextW := soundVPlusItem.Text.Text.BoundsOf(soundVPlusItem.Text.Raw).W()
//	totalWidth := (soundVTextW + soundVMinusTextW + soundVVolumeTextW + soundVPlusTextW) * Size1
//	soundVItem.Transform.Pos = pixel.V((totalWidth - soundVTextW * Size1) * -0.5, Spacing * 1.)
//	soundVMinusItem.Transform.Pos = pixel.V(totalWidth * -0.5 + (soundVTextW + soundVMinusTextW * 0.5) * Size1, Spacing)
//	soundVolumeItem.Transform.Pos = pixel.V(totalWidth * 0.5 - (soundVPlusTextW + soundVVolumeTextW * 0.5) * Size1, Spacing)
//	soundVPlusItem.Transform.Pos = pixel.V((totalWidth - soundVPlusTextW * Size1) * 0.5, Spacing * 1.)
//}
//
//func SetOptionVSyncWidth() {
//	vsyncItem := Options.Items["vsync"]
//	vsyncTextW := vsyncItem.Text.Text.BoundsOf(vsyncItem.Text.Raw).W()
//	vsyncOptionItem := Options.Items["vsync_options"]
//	vsyncOptionTextW := vsyncOptionItem.Text.Text.BoundsOf(vsyncOptionItem.Text.Raw).W()
//	totalWidth := (vsyncTextW + vsyncOptionTextW) * Size1
//	vsyncItem.Transform.Pos = pixel.V((totalWidth - vsyncTextW * Size1) * -0.5, 0.)
//	vsyncOptionItem.Transform.Pos = pixel.V((totalWidth - vsyncOptionTextW * Size1) * 0.5, 0.)
//}
//
//func SetOptionFullscreenWidth() {
//	fullscreenItem := Options.Items["fullscreen"]
//	fullscreenTextW := fullscreenItem.Text.Text.BoundsOf(fullscreenItem.Text.Raw).W()
//	fullscreenOptionItem := Options.Items["fullscreen_options"]
//	fullscreenOptionTextW := fullscreenOptionItem.Text.Text.BoundsOf(fullscreenOptionItem.Text.Raw).W()
//	totalWidth := (fullscreenTextW + fullscreenOptionTextW) * Size1
//	fullscreenItem.Transform.Pos = pixel.V((totalWidth - fullscreenTextW * Size1) * -0.5, Spacing * -1.)
//	fullscreenOptionItem.Transform.Pos = pixel.V((totalWidth - fullscreenOptionTextW * Size1) * 0.5, Spacing * -1.)
//}
//
//func SetOptionResolutionWidth() {
//	resolutionItem := Options.Items["resolution"]
//	resolutionTextW := resolutionItem.Text.Text.BoundsOf(resolutionItem.Text.Raw).W()
//	resolutionOptItem := Options.Items["resolution_options"]
//	resolutionOptTextW := resolutionOptItem.Text.Text.BoundsOf(resolutionOptItem.Text.Raw).W()
//	totalWidth := (resolutionTextW + resolutionOptTextW) * Size1
//	resolutionItem.Transform.Pos = pixel.V((totalWidth - resolutionTextW * Size1) * -0.5, Spacing * -2.)
//	resolutionOptItem.Transform.Pos = pixel.V((totalWidth - resolutionOptTextW * Size1) * 0.5, Spacing * -2.)
//}

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
		newState = 4
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

//func InitializeEnchantShopMenu() bool {
//	EnchantShop.Items = map[string]*menu.Item{}
//	pe := dungeon.Dungeon.GetPlayer().Enchants
//	list := enchants.Enchantments
//	for _, i := range pe {
//		if len(list) > 1 {
//			list = append(list[:i], list[i+1:]...)
//		} else {
//			return false
//		}
//	}
//	choices := util.RandomSampleRange(util.Min(len(list), 3), 0, len(list), random.CaveGen)
//	var e1, e2, e3 *data.Enchantment
//	for i, c := range choices {
//		if i == 0 {
//			e1 = list[c]
//		} else if i == 1 {
//			e2 = list[c]
//		} else if i == 2 {
//			e3 = list[c]
//		}
//	}
//	if e1 == nil {
//		return false
//	}
//	enchant1S := e1.Title
//	enchant1Text := menu.NewItemText(enchant1S, colornames.Aliceblue, pixel.V(1.1, 1.1), menu.Center, menu.Top)
//	enchant1Text.HoverColor = colornames.Mediumblue
//	enchant1R := pixel.R(0., 0., enchant1Text.Text.BoundsOf(enchant1S).W() * 1.2, enchant1Text.Text.BoundsOf(enchant1S).H() * 1.2)
//	enchant1Item := menu.NewItem(enchant1Text, enchant1R, EnchantShop.Canvas.Bounds())
//	enchant1Item.Transform.Pos = pixel.V(0., 50.)
//	enchant1Item.Transform.Anchor = transform.Anchor{
//		H: transform.Center,
//		V: transform.Bottom,
//	}
//	enchant1Item.SetOnHoverFn(func() {
//		sfx.SoundPlayer.PlaySound("click", 2.0)
//	})
//	enchant1Item.SetClickFn(func() {
//		enchants.AddEnchantment(e1)
//		newState = 0
//	})
//	EnchantShop.Items["enchant1"] = enchant1Item
//
//	if e2 != nil {
//		enchant2S := e2.Title
//		enchant2Text := menu.NewItemText(enchant2S, colornames.Aliceblue, pixel.V(1.1, 1.1), menu.Center, menu.Top)
//		enchant2Text.HoverColor = colornames.Mediumblue
//		enchant2R := pixel.R(0., 0., enchant2Text.Text.BoundsOf(enchant2S).W()*1.2, enchant2Text.Text.BoundsOf(enchant2S).H()*1.2)
//		enchant2Item := menu.NewItem(enchant2Text, enchant2R, EnchantShop.Canvas.Bounds())
//		enchant2Item.Transform.Pos = pixel.V(0., 100.)
//		enchant2Item.Transform.Anchor = transform.Anchor{
//			H: transform.Center,
//			V: transform.Bottom,
//		}
//		enchant2Item.SetOnHoverFn(func() {
//			sfx.SoundPlayer.PlaySound("click", 2.0)
//		})
//		enchant2Item.SetClickFn(func() {
//			enchants.AddEnchantment(e2)
//			newState = 0
//		})
//		EnchantShop.Items["enchant2"] = enchant2Item
//	}
//
//	if e3 != nil {
//		enchant3S := e3.Title
//		enchant3Text := menu.NewItemText(enchant3S, colornames.Aliceblue, pixel.V(1.1, 1.1), menu.Center, menu.Top)
//		enchant3Text.HoverColor = colornames.Mediumblue
//		enchant3R := pixel.R(0., 0., enchant3Text.Text.BoundsOf(enchant3S).W()*1.2, enchant3Text.Text.BoundsOf(enchant3S).H()*1.2)
//		enchant3Item := menu.NewItem(enchant3Text, enchant3R, EnchantShop.Canvas.Bounds())
//		enchant3Item.Transform.Pos = pixel.V(0., 150.)
//		enchant3Item.Transform.Anchor = transform.Anchor{
//			H: transform.Center,
//			V: transform.Bottom,
//		}
//		enchant3Item.SetOnHoverFn(func() {
//			sfx.SoundPlayer.PlaySound("click", 2.0)
//		})
//		enchant3Item.SetClickFn(func() {
//			enchants.AddEnchantment(e3)
//			newState = 0
//		})
//		EnchantShop.Items["enchant3"] = enchant3Item
//	}
//	return true
//}

func InitMainMenu(win *pixelgl.Window) {
	MainMenu = menus.New("main", pixel.R(0., 0., 64., 64.), camera.Cam)
	start := MainMenu.AddItem("start", "Start Game", "")
	options := MainMenu.AddItem("options", "Options", "")
	credit := MainMenu.AddItem("credits", "Credits", "")
	quit := MainMenu.AddItem("quit", "Quit", "")

	start.SetClickFn(func() {
		MainMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		newState = 4
	})
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
}

func InitOptionsMenu() {
	OptionsMenu = menus.New("options", pixel.R(0., 0., 64., 64.), camera.Cam)
	optionsTitle := OptionsMenu.AddItem("title", "Options", "")
	volume := OptionsMenu.AddItem("s_volume", "Volume", strconv.Itoa(sfx.GetSoundVolume()))
	vsync := OptionsMenu.AddItem("vsync", "VSync", "On")
	fullscreen := OptionsMenu.AddItem("fullscreen", "Fullscreen", "Off")
	resolution := OptionsMenu.AddItem("resolution", "Resolution", cfg.ResStrings[cfg.ResIndex])
	back := OptionsMenu.AddItem("back", "Back", "")

	optionsTitle.NoHover = true
	volume.SetRightFn(func() {
		n := sfx.GetSoundVolume() + 5
		if n > 100 {
			n = 100
		}
		sfx.SetSoundVolume(n)
		volume.SRaw = strconv.Itoa(n)
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
	})
	volume.SetLeftFn(func() {
		n := sfx.GetSoundVolume() - 5
		if n < 0 {
			n = 0
		}
		sfx.SetSoundVolume(n)
		volume.SRaw = strconv.Itoa(n)
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
	})
	vsync.SetClickFn(func() {
		s := "On"
		if cfg.VSync {
			s = "Off"
		}
		cfg.VSync = !cfg.VSync
		vsync.SRaw = s
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	fullscreen.SetClickFn(func() {
		s := "On"
		if cfg.FullScreen {
			s = "Off"
		}
		cfg.FullScreen = !cfg.FullScreen
		cfg.ChangeScreenSize = true
		fullscreen.SRaw = s
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	fn := func() {
		cfg.ResIndex += 1
		cfg.ResIndex %= len(cfg.Resolutions)
		cfg.ChangeScreenSize = true
		resolution.SRaw = cfg.ResStrings[cfg.ResIndex]
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
	resolution.SetClickFn(fn)
	resolution.SetRightFn(fn)
	resolution.SetLeftFn(func() {
		cfg.ResIndex += len(cfg.Resolutions) - 1
		cfg.ResIndex %= len(cfg.Resolutions)
		cfg.ChangeScreenSize = true
		resolution.SRaw = cfg.ResStrings[cfg.ResIndex]
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	back.SetClickFn(func() {
		OptionsMenu.Close()
	})
}

func InitPauseMenu(win *pixelgl.Window) {
	PauseMenu = menus.New("pause", pixel.R(0., 0., 64., 64.), camera.Cam)
	pauseTitle := PauseMenu.AddItem("title", "Paused", "")
	resume := PauseMenu.AddItem("resume", "Resume", "")
	options := PauseMenu.AddItem("options", "Options", "")
	mainMenu := PauseMenu.AddItem("main_menu", "Abandon Run", "")
	quit := PauseMenu.AddItem("quit", "Quit Game", "")

	pauseTitle.NoHover = true
	resume.SetClickFn(func() {
		PauseMenu.Close()
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
	EnchantMenu = menus.New("enchant", pixel.R(0., 0., 64., 64.), camera.Cam)
	chooseTitle := EnchantMenu.AddItem("title", "Choose an Enchantment", "")
	skip := EnchantMenu.AddItem("skip", "Skip", "")

	chooseTitle.NoHover = true
	skip.SetClickFn(func() {
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
	option1 := EnchantMenu.InsertItem("option1", e1.Title, "", 1)
	option1.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		enchants.AddEnchantment(e1)
		EnchantMenu.CloseInstant()
		ClearEnchantMenu()
		newState = 0
	})
	if e2 != nil {
		option2 := EnchantMenu.InsertItem("option2", e2.Title, "", 2)
		option2.SetClickFn(func() {
			sfx.SoundPlayer.PlaySound("click", 2.0)
			enchants.AddEnchantment(e2)
			EnchantMenu.CloseInstant()
			ClearEnchantMenu()
			newState = 0
		})
	}
	if e3 != nil {
		option3 := EnchantMenu.InsertItem("option3", e3.Title, "", 3)
		option3.SetClickFn(func() {
			sfx.SoundPlayer.PlaySound("click", 2.0)
			enchants.AddEnchantment(e3)
			EnchantMenu.CloseInstant()
			ClearEnchantMenu()
			newState = 0
		})
	}
	return true
}