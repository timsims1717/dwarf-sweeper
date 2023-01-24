package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"math"
)

type GnomeBossState int

const (
	GBWaiting = iota
	GBSearching
	GBEmerge
	GBRoar
	GBCharge
	GBIdle
	GBDig
	GBFlee
)

const (
	GBSpeed = 115.
	GBAcc   = 10.
)

var gbImmunity = map[data.DamageType]data.Immunity{
	data.Enemy: {
		KB:    true,
		DMG:   true,
		Dazed: true,
	},
	data.Shovel: {
		KB:  true,
		DMG: true,
	},
}

type GnomeBoss struct {
	Transform  *transform.Transform
	Physics    *physics.Physics
	Reanimator *reanimator.Tree
	Entity     *ecs.Entity
	Health     *data.Health
	Collider   *data.Collider

	State    GnomeBossState
	Charge   bool
	faceLeft bool

	emergeTimer *timing.Timer
	sfxTimer    *timing.Timer
	currLayer   int
	onDamageFn  func()
	onFleeFn    func()
	fleeing     bool
	tries       int
}

func (gnome *GnomeBoss) Update() bool {
	gnome.sfxTimer.Update()
	if gnome.Health.Dazed {
		gnome.Collider.ThroughWalls = false
		gnome.Physics.GravityOff = false
		gnome.Health.Immune = data.FullImmunity
		gnome.Collider.Damage = nil
	} else {
		debug.AddText(fmt.Sprintf("GnomeBoss Pos: (%d,%d)", int(gnome.Transform.Pos.X), int(gnome.Transform.Pos.Y)))
		if gnome.State == GBCharge || gnome.State == GBEmerge || gnome.State == GBRoar || gnome.State == GBIdle {
			gnome.Health.Immune = gbImmunity
		} else {
			gnome.Health.Immune = data.FullImmunity
		}
		gnome.Collider.Damage = nil
		if gnome.State != GBDig && gnome.Health.Curr < gnome.currLayer && gnome.onDamageFn != nil {
			gnome.onDamageFn()
		}
		if gnome.State != GBWaiting {
			gnome.Physics.GravityOff = false
			gnome.Collider.ThroughWalls = true
		}
		if gnome.State != GBCharge && gnome.State != GBSearching && gnome.State != GBWaiting {
			p := Descent.GetClosestPlayer(gnome.Transform.Pos)
			if gnome.Transform.Pos.X > p.Transform.Pos.X {
				gnome.faceLeft = true
			} else if gnome.Transform.Pos.X < p.Transform.Pos.X {
				gnome.faceLeft = false
			}
		}
		switch gnome.State {
		case GBSearching:
			if gnome.emergeTimer.UpdateDone() {
				if res, found := gnome.EmergeCoords(); found {
					gnome.Transform.Pos = res
					gnome.State = GBEmerge
				}
			}
		case GBCharge:
			if gnome.faceLeft {
				gnome.Physics.SetVelX(-GBSpeed, GBAcc)
			} else {
				gnome.Physics.SetVelX(GBSpeed, GBAcc)
			}
			half := world.TileSize * 0.501
			tul := Descent.Cave.GetTile(pixel.V(gnome.Transform.Pos.X-half, gnome.Transform.Pos.Y+half))
			tur := Descent.Cave.GetTile(pixel.V(gnome.Transform.Pos.X+half, gnome.Transform.Pos.Y+half))
			tdl := Descent.Cave.GetTile(pixel.V(gnome.Transform.Pos.X-half, gnome.Transform.Pos.Y-half))
			tdr := Descent.Cave.GetTile(pixel.V(gnome.Transform.Pos.X+half, gnome.Transform.Pos.Y-half))
			if tul.Solid() && tur.Solid() && tdl.Solid() && tdr.Solid() {
				gnome.State = GBSearching
				gnome.emergeTimer = timing.New(3.0)
			}
			gnome.Collider.Damage = &data.Damage{
				SourceID:  gnome.Transform.ID,
				Amount:    1,
				Dazed:     3.,
				Knockback: 8.,
				Angle:     &mmAngle,
				Source:    gnome.Transform.Pos,
				Type:      data.Enemy,
			}
		case GBDig:
			currT := Descent.Cave.GetTile(gnome.Transform.Pos)
			if math.Abs(currT.Transform.Pos.X-gnome.Transform.Pos.X) < 0.01 {
				gnome.Transform.Pos.X -= 0.5
			}
			edge := world.TileSize
			next := world.TileSize * 1.25
			tdl := Descent.Cave.GetTile(pixel.V(gnome.Transform.Pos.X-edge, gnome.Transform.Pos.Y-next))
			td := Descent.Cave.GetTile(pixel.V(gnome.Transform.Pos.X, gnome.Transform.Pos.Y-next))
			tdr := Descent.Cave.GetTile(pixel.V(gnome.Transform.Pos.X+edge, gnome.Transform.Pos.Y-next))
			if tdl.RCoords.Y == Descent.CoordsMap["current_layer"].Y ||
				td.RCoords.Y == Descent.CoordsMap["current_layer"].Y ||
				tdr.RCoords.Y == Descent.CoordsMap["current_layer"].Y {
				gnome.State = GBSearching
				gnome.currLayer = gnome.Health.Curr
			} else if gnome.Physics.Grounded || gnome.Physics.NearGround {
				if tdl.Solid() {
					tdl.DestroySpecial(true, true, true)
				}
				if td.Solid() {
					td.DestroySpecial(true, true, true)
				}
				if tdr.Solid() {
					tdr.DestroySpecial(true, true, true)
				}
			}
		}
		gnome.Transform.Flip = gnome.faceLeft
	}
	return false
}

