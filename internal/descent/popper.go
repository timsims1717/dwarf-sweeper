package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/beefsack/go-astar"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"math"
)

const (
	popperSpeed = 20.
	popperAcc   = 15.
	seekDist    = 8.
	runDist     = 2.
	effectSec   = 0.5
	fireSec     = 4.
	fireVel     = 275.
)

type PopperAction int

const (
	Wait = iota
	Seek
	Pop
	Unpop
	Aim
	Fire
	Dazed
	Crawl
	Dead
)

type Popper struct {
	Transform   *transform.Transform
	Physics     *physics.Physics
	Reanimator  *reanimator.Tree
	Entity      *ecs.Entity
	created     bool
	Health      *data.Health
	faceLeft    bool
	angle       float64
	aimNorm     pixel.Vec
	action      PopperAction
	path        []*cave.Tile
	target      *cave.Tile
	popFrame    int
	rootPos     pixel.Vec
	poppedPos   pixel.Vec
	effectTimer *timing.Timer
	fireTimer   *timing.Timer
}

func (p *Popper) Update() {
	var relaventTile *cave.Tile
	if p.action == Wait || p.action == Seek {
		relaventTile = Descent.GetCave().GetTile(p.Transform.Pos)
	} else if p.action != Dazed {
		relaventTile = Descent.GetCave().GetTile(p.rootPos)
	}
	if relaventTile != nil && !relaventTile.Solid() {
		p.Health.Dazed = true
		p.Health.DazedTimer = timing.New(3.)
	}
	if !p.Health.Dazed && !p.Health.Dead {
		d := Descent.GetClosestPlayer(p.Transform.Pos)
		action := p.action
		var distance float64
		if p.action == Wait || p.action == Seek {
			distance = util.Magnitude(d.Transform.Pos.Sub(p.Transform.Pos))
		} else {
			distance = util.Magnitude(d.Transform.Pos.Sub(p.rootPos))
		}
		tarDist := 0.
		if p.target != nil {
			tarDist = util.Magnitude(d.Transform.Pos.Sub(p.target.Transform.Pos))
		}
		wait := distance > world.TileSize*seekDist
		run := distance < world.TileSize*runDist
		switch p.action {
		case Wait:
			p.Physics.SetVelX(0., popperAcc)
			p.Physics.SetVelY(0., popperAcc)
			if !wait {
				action = Seek
			}
		case Seek:
			if wait {
				action = Wait
			} else {
				if p.target == nil || len(p.path) == 0 || !p.target.Solid() || tarDist < world.TileSize*runDist || tarDist > world.TileSize*seekDist {
					tries := 0
					targetFound := false
					for tries < 5 && !targetFound {
						var target *cave.Tile
						pos := Descent.GetTile(p.Transform.Pos).RCoords
						pos.X += random.Effects.Intn(11) - 5
						pos.Y += random.Effects.Intn(11) - 5
						target = Descent.GetCave().GetTileInt(pos.X, pos.Y)
						// if the target is nil or solid, try again
						if target == nil || target.Solid() {
							tries++
							continue
						}
						// if the dwarf is unreachable or out of range, try again
						Descent.Cave.PathRule = cave.MakePathRule(target.RCoords, cave.EmptyTypes, false, false, true, true)
						_, dist, found := astar.Path(target, Descent.Cave.GetTile(d.Transform.Pos))
						if !found || dist > seekDist || dist < runDist {
							tries++
							continue
						}
						// if the tile isn't next to a solid tile that we can reach, try again
						found = false
						for _, n := range target.RCoords.Neighbors() {
							if n.X == target.RCoords.X || n.Y == target.RCoords.Y {
								nt := Descent.Cave.GetTileInt(n.X, n.Y)
								if nt.Breakable() {
									target = nt
									found = true
									break
								}
							}
						}
						if !found {
							tries++
							continue
						}

						Descent.Cave.PathRule = cave.MakePathRule(target.RCoords, cave.NonWallTypes, false, false, false, false)
						path, d, found := astar.Path(target, Descent.GetTile(p.Transform.Pos))
						if !found || d > 8 {
							tries++
							continue
						}
						// success!
						p.path = []*cave.Tile{}
						for _, ip := range path {
							i := ip.(*cave.Tile)
							p.path = append(p.path, i)
						}
						targetFound = true
						p.target = target
					}
					if !targetFound {
						action = Wait
						p.path = nil
						p.target = nil
					}
				}
				// follow path
				legal := true
				for _, t := range p.path {
					if !t.Solid() {
						legal = false
						break
					}
				}
				if !legal {
					p.path = nil
				}
				if len(p.path) > 0 {
					next := p.path[0]
					if util.Magnitude(next.Transform.Pos.Sub(p.Transform.Pos)) < 1.5 {
						if len(p.path) == 1 {
							p.Transform.Pos = next.Transform.Pos
							p.path = nil
							next = nil
						} else {
							p.path = p.path[1:]
							next = p.path[0]
						}
					}
					if next == nil {
						var empty *cave.Tile
						curr := Descent.Cave.GetTile(p.Transform.Pos).RCoords
						for _, n := range curr.Neighbors() {
							if n.Y == curr.Y || n.X == curr.X {
								t := Descent.Cave.GetTileInt(n.X, n.Y)
								Descent.Cave.PathRule = cave.MakePathRule(t.RCoords, cave.EmptyTypes, false, false, true, true)
								_, _, found := astar.Path(t, Descent.Cave.GetTile(d.Transform.Pos))
								if found && !t.Solid() {
									empty = t
									break
								}
							}
						}
						if empty != nil {
							action = Pop
							p.rootPos = p.target.Transform.Pos
							p.poppedPos = empty.Transform.Pos
							p.Entity.AddComponent(myecs.Collision, data.NewCollider(pixel.R(0., 0., 16., 16.), true, true))
							p.Health.Immune[data.Shovel] = data.Immunity{
								KB:    false,
								DMG:   false,
								Dazed: false,
							}
							p.Health.Immune[data.Blast] = data.Immunity{
								KB:    false,
								DMG:   false,
								Dazed: false,
							}
						}
						p.Physics.CancelMovement()
						p.path = nil
						p.target = nil
					} else {
						if p.effectTimer == nil || p.effectTimer.UpdateDone() {
							e := myecs.Manager.NewEntity()
							trans := transform.New()
							trans.Pos = p.Transform.Pos
							e.AddComponent(myecs.Transform, trans).
								AddComponent(myecs.Temp, timing.New(1.5)).
								AddComponent(myecs.Drawable, img.Batchers[constants.ParticleKey].GetSprite("dig_thru")).
								AddComponent(myecs.Batch, constants.ParticleKey)
							myecs.AddEffect(e, data.NewFadeOut(colornames.White, 1.5))
							p.effectTimer = timing.New(effectSec)
						}
						move := util.Normalize(next.Transform.Pos.Sub(p.Transform.Pos))
						p.Physics.SetVelX(move.X*popperSpeed, popperAcc)
						p.Physics.SetVelY(move.Y*popperSpeed, popperAcc)
						if debug.Debug {
							moveDir := p.Transform.Pos
							moveDir.X += move.X * 8.
							moveDir.Y += move.Y * 8.
							debug.AddLine(colornames.Green, imdraw.SharpEndShape, p.Transform.Pos, moveDir, 2.)
						}

						var n *cave.Tile
						col := colornames.Green
						for _, t := range p.path {
							if n != nil {
								debug.AddLine(col, imdraw.SharpEndShape, n.Transform.Pos, t.Transform.Pos, 2.)
								col.R += 25
							}
							n = t
						}
					}
				}
			}
		case Pop:
			// if too far away, go to unpop
			if wait || run {
				action = Unpop
			} else {
				p.Transform.Pos = p.poppedPos
				if p.rootPos.X > p.poppedPos.X {
					p.Transform.Rot = 0.5
				} else if p.rootPos.X < p.poppedPos.X {
					p.Transform.Rot = -0.5
				} else if p.rootPos.Y > p.poppedPos.Y {
					p.Transform.Rot = 1.
				} else {
					p.Transform.Rot = 0.
				}
			}
		case Unpop:
			// if close again, go to pop
			if !wait && !run {
				action = Pop
			} else {
				p.Transform.Pos = p.poppedPos
				if p.rootPos.X > p.poppedPos.X {
					p.Transform.Rot = 0.5
				} else if p.rootPos.X < p.poppedPos.X {
					p.Transform.Rot = -0.5
				} else if p.rootPos.Y > p.poppedPos.Y {
					p.Transform.Rot = 1.
				} else {
					p.Transform.Rot = 0.
				}
			}
		case Aim:
			if wait || run {
				action = Unpop
			} else {
				shoot := true
				pPos := d.Transform.Pos
				ray := pPos.Sub(p.poppedPos)
				ray.Y += math.Abs(ray.X) * 0.2
				norm := util.Normalize(ray)
				p.Transform.Pos = p.poppedPos
				if p.rootPos.X > p.poppedPos.X { // root to right
					p.Transform.Rot = 0.5
					p.angle = ray.Rotated(math.Pi * -0.5).Angle()
					p.Transform.Flip = pPos.Y > p.poppedPos.Y
					shoot = pPos.X < p.poppedPos.X+world.TileSize*0.25
				} else if p.rootPos.X < p.poppedPos.X { // root to left
					p.Transform.Rot = -0.5
					p.angle = ray.Rotated(math.Pi * 0.5).Angle()
					p.Transform.Flip = pPos.Y < p.poppedPos.Y
					shoot = pPos.X > p.poppedPos.X-world.TileSize*0.25
				} else if p.rootPos.Y > p.poppedPos.Y { // root above
					p.Transform.Rot = 1.
					p.angle = ray.Rotated(ray.Angle() * -2.).Angle()
					p.Transform.Flip = pPos.X < p.poppedPos.X
					shoot = pPos.Y < p.poppedPos.Y+world.TileSize*0.25
				} else { // root below
					p.Transform.Rot = 0.
					p.angle = ray.Angle()
					p.Transform.Flip = pPos.X > p.poppedPos.X
					shoot = pPos.Y > p.poppedPos.Y-world.TileSize*0.25
				}
				if debug.Debug {
					aimDir := p.poppedPos
					aimDir.X += norm.X * 8.
					aimDir.Y += norm.Y * 8.
					debug.AddLine(colornames.Orange, imdraw.SharpEndShape, p.poppedPos, aimDir, 2.)
				}
				p.aimNorm = norm
				if shoot && p.fireTimer.UpdateDone() {
					action = Fire
				}
			}
		case Fire:
			pPos := d.Transform.Pos
			ray := pPos.Sub(p.poppedPos)
			ray.Y += math.Abs(ray.X) * 0.2
			norm := util.Normalize(ray)
			p.Transform.Pos = p.poppedPos
			if p.rootPos.X > p.poppedPos.X { // root to right
				p.Transform.Rot = 0.5
				p.angle = ray.Rotated(math.Pi * -0.5).Angle()
				p.Transform.Flip = pPos.Y > p.poppedPos.Y
			} else if p.rootPos.X < p.poppedPos.X { // root to left
				p.Transform.Rot = -0.5
				p.angle = ray.Rotated(math.Pi * 0.5).Angle()
				p.Transform.Flip = pPos.Y < p.poppedPos.Y
			} else if p.rootPos.Y > p.poppedPos.Y { // root above
				p.Transform.Rot = 1.
				p.angle = ray.Rotated(ray.Angle() * -2.).Angle()
				p.Transform.Flip = pPos.X < p.poppedPos.X
			} else { // root below
				p.Transform.Rot = 0.
				p.angle = ray.Angle()
				p.Transform.Flip = pPos.X > p.poppedPos.X
			}
			if debug.Debug {
				aimDir := p.poppedPos
				aimDir.X += norm.X * 8.
				aimDir.Y += norm.Y * 8.
				debug.AddLine(colornames.Orange, imdraw.SharpEndShape, p.poppedPos, aimDir, 2.)
			}
			p.aimNorm = norm
		case Dazed:
			if p.Physics.Grounded {
				p.Physics.GravityOff = true
				pos := p.Transform.Pos
				pos.Y -= world.TileSize
				b := Descent.Cave.GetTile(pos)
				if b.Solid() && b.Type != cave.Wall {
					p.rootPos = b.Transform.Pos
					p.poppedPos = p.Transform.Pos
					action = Unpop
				} else {
					action = Crawl
				}
			}
		case Crawl:

		}
		p.action = action
	} else if p.Health.Dazed && !p.Health.Dead {
		p.action = Dazed
		p.Physics.GravityOff = false
		p.Entity.AddComponent(myecs.Collision, data.NewCollider(pixel.R(0., 0., 16., 16.), true, true))
		p.Health.Immune[data.Shovel] = data.Immunity{
			KB:    false,
			DMG:   false,
			Dazed: false,
		}
		p.Health.Immune[data.Blast] = data.Immunity{
			KB:    false,
			DMG:   false,
			Dazed: false,
		}
		p.Transform.Rot = 0.
	} else if p.Health.Dead {
		if p.action != Dead {
			p.action = Dead
			p.Entity.AddComponent(myecs.Temp, timing.New(4.))
			p.Entity.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
				myecs.AddEffect(p.Entity, data.NewBlink(2.))
				return true
			}, 2.))
		}
		p.Physics.GravityOff = false
		p.Transform.Rot = 0.
	}
	if p.action == Wait {
		debug.AddCircle(colornames.Gray, p.Transform.Pos, 4., 0.)
	} else if p.action == Seek {
		debug.AddCircle(colornames.Orange, p.Transform.Pos, 4., 0.)
	}
	if p.target != nil {
		debug.AddCircle(colornames.Yellow, p.target.Transform.Pos, 4., 0.)
	}
}

