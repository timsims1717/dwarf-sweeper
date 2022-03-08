package sfx

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

var (
	SoundPlayer *soundPlayer
	RepeatDelay = 0.05
)

type soundPlayer struct {
	volumes     map[uuid.UUID]*effects.Volume
	sounds      map[string]*beep.Buffer
	soundTimers map[string]time.Time
}

func init() {
	SoundPlayer = &soundPlayer{
		volumes:     make(map[uuid.UUID]*effects.Volume),
		sounds:      make(map[string]*beep.Buffer),
		soundTimers: make(map[string]time.Time),
	}
}

func (p *soundPlayer) RegisterSound(path, key string) error {
	errMsg := fmt.Sprintf("register sound %s", key)
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

func (p *soundPlayer) PlaySound(key string, vol float64) *uuid.UUID {
	if !soundMuted && !masterMuted {
		if sound, ok := p.sounds[key]; ok {
			if t, ok := p.soundTimers[key]; !ok || time.Since(t).Seconds() > RepeatDelay {
				volume := &effects.Volume{
					Streamer: sound.Streamer(0, sound.Len()),
					Base:     2,
					Volume:   getSoundVolume() + vol,
					Silent:   false,
				}
				id := uuid.New()
				p.volumes[id] = volume
				speaker.Play(beep.Seq(
					volume,
					beep.Callback(func() {
						delete(p.volumes, id)
					}),
				))
				//speaker.Play(volume)
				p.soundTimers[key] = time.Now()
				return &id
			}
		} else {
			fmt.Printf("WARNING: SoundPlayer key %s not registered\n", key)
		}
	}
	return nil
}

func (p *soundPlayer) PlayRandomSound(keys []string) {
	p.PlaySound(keys[random.Intn(len(keys))], 0.)
}

func (p *soundPlayer) KillSound(id *uuid.UUID) {
	if id != nil {
		if s, ok := p.volumes[*id]; ok {
			s.Silent = true
		}
	}
}

func (p *soundPlayer) KillAll() {
	for id, vol := range p.volumes {
		vol.Silent = true
		delete(p.volumes, id)
	}
}