package sfx

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"strings"
	"time"
)

const sampleRate = beep.SampleRate(44100)

// Volumes are stored as integers from 0 to 100.
var (
	random       *rand.Rand
	masterVolume = 100
	masterMuted  = false
	musicVolume  = 100
	musicMuted   = false
	soundVolume  = 100
	soundMuted   = false
	//sfxVolume    = map[string]int{}
	//sfxMuted     = map[string]bool{}
)

func init() {
	random = rand.New(rand.NewSource(time.Now().Unix()))
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/25))
	if err != nil {
		panic(err)
	}
}

/* INTERNAL GETTERS */

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

func getMasterMuted() bool {
	return masterMuted
}

func getMusicMuted() bool {
	return masterMuted || musicMuted
}

func getSoundMuted() bool {
	return masterMuted || soundMuted
}

//func getSfxVolume(key string) float64 {
//	if sfxMuted[key] || masterMuted {
//		return -5.
//	} else {
//		return float64(sfxVolume[key]*masterVolume)/2000. - 5.
//	}
//}

/* EXTERNAL GETTERS */

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

func GetMasterMuted() bool {
	return masterMuted
}

func GetMusicMuted() bool {
	return musicMuted
}

func GetSoundMuted() bool {
	return soundMuted
}

//func GetSfxVolume(key string) int {
//	if sfxMuted[key] {
//		return 0
//	} else {
//		return sfxVolume[key]
//	}
//}

/* SET MUTE */

func MuteMaster(muted bool) {
	speaker.Lock()
	if muted {
		masterMuted = true
		for _, set := range MusicPlayer.sets {
			if set.volume != nil {
				set.volume.Silent = true
			}
		}
		for _, vol := range SoundPlayer.volumes {
			vol.Silent = true
		}
	} else {
		masterMuted = false
		if !musicMuted {
			for _, set := range MusicPlayer.sets {
				if set.volume != nil {
					set.volume.Silent = false
				}
			}
		}
		if !soundMuted {
			for _, vol := range SoundPlayer.volumes {
				vol.Silent = false
			}
		}
	}
	speaker.Unlock()
}

func MuteMusic(muted bool) {
	speaker.Lock()
	if muted {
		musicMuted = true
		for _, set := range MusicPlayer.sets {
			if set.volume != nil {
				set.volume.Silent = true
			}
		}
	} else {
		musicMuted = false
		if !masterMuted {
			for _, set := range MusicPlayer.sets {
				if set.volume != nil {
					set.volume.Silent = false
				}
			}
		}
	}
	speaker.Unlock()
}

func MuteSound(muted bool) {
	speaker.Lock()
	if muted {
		soundMuted = true
		for _, vol := range SoundPlayer.volumes {
			vol.Silent = true
		}
	} else {
		soundMuted = false
		if !masterMuted {
			for _, vol := range SoundPlayer.volumes {
				vol.Silent = false
			}
		}
	}
	speaker.Unlock()
}

/* SET VOLUME */

func SetMasterVolume(v int) {
	speaker.Lock()
	if v == 0 {
		masterMuted = true
		for _, set := range MusicPlayer.sets {
			if set.volume != nil {
				set.volume.Silent = true
			}
		}
		for _, vol := range SoundPlayer.volumes {
			vol.Silent = true
		}
	} else {
		masterMuted = false
		if !musicMuted {
			for _, set := range MusicPlayer.sets {
				if set.volume != nil {
					set.volume.Silent = false
				}
			}
		}
		if !soundMuted {
			for _, vol := range SoundPlayer.volumes {
				vol.Silent = false
			}
		}
	}
	masterVolume = v
	speaker.Unlock()
}

func SetMusicVolume(v int) {
	speaker.Lock()
	if v == 0 {
		musicMuted = true
		for _, set := range MusicPlayer.sets {
			if set.volume != nil {
				set.volume.Silent = true
			}
		}
	} else {
		musicMuted = false
		if !masterMuted {
			for _, set := range MusicPlayer.sets {
				if set.volume != nil {
					set.volume.Silent = false
				}
			}
		}
	}
	musicVolume = v
	speaker.Unlock()
}

func SetSoundVolume(v int) {
	speaker.Lock()
	if v == 0 {
		soundMuted = true
		for _, vol := range SoundPlayer.volumes {
			vol.Silent = true
		}
	} else {
		soundMuted = false
		if !masterMuted {
			for _, vol := range SoundPlayer.volumes {
				vol.Silent = false
			}
		}
	}
	soundVolume = v
	speaker.Unlock()
}

//func SetSfxVolume(v int, key string) {
//	if v == 0 {
//		sfxMuted[key] = true
//	} else {
//		sfxMuted[key] = false
//	}
//	sfxVolume[key] = v
//}

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
