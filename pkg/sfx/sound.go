package sfx

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/pkg/errors"
	"time"
)

var (
	SoundPlayer *soundPlayer
	RepeatDelay = 0.25
)

type soundPlayer struct {
	sounds     map[string]*beep.Buffer
	currSounds map[string]time.Time
}

func init() {
	SoundPlayer = &soundPlayer{
		sounds: make(map[string]*beep.Buffer),
		currSounds: make(map[string]time.Time),
	}
}

func (p *soundPlayer) RegisterSound(path, key string) error {
	errMsg := "register sound"
	streamer, format, err := loadSoundFile(path)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	resampled := beep.Resample(4, format.SampleRate, sampleRate, streamer)

	buffer := beep.NewBuffer(format)
	buffer.Append(resampled)
	err = streamer.Close()
	if err != nil {
		fmt.Println(errors.Wrap(err, errMsg))
	}
	p.sounds[key] = buffer
	return nil
}

func (p *soundPlayer) PlaySound(key string, vol float64) {
	if sound, ok := p.sounds[key]; ok {
		if t, ok := p.currSounds[key]; !ok || time.Since(t).Seconds() > RepeatDelay {
			volume := &effects.Volume{
				Streamer: sound.Streamer(0, sound.Len()),
				Base:     2,
				Volume:   getSoundVolume() + vol,
				Silent:   false,
			}
			speaker.Play(volume)
			p.currSounds[key] = time.Now()
		}
	} else {
		fmt.Printf("WARNING: SoundPlayer key %s not registered\n", key)
	}
}

func (p *soundPlayer) PlayRandomSound(keys []string) {
	p.PlaySound(keys[random.Intn(len(keys))], 0.)
}
