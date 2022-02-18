package states

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/enchants"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/sfx"
	"github.com/faiface/pixel/pixelgl"
)

func InitPauseMenu(win *pixelgl.Window) {
	PauseMenu = menus.New("pause", camera.Cam)
	PauseMenu.Title = true
	PauseMenu.SetCloseFn(func() {
		sfx.MusicPlayer.Pause("pause", true)
		sfx.MusicPlayer.Resume(constants.GameMusic)
	})
	pauseTitle := PauseMenu.AddItem("title", "Paused", false)
	resume := PauseMenu.AddItem("resume", "Resume", false)
	options := PauseMenu.AddItem("options", "Options", false)
	mainMenu := PauseMenu.AddItem("main_menu", "Abandon Run", false)
	quit := PauseMenu.AddItem("quit", "Quit Game", false)

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
		SwitchState(MenuStateKey)
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
	chooseTitle := EnchantMenu.AddItem("title", "Enchant!", false)
	skip := EnchantMenu.AddItem("skip", "Skip", false)

	chooseTitle.NoHover = true
	skip.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		EnchantMenu.CloseInstant()
		ClearEnchantMenu()
		SwitchState(DescentStateKey)
	})
}

func ClearEnchantMenu() {
	EnchantMenu.RemoveItem("option1")
	EnchantMenu.RemoveItem("option2")
	EnchantMenu.RemoveItem("option3")
}

func FillEnchantMenu() bool {
	ClearEnchantMenu()
	choices := enchants.PickEnchantments(descent.Descent.Player.Enchants)
	if len(choices) == 0 {
		return false
	}
	e1 := choices[0]
	option1 := EnchantMenu.InsertItem("option1", e1.Title, 1, false)
	option1.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		enchants.AddEnchantment(e1)
		EnchantMenu.CloseInstant()
		ClearEnchantMenu()
		SwitchState(DescentStateKey)
	})
	option1.Hint = e1.Desc
	if len(choices) > 1 {
		e2 := choices[1]
		option2 := EnchantMenu.InsertItem("option2", e2.Title, 2, false)
		option2.SetClickFn(func() {
			sfx.SoundPlayer.PlaySound("click", 2.0)
			enchants.AddEnchantment(e2)
			EnchantMenu.CloseInstant()
			ClearEnchantMenu()
			SwitchState(DescentStateKey)
		})
		option2.Hint = e2.Desc
	}
	if len(choices) > 2 {
		e3 := choices[2]
		option3 := EnchantMenu.InsertItem("option3", e3.Title, 3, false)
		option3.SetClickFn(func() {
			sfx.SoundPlayer.PlaySound("click", 2.0)
			enchants.AddEnchantment(e3)
			EnchantMenu.CloseInstant()
			ClearEnchantMenu()
			SwitchState(DescentStateKey)
		})
		option3.Hint = e3.Desc
	}
	return true
}

func InitPostGameMenu() {
	PostMenu = menus.New("post", camera.Cam)
	PostMenu.Title = true
	PostMenu.SetBackFn(func() {})
	blocksDug := PostMenu.AddItem("blocks", "Blocks Dug", false)
	blocksDugS := PostMenu.AddItem("blocks_s", "", true)
	gems := PostMenu.AddItem("gem_count", "Gems Found", false)
	gemsS := PostMenu.AddItem("gem_count_s", "", true)
	bombs := PostMenu.AddItem("bombs_flagged", "Bombs Flagged", false)
	bombsS := PostMenu.AddItem("bombs_flagged_s", "", true)
	wrongFlags := PostMenu.AddItem("wrong_flags", "Incorrect Flags", false)
	wrongFlagsS := PostMenu.AddItem("wrong_flags_s", "", true)
	totalScore := PostMenu.AddItem("total_score", "Total Score", false)
	totalScoreS := PostMenu.AddItem("total_score_s", "", true)
	retry := PostMenu.AddItem("retry", "Retry", false)
	backToMenu := PostMenu.AddItem("menu", "Main Menu", true)

	blocksDug.NoHover = true
	blocksDug.NoDraw = true
	blocksDugS.NoHover = true
	blocksDugS.NoDraw = true
	gems.NoHover = true
	gems.NoDraw = true
	gemsS.NoHover = true
	gemsS.NoDraw = true
	bombs.NoHover = true
	bombs.NoDraw = true
	bombsS.NoHover = true
	bombsS.NoDraw = true
	wrongFlags.NoHover = true
	wrongFlags.NoDraw = true
	wrongFlagsS.NoHover = true
	wrongFlagsS.NoDraw = true
	totalScore.NoHover = true
	totalScore.NoDraw = true
	totalScoreS.NoHover = true
	totalScoreS.NoDraw = true
	retry.SetClickFn(func() {
		PostMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		DescentState.start = true
		SwitchState(DescentStateKey)
	})
	backToMenu.SetClickFn(func() {
		PostMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		SwitchState(MenuStateKey)
	})
}
