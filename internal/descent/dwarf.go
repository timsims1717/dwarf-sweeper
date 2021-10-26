package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

const (
	JumpDelay          = 0.05
	stepTime           = 0.4
	GroundAcceleration = 5.
	AirAcceleration    = 10.
	JumpVel            = 150.
	DigRange           = 1.5
	selectTimerSec     = 2.0
	AngleSec           = 0.1
	FacingSec          = 1.0
	angleDiff          = 0.15
)

var (
	ClimbSpeed      = 50.
	Speed           = 80.
	LedgeClimbSpd   = 100.
	MaxJump         = 4
	ShovelKnockback = 6.
	ShovelDazed     = 2.
	ShovelDamage    = 0
)

type DwarfStats struct {
	ClimbSpeed      float64
	Speed           float64
	MaxJump         int
	ShovelKnockback float64
	ShovelDazed     float64
	ShovelDamage    int
}

func DefaultStats() DwarfStats {
	return DwarfStats{
		ClimbSpeed:      ClimbSpeed,
		Speed:           Speed,
		MaxJump:         MaxJump,
		ShovelKnockback: ShovelKnockback,
		ShovelDazed:     ShovelDazed,
		ShovelDamage:    ShovelDamage,
	}
}

type Dwarf struct {
	DwarfStats
	Physics    *physics.Physics
	Transform  *transform.Transform
	Collider   *data.Collider
	Reanimator *reanimator.Tree
	Entity     *ecs.Entity
	Enchants   []string
	EnchantMax int
	faceLeft   bool

	selectLegal bool
	selectTimer *timing.FrameTimer
	angleTimer  *timing.FrameTimer
	facingTimer *timing.FrameTimer

	facing      pixel.Vec
	hovered     *cave.Tile
	relative    pixel.Vec
	isRelative  bool
	digTile     *cave.Tile
	attackPoint pixel.Vec
	tileQueue   []struct{
		a int
		t *cave.Tile
		f pixel.Vec
	}

	walkTimer *timing.FrameTimer
	walking   bool

	jumping    bool
	jumpOrigY  float64
	jumpTarget float64
	jumpTimer  *timing.FrameTimer
	toJump     bool
	jumpEnd    bool
	distFell   float64

	digging    bool
	attacking  bool
	marking    bool
	climbing   bool
	airDig     bool
	digHold    bool
	markHold   bool

	Health    *data.Health
	DeadStop  bool

	Bubble *Bubble
}

