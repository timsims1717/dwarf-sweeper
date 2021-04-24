package sfx

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/pkg/errors"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"math/rand"
)

var MusicPlayer *musicPlayer

type musicPlayer struct {
	currSet  []string
	curr     string
	next     string
	tracks   map[string]string
	stream   beep.StreamSeekCloser
	ctrl     *beep.Ctrl
	volume   *effects.Volume
	interV   *gween.Tween
	format   beep.Format
	silent   bool
	loading  bool
}

func init() {
	MusicPlayer = &musicPlayer{
		tracks: make(map[string]string),
	}
}

func (p *musicPlayer) Update() {
	if !p.loading {
		if p.next != "" {
			if p.volume == nil || p.volume.Silent {
				p.loading = true
				go func() {
					if err := p.loadTrack(p.next); err != nil {
						fmt.Printf("music player error %s: %s\n", p.next, err)
					}
					p.loading = false
				}()
			}
		}
		if p.volume != nil {
			if p.interV != nil {
				v, fin := p.interV.Update(timing.DT)
				if fin {
					p.volume.Silent = true
					p.silent = true
					p.interV = nil
				} else {
					speaker.Lock()
					p.volume.Volume = v
					speaker.Unlock()
				}
			} else {
				speaker.Lock()
				p.volume.Silent = musicMuted || p.silent
				p.volume.Volume = getMusicVolume()
				speaker.Unlock()
			}
		}
	}
}

func (p *musicPlayer) RegisterMusicTrack(path, key string) {
	p.tracks[key] = path
}

func (p *musicPlayer) SetCurrentTracks(keys []string) {
	p.currSet = keys
}

func (p *musicPlayer) PlayTrack(key string, fadeOut float64) {
	p.next = key
	if p.volume != nil {
		p.interV = gween.New(p.volume.Volume, -8., fadeOut, ease.Linear)
	}
}

func (p *musicPlayer) PlayNextTrack(fadeOut float64, mustSwitch bool) {
	if len(p.currSet) > 0 && (mustSwitch || !util.ContainsStr(p.curr, p.currSet)) {
		p.PlayTrack(p.currSet[rand.Intn(len(p.currSet))], fadeOut)
	}
}

func (p *musicPlayer) FadeOut(fade float64) {
	if p.volume != nil {
		p.interV = gween.New(p.volume.Volume, -8., fade, ease.Linear)
	}
}

func (p *musicPlayer) loadTrack(key string) error {
	errMsg := "load track"
	if path, ok := p.tracks[key]; ok {
		streamer, format, err := loadSoundFile(path)
		if err != nil {
			return errors.Wrap(err, errMsg)
		}
		speaker.Lock()
		if p.stream != nil {
			err = p.stream.Close()
			if err != nil {
				fmt.Println(errors.Wrap(err, errMsg))
			}
		}
		if p.ctrl != nil {
			p.ctrl.Paused = true
		}
		if p.volume != nil {
			p.volume.Silent = true
		}
		p.stream = streamer
		p.ctrl = &beep.Ctrl{
			Streamer: p.stream,
			Paused:   false,
		}
		p.volume = &effects.Volume{
			Streamer: p.ctrl,
			Base:     2,
			Volume:   getMusicVolume(),
			Silent:   false,
		}
		p.silent = false
		p.format = format
		p.curr = p.next
		p.next = ""
		p.interV = nil
		speaker.Unlock()
		speaker.Play(beep.Seq(
			beep.Resample(4, format.SampleRate, sampleRate, p.volume),
			beep.Callback(func() {
				if len(p.currSet) > 0 {
					p.PlayTrack(p.currSet[rand.Intn(len(p.currSet)-1)], 0.)
				}
			}),
		))
		return nil
	}
	return errors.Wrap(fmt.Errorf("key %s is not a registered track", key), errMsg)
}

func (p *musicPlayer) stopMusic() {
	speaker.Clear()
	p.ctrl = nil
	p.volume = nil
}