func CreateGnomeBoss(maxHP int) *GnomeBoss {
	e := myecs.Manager.NewEntity()
	trans := transform.New().WithID("gnome-boss")
	trans.KeepLoaded = true
	trans.Load = true
	phys := physics.New()
	phys.GravityOff = true
	phys.Bounciness = 0.
	hp := &data.Health{
		Max:          maxHP,
		Curr:         maxHP,
		TempInvTimer: timing.New(0.5),
		TempInvSec:   5.,
		DazedTime:    1.5,
		Immune:       data.FullImmunity,
	}
	coll := data.NewCollider(pixel.R(0., 0., 32., 32.), data.Critter)
	coll.Debug = true
	coll.ThroughWalls = true
	gb := &GnomeBoss{
		Transform: trans,
		Physics:   phys,
		Entity:    e,
		Health:    hp,
		Collider:  coll,
		currLayer: maxHP,
		sfxTimer:  timing.New(0.25),
	}
	emergePartFn := func() {
		exit := gb.Transform.Pos
		exit.Y -= 16.
		particles.BiomeParticles(exit, Descent.Cave.Biome, 10, 16, 12., 0., math.Pi*0.5, 0.5, 130., 15., 0.75, 0.1, true)
	}
	runFXFn := func() {
		// check if inside any blocks
		half := world.TileSize * 0.501
		tul := Descent.Cave.GetTile(pixel.V(gb.Transform.Pos.X-half, gb.Transform.Pos.Y+half))
		tur := Descent.Cave.GetTile(pixel.V(gb.Transform.Pos.X+half, gb.Transform.Pos.Y+half))
		tdl := Descent.Cave.GetTile(pixel.V(gb.Transform.Pos.X-half, gb.Transform.Pos.Y-half))
		tdr := Descent.Cave.GetTile(pixel.V(gb.Transform.Pos.X+half, gb.Transform.Pos.Y-half))
		if tul.Solid() || tur.Solid() || tdl.Solid() || tdr.Solid() {
			if gb.sfxTimer.Done() {
				PlayRocks(-1.0)
				gb.sfxTimer = timing.New(0.25)
			}
			if tul.Solid() {
				if !tdl.Solid() {
					orig := tul.Transform.Pos
					orig.Y -= half
					particles.BiomeParticles(orig, Descent.Cave.Biome, 5, 7, half, 0., math.Pi*-0.5, 0.5, 80., 10., 0.75, 0.1, true)
				}
				if gb.faceLeft && !tur.Solid() {
					orig := tul.Transform.Pos
					orig.X += half
					particles.BiomeParticles(orig, Descent.Cave.Biome, 8, 9, 0., half, 0., 0.5, 80., 10., 0.75, 0.1, true)
				}
			}
			if tur.Solid() {
				if !tdr.Solid() {
					orig := tur.Transform.Pos
					orig.Y -= half
					particles.BiomeParticles(orig, Descent.Cave.Biome, 5, 7, half, 0., math.Pi*-0.5, 0.5, 80., 10., 0.75, 0.1, true)
				}
				if !gb.faceLeft && !tul.Solid() {
					orig := tur.Transform.Pos
					orig.X -= half
					particles.BiomeParticles(orig, Descent.Cave.Biome, 8, 9, 0., half, math.Pi, 0.5, 80., 10., 0.75, 0.1, true)
				}
			}
			if tdl.Solid() {
				if !tul.Solid() {
					orig := tdl.Transform.Pos
					orig.Y += half
					particles.BiomeParticles(orig, Descent.Cave.Biome, 5, 7, half, 0., math.Pi*0.5, 0.5, 80., 10., 0.75, 0.1, true)
				}
				if gb.faceLeft && !tdr.Solid() {
					orig := tdl.Transform.Pos
					orig.X += half
					particles.BiomeParticles(orig, Descent.Cave.Biome, 8, 9, 0., half, 0., 0.5, 80., 10., 0.75, 0.1, true)
				}
			}
			if tdr.Solid() {
				if !tur.Solid() {
					orig := tdr.Transform.Pos
					orig.Y += half
					particles.BiomeParticles(orig, Descent.Cave.Biome, 5, 7, half, 0., math.Pi*0.5, 0.5, 80., 10., 0.75, 0.1, true)
				}
				if !gb.faceLeft && !tdl.Solid() {
					orig := tdr.Transform.Pos
					orig.X -= half
					particles.BiomeParticles(orig, Descent.Cave.Biome, 8, 9, 0., half, math.Pi, 0.5, 80., 10., 0.75, 0.1, true)
				}
			}
		}
		// if the sfx timer is done, play a crunchy sound and reset it
		// for each face of the block that faces up or back (the opposite of the run)
		// blast a set of particles out of it
	}
	fleeFxFn := func() {
		if gb.sfxTimer.Done() {
			PlayRocks(-1.0)
			gb.sfxTimer = timing.New(0.25)
		}
		particles.BiomeParticles(gb.Transform.Pos, Descent.Cave.Biome, 3, 5, 8., 8., math.Pi*0.5, 0.5, 80., 10., 0.75, 0.1, true)
	}
	gb.Reanimator = reanimator.New(reanimator.NewSwitch().
		AddNull().
		AddAnimation(reanimator.NewAnimFromSprites("gnome_emerge", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_emerge").S, reanimator.Tran).
			SetTrigger(0, func() {
				emergePartFn()
				sfx.SoundPlayer.PlaySound("emerge", 0.)
			}).
			SetTrigger(1, func() {
				emergePartFn()
			}).
			SetTrigger(2, func() {
				emergePartFn()
			}).
			SetTrigger(3, func() {
				emergePartFn()
			}).
			SetTrigger(6, func() {
				gb.State = GBIdle
				gb.Entity.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
					gb.State = GBRoar
					return true
				}, 0.3))
			})).
		AddAnimation(reanimator.NewAnimFromSprites("gnome_roar", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_roar").S, reanimator.Hold).
			SetTrigger(0, func() {
				sfx.SoundPlayer.PlaySound("roar", 0.)
			}).
			SetTrigger(2, func() {
				for _, d := range Descent.GetPlayers() {
					ShakeCam(d, 2.5, random.Effects.Float64()*4.+8.)
				}
				gb.Entity.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
					if gb.Charge {
						gb.State = GBCharge
					} else {
						gb.State = GBIdle
					}
					return true
				}, 1.5))
			})).
		AddAnimation(reanimator.NewAnimFromSprites("gnome_run", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_run").S, reanimator.Loop).
			SetTrigger(0, func() {
				runFXFn()
			}).
			SetTrigger(1, func() {
				runFXFn()
			}).
			SetTrigger(2, func() {
				runFXFn()
			}).
			SetTrigger(3, func() {
				runFXFn()
				sfx.SoundPlayer.PlaySound("gnomestep", random.Effects.Float64()-1.)
			})).
		AddAnimation(reanimator.NewAnimFromSprites("gnome_idle", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_idle").S, reanimator.Loop)).
		AddAnimation(reanimator.NewAnimFromSprites("gnome_hurt", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_hurt").S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("gnome_dig", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_dig").S, reanimator.Loop)).
		AddAnimation(reanimator.NewAnimFromSprites("gnome_flee", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_dig").S, reanimator.Loop).
			SetTrigger(0, func() {
				fleeFxFn()
				if !gb.fleeing {
					myecs.AddEffect(gb.Entity, data.NewFadeBlack(colornames.White, 2.0))
					gb.Entity.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
						if gb.onFleeFn != nil {
							gb.onFleeFn()
						}
						myecs.Manager.DisposeEntity(gb.Entity)
						return true
					}, 2.0))
					gb.fleeing = true
				}
			}).
			SetTrigger(1, func() {
				fleeFxFn()
			}).
			SetTrigger(2, func() {
				fleeFxFn()
			}).
			SetTrigger(3, func() {
				fleeFxFn()
			}).
			SetTrigger(4, func() {
				fleeFxFn()
			})).
		SetChooseFn(func() int {
			if gb.State == GBWaiting || gb.State == GBSearching {
				return 0
			} else if gb.Health.Dazed {
				return 5
			} else if gb.State == GBEmerge {
				return 1
			} else if gb.State == GBRoar {
				return 2
			} else if gb.State == GBCharge {
				return 3
			} else if gb.State == GBDig {
				return 6
			} else if gb.State == GBFlee {
				gb.Health.Immune = data.FullImmunity
				return 7
			} else {
				return 4
			}
		}), "")
	e.AddComponent(myecs.Animation, gb.Reanimator).
		AddComponent(myecs.Drawable, gb.Reanimator).
		AddComponent(myecs.Transform, gb.Transform).
		AddComponent(myecs.Physics, gb.Physics).
		AddComponent(myecs.Health, gb.Health).
		AddComponent(myecs.Collision, gb.Collider).
		AddComponent(myecs.Update, data.NewFrameFunc(gb.Update)).
		AddComponent(myecs.Batch, constants.BigEntityKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	return gb
}

