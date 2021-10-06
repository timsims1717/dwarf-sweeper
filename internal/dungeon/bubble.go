package dungeon

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

const (
	BubbleVel = 45.
	BubbleAcc = 5.
)

type Bubble struct {
	Physics    *physics.Physics
	Transform  *transform.Transform
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
	popped     bool
	created    bool
	Timer      *timing.FrameTimer
}

func (b *Bubble) Update() {
	if b.created && b.Timer.UpdateDone() && !b.popped {
		b.Pop()
	}
}

func (b *Bubble) Create(_ pixel.Vec) {
	b.Timer = timing.New(12.)
	b.Transform = transform.NewTransform()
	b.Physics = physics.New()
	b.Physics.Gravity = 50.
	b.Physics.Friction = 200.
	Dungeon.Player.Entity.AddComponent(myecs.Physics, b.Physics)
	Dungeon.Player.Bubble = b
	b.created = true
	b.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("bubble_vibe", img.Batchers[constants.BigEntityKey].Animations["bubble_vibe"].S, reanimator.Loop, nil),
			reanimator.NewAnimFromSprites("bubble_pop", img.Batchers[constants.BigEntityKey].Animations["bubble_pop"].S, reanimator.Tran, map[int]func() {
				2: func() {
					b.Delete()
				},
			}),
		),
		Check: func() int {
			if !b.popped {
				return 0
			} else {
				return 1
			}
		},
	}, "bubble_vibe")
	b.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, b).
		AddComponent(myecs.Transform, b.Transform).
		AddComponent(myecs.Parent, Dungeon.Player.Transform).
		AddComponent(myecs.Animation, b.Reanimator).
		AddComponent(myecs.Batch, constants.BigEntityKey)
}

func (b *Bubble) Delete() {
	if !b.popped {
		b.Pop()
	}
	myecs.Manager.DisposeEntity(b.entity)
}

func (b *Bubble) Pop() {
	b.popped = true
	b.entity.RemoveComponent(myecs.Parent)
	b.Transform.Pos = Dungeon.Player.Transform.Pos
	Dungeon.Player.Physics.Velocity = b.Physics.Velocity
	Dungeon.Player.Entity.AddComponent(myecs.Physics, Dungeon.Player.Physics)
	Dungeon.Player.Bubble = nil
}