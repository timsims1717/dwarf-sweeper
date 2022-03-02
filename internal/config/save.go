package config

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/sfx"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
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
		InputP1:  DefaultInput,
		InputP2:  DefaultInput,
		InputP3:  DefaultInput,
		InputP4:  DefaultInput,
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
	loadInput(&conf)
	mu.Unlock()
}

func SaveAsync() {
	go SaveConfig()
}

func SaveConfig() {
	mu.Lock()
	os.Remove(constants.ConfigFile)
	file, err := os.Create(constants.ConfigFile)
	if err != nil {
		panic(err)
	}
	var conf config
	conf.Audio.SoundVolume = sfx.GetSoundVolume()
	conf.Audio.MusicVolume = sfx.GetMusicVolume()
	conf.Graphics.VSync = constants.VSync
	conf.Graphics.FullS = constants.FullScreen
	conf.Graphics.ResIn = constants.ResIndex

	saveInput(&conf)

	encode := toml.NewEncoder(file)
	err = encode.Encode(conf)
	if err != nil {
		fmt.Printf("couldn't save configuration: %s\n", err)
	}
	mu.Unlock()
}