func (gnome *GnomeBoss) EmergeCoords() (pixel.Vec, bool) {
	tries := 0
	for tries < 5 {
		dt := Descent.GetRandomPlayerTile()
		outline := Descent.Cave.GetOutline(dt.RCoords, 7.5)
		inline := Descent.Cave.GetOutline(dt.RCoords, 5.25)
		candidates := world.NotIn(outline, inline)
		if len(candidates) > 0 {
			i := random.Effects.Intn(len(candidates))
			tc := candidates[i]
			if tc.Y >= Descent.CoordsMap["current_layer"].Y {
				tc1x := tc.X
				if random.Effects.Intn(2) == 0 {
					tc1x++
				} else {
					tc1x--
				}
				bt := Descent.Cave.GetTileInt(tc.X, tc.Y)
				bt1 := Descent.Cave.GetTileInt(tc1x, tc.Y)
				nt := Descent.Cave.GetTileInt(tc.X, tc.Y-1)
				nt1 := Descent.Cave.GetTileInt(tc1x, tc.Y-1)
				ut := Descent.Cave.GetTileInt(tc.X, tc.Y-2)
				ut1 := Descent.Cave.GetTileInt(tc1x, tc.Y-2)
				if !nt.Solid() && !nt1.Solid() && !ut.Solid() && !ut1.Solid() && bt.Solid() && bt1.Solid() {
					result := nt.Transform.Pos
					if i > 0 {
						result.X += world.TileSize * 0.5
					} else {
						result.X -= world.TileSize * 0.5
					}
					result.Y += world.TileSize * 0.5
					return result, true
				}
			}
		}
		tries++
	}
	gnome.tries++
	return pixel.ZV, false
}

func (gnome *GnomeBoss) SetOnDamageFn(fn func()) {
	gnome.onDamageFn = fn
}

func (gnome *GnomeBoss) SetOnFleeFn(fn func()) {
	gnome.onFleeFn = fn
}