func NewDwarf(start pixel.Vec) *Dwarf {
	tran := transform.NewTransform()
	tran.Pos = start
	d := &Dwarf{
		DwarfStats: DefaultStats(),
		Physics:    physics.New(),
		Transform:  tran,
		Health: &data.Health{
			Max:          3,
			Curr:         3,
			TempInvSec:   3.,
			DazeOverride: true,
			Immune:       []data.DamageType{data.Shovel},
		},
	}
	batcher := img.Batchers[constants.DwarfKey]
	idleAnim := batcher.GetAnimation("idle")
	climbAnim := batcher.GetAnimation("climb")
	d.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewSwitch(&reanimator.Switch{
				Elements: reanimator.NewElements(
					reanimator.NewSwitch(&reanimator.Switch{
						Elements: reanimator.NewElements(
								reanimator.NewAnimFromSprite("hit_front", batcher.GetSprite("hit_front"), reanimator.Hold, nil), // hit_front
								reanimator.NewAnimFromSprite("hit_back", batcher.GetSprite("hit_front"), reanimator.Hold, nil),  // hit_back
							),
						Check: func() int {
							if d.faceLeft {
								if d.Physics.Velocity.X > 0 {
									return 0
								} else {
									return 1
								}
							} else {
								if d.Physics.Velocity.X > 0 {
									return 1
								} else {
									return 0
								}
							}
						},
					}),
					reanimator.NewAnimFromSprite("flat", batcher.GetSprite("flat"), reanimator.Hold, nil), // flat
				),
				Check: func() int {
					if !d.Physics.Grounded {
						return 0
					} else {
						return 1
					}
				},
			}),
			reanimator.NewAnimFromSprite("dig_hold", batcher.GetSprite("dig"), reanimator.Hold, nil), // dig hold
			reanimator.NewAnimFromSprite("mark_hold", batcher.GetSprite("mark"), reanimator.Hold, nil), // mark hold
			reanimator.NewAnimFromSprites("dig", batcher.GetAnimation("dig").S, reanimator.Tran, map[int]func() {
				0: func() {
					if d.digTile == nil || !d.digTile.Solid() {
						trans := transform.NewTransform()
						trans.Pos = d.attackPoint
						trans.Flip = d.faceLeft
						key := "shovel_swipe"
						if d.facing.X == 0. {
							if d.faceLeft {
								if d.facing.Y > 0. {
									trans.Rot = -0.5
								} else {
									trans.Rot = 0.5
								}
							} else {
								if d.facing.Y > 0. {
									trans.Rot = 0.5
								} else {
									trans.Rot = -0.5
								}
							}
						} else if d.facing.Y != 0. && d.facing.X != 0. {
							key = "shovel_swipe_diag"
							if d.facing.Y < 0. {
								if d.faceLeft {
									trans.Rot = 0.5
								} else {
									trans.Rot = -0.5
								}
							}
						}
						anim := img.Batchers[constants.ParticleKey].GetAnimation(key).S
						e := myecs.Manager.NewEntity()
						e.AddComponent(myecs.Animation, reanimator.NewSimple(
								reanimator.NewAnimFromSprites("swipe", anim, reanimator.Done, map[int]func(){
									3: func() {
										myecs.Manager.DisposeEntity(e)
									},
								}).Anim,
							)).
							AddComponent(myecs.Transform, trans).
							AddComponent(myecs.Batch, constants.ParticleKey).
							AddComponent(myecs.Temp, myecs.ClearFlag(false))
					}
				},
				1: func() {
					if d.digTile != nil && d.digTile.Solid() {
						CaveBlocksDug++
						d.digTile.Destroy(true)
						d.digTile = nil
					} else {
						myecs.Manager.NewEntity().
							AddComponent(myecs.AreaDmg, &data.AreaDamage{
								Center:    d.attackPoint,
								Radius:    world.TileSize * 1.2,
								Amount:    d.ShovelDamage,
								Dazed:     d.ShovelDazed,
								Knockback: d.ShovelKnockback,
								Source:    d.Transform.Pos,
								Type:      data.Shovel,
							})
					}
					sfx.SoundPlayer.PlaySound("shovel", 1.0)
				},
				3: func() {
					d.digging = false
					d.attacking = false
				},
			}), // digging
			reanimator.NewAnimFromSprite("mark", batcher.GetSprite("mark"), reanimator.Tran, map[int]func() {
				1: func() {
					d.marking = false
				},
			}), // marking
			reanimator.NewSwitch(&reanimator.Switch{
				Elements: reanimator.NewElements(
					reanimator.NewSwitch(&reanimator.Switch{
						Elements: reanimator.NewElements(
							reanimator.NewAnimFromSprites("run", batcher.GetAnimation("run").S, reanimator.Loop, nil), // run
							reanimator.NewSwitch(&reanimator.Switch{
								Elements: reanimator.NewElements(
									reanimator.NewAnimFromSprite("flat", batcher.GetSprite("flat"), reanimator.Hold, nil), // flat
									reanimator.NewAnimFromSprites("idle", []*pixel.Sprite{
										idleAnim.S[0], idleAnim.S[0],
										idleAnim.S[1], idleAnim.S[1],
									}, reanimator.Loop, nil), // idle
								),
								Check: func() int {
									if d.distFell > 100. {
										return 0
									} else {
										return 1
									}
								},
							}),
						),
						Check: func() int {
							if d.Physics.IsMovingX() || (d.Bubble != nil && d.Bubble.Physics.IsMovingX()) {
								return 0
							} else {
								return 1
							}
						},
					}),
					reanimator.NewSwitch(&reanimator.Switch{
						Elements: reanimator.NewElements(
							reanimator.NewAnimFromSprites("climb_up", climbAnim.S, reanimator.Loop, nil), // climb_up
							reanimator.NewAnimFromSprites("climb_dwn", img.Reverse(climbAnim.S), reanimator.Loop, nil), // climb_dwn
							reanimator.NewAnimFromSprite("climb_still", climbAnim.S[3], reanimator.Hold, nil), // climb_still
						),
						Check: func() int {
							if d.Physics.Velocity.Y > 0. {
								return 0
							} else if d.Physics.Velocity.Y < 0. {
								return 1
							} else {
								return 2
							}
						},
					}),
					reanimator.NewSwitch(&reanimator.Switch{
						Elements: reanimator.NewElements(
							reanimator.NewAnimFromSprites("jump", batcher.GetAnimation("jump").S, reanimator.Hold, nil), // jump
							reanimator.NewAnimFromSprite("fall", batcher.GetSprite("fall"), reanimator.Hold, nil),   // fall
						),
						Check: func() int {
							if d.Physics.Velocity.Y > 0. || d.jumping || d.toJump || d.jumpEnd ||
								(d.Bubble != nil && d.Bubble.Physics.Velocity.Y > 0.) {
								return 0
							} else {
								return 1
							}
						},
					}),
				),
				Check: func() int {
					if (d.Physics.Grounded && !d.jumping && !d.toJump && d.Bubble == nil) ||
						(d.Bubble != nil && d.Bubble.Physics.Grounded) {
						return 0
					} else if d.climbing {
						return 1
					} else {
						return 2
					}
				},
			}),
		),
		Check: func() int {
			if d.Health.Dazed || d.Health.Dead {
				return 0
			} else if d.digHold {
				return 1
			} else if d.markHold {
				return 2
			} else if d.digging || d.attacking {
				return 3
			} else if d.marking {
				return 4
			} else {
				return 5
			}
		},
	}, "idle")
	d.Collider = &data.Collider{
		Hitbox: pixel.R(0., 0., 16., 16.),
		CanPass: true,
	}
	d.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, tran).
		AddComponent(myecs.Physics, d.Physics).
		AddComponent(myecs.Collision, d.Collider).
		AddComponent(myecs.Animation, d.Reanimator).
		AddComponent(myecs.Health, d.Health)
	return d
}

