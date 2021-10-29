package state

import (
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
		sfx.MusicPlayer.PauseMusic("pause", true)
		sfx.MusicPlayer.UnpauseOrNext("game")
	})
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
		SwitchState(1)
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
		SwitchState(0)
	})
}

func ClearEnchantMenu() {
	EnchantMenu.RemoveItem("option1")
	EnchantMenu.RemoveItem("option2")
	EnchantMenu.RemoveItem("option3")
}

func FillEnchantMenu() bool {
	ClearEnchantMenu()
	choices := enchants.PickEnchantments()
	if len(choices) == 0 {
		return false
	}
	e1 := choices[0]
	option1 := EnchantMenu.InsertItem("option1", e1.Title, 1)
	option1.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		enchants.AddEnchantment(e1)
		EnchantMenu.CloseInstant()
		ClearEnchantMenu()
		SwitchState(0)
	})
	option1.Hint = e1.Desc
	if len(choices) > 1 {
		e2 := choices[1]
		option2 := EnchantMenu.InsertItem("option2", e2.Title, 2)
		option2.SetClickFn(func() {
			sfx.SoundPlayer.PlaySound("click", 2.0)
			enchants.AddEnchantment(e2)
			EnchantMenu.CloseInstant()
			ClearEnchantMenu()
			SwitchState(0)
		})
		option2.Hint = e2.Desc
	}
	if len(choices) > 2 {
		e3 := choices[2]
		option3 := EnchantMenu.InsertItem("option3", e3.Title, 3)
		option3.SetClickFn(func() {
			sfx.SoundPlayer.PlaySound("click", 2.0)
			enchants.AddEnchantment(e3)
			EnchantMenu.CloseInstant()
			ClearEnchantMenu()
			SwitchState(0)
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
	blocksDug.NoDraw = true
	blocksDugS.Right = true
	blocksDugS.NoHover = true
	blocksDugS.NoDraw = true
	gems.NoHover = true
	gems.NoDraw = true
	gemsS.Right = true
	gemsS.NoHover = true
	gemsS.NoDraw = true
	bombs.NoHover = true
	bombs.NoDraw = true
	bombsS.Right = true
	bombsS.NoHover = true
	bombsS.NoDraw = true
	wrongMarks.NoHover = true
	wrongMarks.NoDraw = true
	wrongMarksS.Right = true
	wrongMarksS.NoHover = true
	wrongMarksS.NoDraw = true
	totalScore.NoHover = true
	totalScore.NoDraw = true
	totalScoreS.Right = true
	totalScoreS.NoHover = true
	totalScoreS.NoDraw = true
	retry.SetClickFn(func() {
		PostMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		SwitchState(4)
	})
	backToMenu.Right = true
	backToMenu.SetClickFn(func() {
		PostMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		SwitchState(1)
	})
}