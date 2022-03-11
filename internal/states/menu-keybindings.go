package states

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/typeface"
	"fmt"
	"strings"
)

func InitKeybindingMenu() {
	KeybindingMenu = menus.New("keybinding", camera.Cam)
	KeybindingMenu.HideArrow = true
	KeybindingMenu.SetCloseFn(func() {
		UpdateKeybindings(data.CurrInput)
		RegisterPlayerSymbols(data.CurrInput.Key, data.CurrInput)
	})
	keybindingA := KeybindingMenu.AddItem("line_a", "Set key/button ", false)
	keybindingA.NoHover = true
	keybindingB := KeybindingMenu.AddItem("line_b", "", false)
	keybindingB.NoHover = true
}

func OpenKeybindingMenu(name, key string) {
	KeybindingMenu.ItemMap["line_b"].SetText(fmt.Sprintf("for %s", name))
	KeyString = key
	OpenMenu(KeybindingMenu)
	sfx.SoundPlayer.PlaySound("click", 2.0)
}

func UpdateKeybindings(in *input.Input) {
	UpdateKeybinding("left", in)
	UpdateKeybinding("right", in)
	UpdateKeybinding("up", in)
	UpdateKeybinding("down", in)
	UpdateKeybinding("jump", in)
	UpdateKeybinding("dig", in)
	UpdateKeybinding("flag", in)
	UpdateKeybinding("interact", in)
	UpdateKeybinding("use", in)
	UpdateKeybinding("prev", in)
	UpdateKeybinding("next", in)
	UpdateKeybinding("mine_puzz_bomb", in)
	UpdateKeybinding("mine_puzz_safe", in)
}

func UpdateKeybinding(key string, in *input.Input) {
	r := InputMenu.ItemMap[fmt.Sprintf("%s_r", key)]
	bs := in.Buttons[key]
	builder := strings.Builder{}
	first := true
	if in.Mode != input.Gamepad {
		for _, k := range bs.Keys {
			if first {
				first = false
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString(fmt.Sprintf("{symbol:%s}", k.String()))
		}
		if bs.Scroll > 0 {
			if first {
				first = false
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString("{symbol:MouseScrollUp}")
		} else if bs.Scroll < 0 {
			if first {
				first = false
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString("{symbol:MouseScrollDown}")
		}
	}
	if in.Mode != input.KeyboardMouse {
		for _, b := range bs.Buttons {
			if first {
				first = false
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString(fmt.Sprintf("{symbol:%s}", input.GamepadString(b)))
		}
		if bs.AxisV != 0 {
			if first {
				first = false
			} else {
				builder.WriteString(" ")
			}
			builder.WriteString(fmt.Sprintf("{symbol:%s}", input.AxisDirString(bs.Axis, bs.AxisV > 0)))
		}
	}
	r.SetText(builder.String())
}

func RegisterPlayerSymbols(pCode string, in *input.Input) {
	RegisterPlayerSymbol("left", pCode, in)
	RegisterPlayerSymbol("right", pCode, in)
	RegisterPlayerSymbol("up", pCode, in)
	RegisterPlayerSymbol("down", pCode, in)
	RegisterPlayerSymbol("jump", pCode, in)
	RegisterPlayerSymbol("dig", pCode, in)
	RegisterPlayerSymbol("flag", pCode, in)
	RegisterPlayerSymbol("interact", pCode, in)
	RegisterPlayerSymbol("use", pCode, in)
	RegisterPlayerSymbol("prev", pCode, in)
	RegisterPlayerSymbol("next", pCode, in)
	RegisterPlayerSymbol("mine_puzz_bomb", pCode, in)
	RegisterPlayerSymbol("mine_puzz_safe", pCode, in)
}

func RegisterPlayerSymbol(key, pCode string, in *input.Input) {
	fullKey := fmt.Sprintf("%s-%s", pCode, key)
	bs := in.Buttons[key]
	if in.Mode != input.Gamepad {
		for _, k := range bs.Keys {
			typeface.RegisterSymbol(fullKey, img.Batchers[constants.MenuSprites].GetSprite(k.String()), 1.)
			return
		}
		if bs.Scroll > 0 {
			typeface.RegisterSymbol(fullKey, img.Batchers[constants.MenuSprites].GetSprite("MouseScrollUp"), 1.)
			return
		} else if bs.Scroll < 0 {
			typeface.RegisterSymbol(fullKey, img.Batchers[constants.MenuSprites].GetSprite("MouseScrollDown"), 1.)
			return
		}
	}
	if in.Mode != input.KeyboardMouse {
		for _, b := range bs.Buttons {
			typeface.RegisterSymbol(fullKey, img.Batchers[constants.MenuSprites].GetSprite(input.GamepadString(b)), 1.)
			return
		}
		if bs.AxisV != 0 {
			typeface.RegisterSymbol(fullKey, img.Batchers[constants.MenuSprites].GetSprite(input.AxisDirString(bs.Axis, bs.AxisV > 0)), 1.)
			return
		}
	}
}