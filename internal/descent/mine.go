package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	player2 "dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

const (
	MineFuse = 0.4
)

func CreateMine(pos pixel.Vec) {
	e := myecs.Manager.NewEntity()
	trans := transform.New()
	trans.Pos = pos
	fuse := timing.New(MineFuse)
	anim := reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("mine_1", img.Batchers[constants.EntityKey].Animations["mine_1"].S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("mine_2", img.Batchers[constants.EntityKey].Animations["mine_2"].S, reanimator.Tran).
			SetTrigger(1, func(_ *reanimator.Anim, _ string, _ int) {
				e.AddComponent(myecs.Func, data.NewFrameFunc(func() bool {
					player2.CaveBombsLeft--
					tile := Descent.GetCave().GetTile(trans.Pos)
					for _, n := range tile.SubCoords.Neighbors() {
						t := tile.Chunk.Get(n)
						t.Destroy(nil, false)
					}
					myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
						SourceID:       trans.ID,
						Center:         trans.Pos,
						Radius:         MineBaseRadius * world.TileSize,
						Amount:         1,
						Dazed:          3.,
						Knockback:      MineBaseKnockback,
						KnockbackDecay: true,
					})
					particles.CreateSmallExplosion(trans.Pos)
					sfx.SoundPlayer.PlaySound("blast1", 0.0)
					//camera.Cam.Shake(0.5, 10.)
					myecs.Manager.DisposeEntity(e)
					return false
				}))
			})).
		SetChooseFn(func() int {
			if MineFuse*0.5 > fuse.Elapsed() {
				return 0
			} else {
				return 1
			}
		}), "mine_1")
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Func, data.NewFrameFunc(func() bool {
			fuse.Update()
			return false
		})).
		AddComponent(myecs.Animation, anim).
		AddComponent(myecs.Drawable, anim).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
}

//type Mine struct {
//	Transform  *transform.Transform
//	Timer      *timing.FrameTimer
//	Tile       *cave.Tile
//	created    bool
//	explode    bool
//	Reanimator *reanimator.Tree
//	entity     *ecs.Entity
//}
//
//func (m *Mine) Update() {
//	if m.created {
//		if m.Timer.UpdateDone() || m.explode {
//			player.CaveBombsLeft--
//			player.CaveBlownUpBombs++
//			area := []pixel.Vec{m.Transform.Pos}
//			for _, n := range m.Tile.SubCoords.Neighbors() {
//				t := m.Tile.Chunk.Get(n)
//				t.Destroy(false)
//				area = append(area, t.Transform.Pos)
//			}
//			myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
//				SourceID:       m.Transform.ID,
//				Center:         m.Transform.Pos,
//				Radius:         MineBaseRadius * world.TileSize,
//				Amount:         1,
//				Dazed:          3.,
//				Knockback:      MineBaseKnockback,
//				KnockbackDecay: true,
//			})
//			particles.CreateSmallExplosion(m.Tile.Transform.Pos)
//			sfx.SoundPlayer.PlaySound("blast1", -1.0)
//			m.Delete()
//		}
//	}
//}
//
//func (m *Mine) Create(pos pixel.Vec) {
//	m.Transform = transform.New()
//	m.Transform.Pos = pos
//	m.created = true
//	m.Timer = timing.New(MineFuse)
//	m.Reanimator = reanimator.New(reanimator.NewSwitch().
//		AddAnimation(reanimator.NewAnimFromSprites("mine_1", img.Batchers[constants.EntityKey].Animations["mine_1"].S, reanimator.Hold)).
//		AddAnimation(reanimator.NewAnimFromSprites("mine_2", img.Batchers[constants.EntityKey].Animations["mine_2"].S, reanimator.Tran).
//			SetTrigger(2, func(_ *reanimator.Anim, _ string, _ int) {
//				m.explode = true
//			}),
//		).
//		SetChooseFn(func() int {
//			if MineFuse*0.5 > m.Timer.Elapsed() {
//				return 0
//			} else {
//				return 1
//			}
//		}), "mine_1")
//	m.entity = myecs.Manager.NewEntity().
//		AddComponent(myecs.Entity, m).
//		AddComponent(myecs.Transform, m.Transform).
//		AddComponent(myecs.Animation, m.Reanimator).
//		AddComponent(myecs.Drawable, m.Reanimator).
//		AddComponent(myecs.Batch, constants.EntityKey)
//}
//
//func (m *Mine) Delete() {
//	myecs.Manager.DisposeEntity(m.entity)
//}
