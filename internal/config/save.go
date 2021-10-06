package config

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/player"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/sfx"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/faiface/pixel/pixelgl"
	"os"
	"os/user"
	"runtime"
	"sync"
)

var (
	mu *sync.Mutex
)

func init() {
	mu = &sync.Mutex{}
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	constants.HomeDir = usr.HomeDir
	constants.ConfigDir = constants.HomeDir
	constants.System = runtime.GOOS
	switch constants.System {
	case "windows":
		fmt.Println("Windows")
		constants.ConfigDir += constants.WinDir
	case "darwin":
		fmt.Println("Mac")
		constants.ConfigDir += constants.MacDir
	case "linux":
		fmt.Println("Linux")
		constants.ConfigDir += constants.LinuxDir
	default:
		fmt.Printf("Unknown: %s.\n", constants.System)
		constants.ConfigDir += constants.LinuxDir
	}
	err = os.MkdirAll(constants.ConfigDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	constants.ConfigFile = constants.ConfigDir + "/config.toml"
	_, err = os.Open(constants.ConfigFile)
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(constants.ConfigFile)
		if err != nil {
			panic(err)
		}
		conf := config{
			Audio: audio{
				SoundVolume: 75,
				MusicVolume: 75,
			},
			Graphics: graphics{
				VSync: true,
				FullS: false,
				ResIn: 2,
			},
			Inputs: inputs{
				Gamepad:   -1,
				DigMode:   2,
				Deadzone:  0.25,
				LeftStick: true,
				Left:  input.New(pixelgl.KeyA, pixelgl.ButtonDpadLeft),
				Right: input.New(pixelgl.KeyD, pixelgl.ButtonDpadRight),
				Up:    input.New(pixelgl.KeyW, pixelgl.ButtonDpadUp),
				Down:  input.New(pixelgl.KeyS, pixelgl.ButtonDpadDown),
				Jump:  input.New(pixelgl.KeySpace, pixelgl.ButtonA),
				Dig: &input.ButtonSet{
					Keys:    []pixelgl.Button{pixelgl.MouseButtonLeft, pixelgl.KeyLeftShift},
					Buttons: []pixelgl.GamepadButton{pixelgl.ButtonX},
					Axis:    pixelgl.AxisRightTrigger,
					AxisV:   1,
				},
				Mark: &input.ButtonSet{
					Keys:    []pixelgl.Button{pixelgl.MouseButtonRight, pixelgl.KeyLeftControl},
					Buttons: []pixelgl.GamepadButton{pixelgl.ButtonY},
					Axis:    pixelgl.AxisLeftTrigger,
					AxisV:   1,
				},
				Use: input.New(pixelgl.KeyF, pixelgl.ButtonB),
				Prev: &input.ButtonSet{
					Keys:    []pixelgl.Button{pixelgl.KeyQ},
					Buttons: []pixelgl.GamepadButton{pixelgl.ButtonLeftBumper},
					Scroll:  -1,
				},
				Next: &input.ButtonSet{
					Keys:    []pixelgl.Button{pixelgl.KeyE},
					Buttons: []pixelgl.GamepadButton{pixelgl.ButtonRightBumper},
					Scroll:  1,
				},
			},
		}
		encode := toml.NewEncoder(file)
		err = encode.Encode(conf)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
}

func LoadAsync() {
	go LoadConfig()
}

func LoadConfig() {
	mu.Lock()
	var conf config
	if _, err := toml.DecodeFile(constants.ConfigFile, &conf); err != nil {
		fmt.Printf("couldn't decode configuration file: %s\n", err)
		return
	}
	sfx.SetSoundVolume(conf.Audio.SoundVolume)
	sfx.SetMusicVolume(conf.Audio.MusicVolume)
	constants.VSync = conf.Graphics.VSync
	constants.FullScreen = conf.Graphics.FullS
	constants.ResIndex = conf.Graphics.ResIn
	if conf.Inputs.Gamepad < 0 {
		player.GameInput.Mode = input.KeyboardMouse
	} else {
		player.GameInput.Mode = input.Gamepad
		player.GameInput.Joystick = pixelgl.Joystick(conf.Inputs.Gamepad)
	}
	constants.DigMode = conf.Inputs.DigMode
	player.GameInput.StickD = conf.Inputs.LeftStick
	input.Deadzone = conf.Inputs.Deadzone
	player.GameInput.Buttons["left"] = conf.Inputs.Left
	player.GameInput.Buttons["right"] = conf.Inputs.Right
	player.GameInput.Buttons["up"] = conf.Inputs.Up
	player.GameInput.Buttons["down"] = conf.Inputs.Down
	player.GameInput.Buttons["jump"] = conf.Inputs.Jump
	player.GameInput.Buttons["dig"] = conf.Inputs.Dig
	player.GameInput.Buttons["mark"] = conf.Inputs.Mark
	player.GameInput.Buttons["use"] = conf.Inputs.Use
	player.GameInput.Buttons["prev"] = conf.Inputs.Prev
	player.GameInput.Buttons["next"] = conf.Inputs.Next
	mu.Unlock()
}

func SaveAsync() {
	go SaveConfig()
}

func SaveConfig() {
	mu.Lock()
	file, err := os.OpenFile(constants.ConfigFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("couldn't open configuration file: %s\n", err)
		return
	}
	var conf config
	conf.Audio.SoundVolume = sfx.GetSoundVolume()
	conf.Audio.MusicVolume = sfx.GetMusicVolume()
	conf.Graphics.VSync = constants.VSync
	conf.Graphics.FullS = constants.FullScreen
	conf.Graphics.ResIn = constants.ResIndex

	if player.GameInput.Mode == input.KeyboardMouse {
		conf.Inputs.Gamepad = -1
	} else {
		conf.Inputs.Gamepad = int(player.GameInput.Joystick)
	}
	conf.Inputs.DigMode = constants.DigMode
	conf.Inputs.LeftStick = player.GameInput.StickD
	conf.Inputs.Deadzone = input.Deadzone
	conf.Inputs.Left = player.GameInput.Buttons["left"]
	conf.Inputs.Right = player.GameInput.Buttons["right"]
	conf.Inputs.Up = player.GameInput.Buttons["up"]
	conf.Inputs.Down = player.GameInput.Buttons["down"]
	conf.Inputs.Jump = player.GameInput.Buttons["jump"]
	conf.Inputs.Dig = player.GameInput.Buttons["dig"]
	conf.Inputs.Mark = player.GameInput.Buttons["mark"]
	conf.Inputs.Use = player.GameInput.Buttons["use"]
	conf.Inputs.Prev = player.GameInput.Buttons["prev"]
	conf.Inputs.Next = player.GameInput.Buttons["next"]
	encode := toml.NewEncoder(file)
	err = encode.Encode(conf)
	if err != nil {
		fmt.Printf("couldn't save configuration: %s\n", err)
	}
	mu.Unlock()
}