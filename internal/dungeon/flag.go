package dungeon

import (
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/transform"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

type Flag struct {
	Transform  *transform.Transform
	Tile       *Tile
	created    bool
	done       bool
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
}

func (f *Flag) Update() {
	if f.created && !f.done {
		if !f.Tile.Solid || f.Tile.destroyed || !f.Tile.marked {
			f.done = true
			// todo: particles?
		}
	}
}

func (f *Flag) Draw(target pixel.Target) {
	if f.created && !f.done {
		f.Reanimator.CurrentSprite().Draw(target, f.Transform.Mat)
	}
}

func (f *Flag) Create(from pixel.Vec, batcher *img.Batcher) {
	f.Transform = transform.NewTransform()
	f.Transform.Pos = f.Tile.Transform.Pos
	f.created = true
	f.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("flag_hang", batcher.Animations["flag_hang"].S, reanimator.Loop, nil),
		),
		Check: func() int {
			return 0
		},
	}, "flag_hang")
	f.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, f.Transform).
		AddComponent(myecs.Animation, f.Reanimator)
}

func (f *Flag) Done() bool {
	return f.done
}

func (f *Flag) Delete() {
	myecs.Manager.DisposeEntity(f.entity)
}