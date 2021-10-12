package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

const (
	MineBaseKnockback = 25.
	MineBaseRadius = 1.99
)

type Bomb struct {
	Transform  *transform.Transform
	Timer      *timing.FrameTimer
	FuseLength float64
	Tile       *cave.Tile
	created    bool
	explode    bool
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
}

func (b *Bomb) Update() {
	if b.created {
		if b.explode{
			CaveBombsLeft--
			CaveBlownUpBombs++
			area := []pixel.Vec{b.Transform.Pos}
			for _, n := range b.Tile.SubCoords.Neighbors(){
				t := b.Tile.Chunk.Get(n)
				t.Destroy(false)
				area = append(area, t.Transform.Pos)
			}
			myecs.Manager.NewEntity().
			AddComponent(myecs.AreaDmg, &data.AreaDamage{
				Center:         b.Transform.Pos,
				Radius:         MineBaseRadius * world.TileSize,
				Amount:         1,
				Dazed:          3.,
				Knockback:      MineBaseKnockback,
				KnockbackDecay: true,
				Source:         b.Transform.Pos,
				Override:       true,
			})
			vfx.CreateExplosion(b.Tile.Transform.Pos)
			sfx.SoundPlayer.PlaySound("blast1", 0.0)
			camera.Cam.Shake()
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
			reanimator.NewAnimFromSprites("bomb_fuse", img.Batchers[constants.EntityKey].Animations["bomb_fuse"].S, reanimator.Loop, nil),
			reanimator.NewAnimFromSprites("bomb_blow", img.Batchers[constants.EntityKey].Animations["bomb_blow"].S, reanimator.Tran, map[int]func() {
				2: func() {
					b.explode = true
				},
			}),
		),
		Check: func() int {
			if b.FuseLength - b.Timer.Elapsed() > 0.3 {
				return 0
			} else {
				return 1
			}
		},
	}, "bomb_unlit")
	b.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, b).
		AddComponent(myecs.Transform, b.Transform).
		AddComponent(myecs.Animation, b.Reanimator).
		AddComponent(myecs.Batch, constants.EntityKey)
}

func (b *Bomb) Delete() {
	myecs.Manager.DisposeEntity(b.entity)
}