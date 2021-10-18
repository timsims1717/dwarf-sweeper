package config

import (
	"dwarf-sweeper/pkg/input"
)

type config struct {
	Audio    audio    `toml:"audio"`
	Graphics graphics `toml:"graphics"`
	Inputs   inputs   `toml:"inputs"`
}

type audio struct {
	SoundVolume int `toml:"sound_volume"`
	MusicVolume int `toml:"music_volume"`
}

type graphics struct {
	VSync bool `toml:"vsync"`
	FullS bool `toml:"fullscreen"`
	ResIn int  `toml:"resolution"`
}

type inputs struct {
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
	Mark         *input.ButtonSet `toml:"mark"`
	Use          *input.ButtonSet `toml:"use"`
	Interact     *input.ButtonSet `toml:"interact"`
	Prev         *input.ButtonSet `toml:"prev"`
	Next         *input.ButtonSet `toml:"next"`
}