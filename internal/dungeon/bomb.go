package dungeon

import (
	"dwarf-sweeper/internal/character"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

type Bomb struct {
	EID        int
	Transform  *transform.Transform
	Timer      *timing.FrameTimer
	FuseLength float64
	Tile       *Tile
	created    bool
	first      bool
	explode    bool
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
}

func (b *Bomb) Update() {
	if b.created {
		if b.explode{
			area := []pixel.Vec{b.Transform.Pos}
			for _, n := range b.Tile.SubCoords.Neighbors(){
				t := b.Tile.Chunk.Get(n)
				t.Destroy()
				area = append(area, t.Transform.Pos)
			}
			myecs.Manager.NewEntity().
			AddComponent(myecs.AreaDmg, &character.AreaDamage{
				Area:           area,
				Amount:         1,
				Dazed:          3.,
				Knockback:      MineKnockback,
				KnockbackDecay: true,
				Source:         b.Transform.Pos,
				Override:       true,
			})
			vfx.CreateExplosion(b.Tile.Transform.Pos)
			b.Delete()
		} else {
			b.Timer.Update()
		}
	}
}

func (b *Bomb) Create(pos pixel.Vec) {
	b.Transform = transform.NewTransform()
	b.Transform.Pos = pos
	b.created = true
	b.Timer = timing.New(b.FuseLength)
	b.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewAnimFromSprites("bomb_unlit", img.Batchers[entityBKey].Animations["bomb_unlit"].S, reanimator.Tran, map[int]func() {
				1: func() {
					b.first = false
				},
			}),
			reanimator.NewAnimFromSprites("bomb_fuse", img.Batchers[entityBKey].Animations["bomb_fuse"].S, reanimator.Loop, nil),
			reanimator.NewAnimFromSprites("bomb_blow", img.Batchers[entityBKey].Animations["bomb_blow"].S, reanimator.Tran, map[int]func() {
				2: func() {
					b.explode = true
				},
			}),
		),
		Check: func() int {
			if b.first {
				return 0
			} else if b.FuseLength - b.Timer.Elapsed() > 0.3 {
				return 1
			} else {
				return 2
			}
		},
	}, "bomb_unlit")
	b.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, b).
		AddComponent(myecs.Transform, b.Transform).
		AddComponent(myecs.Animation, b.Reanimator).
		AddComponent(myecs.Batch, entityBKey)
	Dungeon.AddEntity(b)
}

func (b *Bomb) Delete() {
	myecs.Manager.DisposeEntity(b.entity)
	Dungeon.RemoveEntity(b.EID)
}

func (b *Bomb) SetId(i int) {
	b.EID = i
}