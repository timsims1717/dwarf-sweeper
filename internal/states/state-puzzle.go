package states

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/puzzles"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/state"
	"github.com/faiface/pixel/pixelgl"
)

type puzzleState struct {
	*state.AbstractState
	minePuzzle *puzzles.MinePuzzle
	puzLevel   int
}

func (s *puzzleState) Unload() {
	sfx.MusicPlayer.Stop("pause")
}

func (s *puzzleState) Load(done chan struct{}) {
	s.puzLevel++
	reanimator.SetFrameRate(10)
	reanimator.Reset()
	s.minePuzzle = &puzzles.MinePuzzle{}
	s.minePuzzle.Create(camera.Cam, s.puzLevel)
	s.minePuzzle.Open()
	sfx.MusicPlayer.PlayMusic("pause")
	done <- struct{}{}
}

func (s *puzzleState) Update(win *pixelgl.Window) {
	if s.minePuzzle != nil {
		reanimator.Update()
		s.minePuzzle.Update(data.GameInput)
		if s.minePuzzle.Solved() {
			s.minePuzzle.Close()
		}
		if s.minePuzzle.IsClosed() && s.minePuzzle.Solved() {
			s.puzLevel++
			s.minePuzzle.Create(camera.Cam, s.puzLevel)
			s.minePuzzle.Open()
		}
	}
}

func (s *puzzleState) Draw(win *pixelgl.Window) {
	if s.minePuzzle != nil {
		s.minePuzzle.Draw(win)
	}
}

func (s *puzzleState) SetAbstract(aState *state.AbstractState) {
	s.AbstractState = aState
}