package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
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
	HighYieldRadius   = 2.49
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
	trans := transform.New().WithID("bomb")
	trans.Pos = pos
	phys := physics.New()
	phys.GravityOff = true
	phys.FrictionOff = true
	phys.Bounciness = 0.4
	phys.Friction = 300.
	fuse := timing.New(BombFuse)
	fuseSFX := sfx.SoundPlayer.PlaySound("fuseshort", 0.0)
	anim := reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("bomb_fuse", img.Batchers[constants.EntityKey].Animations["bomb_fuse"].S, reanimator.Loop)).
		AddAnimation(reanimator.NewAnimFromSprites("bomb_blow", img.Batchers[constants.EntityKey].Animations["bomb_blow"].S, reanimator.Tran).
		SetTrigger(2, func() {
				sfx.SoundPlayer.PlaySound("doubleblast", 0.0)
			}).
		SetTrigger(3, func() {
			sfx.SoundPlayer.KillSound(fuseSFX)
			e.AddComponent(myecs.Func, data.NewFrameFunc(func() bool {
				Descent.GetCave().BombsLeft--
				tile := Descent.GetCave().GetTile(trans.Pos)
				for _, n := range tile.RCoords.Neighbors() {
					t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
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
				//camera.Cam.Shake(0.5, 10.)
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
	collider := data.NewCollider(pixel.R(0., 0., 16., 16.), data.GroundOnly)
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Collision, collider).
		AddComponent(myecs.Health, &data.SimpleHealth{
			Immune: BombImmunity,
		}).
		AddComponent(myecs.Func, data.NewFrameFunc(func() bool {
			fuse.Update()
			return false
		})).
		AddComponent(myecs.Animation, anim).
		AddComponent(myecs.Drawable, anim).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	//sfx.SoundPlayer.PlaySound("fuseshort", 0.0)
}

func CreateHighYieldBomb(pos pixel.Vec) {
	e := myecs.Manager.NewEntity()
	hy := myecs.Manager.NewEntity()
	trans := transform.New().WithID("hy-bomb")
	trans.Pos = pos
	phys := physics.New()
	phys.GravityOff = true
	phys.FrictionOff = true
	phys.Bounciness = 0.4
	phys.Friction = 300.
	fuse := timing.New(BombFuse)
	fuseSFX := sfx.SoundPlayer.PlaySound("fuseshort", 0.0)
	anim := reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprites("bomb_fuse", img.Batchers[constants.EntityKey].Animations["bomb_fuse"].S, reanimator.Loop)).
		AddAnimation(reanimator.NewAnimFromSprites("bomb_blow", img.Batchers[constants.EntityKey].Animations["bomb_blow"].S, reanimator.Tran).
		SetTrigger(2, func() {
				sfx.SoundPlayer.PlaySound("doubleblast", -1.0)
				hy.AddComponent(myecs.Drawable, img.Batchers[constants.EntityKey].GetSprite("bomb_high_yield_burst"))
			}).
		SetTrigger(3, func() {
			sfx.SoundPlayer.KillSound(fuseSFX)
			e.AddComponent(myecs.Func, data.NewFrameFunc(func() bool {
				Descent.GetCave().BombsLeft--
				tile := Descent.GetCave().GetTile(trans.Pos)
				for _, n := range tile.RCoords.Neighbors() {
					t := tile.Chunk.Cave.GetTileInt(n.X, n.Y)
					for i, sn := range t.RCoords.Neighbors() {
						if i % 2 == 0 {
							tile.Chunk.Cave.GetTileInt(sn.X, sn.Y).Destroy(nil, false)
						}
					}
					t.Destroy(nil, false)
				}
				myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
					SourceID:       trans.ID,
					Center:         trans.Pos,
					Radius:         HighYieldRadius * world.TileSize,
					Amount:         1,
					Dazed:          3.,
					Knockback:      MineBaseKnockback,
					KnockbackDecay: true,
				})
				particles.CreateBigExplosion(trans.Pos)
				//camera.Cam.Shake(0.5, 10.)
				myecs.Manager.DisposeEntity(e)
				myecs.Manager.DisposeEntity(hy)
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
	collider := data.NewCollider(pixel.R(0., 0., 16., 16.), data.GroundOnly)
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Collision, collider).
		AddComponent(myecs.Health, &data.SimpleHealth{
			Immune: BombImmunity,
		}).
		AddComponent(myecs.Func, data.NewFrameFunc(func() bool {
			fuse.Update()
			return false
		})).
		AddComponent(myecs.Animation, anim).
		AddComponent(myecs.Drawable, anim).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	hy.AddComponent(myecs.Transform, transform.New()).
		AddComponent(myecs.Parent, trans).
		AddComponent(myecs.Drawable, img.Batchers[constants.EntityKey].GetSprite("bomb_high_yield")).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
}