func (d *Dwarf) Update(in *input.Input) {
	d.angleTimer.Update()
	if d.Physics.Grounded || d.climbing || d.Bubble != nil {
		d.airDig = false
	}
	if d.Health.Dazed || d.Health.Dead {
		if d.Bubble != nil {
			d.Bubble.Pop()
		}
		d.tileQueue = []struct {
			a int
			t *cave.Tile
			f pixel.Vec
		}{}
		d.digging = false
		d.attacking = false
		d.jumping = false
		d.walking = false
		d.climbing = false
		if d.Health.Dead {
			if d.Physics.Grounded && !d.Physics.IsMovingX() && !d.DeadStop {
				d.Physics.CancelMovement()
				//d.distFell = 150.
				d.DeadStop = true
			}
		} else if d.Physics.Grounded && !d.Physics.IsMovingX() &&
			(in.Get("left").Pressed() || in.Get("right").Pressed() ||
				in.Get("up").Pressed() || in.Get("down").Pressed() ||
				in.Get("jump").JustPressed() || in.Get("dig").JustPressed() ||
				in.Get("mark").JustPressed()) {
			d.Health.DazedO = false
		}
	}
	if !d.Health.Dazed && !d.Health.Dead {
		if constants.AimDedicated {
			if in.Mode != input.KeyboardMouse &&
				(in.Axes["targetX"].F > 0. || in.Axes["targetX"].F < 0. ||
				in.Axes["targetY"].F > 0. || in.Axes["targetY"].F < 0.) {
				x := in.Axes["targetX"].R
				y := in.Axes["targetY"].R
				xA := math.Abs(x)
				yA := math.Abs(y)
				diff := math.Abs(x - y)
				if xA > input.Deadzone ||
					(yA > input.Deadzone && diff < angleDiff) {
					if x > 0. {
						x = world.TileSize
					} else {
						x = -world.TileSize
					}
				}
				if yA > input.Deadzone ||
					(xA > input.Deadzone && diff < angleDiff) {
					if y > 0. {
						y = -world.TileSize
					} else {
						y = world.TileSize
					}
				}
				d.relative = pixel.V(x, y)
				d.isRelative = true
			} else if in.Mode != input.Gamepad {
				d.hovered = Descent.GetCave().GetTile(in.World)
				d.facingTimer = nil
				d.isRelative = false
			}
		} else {
			moveSelecting := in.Get("left").Pressed() || in.Get("right").Pressed() || in.Get("up").Pressed() || in.Get("down").Pressed()
			jpMS := in.Get("left").JustPressed() || in.Get("right").JustPressed() || in.Get("up").JustPressed() || in.Get("down").JustPressed()
			if in.Get("left").JustReleased() || in.Get("right").JustReleased() || in.Get("up").JustReleased() || in.Get("down").JustReleased() {
				d.angleTimer = timing.New(AngleSec)
			}
			if moveSelecting && (jpMS || d.angleTimer.Done()) {
				x := 0.
				y := 0.
				if in.Get("left").Pressed() {
					x = -world.TileSize
				} else if in.Get("right").Pressed() {
					x = world.TileSize
				}
				if in.Get("down").Pressed() {
					y = -world.TileSize
				} else if in.Get("up").Pressed() {
					y = world.TileSize
				}
				d.relative = pixel.V(x, y)
				d.isRelative = true
			}
		}
		if d.isRelative {
			p := d.Transform.Pos
			p.X += d.relative.X
			p.Y += d.relative.Y
			d.hovered = Descent.GetCave().GetTile(p)
			d.facingTimer = timing.New(FacingSec)
		}
		if d.hovered != nil && !d.airDig && len(d.tileQueue) < 3 {
			d.selectLegal = math.Abs(d.Transform.Pos.X-d.hovered.Transform.Pos.X) < world.TileSize*DigRange && math.Abs(d.Transform.Pos.Y-d.hovered.Transform.Pos.Y) < world.TileSize*DigRange
			facing := pixel.ZV
			if (in.Get("dig").JustPressed() && !constants.DigOnRelease) || (in.Get("dig").JustReleased() && constants.DigOnRelease) {
				if in.Mode != input.Gamepad && constants.AimDedicated {
					angle := d.Transform.Pos.Sub(in.World).Angle()
					if angle > math.Pi*(5./8.) || angle < math.Pi*-(5./8.) {
						facing.X = 1
					} else if angle < math.Pi*(3./8.) && angle > math.Pi*-(3./8.) {
						facing.X = -1
					} else {
						facing.X = 0
					}
					if angle > math.Pi/8. && angle < math.Pi*(7./8.) {
						facing.Y = -1
					} else if angle < math.Pi/-8. && angle > math.Pi*-(7./8.) {
						facing.Y = 1
					} else {
						facing.Y = 0
					}
				} else if d.isRelative {
					if d.relative.X < 0 {
						facing.X = -1
					} else if d.relative.X > 0 {
						facing.X = 1
					} else {
						facing.X = 0
					}
					if d.relative.Y < 0 {
						facing.Y = -1
					} else if d.relative.Y > 0 {
						facing.Y = 1
					} else {
						facing.Y = 0
					}
				}
				if d.hovered.Solid() && d.selectLegal {
					d.tileQueue = append(d.tileQueue, struct {
						a int
						t *cave.Tile
						f pixel.Vec
					}{
						a: 0,
						t: d.hovered,
						f: facing,
					})
				} else {
					d.tileQueue = append(d.tileQueue, struct {
						a int
						t *cave.Tile
						f pixel.Vec
					}{
						a: 2,
						t: nil,
						f: facing,
					})
				}
			} else if ((in.Get("mark").JustPressed() && !constants.DigOnRelease) || (in.Get("mark").JustReleased() && constants.DigOnRelease)) && d.hovered.Solid() && d.selectLegal {
				d.tileQueue = append(d.tileQueue, struct {
					a int
					t *cave.Tile
					f pixel.Vec
				}{
					a: 1,
					t: d.hovered,
					f: facing,
				})
			}
		}
		if len(d.tileQueue) > 0 && !d.digging && !d.attacking && !d.marking {
			next := d.tileQueue[0]
			d.tileQueue = d.tileQueue[1:]
			d.facing = next.f
			if next.t != nil {
				if next.t.Transform.Pos.X < d.Transform.Pos.X {
					d.faceLeft = true
				} else if next.t.Transform.Pos.X > d.Transform.Pos.X {
					d.faceLeft = false
				}
			}
			digLegal := next.t != nil &&
				math.Abs(d.Transform.Pos.X-next.t.Transform.Pos.X) < world.TileSize*DigRange &&
				math.Abs(d.Transform.Pos.Y-next.t.Transform.Pos.Y) < world.TileSize*DigRange
			if next.a == 0 && next.t.Solid() && digLegal {
				d.digging = true
				d.attacking = false
				d.jumping = false
				d.walking = false
				d.distFell = 0.
				d.digTile = next.t
			} else if next.a == 1 && next.t.Solid() && digLegal {
				d.marking = true
				d.distFell = 0.
				Mark(next.t)
			} else {
				var x, y float64
				if d.facing.Y < 0 {
					y = d.Transform.Pos.Y - world.TileSize * 0.6
				} else if d.facing.Y > 0 {
					y = d.Transform.Pos.Y + world.TileSize * 0.6
				} else {
					y = d.Transform.Pos.Y
				}
				if d.facing.X < 0 {
					x = d.Transform.Pos.X - world.TileSize * 0.6
				} else if d.facing.X > 0 {
					x = d.Transform.Pos.X + world.TileSize * 0.6
				} else {
					x = d.Transform.Pos.X
				}
				d.attackPoint = pixel.V(x, y)
				d.digging = false
				d.attacking = true
				d.jumping = false
				d.walking = false
				d.distFell = 0.
				d.digTile = nil
			}
		}
		d.digHold = constants.DigOnRelease && in.Get("dig").Pressed() && !d.markHold
		d.markHold = constants.DigOnRelease && in.Get("mark").Pressed() && !d.digHold
		if d.digHold || d.markHold {
			if d.climbing {
				d.Physics.CancelMovement()
			}
			d.digging = false
			d.attacking = false
			d.jumping = false
			d.walking = false
			if d.hovered.Transform.Pos.X < d.Transform.Pos.X {
				d.facing.X = -1
			} else {
				d.facing.X = 1
			}
		} else {
			if d.digging || d.attacking || d.marking {
				if !d.Physics.Grounded && !d.airDig && d.Bubble == nil {
					d.airDig = true
				}
				d.Physics.CancelMovement()
			} else if d.Bubble != nil {
				if in.Get("left").Pressed() && !in.Get("right").Pressed() {
					d.faceLeft = true
					d.Bubble.Physics.SetVelX(-BubbleVel, BubbleAcc)
				} else if in.Get("right").Pressed() && !in.Get("left").Pressed() {
					d.faceLeft = false
					d.Bubble.Physics.SetVelX(BubbleVel, BubbleAcc)
				}
				if in.Get("up").Pressed() && !in.Get("down").Pressed() {
					d.Bubble.Physics.SetVelY(BubbleVel, BubbleAcc)
				} else if in.Get("down").Pressed() && !in.Get("up").Pressed() {
					d.Bubble.Physics.SetVelY(-BubbleVel, BubbleAcc)
				}
				d.walking = false
				d.jumping = false
				d.climbing = false
				d.distFell = 0.
			} else {
				canJump := d.Physics.NearGround || d.Physics.Grounded
				canClimb := d.Physics.RightWalled || d.Physics.LeftWalled

				xDir := 0
				if in.Get("left").Pressed() && !in.Get("right").Pressed() {
					xDir = 1
				} else if in.Get("right").Pressed() && !in.Get("left").Pressed() {
					xDir = 2
				}

				switch xDir {
				case 1:
					if d.Physics.Grounded {
						d.faceLeft = true
						d.Physics.SetVelX(-d.Speed, GroundAcceleration)
					} else {
						d.Physics.SetVelX(-d.Speed, AirAcceleration)
					}
				case 2:
					if d.Physics.Grounded {
						d.faceLeft = false
						d.Physics.SetVelX(d.Speed, GroundAcceleration)
					} else {
						d.Physics.SetVelX(d.Speed, AirAcceleration)
					}
				}
				// Ground test, considered on the ground for jumping purposes until half a tile out
				if !d.jumping && canJump && in.Get("jump").JustPressed() {
					d.toJump = true
					d.climbing = false
					d.walking = false
					d.distFell = 0.
					d.jumpTimer = timing.New(JumpDelay)
				} else if d.toJump && d.jumpTimer.UpdateDone() {
					d.climbing = false
					d.toJump = false
					d.walking = false
					d.jumping = true
					d.jumpOrigY = d.Transform.Pos.Y
					sfx.SoundPlayer.PlaySound(fmt.Sprintf("step%d", random.Effects.Intn(4)+1), 0.)
					d.distFell = 0.
					d.Physics.SetVelY(JumpVel, 0.)
				} else if d.climbing {
					if in.Get("jump").JustPressed() {
						d.climbing = false
					} else if canClimb {
						d.distFell = 0.
						if in.Get("up").Pressed() && !in.Get("down").Pressed() {
							d.Physics.SetVelY(d.ClimbSpeed, 0.)
						} else if in.Get("down").Pressed() && !in.Get("up").Pressed() {
							d.Physics.SetVelY(-d.ClimbSpeed, 0.)
						} else {
							d.Physics.SetVelY(0., 0.)
						}
						if d.Physics.RightWalled && (!d.Physics.LeftWalled || in.Get("right").Pressed()) {
							d.faceLeft = false
						} else if d.Physics.LeftWalled && (!d.Physics.RightWalled || in.Get("left").Pressed()) {
							d.faceLeft = true
						}
					} else if in.Get("up").Pressed() && !d.Physics.Grounded &&
						!(in.Get("left").Pressed() || in.Get("right").Pressed()) &&
						((d.Collider.CDL && d.faceLeft) || (d.Collider.CDR && !d.faceLeft)) {
						d.Physics.SetVelY(d.ClimbSpeed, 0.)
						if d.faceLeft {
							d.Physics.SetVelX(-Speed, AirAcceleration)
						} else {
							d.Physics.SetVelX(Speed, AirAcceleration)
						}
					} else {
						d.climbing = false
					}
				} else if canClimb && !d.toJump && in.Get("up").Pressed() {
					d.climbing = true
					d.walking = false
					d.jumping = false
					d.toJump = false
					d.distFell = 0.
					d.Physics.SetVelY(d.ClimbSpeed, 0.)
					if d.Physics.RightWalled && !d.Physics.LeftWalled {
						d.faceLeft = false
					} else if d.Physics.LeftWalled && !d.Physics.RightWalled {
						d.faceLeft = true
					}
				} else if !d.jumping && !d.toJump && d.Physics.Grounded {
					wasWalking := d.walking
					if math.Abs(d.Physics.Velocity.X) < 20.0 {
						if in.Get("up").Pressed() && !in.Get("down").Pressed() {
							d.distFell = 0.
							camera.Cam.Up()
						} else if in.Get("down").Pressed() && !in.Get("up").Pressed() {
							d.distFell = 0.
							camera.Cam.Down()
						}
						d.walking = false
						d.climbing = false
					} else if d.Physics.Velocity.X > 0. {
						d.faceLeft = false
						d.walking = true
					} else if d.Physics.Velocity.X < 0. {
						d.faceLeft = true
						d.walking = true
					}
					if d.walking {
						d.distFell = 0.
						if !wasWalking {
							d.walkTimer = timing.New(stepTime)
						}
					}
				} else {
					d.walking = false
					d.climbing = false
					if d.jumping || d.jumpEnd {
						height := int(((d.Transform.Pos.Y - d.jumpOrigY) + world.TileSize*1.0) / world.TileSize)
						if d.Physics.Velocity.Y <= 0. {
							d.jumping = false
							d.jumpEnd = false
						} else if height < d.MaxJump-1 && in.Get("jump").Pressed() {
							d.Physics.SetVelY(JumpVel, 0.)
						} else if !d.jumpEnd {
							in.Get("jump").Consume()
							d.jumping = false
							d.jumpEnd = true
							d.jumpTarget = Descent.GetCave().GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y+world.TileSize*1.0)).Transform.Pos.Y
						}
						if d.jumpEnd {
							if d.jumpTarget > d.Transform.Pos.Y {
								d.Physics.SetVelY(0., 0.5)
							} else {
								d.Physics.SetVelY(0., 0.)
								d.jumpEnd = false
							}
						}
					}
					if d.Physics.Velocity.Y < 0. {
						d.distFell += math.Abs(d.Physics.Velocity.Y * timing.DT)
					}
				}
			}
			if in.Get("prev").JustPressed() {
				PrevItem()
			} else if in.Get("next").JustPressed() {
				NextItem()
			} else if in.Get("use").JustPressed() {
				UseEquipped()
			}
		}
	}
	if d.facingTimer != nil && d.facingTimer.UpdateDone() {
		if d.faceLeft {
			d.facing.X = -1
		} else {
			d.facing.X = 1
		}
		d.facing.Y = 0
	} else if d.facing.X < 0 && (d.digHold || d.markHold || d.digging || d.marking || d.attacking) {
		d.faceLeft = true
	} else if d.facing.X > 0 && (d.digHold || d.markHold || d.digging || d.marking || d.attacking) {
		d.faceLeft = false
	}

}

