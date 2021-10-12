package sfx

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/pkg/errors"
)

var MusicPlayer *musicPlayer

type musicPlayer struct {
	tracks   map[string]string
	sets     map[string]*musicSet
	loading  bool
}

func init() {
	MusicPlayer = &musicPlayer{
		tracks: make(map[string]string),
		sets:   make(map[string]*musicSet),
	}
}

func (p *musicPlayer) Update() {
	for _, s := range p.sets {
		s.update()
	}
}

func (p *musicPlayer) RegisterMusicTrack(path, key string) {
	p.tracks[key] = path
}

func (p *musicPlayer) NewSet(key string, set []string, isRandom bool, vol, fade float64) {
	p.sets[key] = &musicSet{
		key:      key,
		isRandom: isRandom,
		fade:     fade,
		vol:      vol,
	}
	p.sets[key].setTracks(set)
}

func (p *musicPlayer) SetTracks(key string, set []string) {
	if s, ok := p.sets[key]; ok {
		s.setTracks(set)
	} else {
		p.sets[key] = &musicSet{
			key: key,
		}
		p.sets[key].setTracks(set)
	}
}

func (p *musicPlayer) UnpauseOrNext(set string) {
	if s, ok := p.sets[set]; ok {
		s.play()
	} else {
		fmt.Printf("music player error: no set '%s'\n", set)
	}
}

func (p *musicPlayer) PlayTrack(set, key string) {
	if s, ok := p.sets[set]; ok {
		s.playTrack(key)
	} else {
		fmt.Printf("music player error: no set '%s'\n", set)
	}
}

func (p *musicPlayer) PlayNext(set string) {
	if s, ok := p.sets[set]; ok {
		s.playNextTrack()
	} else {
		fmt.Printf("music player error: no set '%s'\n", set)
	}
}

func (p *musicPlayer) PauseMusic(set string, pause bool) {
	if s, ok := p.sets[set]; ok {
		s.pause(pause)
	} else {
		fmt.Printf("music player error: no set '%s'\n", set)
	}
}

func (p *musicPlayer) SetVolume(set string, vol float64) {
	if s, ok := p.sets[set]; ok {
		s.setVolume(vol)
	} else {
		fmt.Printf("music player error: no set '%s'\n", set)
	}
}

func (p *musicPlayer) SetFade(set string, fade float64) {
	if s, ok := p.sets[set]; ok {
		s.setFade(fade)
	} else {
		fmt.Printf("music player error: no set '%s'\n", set)
	}
}

func (p *musicPlayer) loadTrack(set *musicSet) {
	p.loading = true
	if err := p.loadTrackInner(set); err != nil {
		fmt.Printf("music player error: %s\n", err)
	} else {
		set.playNext = false
	}
	p.loading = false
}

func (p *musicPlayer) loadTrackInner(set *musicSet) error {
	errMsg := fmt.Sprintf("load track %s", set.next)
	if path, ok := p.tracks[set.next]; ok {
		streamer, format, err := loadSoundFile(path)
		if err != nil {
			return errors.Wrap(err, errMsg)
		}
		speaker.Lock()
		if set.stream != nil {
			err = set.stream.Close()
			if err != nil {
				fmt.Println(errors.Wrap(err, errMsg))
			}
		}
		if set.ctrl != nil {
			set.ctrl.Paused = true
		}
		if set.volume != nil {
			set.volume.Silent = true
		}
		set.stream = streamer
		set.ctrl = &beep.Ctrl{
			Streamer: set.stream,
			Paused:   false,
		}
		set.volume = &effects.Volume{
			Streamer: set.ctrl,
			Base:     2,
			Volume:   getMusicVolume(),
			Silent:   false,
		}
		set.paused = false
		set.interV = nil
		fmt.Printf("playing track %s\n", set.next)
		speaker.Unlock()
		speaker.Play(beep.Seq(
			beep.Resample(4, format.SampleRate, sampleRate, set.volume),
			beep.Callback(func() {
				set.playNextTrack()
			}),
		))
		return nil
	}
	return errors.Wrap(fmt.Errorf("key %s is not a registered track", set.next), errMsg)
}

func (p *musicPlayer) stopAllMusic() {
	speaker.Clear()
	for _, s := range p.sets {
		if s.stream != nil {
			s.stream.Close()
		}
		s.ctrl = nil
		s.volume = nil
		s.interV = nil
		s.paused = true
	}
}

func (p *musicPlayer) PauseAllMusic() {
	for _, s := range p.sets {
		s.pause(true)
	}
}