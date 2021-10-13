package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

type Flag struct {
	Transform  *transform.Transform
	Tile       *cave.Tile
	created    bool
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
	correct    bool
}

func (f *Flag) Update() {
	if f.created {
		if !f.Tile.Solid() || f.Tile.Destroyed || !f.Tile.Marked {
			f.Delete()
			// todo: particles?
		}
	}
}

func (f *Flag) Create(_ pixel.Vec) {
	f.Transform = transform.NewTransform()
	f.Transform.Pos = f.Tile.Transform.Pos
	f.created = true
	f.correct = f.Tile.Bomb
	if f.correct {
		CaveBombsMarked++
		CaveCorrectMarks++
	} else {
		CaveWrongMarks++
	}
	f.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("flag_hang", img.Batchers[constants.EntityKey].Animations["flag_hang"].S, reanimator.Loop, nil),
		),
		Check: func() int {
			return 0
		},
	}, "flag_hang")
	f.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, f).
		AddComponent(myecs.Transform, f.Transform).
		AddComponent(myecs.Animation, f.Reanimator).
		AddComponent(myecs.Batch, constants.EntityKey)
}

func (f *Flag) Delete() {
	f.Tile.Marked = false
	if f.Tile.Solid() {
		if f.correct {
			CaveBombsMarked--
			CaveCorrectMarks--
		} else {
			CaveWrongMarks--
		}
	} else if f.correct {
		CaveBombsMarked--
	} else {
		CaveWrongMarks--
	}
	myecs.Manager.DisposeEntity(f.entity)
}