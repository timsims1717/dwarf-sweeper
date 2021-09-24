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

var (
	PostGame    *menu.Menu
	MainMenu    *menus.DwarfMenu
	PauseMenu   *menus.DwarfMenu
	OptionsMenu *menus.DwarfMenu
	EnchantMenu *menus.DwarfMenu
)

func InitializeMenus(win *pixelgl.Window) {
	PostGame    = menu.NewMenu(pixel.R(0,0, cfg.BaseW, cfg.BaseH), camera.Cam)
	InitMainMenu(win)
	InitOptionsMenu()
	InitPauseMenu(win)
	InitEnchantMenu()
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

func InitMainMenu(win *pixelgl.Window) {
	MainMenu = menus.New("main", pixel.R(0., 0., 64., 64.), camera.Cam)
	start := MainMenu.AddItem("start", "Start Game")
	options := MainMenu.AddItem("options", "Options")
	credit := MainMenu.AddItem("credits", "Credits")
	quit := MainMenu.AddItem("quit", "Quit")

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
	optionsTitle := OptionsMenu.AddItem("title", "Options")
	volume := OptionsMenu.AddItem("s_volume", "Volume")
	volumeR := OptionsMenu.AddItem("s_volume_r", strconv.Itoa(sfx.GetSoundVolume()))
	vsync := OptionsMenu.AddItem("vsync", "VSync")
	vsyncR := OptionsMenu.AddItem("vsync_r", "On")
	fullscreen := OptionsMenu.AddItem("fullscreen", "Fullscreen")
	fullscreenR := OptionsMenu.AddItem("fullscreen_r", "Off")
	resolution := OptionsMenu.AddItem("resolution", "Resolution")
	resolutionR := OptionsMenu.AddItem("resolution_r", cfg.ResStrings[cfg.ResIndex])
	back := OptionsMenu.AddItem("back", "Back")

	optionsTitle.NoHover = true
	volume.SetRightFn(func() {
		n := sfx.GetSoundVolume() + 5
		if n > 100 {
			n = 100
		}
		sfx.SetSoundVolume(n)
		volumeR.Raw = strconv.Itoa(n)
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
	})
	volume.SetLeftFn(func() {
		n := sfx.GetSoundVolume() - 5
		if n < 0 {
			n = 0
		}
		sfx.SetSoundVolume(n)
		volumeR.Raw = strconv.Itoa(n)
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
	})
	volumeR.NoHover = true
	volumeR.Right = true
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
		OptionsMenu.Close()
	})
}

func InitPauseMenu(win *pixelgl.Window) {
	PauseMenu = menus.New("pause", pixel.R(0., 0., 64., 64.), camera.Cam)
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
	EnchantMenu = menus.New("enchant", pixel.R(0., 0., 64., 64.), camera.Cam)
	chooseTitle := EnchantMenu.AddItem("title", "Choose an Enchantment")
	skip := EnchantMenu.AddItem("skip", "Skip")

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
	if e2 != nil {
		option2 := EnchantMenu.InsertItem("option2", e2.Title, 2)
		option2.SetClickFn(func() {
			sfx.SoundPlayer.PlaySound("click", 2.0)
			enchants.AddEnchantment(e2)
			EnchantMenu.CloseInstant()
			ClearEnchantMenu()
			newState = 0
		})
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
	}
	return true
}