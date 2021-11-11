package state

import (
	"dwarf-sweeper/internal/credits"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/sfx"
	"github.com/faiface/pixel/pixelgl"
)

func InitMainMenu(win *pixelgl.Window) {
	MainMenu = menus.New("main", camera.Cam)
	MainMenu.Title = true
	start := MainMenu.AddItem("start", "Start Game")
	options := MainMenu.AddItem("options", "Options")
	credit := MainMenu.AddItem("credits", "Credits")
	quit := MainMenu.AddItem("quit", "Quit")

	start.SetClickFn(func() {
		OpenMenu(StartMenu)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	options.SetClickFn(func() {
		OpenMenu(OptionsMenu)
		sfx.SoundPlayer.PlaySound("click", 2.0)
	})
	credit.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		credits.Open()
	})
	quit.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		win.SetClosed(true)
	})
	quit.Hint = "You're going to leave?"
}

func InitStartMenu() {
	StartMenu = menus.New("start", camera.Cam)
	StartMenu.Title = true
	normal := StartMenu.AddItem("normal", "Normal Descent")
	infinite := StartMenu.AddItem("infinite", "Infinite Cave")
	back := StartMenu.AddItem("back", "Back")

	normal.SetClickFn(func() {
		StartMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		descent.Descent.Type = cave.Normal
		SwitchState(4)
	})
	normal.Hint = "Start a new run through a variety of caves!"
	infinite.SetClickFn(func() {
		StartMenu.CloseInstant()
		sfx.SoundPlayer.PlaySound("click", 2.0)
		descent.Descent.Type = cave.Infinite
		SwitchState(4)
	})
	infinite.Hint = "Survive in a cave that never ends!"
	back.SetClickFn(func() {
		sfx.SoundPlayer.PlaySound("click", 2.0)
		StartMenu.Close()
	})
}