func (p *Popper) Create(pos pixel.Vec) {
	p.Transform = transform.New()
	p.Transform.Pos = pos
	p.Physics = physics.New()
	p.Physics.GravityOff = true
	p.Health = &data.Health{
		Max:          2,
		Curr:         2,
		TempInvTimer: timing.New(0.5),
		TempInvSec:   0.5,
		Immune: map[data.DamageType]data.Immunity{
			data.Enemy: {
				KB:    true,
				DMG:   true,
				Dazed: true,
			},
			data.Blast: {
				KB: true,
			},
			data.Shovel: {
				KB:    true,
				DMG:   true,
				Dazed: true,
			},
		},
	}
	p.created = true
	p.Reanimator = reanimator.New(reanimator.NewSwitch().
		AddNull().
		AddAnimation(reanimator.NewAnimFromSprites("popper_out", img.Batchers[constants.EntityKey].GetAnimation("popper_out").S, reanimator.Tran).
			SetTrigger(0, func(a *reanimator.Anim, pKey string, pFrame int) {
				if pKey == "popper_in" {
					a.Step = 4 - pFrame
				} else {
					exit := p.poppedPos
					var varX, varY, angle float64
					if p.poppedPos.X > p.rootPos.X {
						exit.X -= world.TileSize * 0.4
						varY = 2.
						angle = 0.
					} else if p.poppedPos.X < p.rootPos.X {
						exit.X += world.TileSize * 0.4
						varY = 2.
						angle = math.Pi
					} else if p.poppedPos.Y > p.rootPos.Y {
						exit.Y -= world.TileSize * 0.4
						varX = 2.
						angle = math.Pi * 0.5
					} else {
						exit.Y += world.TileSize * 0.4
						varX = 2.
						angle = math.Pi * -0.5
					}
					particles.BiomeParticles(exit, Descent.Cave.Biome, 4, 6, varX, varY, angle, 0.5, 100., 15., 0.75, 0.1, true)
				}
			}).
			SetTrigger(5, func(_ *reanimator.Anim, _ string, _ int) {
				p.action = Aim
				p.fireTimer = timing.New(2.)
			})).
		AddAnimation(reanimator.NewAnimFromSprites("popper_in", img.Reverse(img.Batchers[constants.EntityKey].GetAnimation("popper_out").S), reanimator.Tran).
			SetTrigger(0, func(a *reanimator.Anim, pKey string, pFrame int) {
				if pKey == "popper_out" {
					a.Step = 4 - pFrame
				}
			}).
			SetTrigger(5, func(_ *reanimator.Anim, _ string, _ int) {
				p.action = Seek
				p.Transform.Pos = p.rootPos
				p.Entity.RemoveComponent(myecs.Collision)
				p.Health.Immune[data.Shovel] = data.Immunity{
					KB:    true,
					DMG:   true,
					Dazed: true,
				}
				p.Health.Immune[data.Blast] = data.Immunity{
					KB:    true,
					DMG:   false,
					Dazed: false,
				}
			})).
		AddAnimation(reanimator.NewAnimFromSprites("popper_side", []*pixel.Sprite{img.Batchers[constants.EntityKey].GetSprite("popper_side")}, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("popper_diag", []*pixel.Sprite{img.Batchers[constants.EntityKey].GetSprite("popper_diag")}, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("popper_up", []*pixel.Sprite{img.Batchers[constants.EntityKey].GetSprite("popper_up")}, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("popper_side_fire", img.Batchers[constants.EntityKey].GetAnimation("popper_side_fire").S, reanimator.Tran).
			SetTrigger(0, func(a *reanimator.Anim, pKey string, pFrame int) {
				if pKey == "popper_diag_fire" || pKey == "popper_up_fire" {
					a.Step = pFrame
				}
			}).
			SetTrigger(3, func(_ *reanimator.Anim, _ string, _ int) {
				p.CreateProjectile(p.aimNorm)
			}).
			SetTrigger(4, func(_ *reanimator.Anim, _ string, _ int) {
				p.fireTimer = timing.New(fireSec)
				p.action = Aim
			})).
		AddAnimation(reanimator.NewAnimFromSprites("popper_diag_fire", img.Batchers[constants.EntityKey].GetAnimation("popper_diag_fire").S, reanimator.Tran).
			SetTrigger(0, func(a *reanimator.Anim, pKey string, pFrame int) {
				if pKey == "popper_side_fire" || pKey == "popper_up_fire" {
					a.Step = pFrame
				}
			}).
			SetTrigger(3, func(_ *reanimator.Anim, _ string, _ int) {
				p.CreateProjectile(p.aimNorm)
			}).
			SetTrigger(4, func(_ *reanimator.Anim, _ string, _ int) {
				p.fireTimer = timing.New(fireSec)
				p.action = Aim
			})).
		AddAnimation(reanimator.NewAnimFromSprites("popper_up_fire", img.Batchers[constants.EntityKey].GetAnimation("popper_up_fire").S, reanimator.Tran).
			SetTrigger(0, func(a *reanimator.Anim, pKey string, pFrame int) {
				if pKey == "popper_side_fire" || pKey == "popper_diag_fire" {
					a.Step = pFrame
				}
			}).
			SetTrigger(3, func(_ *reanimator.Anim, _ string, _ int) {
				p.CreateProjectile(p.aimNorm)
			}).
			SetTrigger(4, func(_ *reanimator.Anim, _ string, _ int) {
				p.fireTimer = timing.New(fireSec)
				p.action = Aim
			})).
		SetChooseFn(func() int {
			if p.Health.Dazed {
				return 5
			} else if p.action == Wait || p.action == Seek {
				return 0
			} else if p.action == Pop {
				return 1
			} else if p.action == Unpop {
				return 2
			} else if p.action == Aim {
				if p.angle > math.Pi*0.8 || p.angle < math.Pi*0.2 {
					return 3
				} else if p.angle > math.Pi*0.6 || p.angle < math.Pi*0.4 {
					return 4
				} else {
					return 5
				}
			} else if p.action == Fire {
				if p.angle > math.Pi*0.8 || p.angle < math.Pi*0.2 {
					return 6
				} else if p.angle > math.Pi*0.6 || p.angle < math.Pi*0.4 {
					return 7
				} else {
					return 8
				}
			} else {
				return 5
			}
		}), "")
	p.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Entity, p).
		AddComponent(myecs.Animation, p.Reanimator).
		AddComponent(myecs.Drawable, p.Reanimator).
		AddComponent(myecs.Transform, p.Transform).
		AddComponent(myecs.Physics, p.Physics).
		AddComponent(myecs.Health, p.Health).
		AddComponent(myecs.Batch, constants.EntityKey)
}

func (p *Popper) Delete() {
	myecs.Manager.DisposeEntity(p.Entity)
}

func (p *Popper) CreateProjectile(norm pixel.Vec) {
	e := myecs.Manager.NewEntity()
	trans := transform.New()
	trans.Pos = p.poppedPos
	trans.Pos.X += norm.X * 8.
	trans.Pos.Y += norm.Y * 8.
	phys := physics.New()
	phys.SetVelX(norm.X*fireVel, 0.)
	phys.SetVelY(norm.Y*fireVel, 0.)
	phys.RagDollX = true
	phys.RagDollY = true
	spr := img.Batchers[constants.ParticleKey].GetSprite("dirt_shot")
	coll := data.NewCollider(pixel.R(0., 0., spr.Frame().W(), spr.Frame().H()), false, false)
	coll.Damage = &data.Damage{
		SourceID:  p.Transform.ID,
		Amount:    1,
		Dazed:     1.,
		Knockback: 8.,
		Type:      data.Projectile,
	}
	hp := &data.SimpleHealth{
		Immune: data.EnemyImmunity,
	}
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Collision, coll).
		AddComponent(myecs.Health, hp).
		AddComponent(myecs.Func, data.NewFrameFunc(func() bool {
			coll.Damage.Source = trans.Pos
			if coll.Collided {
				myecs.Manager.DisposeEntity(e)
				particles.CreateRandomParticles(4, 6, []string{"dirt_shot_0", "dirt_shot_1", "dirt_shot_2", "dirt_shot_3", "dirt_shot_4"}, trans.Pos, 0., 0., phys.Velocity.Rotated(math.Pi).Angle(), math.Pi*0.25, math.Min(util.Magnitude(phys.Velocity), 120.), 10.0, 2., 0.5, true)
			}
			return false
		})).
		AddComponent(myecs.Drawable, spr).
		AddComponent(myecs.Batch, constants.ParticleKey)
}
