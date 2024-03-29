package states

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/menus"
	"github.com/faiface/pixel/pixelgl"
	pxginput "github.com/timsims1717/pixel-go-input"
)

var (
	MainMenu       *menus.DwarfMenu
	StartMenu      *menus.DwarfMenu
	AudioMenu      *menus.DwarfMenu
	GameplayMenu   *menus.DwarfMenu
	GraphicsMenu   *menus.DwarfMenu
	InputMenu      *menus.DwarfMenu
	KeybindingMenu *menus.DwarfMenu
	PauseMenu      *menus.DwarfMenu
	QuestMenu      *menus.DwarfMenu
	OptionsMenu    *menus.DwarfMenu
	PostMenu       *menus.DwarfMenu
	DebugMenu      *menus.DwarfMenu
	KeyString      string
	NumPlayers     = 1
	BiomeIndex     = 0
	focused        bool
)

func InitializeMenus(win *pixelgl.Window) {
	InitMainMenu(win)
	InitStartMenu()
	InitOptionsMenu()
	// todo: accessibility
	InitAudioMenu()
	InitGameplayMenu()
	InitGraphicsMenu()
	InitInputMenu(win)
	InitKeybindingMenu()
	InitPauseMenu(win)
	InitQuestMenu()
	InitPostGameMenu()
	InitDebugMenu()
	UpdateKeybindings(data.CurrInput)
	RegisterPlayerSymbols(data.GameInputP1.Key, data.GameInputP1)
	RegisterPlayerSymbols(data.GameInputP2.Key, data.GameInputP2)
	RegisterPlayerSymbols(data.GameInputP3.Key, data.GameInputP3)
	RegisterPlayerSymbols(data.GameInputP4.Key, data.GameInputP4)
}

func UpdateMenus(win *pixelgl.Window) {
	if win.Focused() && !focused {
		focused = true
	} else if !win.Focused() && focused {
		focused = false
	}
	for i, me := range menuStack {
		if i == len(menuStack)-1 {
			if !win.Focused() {
				me.UnhoverAll()
			}
			me.Update(menuInput)
			if me.IsClosed() {
				if len(menuStack) > 1 {
					menuStack = menuStack[:len(menuStack)-1]
				} else {
					menuStack = []*menus.DwarfMenu{}
				}
			} else if me.Key == "keybinding" && me.IsOpen() {
				if menuInput.Get("inputClear").JustPressed() {
					pxginput.ClearInput(data.CurrInput, KeyString)
					menuInput.Get("inputClear").Consume()
					me.Close()
				} else {
					if pxginput.CheckAssign(win, data.CurrInput, KeyString) {
						data.CurrInput.Buttons[KeyString].Consume()
						me.Close()
					}
				}
			}
		} else {
			me.Update(nil)
		}
	}
}

func MenuClosed() bool {
	return len(menuStack) < 1
}

func OpenMenu(menu *menus.DwarfMenu) {
	menu.Open()
	menuStack = append(menuStack, menu)
}

func clearMenus() {
	menuStack = []*menus.DwarfMenu{}
}
