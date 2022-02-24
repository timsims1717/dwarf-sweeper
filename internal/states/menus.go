package states

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/pkg/input"
	"github.com/faiface/pixel/pixelgl"
)

var (
	MainMenu       *menus.DwarfMenu
	StartMenu      *menus.DwarfMenu
	AudioMenu      *menus.DwarfMenu
	GraphicsMenu   *menus.DwarfMenu
	InputMenu      *menus.DwarfMenu
	KeybindingMenu *menus.DwarfMenu
	PauseMenu      *menus.DwarfMenu
	OptionsMenu    *menus.DwarfMenu
	EnchantMenu    *menus.DwarfMenu
	PostMenu       *menus.DwarfMenu
	DebugMenu      *menus.DwarfMenu
	KeyString      string
	focused        bool
)

func InitializeMenus(win *pixelgl.Window) {
	InitMainMenu(win)
	InitStartMenu()
	InitOptionsMenu()
	// todo: accessibility
	InitAudioMenu()
	InitGraphicsMenu()
	InitInputMenu(win)
	InitKeybindingMenu()
	InitPauseMenu(win)
	InitEnchantMenu()
	InitPostGameMenu()
	InitDebugMenu()
	UpdateKeybindings()
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
					input.ClearInput(data.GameInputP1, KeyString)
					menuInput.Get("inputClear").Consume()
					me.Close()
				} else {
					if input.CheckAssign(win, data.GameInputP1, KeyString) {
						data.GameInputP1.Buttons[KeyString].Consume()
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
