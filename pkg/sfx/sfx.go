package sfx

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"github.com/pkg/errors"
	"os"
	"strings"
	"time"
)

const sampleRate = beep.SampleRate(44100 * 4)

// Volumes are stored as integers from 0 to 100.
var (
	masterVolume = 100
	masterMuted  = false
	musicVolume  = 100
	musicMuted   = false
	soundVolume  = 100
	soundMuted   = false
	sfxVolume    = map[string]int{}
	sfxMuted     = map[string]bool{}
)

func init() {
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/100))
	if err != nil {
		panic(err)
	}
}

func getMasterVolume() float64 {
	if masterMuted {
		return -5.
	} else {
		return float64(masterVolume)/20. - 5.
	}
}

func getMusicVolume() float64 {
	if musicMuted || masterMuted {
		return -5.
	} else {
		return float64(musicVolume*masterVolume)/2000. - 5.
	}
}

func getSoundVolume() float64 {
	if soundMuted || masterMuted {
		return -5.
	} else {
		return float64(soundVolume*masterVolume)/2000. - 5.
	}
}

func getSfxVolume(key string) float64 {
	if sfxMuted[key] || masterMuted {
		return -5.
	} else {
		return float64(sfxVolume[key]*masterVolume)/2000. - 5.
	}
}

func GetMasterVolume() int {
	if masterMuted {
		return 0
	} else {
		return masterVolume
	}
}

func GetMusicVolume() int {
	if musicMuted {
		return 0
	} else {
		return musicVolume
	}
}

func GetSoundVolume() int {
	if soundMuted {
		return 0
	} else {
		return soundVolume
	}
}

func GetSfxVolume(key string) int {
	if sfxMuted[key] {
		return 0
	} else {
		return sfxVolume[key]
	}
}

func SetMasterVolume(v int) {
	if v == 0 {
		masterMuted = true
	} else {
		masterMuted = false
	}
	masterVolume = v
}

func SetMusicVolume(v int) {
	if v == 0 {
		musicMuted = true
	} else {
		musicMuted = false
	}
	musicVolume = v
}

func SetSoundVolume(v int) {
	if v == 0 {
		soundMuted = true
	} else {
		soundMuted = false
	}
	soundVolume = v
}

func SetSfxVolume(v int, key string) {
	if v == 0 {
		sfxMuted[key] = true
	} else {
		sfxMuted[key] = false
	}
	sfxVolume[key] = v
}

func loadSoundFile(path string) (beep.StreamSeekCloser, beep.Format, error) {
	errMsg := "load sound file"
	file, err := os.Open(path)
	if err != nil {
		return nil, beep.Format{}, errors.Wrap(err, errMsg)
	}
	p := strings.ToLower(path)
	if strings.Contains(p, "mp3") {
		return mp3.Decode(file)
	} else if strings.Contains(p, "wav") {
		return wav.Decode(file)
	} else if strings.Contains(p, "ogg") {
		return vorbis.Decode(file)
	} else if strings.Contains(p, "flac") {
		return vorbis.Decode(file)
	}
	return nil, beep.Format{}, errors.Wrap(fmt.Errorf("could not determine file type of %s", path), errMsg)
}
