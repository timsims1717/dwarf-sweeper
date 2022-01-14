package descend

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/generate"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/player"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/sfx"
	"fmt"
	"github.com/faiface/pixel"
)

func Generate() {
	descent.New()
	for i := 0; i < descent.Descent.Depth; i++ {
		if i == descent.Descent.Depth - 1 {
			caveBuilders, err := builder.LoadBuilder(fmt.Sprint("assets/bosses.json"))
			if err != nil {
				panic(err)
			}
			choice := random.Effects.Intn(len(caveBuilders))
			descent.Descent.Builders = append(descent.Descent.Builders, []*builder.CaveBuilder{caveBuilders[choice]})
		} else if i % 2 == 0 {
			caveBuilders, err := builder.LoadBuilder(fmt.Sprint("assets/caves.json"))
			if err != nil {
				panic(err)
			}
			choice := random.Effects.Intn(len(caveBuilders))
			descent.Descent.Builders = append(descent.Descent.Builders, []*builder.CaveBuilder{caveBuilders[choice]})
		} else {
			caveBuilders, err := builder.LoadBuilder(fmt.Sprint("assets/puzzles.json"))
			if err != nil {
				panic(err)
			}
			choice := random.Effects.Intn(len(caveBuilders))
			descent.Descent.Builders = append(descent.Descent.Builders, []*builder.CaveBuilder{caveBuilders[choice]})
		}
	}
}

func Descend() {
	if descent.Descent.Start {
		if descent.Descent.Player != nil {
			descent.Descent.Player.Delete()
			descent.Descent.Player = nil
		}
		descent.Descent.SetPlayer(descent.NewDwarf(pixel.Vec{}))
		player.InitHUD()
		descent.Inventory = []*descent.InvItem{}
		descent.ResetStats()
		descent.Descent.Start = false
	} else {
		descent.ResetCaveStats()
		descent.Descent.Level++
	}
	descent.Descent.Builder = descent.Descent.Builders[descent.Descent.Level][0]
	descent.Descent.SetCave(generate.NewCave(descent.Descent.Builder, descent.Descent.Level * descent.Descent.Difficulty))
	if len(descent.Descent.Builder.Tracks) > 0 {
		sfx.MusicPlayer.ChooseNextTrack(constants.GameMusic, descent.Descent.Builder.Tracks)
	} else {
		sfx.MusicPlayer.Stop(constants.GameMusic)
	}
	sfx.MusicPlayer.Resume(constants.GameMusic)

	descent.Descent.Player.Transform.Pos = descent.Descent.GetCave().GetStart().Transform.Pos
	camera.Cam.SnapTo(descent.Descent.GetPlayer().Transform.Pos)
}