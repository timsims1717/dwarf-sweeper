package config

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
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
		CreateConfig()
	} else if err != nil {
		panic(err)
	}
}

func CreateConfig() {
	os.Remove(constants.ConfigFile)
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
			Gamepad:      -1,
			AimDedicated: true,
			DigOnRelease: true,
			Deadzone:     0.25,
			LeftStick:    true,
			Left:         input.New(pixelgl.KeyA, pixelgl.ButtonDpadLeft),
			Right:        input.New(pixelgl.KeyD, pixelgl.ButtonDpadRight),
			Up:           input.New(pixelgl.KeyW, pixelgl.ButtonDpadUp),
			Down:         input.New(pixelgl.KeyS, pixelgl.ButtonDpadDown),
			Jump:         input.New(pixelgl.KeySpace, pixelgl.ButtonA),
			Dig: &input.ButtonSet{
				Keys:    []pixelgl.Button{pixelgl.MouseButtonLeft},
				Buttons: []pixelgl.GamepadButton{pixelgl.ButtonX},
				Axis:    pixelgl.AxisRightTrigger,
				AxisV:   1,
			},
			Flag: &input.ButtonSet{
				Keys:    []pixelgl.Button{pixelgl.MouseButtonRight},
				Axis:    pixelgl.AxisLeftTrigger,
				AxisV:   1,
			},
			Use:      input.New(pixelgl.KeyE, pixelgl.ButtonB),
			Interact: input.New(pixelgl.KeyQ, pixelgl.ButtonY),
			Prev: &input.ButtonSet{
				Buttons: []pixelgl.GamepadButton{pixelgl.ButtonLeftBumper},
				Scroll:  -1,
			},
			Next: &input.ButtonSet{
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
}

func LoadAsync() {
	go LoadConfig()
}

func LoadConfig() {
	mu.Lock()
	var conf config
	if _, err := toml.DecodeFile(constants.ConfigFile, &conf); err != nil {
		fmt.Printf("couldn't decode configuration file: %s\n", err)
		CreateConfig()
		return
	}
	sfx.SetSoundVolume(conf.Audio.SoundVolume)
	sfx.SetMusicVolume(conf.Audio.MusicVolume)
	constants.VSync = conf.Graphics.VSync
	constants.FullScreen = conf.Graphics.FullS
	constants.ResIndex = conf.Graphics.ResIn
	if conf.Inputs.Gamepad < 0 {
		data.GameInput.Mode = input.KeyboardMouse
	} else {
		data.GameInput.Mode = input.Gamepad
		data.GameInput.Joystick = pixelgl.Joystick(conf.Inputs.Gamepad)
	}
	constants.AimDedicated = conf.Inputs.AimDedicated
	constants.DigOnRelease = conf.Inputs.DigOnRelease
	data.GameInput.StickD = conf.Inputs.LeftStick
	input.Deadzone = conf.Inputs.Deadzone
	data.GameInput.Buttons["left"] = conf.Inputs.Left
	data.GameInput.Buttons["right"] = conf.Inputs.Right
	data.GameInput.Buttons["up"] = conf.Inputs.Up
	data.GameInput.Buttons["down"] = conf.Inputs.Down
	data.GameInput.Buttons["jump"] = conf.Inputs.Jump
	data.GameInput.Buttons["dig"] = conf.Inputs.Dig
	data.GameInput.Buttons["flag"] = conf.Inputs.Flag
	data.GameInput.Buttons["use"] = conf.Inputs.Use
	data.GameInput.Buttons["interact"] = conf.Inputs.Interact
	data.GameInput.Buttons["prev"] = conf.Inputs.Prev
	data.GameInput.Buttons["next"] = conf.Inputs.Next
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

	if data.GameInput.Mode == input.KeyboardMouse {
		conf.Inputs.Gamepad = -1
	} else {
		conf.Inputs.Gamepad = int(data.GameInput.Joystick)
	}
	conf.Inputs.AimDedicated = constants.AimDedicated
	conf.Inputs.DigOnRelease = constants.DigOnRelease
	conf.Inputs.LeftStick = data.GameInput.StickD
	conf.Inputs.Deadzone = input.Deadzone
	conf.Inputs.Left = data.GameInput.Buttons["left"]
	conf.Inputs.Right = data.GameInput.Buttons["right"]
	conf.Inputs.Up = data.GameInput.Buttons["up"]
	conf.Inputs.Down = data.GameInput.Buttons["down"]
	conf.Inputs.Jump = data.GameInput.Buttons["jump"]
	conf.Inputs.Dig = data.GameInput.Buttons["dig"]
	conf.Inputs.Flag = data.GameInput.Buttons["flag"]
	conf.Inputs.Use = data.GameInput.Buttons["use"]
	conf.Inputs.Interact = data.GameInput.Buttons["interact"]
	conf.Inputs.Prev = data.GameInput.Buttons["prev"]
	conf.Inputs.Next = data.GameInput.Buttons["next"]
	encode := toml.NewEncoder(file)
	err = encode.Encode(conf)
	if err != nil {
		fmt.Printf("couldn't save configuration: %s\n", err)
	}
	mu.Unlock()
}