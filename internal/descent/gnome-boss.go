package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"math"
)

type GnomeState int

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

	State    GnomeState
	Charge   bool
	faceLeft bool

	emergeTimer *timing.FrameTimer
	sfxTimer    *timing.FrameTimer
	currLayer   int
	onDamageFn  func()
	onFleeFn    func()
	fleeing     bool
}

func (gnome *GnomeBoss) Update() bool {
	if gnome.Health.Dead {
		gnome.Collider.ThroughWalls = false
		gnome.Physics.GravityOff = false
		gnome.State = GBFlee
		gnome.Health.Immune = data.FullImmunity
		gnome.Collider.Damage = nil
	} else if gnome.Health.Dazed {
		gnome.Collider.ThroughWalls = false
		gnome.Physics.GravityOff = false
		gnome.Health.Immune = data.FullImmunity
		gnome.Collider.Damage = nil
	} else {
		if gnome.State == GBCharge || gnome.State == GBEmerge || gnome.State == GBRoar || gnome.State == GBIdle {
			gnome.Health.Immune = gbImmunity
		} else {
			gnome.Health.Immune = data.FullImmunity
		}
		gnome.Collider.Damage = nil
		if gnome.State != GBDig && gnome.Health.Curr != gnome.currLayer && gnome.onDamageFn != nil {
			gnome.onDamageFn()
		}
		if gnome.State != GBWaiting {
			gnome.Physics.GravityOff = false
			gnome.Collider.ThroughWalls = true
		}
		if gnome.State != GBCharge && gnome.State != GBSearching && gnome.State != GBWaiting {
			ownCoords := Descent.GetCave().GetTile(gnome.Transform.Pos).RCoords
			playerCoords := Descent.GetPlayerTile().RCoords
			if ownCoords.X > playerCoords.X {
				gnome.faceLeft = true
			} else if ownCoords.X < playerCoords.X {
				gnome.faceLeft = false
			}
		}
		switch gnome.State {
		case GBSearching:
			if gnome.emergeTimer.UpdateDone() {
				gnome.Emerge(true)
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
	trans := transform.New()
	phys := physics.New()
	phys.GravityOff = true
	hp := &data.Health{
		Max:          maxHP,
		Curr:         maxHP,
		TempInvTimer: timing.New(0.5),
		TempInvSec:   5.,
		DazedTime:    1.5,
		Immune:       data.FullImmunity,
	}
	coll := data.NewCollider(pixel.R(0., 0., 32., 32.), false, true)
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
		gb.sfxTimer.Update()
		// check if inside any blocks
		half := world.TileSize * 0.501
		tul := Descent.Cave.GetTile(pixel.V(gb.Transform.Pos.X-half, gb.Transform.Pos.Y+half))
		tur := Descent.Cave.GetTile(pixel.V(gb.Transform.Pos.X+half, gb.Transform.Pos.Y+half))
		tdl := Descent.Cave.GetTile(pixel.V(gb.Transform.Pos.X-half, gb.Transform.Pos.Y-half))
		tdr := Descent.Cave.GetTile(pixel.V(gb.Transform.Pos.X+half, gb.Transform.Pos.Y-half))
		if tul.Solid() || tur.Solid() || tdl.Solid() || tdr.Solid() {
			if gb.sfxTimer.Done() {
				sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
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
					particles.BiomeParticles(orig, Descent.Cave.Biome, 5, 7, 0., half, 0., 0.5, 80., 10., 0.75, 0.1, true)
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
					particles.BiomeParticles(orig, Descent.Cave.Biome, 5, 7, 0., half, math.Pi, 0.5, 80., 10., 0.75, 0.1, true)
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
					particles.BiomeParticles(orig, Descent.Cave.Biome, 5, 7, 0., half, 0., 0.5, 80., 10., 0.75, 0.1, true)
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
					particles.BiomeParticles(orig, Descent.Cave.Biome, 5, 7, 0., half, math.Pi, 0.5, 80., 10., 0.75, 0.1, true)
				}
			}
		}
		// if the sfx timer is done, play a crunchy sound and reset it
		// for each face of the block that faces up or back (the opposite of the run)
		// blast a set of particles out of it
	}
	fleeFxFn := func() {
		if gb.sfxTimer.UpdateDone() {
			sfx.SoundPlayer.PlaySound(fmt.Sprintf("rocks%d", random.Effects.Intn(5)+1), -1.0)
			gb.sfxTimer = timing.New(0.25)
		}
		particles.BiomeParticles(gb.Transform.Pos, Descent.Cave.Biome, 3, 5, 8., 8., math.Pi*0.5, 0.5, 80., 10., 0.75, 0.1, true)
	}
	gb.Reanimator = reanimator.New(reanimator.NewSwitch().
		AddNull().
		AddAnimation(reanimator.NewAnimFromSprites("gnome_emerge", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_emerge").S, reanimator.Tran).
			SetTrigger(0, func(_ *reanimator.Anim, _ string, _ int) {
				emergePartFn()
				sfx.SoundPlayer.PlaySound("emerge", 0.)
			}).
			SetTrigger(1, func(_ *reanimator.Anim, _ string, _ int) {
				emergePartFn()
			}).
			SetTrigger(2, func(_ *reanimator.Anim, _ string, _ int) {
				emergePartFn()
			}).
			SetTrigger(3, func(_ *reanimator.Anim, _ string, _ int) {
				emergePartFn()
			}).
			SetTrigger(6, func(_ *reanimator.Anim, _ string, _ int) {
				gb.State = GBIdle
				gb.Entity.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
					gb.State = GBRoar
					return true
				}, 0.3))
			})).
		AddAnimation(reanimator.NewAnimFromSprites("gnome_roar", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_roar").S, reanimator.Hold).
			SetTrigger(0, func(_ *reanimator.Anim, _ string, _ int) {
   				sfx.SoundPlayer.PlaySound("roar", 0.)
			}).
			SetTrigger(2, func(_ *reanimator.Anim, _ string, _ int) {
				camera.Cam.ZoomShake(1.4, 30.)
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
			SetTrigger(0, func(_ *reanimator.Anim, _ string, _ int) {
				runFXFn()
			}).
			SetTrigger(1, func(_ *reanimator.Anim, _ string, _ int) {
				runFXFn()
			}).
			SetTrigger(2, func(_ *reanimator.Anim, _ string, _ int) {
				runFXFn()
			}).
			SetTrigger(3, func(_ *reanimator.Anim, _ string, _ int) {
				runFXFn()
			}).
			SetTrigger(4, func(_ *reanimator.Anim, _ string, _ int) {
				runFXFn()
				sfx.SoundPlayer.PlaySound("gnomestep", random.Effects.Float64()-1.)
			})).
		AddAnimation(reanimator.NewAnimFromSprites("gnome_idle", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_idle").S, reanimator.Loop)).
		AddAnimation(reanimator.NewAnimFromSprites("gnome_hurt", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_hurt").S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("gnome_dig", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_dig").S, reanimator.Loop)).
		AddAnimation(reanimator.NewAnimFromSprites("gnome_flee", img.Batchers[constants.BigEntityKey].GetAnimation("gnome_dig").S, reanimator.Loop).
			SetTrigger(0, func(_ *reanimator.Anim, _ string, _ int) {
				fleeFxFn()
				if !gb.fleeing {
					myecs.AddEffect(gb.Entity, data.NewFadeBlack(colornames.White, 3.))
					gb.Entity.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
						if gb.onFleeFn != nil {
							gb.onFleeFn()
						}
						gb.Entity.AddComponent(myecs.Temp, myecs.ClearFlag(true))
						return true
					}, 2.5))
					gb.fleeing = true
				}
			}).
			SetTrigger(1, func(_ *reanimator.Anim, _ string, _ int) {
				fleeFxFn()
			}).
			SetTrigger(2, func(_ *reanimator.Anim, _ string, _ int) {
				fleeFxFn()
			}).
			SetTrigger(3, func(_ *reanimator.Anim, _ string, _ int) {
				fleeFxFn()
			}).
			SetTrigger(4, func(_ *reanimator.Anim, _ string, _ int) {
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

func (gnome *GnomeBoss) Emerge(findCoords bool) {
	// find emerge coords
	found := true
	if findCoords {
		gnome.Transform.Pos, found = EmergeCoords()
	}
	if found {
		gnome.State = GBEmerge
	}
}

func EmergeCoords() (pixel.Vec, bool) {
	x := Descent.GetPlayerTile().RCoords.X
	y := Descent.CoordsMap["current_layer"].Y
	pCos := world.Coords{
		X: x,
		Y: y,
	}
	i := 5
	for i < 10 && i > -10 {
		next := pCos
		next.X += i
		next1x := next.X
		if i > 0 {
			next1x++
		} else {
			next1x--
		}
		nt := Descent.Cave.GetTileInt(next.X, next.Y)
		nt1 := Descent.Cave.GetTileInt(next1x, next.Y)
		bt := Descent.Cave.GetTileInt(next.X, next.Y+1)
		bt1 := Descent.Cave.GetTileInt(next1x, next.Y+1)
		ut := Descent.Cave.GetTileInt(next.X, next.Y-1)
		ut1 := Descent.Cave.GetTileInt(next1x, next.Y-1)
		if !nt.Solid() && !nt1.Solid() && !ut.Solid() && !ut1.Solid() && bt.Solid() && bt1.Solid() && random.CaveGen.Intn(10-util.Abs(i)) == 0 {
			result := nt.Transform.Pos
			if i > 0 {
				result.X += world.TileSize * 0.5
			} else {
				result.X -= world.TileSize * 0.5
			}
			result.Y += world.TileSize * 0.5
			return result, true
		}
		i *= -1
		if i > 0 {
			i++
		}
	}
	return pixel.Vec{}, false
}

func (gnome *GnomeBoss) SetOnDamageFn(fn func()) {
	gnome.onDamageFn = fn
}

func (gnome *GnomeBoss) SetOnFleeFn(fn func()) {
	gnome.onFleeFn = fn
}
