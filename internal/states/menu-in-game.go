package states

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/profile"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/util"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"strings"
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
	quests := PauseMenu.AddItem("quests", "Quests", false)
	options := PauseMenu.AddItem("options", "Options", false)
	mainMenu := PauseMenu.AddItem("main_menu", "Abandon Run", false)
	quit := PauseMenu.AddItem("quit", "Quit Game", false)

	pauseTitle.NoHover = true
	resume.SetClickFn(func() {
		PauseMenu.Close()
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	quests.SetClickFn(func() {
		OpenMenu(QuestMenu)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	quests.Ignore = true
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
	PauseMenu.SetOpenFn(func() {
		quests.Ignore = true
		for _, key := range profile.CurrentProfile.Quests {
			if util.ContainsStr(key, profile.CurrentProfile.QuestsComplete) || util.ContainsStr(key, profile.CurrentProfile.QuestsShown) {
				quests.Ignore = false
			}
		}
	})
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
		DescentState.CurrBiome = "mine"
		descent.New()
		SwitchState(DescentStateKey)
	})
	backToMenu.SetClickFn(func() {
		PostMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		SwitchState(MenuStateKey)
	})
}

func InitQuestMenu() {
	QuestMenu = menus.New("quest", camera.Cam)
	QuestMenu.Title = true
	questTitle := QuestMenu.AddItem("title", "Quests", false)
	completed := QuestMenu.AddItem("completed", "Completed Quests:", false)
	back := QuestMenu.AddItem("back", "Back", false)

	questTitle.NoHover = true
	completed.NoHover = true
	completed.Disabled = true
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		QuestMenu.Close()
	})

	updateQuests := func() {
		for i := len(QuestMenu.Items)-1; i >= 0; i-- {
			item := QuestMenu.Items[i]
			if strings.Contains(item.Key, "quest") {
				QuestMenu.RemoveItem(item.Key)
			}
		}
		for _, key := range profile.CurrentProfile.Quests {
			if !util.ContainsStr(key, profile.CurrentProfile.QuestsComplete) && util.ContainsStr(key, profile.CurrentProfile.QuestsShown) {
				q := data.Quests[key]
				qi := QuestMenu.InsertItem(fmt.Sprintf("quest_%s", key), fmt.Sprintf(" %s", q.Name), "title", false)
				qi.Hint = q.Desc
			}
		}
		for _, key := range profile.CurrentProfile.Quests {
			if util.ContainsStr(key, profile.CurrentProfile.QuestsComplete) {
				q := data.Quests[key]
				qci := QuestMenu.InsertItem(fmt.Sprintf("quest_c_%s", key), fmt.Sprintf(" %s", q.Name), "completed", false)
				qci.Hint = q.Desc
				completed.Ignore = false
			}
		}
		QuestMenu.Hovered = 1
	}

	QuestMenu.SetOpenFn(func() {
		updateQuests()
	})
}