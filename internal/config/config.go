package config

import (
	pxginput "github.com/timsims1717/pixel-go-input"
)

type config struct {
	Audio    audio    `toml:"audio"`
	Gameplay gameplay `toml:"gameplay"`
	Graphics graphics `toml:"graphics"`
	InputP1  inputs   `toml:"inputP1"`
	InputP2  inputs   `toml:"inputP2"`
	InputP3  inputs   `toml:"inputP3"`
	InputP4  inputs   `toml:"inputP4"`
}

type audio struct {
	SoundVolume int  `toml:"sound_volume"`
	MusicVolume int  `toml:"music_volume"`
	MuteUnfocus bool `toml:"mute_on_unfocus"`
}

type graphics struct {
	VSync bool `toml:"vsync"`
	FullS bool `toml:"fullscreen"`
	ResIn int  `toml:"resolution"`
}

type gameplay struct {
	ShowTimer    bool `toml:"show_timer"`
	ScreenShake  bool `toml:"screen_shake"`
	SplitScreenV bool `toml:"split_screen_v"`
}

type inputs struct {
	Key          string           `toml:"name"`
	Gamepad      int              `toml:"gamepad"`
	AimDedicated bool             `toml:"aim_mode"`
	DigOnRelease bool             `toml:"dig_on"`
	Deadzone     float64          `toml:"deadzone"`
	LeftStick    bool             `toml:"left_stick"`
	Left         *pxginput.ButtonSet `toml:"left"`
	Right        *pxginput.ButtonSet `toml:"right"`
	Up           *pxginput.ButtonSet `toml:"up"`
	Down         *pxginput.ButtonSet `toml:"down"`
	Jump         *pxginput.ButtonSet `toml:"jump"`
	Dig          *pxginput.ButtonSet `toml:"dig"`
	Flag         *pxginput.ButtonSet `toml:"flag"`
	Use          *pxginput.ButtonSet `toml:"use"`
	Interact     *pxginput.ButtonSet `toml:"interact"`
	Prev         *pxginput.ButtonSet `toml:"prev"`
	Next         *pxginput.ButtonSet `toml:"next"`
	PuzzLeave    *pxginput.ButtonSet `toml:"puzz_leave"`
	PuzzHelp     *pxginput.ButtonSet `toml:"puzz_help"`
	MinePuzzBomb *pxginput.ButtonSet `toml:"mine_puzz_bomb"`
	MinePuzzSafe *pxginput.ButtonSet `toml:"mine_puzz_safe"`
}
