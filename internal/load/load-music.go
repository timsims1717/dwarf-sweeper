package load

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/pkg/sfx"
)

func Music() {
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Crab Nebula.wav", "crab")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Honeybee.wav", "honey")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Prairie Oyster.wav", "prairie")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Sable.wav", "sable")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Strawberry Jam.wav", "strawberry")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/The Dawn Approaches.wav", "dawn")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/The Hero Approaches.wav", "hero")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Twin Turbo.wav", "turbo")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Voyage.wav", "voyage")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/Ascension II.wav", "ascension2")

	sfx.MusicPlayer.SetTracks("menu", []string{"crab"})
	sfx.MusicPlayer.SetTracks("pause", []string{"sable"})
	sfx.MusicPlayer.NewSet(constants.GameMusic, []string{"honey", "strawberry", "dawn", "hero", "voyage", "prairie"}, sfx.Repeat, 0., 2.)
}