func (d *Dwarf) Update2() {
	d.Transform.Flip = d.faceLeft
	camera.Cam.StayWithin(d.Transform.Pos, world.TileSize * 1.5)
	if d.walking && d.walkTimer.UpdateDone() {
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("step%d", random.Effects.Intn(4) + 1), 0.)
		d.walkTimer = timing.New(stepTime)
	}
}

func (d *Dwarf) Draw(win *pixelgl.Window, in *input.Input) {
	d.Reanimator.CurrentSprite().Draw(win, d.Transform.Mat)
	if d.hovered != nil && !d.Health.Dazed {
		if d.hovered.Solid() && d.selectLegal {
			img.Batchers[constants.ParticleKey].DrawSprite("target", d.hovered.Transform.Mat)
		} else {
			img.Batchers[constants.ParticleKey].DrawSprite("target_blank", d.hovered.Transform.Mat)
		}
	}
	if debug.Text {
		debug.AddText(fmt.Sprintf("world coords: (%d,%d)", int(in.World.X), int(in.World.Y)))
		if d.hovered != nil {
			debug.AddText(fmt.Sprintf("tile coords: (%d,%d)", d.hovered.RCoords.X, d.hovered.RCoords.Y))
			debug.AddText(fmt.Sprintf("chunk coords: (%d,%d)", d.hovered.Chunk.Coords.X, d.hovered.Chunk.Coords.Y))
			debug.AddText(fmt.Sprintf("tile sub coords: (%d,%d)", d.hovered.SubCoords.X, d.hovered.SubCoords.Y))
			debug.AddText(fmt.Sprintf("tile type: '%s'", d.hovered.Type))
			debug.AddText(fmt.Sprintf("tile sprite: '%s'", d.hovered.BGSpriteS))
			debug.AddText(fmt.Sprintf("tile smart str: '%s'", d.hovered.SmartStr))
		}
		debug.AddText(fmt.Sprintf("dwarf position: (%d,%d)", int(d.Transform.APos.X), int(d.Transform.APos.Y)))
		debug.AddText(fmt.Sprintf("dwarf actual position: (%f,%f)", d.Transform.Pos.X, d.Transform.Pos.Y))
		debug.AddText(fmt.Sprintf("dwarf velocity: (%d,%d)", int(d.Physics.Velocity.X), int(d.Physics.Velocity.Y)))
		debug.AddText(fmt.Sprintf("dwarf moving?: (%t,%t)", d.Physics.IsMovingX(), d.Physics.IsMovingY()))
		//debug.AddText(fmt.Sprintf("jump pressed?: %t", input.Input.Jumping.Pressed()))
		debug.AddText(fmt.Sprintf("dwarf grounded?: %t", d.Physics.Grounded))
		debug.AddText(fmt.Sprintf("tile queue len: %d", len(d.tileQueue)))
	}
}

func (d *Dwarf) Delete() {
	d.Health.Delete()
	myecs.Manager.DisposeEntity(d.Entity)
}

func Mark(tile *cave.Tile) {
	if tile != nil && tile.Solid() && !tile.Destroyed && tile.Breakable() {
		if !tile.Marked {
			tile.Marked = true
			f := &Flag{
				Tile: tile,
			}
			f.Create(pixel.ZV)
		} else {
			tile.Marked = false
		}
	}
}