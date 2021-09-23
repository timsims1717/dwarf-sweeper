package dungeon

import (
	"dwarf-sweeper/internal/cfg"
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
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
}

func (f *Flag) Update() {
	if f.created {
		if !f.Tile.Solid || f.Tile.destroyed || !f.Tile.marked {
			f.Delete()
			// todo: particles?
		}
	}
}

func (f *Flag) Create(from pixel.Vec) {
	f.Transform = transform.NewTransform()
	f.Transform.Pos = f.Tile.Transform.Pos
	f.created = true
	f.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("flag_hang", img.Batchers[cfg.EntityKey].Animations["flag_hang"].S, reanimator.Loop, nil),
		),
		Check: func() int {
			return 0
		},
	}, "flag_hang")
	f.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, f).
		AddComponent(myecs.Transform, f.Transform).
		AddComponent(myecs.Animation, f.Reanimator).
		AddComponent(myecs.Batch, cfg.EntityKey)
}

func (f *Flag) Delete() {
	f.Tile.marked = false
	myecs.Manager.DisposeEntity(f.entity)
}