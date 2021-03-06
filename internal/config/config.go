package config

import (
	"dwarf-sweeper/pkg/input"
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
	Left         *input.ButtonSet `toml:"left"`
	Right        *input.ButtonSet `toml:"right"`
	Up           *input.ButtonSet `toml:"up"`
	Down         *input.ButtonSet `toml:"down"`
	Jump         *input.ButtonSet `toml:"jump"`
	Dig          *input.ButtonSet `toml:"dig"`
	Flag         *input.ButtonSet `toml:"flag"`
	Use          *input.ButtonSet `toml:"use"`
	Interact     *input.ButtonSet `toml:"interact"`
	Prev         *input.ButtonSet `toml:"prev"`
	Next         *input.ButtonSet `toml:"next"`
	PuzzLeave    *input.ButtonSet `toml:"puzz_leave"`
	PuzzHelp     *input.ButtonSet `toml:"puzz_help"`
	MinePuzzBomb *input.ButtonSet `toml:"mine_puzz_bomb"`
	MinePuzzSafe *input.ButtonSet `toml:"mine_puzz_safe"`
}
