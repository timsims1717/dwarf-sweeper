package descent

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

var BubbleSec = 12.

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
	b.Timer = timing.New(BubbleSec)
	b.Transform = transform.NewTransform()
	b.Physics = physics.New()
	b.Physics.Gravity = 50.
	b.Physics.Friction = 200.
	Descent.Player.Entity.AddComponent(myecs.Physics, b.Physics)
	Descent.Player.Bubble = b
	b.created = true
	vibe := img.Batchers[constants.ParticleKey].GetAnimation("bubble_vibe")
	b.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("bubble_vibe", []*pixel.Sprite{
				vibe.S[0], vibe.S[0],
				vibe.S[1], vibe.S[1],
				vibe.S[2], vibe.S[2],
				vibe.S[1], vibe.S[1],
			}, reanimator.Loop, nil),
			reanimator.NewAnimFromSprites("bubble_pop", []*pixel.Sprite{
				img.Batchers[constants.ParticleKey].GetSprite("bubble_pop"),
				img.Batchers[constants.ParticleKey].GetSprite("bubble_pop"),
			}, reanimator.Tran, map[int]func() {
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
		AddComponent(myecs.Parent, Descent.Player.Transform).
		AddComponent(myecs.Animation, b.Reanimator).
		AddComponent(myecs.Batch, constants.ParticleKey)
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
	b.Transform.Pos = Descent.Player.Transform.Pos
	Descent.Player.Physics.Velocity = b.Physics.Velocity
	Descent.Player.Entity.AddComponent(myecs.Physics, Descent.Player.Physics)
	Descent.Player.Bubble = nil
}