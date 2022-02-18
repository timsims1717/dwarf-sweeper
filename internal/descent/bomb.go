package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/player"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
)

const (
	MineBaseKnockback = 25.
	MineBaseRadius    = 1.49
	BombFuse          = 1.5
)

var BombImmunity = map[data.DamageType]data.Immunity{
	data.Blast: {
		DMG:   true,
		Dazed: true,
	},
	data.Shovel: {
		DMG:   true,
		Dazed: true,
	},
	data.Enemy: {
		KB:    true,
		DMG:   true,
		Dazed: true,
	},
}

func CreateBomb(pos pixel.Vec) {
	e := myecs.Manager.NewEntity()
	trans := transform.New()
	trans.Pos = pos
	phys := physics.New()
	phys.GravityOff = true
	phys.FrictionOff = true
	phys.Bounciness = 0.4
	phys.Friction = 300.
	fuse := timing.New(BombFuse)
	anim := reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("bomb_fuse", img.Batchers[constants.EntityKey].Animations["bomb_fuse"].S, reanimator.Loop)).
		AddAnimation(reanimator.NewAnimFromSprites("bomb_blow", img.Batchers[constants.EntityKey].Animations["bomb_blow"].S, reanimator.Tran).
			SetTrigger(2, func(_ *reanimator.Anim, _ string, _ int) {
				e.AddComponent(myecs.Func, data.NewFrameFunc(func() bool {
					player.CaveBombsLeft--
					player.CaveBlownUpBombs++
					tile := Descent.GetCave().GetTile(trans.Pos)
					for _, n := range tile.SubCoords.Neighbors() {
						t := tile.Chunk.Get(n)
						t.Destroy(false)
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
					vfx.CreateExplosion(trans.Pos)
					sfx.SoundPlayer.PlaySound("blast1", 0.0)
					camera.Cam.Shake(0.5, 10.)
					myecs.Manager.DisposeEntity(e)
					return false
				}))
			})).
		SetChooseFn(func() int {
			if BombFuse-fuse.Elapsed() > 0.3 {
				return 0
			} else {
				return 1
			}
		}), "bomb_fuse")
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Collision, data.NewCollider(pixel.R(0., 0., 16., 16.), true, false)).
		AddComponent(myecs.Health, &data.SimpleHealth{
			Immune: BombImmunity,
		}).
		AddComponent(myecs.Func, data.NewFrameFunc(func() bool {
			fuse.Update()
			return false
		})).
		AddComponent(myecs.Animation, anim).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
}
