package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/internal/profile"
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
	pxginput "github.com/timsims1717/pixel-go-input"
	"math"
)

const (
	JumpDelay          = 0.05
	stepTime           = 0.4
	GroundAcceleration = 5.
	AirAcceleration    = 10.
	JumpVel            = 150.
	DigRange           = 1.5
	AngleSec           = 0.1
	FacingSec          = 1.0
	angleDiff          = 0.15
)

var (
	ClimbSpeed      = 50.
	Speed           = 80.
	MaxJump         = 4
	ShovelKnockback = 8.
	ShovelDazed     = 2.
	ShovelDamage    = 0
	AirDigs         = 2
)

type DwarfStats struct {
	ClimbSpeed      float64
	Speed           float64
	MaxJump         int
	ShovelKnockback float64
	ShovelDazed     float64
	ShovelDamage    int
	AirDigs         int
}

func DefaultStats() DwarfStats {
	return DwarfStats{
		ClimbSpeed:      ClimbSpeed,
		Speed:           Speed,
		MaxJump:         MaxJump,
		ShovelKnockback: ShovelKnockback,
		ShovelDazed:     ShovelDazed,
		ShovelDamage:    ShovelDamage,
		AirDigs:         AirDigs,
	}
}

type Dwarf struct {
	DwarfStats
	Player     *data.Player
	Health     *data.Health
	Physics    *physics.Physics
	Transform  *transform.Transform
	Collider   *data.Collider
	Reanimator *reanimator.Tree
	Entity     *ecs.Entity
	Enchants   []string
	EnchantMax int
	faceLeft   bool

	SelectLegal bool
	selectTimer *timing.Timer
	angleTimer  *timing.Timer
	facingTimer *timing.Timer

	facing      pixel.Vec
	Hovered     *cave.Tile
	relative    pixel.Vec
	isRelative  bool
	digTile     *cave.Tile
	attackPoint pixel.Vec
	tileQueue   []struct {
		a int
		t *cave.Tile
		f pixel.Vec
	}

	walking    bool
	jumping    bool
	jumpOrigY  float64
	jumpTarget float64
	jumpTimer  *timing.Timer
	toJump     bool
	jumpEnd    bool
	distFell   float64

	digging   bool
	attacking bool
	flagging  bool
	climbing  bool
	airDig    int
	digHold   bool
	flagHold  bool
	dropTimer *timing.Timer
	stopDrop  bool

	DeadStop bool
}

func NewDwarf(p *data.Player) *Dwarf {
	tran := transform.New().WithID(fmt.Sprintf("dwarf-%s", p.Code))
	d := &Dwarf{
		DwarfStats: DefaultStats(),
		Health:     &data.Health{
			Max:        profile.CurrentProfile.StartingAttr.MaxHealth,
			Curr:       profile.CurrentProfile.StartingAttr.MaxHealth,
			TempInvSec: 3.,
			DazedTime:  10.,
			Immune:     data.ShovelImmunity,
		},
		Physics:   physics.New(),
		Transform: tran,
		Player:    p,
		relative:  pixel.V(world.TileSize, 0.),
	}
	batcher := img.Batchers[constants.DwarfKey]
	climbAnim := batcher.GetAnimation("climb")
	d.Reanimator = reanimator.New(reanimator.NewSwitch().
		AddSubSwitch(reanimator.NewSwitch().
			AddSubSwitch(reanimator.NewSwitch().
				AddAnimation(reanimator.NewAnimFromSprite("hit_front", batcher.GetSprite("hit_front"), reanimator.Hold)). // hit_front
				AddAnimation(reanimator.NewAnimFromSprite("hit_back", batcher.GetSprite("hit_front"), reanimator.Hold)).  // hit_back
				SetChooseFn(func() int {
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
				})).
			AddAnimation(reanimator.NewAnimFromSprite("flat", batcher.GetSprite("flat"), reanimator.Hold)). // flat
			SetChooseFn(func() int {
				if !d.Physics.Grounded {
					return 0
				} else {
					return 1
				}
			})).
		AddAnimation(reanimator.NewAnimFromSprite("dig_hold", batcher.GetSprite("dig"), reanimator.Hold)).   // dig hold
		AddAnimation(reanimator.NewAnimFromSprite("flag_hold", batcher.GetSprite("flag"), reanimator.Hold)). // flag hold
		AddAnimation(reanimator.NewAnimFromSprites("dig", batcher.GetAnimation("dig").S, reanimator.Tran).
		SetTrigger(0, func() {
			if d.digTile == nil || !d.digTile.Solid() || d.digTile.IsDeco() {
				trans := transform.New()
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
				tree := reanimator.NewSimple(
					reanimator.NewAnimFromSprites("swipe", anim, reanimator.Done).
						SetTrigger(3, func() {
							myecs.Manager.DisposeEntity(e)
						}),
				)
				e.AddComponent(myecs.Animation, tree).
					AddComponent(myecs.Drawable, tree).
					AddComponent(myecs.Transform, trans).
					AddComponent(myecs.Batch, constants.ParticleKey).
					AddComponent(myecs.Temp, myecs.ClearFlag(false))
			}
		}).
		SetTrigger(1, func() {
			if d.digTile != nil && d.digTile.Solid() && !d.digTile.IsDeco() {
				Dig(d.digTile, d.Player)
				d.digTile = nil
			} else {
				sub := d.attackPoint.Sub(d.Transform.Pos)
				sub.Y += world.TileSize * 0.5
				angle := sub.Angle()
				myecs.Manager.NewEntity().
					AddComponent(myecs.AreaDmg, &data.AreaDamage{
						SourceID:  d.Transform.ID,
						Center:    d.attackPoint,
						Radius:    world.TileSize * 1.2,
						Amount:    d.ShovelDamage,
						Dazed:     d.ShovelDazed,
						Knockback: d.ShovelKnockback,
						Angle:     &angle,
						Type:      data.Shovel,
					})
			}
			sfx.SoundPlayer.PlaySound("shovel", 1.0)
		}).
		SetTrigger(3, func() {
			d.digging = false
			d.attacking = false
		})). // digging
		AddAnimation(reanimator.NewAnimFromSprite("flag", batcher.GetSprite("flag"), reanimator.Tran).
			SetTrigger(1, func() {
				d.flagging = false
			})). // flagging
		AddSubSwitch(reanimator.NewSwitch().
			AddSubSwitch(reanimator.NewSwitch().
				AddAnimation(reanimator.NewAnimFromSprites("run", batcher.GetAnimation("run").S, reanimator.Loop).
			SetTrigger(0, func() {
					PlayStep(0.)
				}).
			SetTrigger(4, func() {
					PlayStep(0.)
				})). // run
				AddSubSwitch(reanimator.NewSwitch().
					AddAnimation(reanimator.NewAnimFromSprite("flat", batcher.GetSprite("flat"), reanimator.Hold)).       // flat
					AddAnimation(reanimator.NewAnimFromSprites("idle", batcher.GetAnimation("idle").S, reanimator.Loop)). // idle
					SetChooseFn(func() int {
						if d.distFell > 100. {
							return 0
						} else {
							return 1
						}
					})).
				SetChooseFn(func() int {
					if d.Physics.IsMovingX() {
						return 0
					} else {
						return 1
					}
				})).
			AddSubSwitch(reanimator.NewSwitch().
				AddAnimation(reanimator.NewAnimFromSprites("climb_up", climbAnim.S, reanimator.Loop)).               // climb_up
				AddAnimation(reanimator.NewAnimFromSprites("climb_dwn", img.Reverse(climbAnim.S), reanimator.Loop)). // climb_dwn
				AddAnimation(reanimator.NewAnimFromSprite("climb_still", climbAnim.S[3], reanimator.Hold)).          // climb_still
				SetChooseFn(func() int {
					if d.Physics.Velocity.Y > 0. {
						return 0
					} else if d.Physics.Velocity.Y < 0. {
						return 1
					} else {
						return 2
					}
				})).
			AddSubSwitch(reanimator.NewSwitch().
				AddAnimation(reanimator.NewAnimFromSprites("jump", batcher.GetAnimation("jump").S, reanimator.Hold)). // jump
				AddAnimation(reanimator.NewAnimFromSprite("fall", batcher.GetSprite("fall"), reanimator.Hold)).       // fall
				SetChooseFn(func() int {
					if d.Physics.Velocity.Y > 0. || d.jumping || d.toJump || d.jumpEnd {
						return 0
					} else {
						return 1
					}
				})).
			SetChooseFn(func() int {
				if d.Physics.Grounded && !d.jumping && !d.toJump {
					return 0
				} else if d.climbing {
					return 1
				} else {
					return 2
				}
			})).
		SetChooseFn(func() int {
			if d.Health.Dazed || d.Health.Dead {
				return 0
			} else if d.digHold {
				return 1
			} else if d.flagHold {
				return 2
			} else if d.digging || d.attacking {
				return 3
			} else if d.flagging {
				return 4
			} else {
				return 5
			}
		}), "idle")
	d.Collider = data.NewCollider(pixel.R(0., 0., 16., 16.), data.PlayerC)
	d.Collider.Debug = true
	d.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, tran).
		AddComponent(myecs.Physics, d.Physics).
		AddComponent(myecs.Collision, d.Collider).
		AddComponent(myecs.Animation, d.Reanimator).
		AddComponent(myecs.Drawable, d.Reanimator).
		AddComponent(myecs.Batch, constants.DwarfKey).
		AddComponent(myecs.Health, d.Health).
		AddComponent(myecs.Player, d.Player)
	return d
}

func (d *Dwarf) SetStart(pos pixel.Vec) {
	d.Transform.Pos = pos
	d.Player.CamPos = pos
	hPos := pos
	hPos.X += world.TileSize
	d.Hovered = Descent.Cave.GetTile(hPos)
	d.isRelative = true
}

func (d *Dwarf) Update(in *pxginput.Input) {
	cameraMoveX := false
	cameraMoveY := false
	useItem := false
	d.angleTimer.Update()
	if d.Physics.Grounded || d.climbing {
		d.airDig = 0
	}
	if d.Health.Dazed || d.Health.Dead {
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
		} else if d.Physics.Grounded &&
			d.Health.DazedTimer.Elapsed() > 1. &&
			(in.Get("left").Pressed() || in.Get("right").Pressed() ||
				in.Get("up").Pressed() || in.Get("down").Pressed() ||
				in.Get("jump").JustPressed() || in.Get("dig").JustPressed() ||
				in.Get("flag").JustPressed()) {
			d.Health.Dazed = false
		}
	}
	if !d.Health.Dazed && !d.Health.Dead && in != nil {
		if in.OptFlags["AimDedicated"] {
			if in.Mode != pxginput.KeyboardMouse &&
				(in.Axes["targetX"].F > 0. || in.Axes["targetX"].F < 0. ||
					in.Axes["targetY"].F > 0. || in.Axes["targetY"].F < 0.) {
				x := in.Axes["targetX"].R
				y := in.Axes["targetY"].R
				xA := math.Abs(x)
				yA := math.Abs(y)
				diff := math.Abs(x - y)
				if xA > in.Deadzone ||
					(yA > in.Deadzone && diff < angleDiff) {
					if x > 0. {
						x = world.TileSize
					} else {
						x = -world.TileSize
					}
				}
				if yA > in.Deadzone ||
					(xA > in.Deadzone && diff < angleDiff) {
					if y > 0. {
						y = -world.TileSize
					} else {
						y = world.TileSize
					}
				}
				d.relative = pixel.V(x, y)
				d.isRelative = true
			} else if in.Mode != pxginput.Gamepad {
				d.Hovered = Descent.GetCave().GetTile(d.Player.CamPos.Sub(d.Player.CanvasPos.Sub(in.World)))
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
			d.Hovered = Descent.GetCave().GetTile(p)
			d.facingTimer = timing.New(FacingSec)
		}
		if d.Hovered != nil && d.airDig < d.AirDigs && len(d.tileQueue) < 3 {
			d.SelectLegal = math.Abs(d.Transform.Pos.X-d.Hovered.Transform.Pos.X) < world.TileSize*DigRange && math.Abs(d.Transform.Pos.Y-d.Hovered.Transform.Pos.Y) < world.TileSize*DigRange
			if (in.Get("dig").JustPressed() && !in.OptFlags["DigOnRelease"]) || (in.Get("dig").JustReleased() && in.OptFlags["DigOnRelease"]) {
				if !d.Physics.Grounded && d.airDig < d.AirDigs {
					d.airDig++
				}
				facing := util.Cardinal(d.Transform.Pos, d.Hovered.Transform.Pos)
				if d.Hovered.Solid() && d.SelectLegal {
					d.tileQueue = append(d.tileQueue, struct {
						a int
						t *cave.Tile
						f pixel.Vec
					}{
						a: 0,
						t: d.Hovered,
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
			} else if ((in.Get("flag").JustPressed() && !in.OptFlags["DigOnRelease"]) || (in.Get("flag").JustReleased() && in.OptFlags["DigOnRelease"])) && d.Hovered.Solid() && d.SelectLegal {
				if !d.Physics.Grounded && d.airDig < d.AirDigs {
					d.airDig++
				}
				facing := util.Cardinal(d.Transform.Pos, d.Hovered.Transform.Pos)
				d.tileQueue = append(d.tileQueue, struct {
					a int
					t *cave.Tile
					f pixel.Vec
				}{
					a: 1,
					t: d.Hovered,
					f: facing,
				})
			}
		}
		if len(d.tileQueue) > 0 && !d.digging && !d.attacking && !d.flagging {
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
				d.flagging = true
				d.distFell = 0.
				FlagTile(d.Player, next.t)
			} else {
				var x, y float64
				if d.facing.Y < 0 {
					y = d.Transform.Pos.Y - world.TileSize*0.6
				} else if d.facing.Y > 0 {
					y = d.Transform.Pos.Y + world.TileSize*0.6
				} else {
					y = d.Transform.Pos.Y
				}
				if d.facing.X < 0 {
					x = d.Transform.Pos.X - world.TileSize*0.6
				} else if d.facing.X > 0 {
					x = d.Transform.Pos.X + world.TileSize*0.6
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
		d.digHold = in.OptFlags["DigOnRelease"] && in.Get("dig").Pressed() && !d.flagHold
		d.flagHold = in.OptFlags["DigOnRelease"] && in.Get("flag").Pressed() && !d.digHold
		if d.digHold || d.flagHold {
			if d.climbing {
				d.Physics.CancelMovement()
			}
			d.digging = false
			d.attacking = false
			d.jumping = false
			d.walking = false
			if d.Hovered.Transform.Pos.X < d.Transform.Pos.X {
				d.facing.X = -1
			} else {
				d.facing.X = 1
			}
		} else {
			if d.digging || d.attacking || d.flagging {
				d.Physics.CancelMovement()
			} else {
				canJump := d.Physics.NearGround || d.Physics.Grounded

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
					if !d.Physics.LeftBound {
						d.Player.CamVel.X = math.Max(d.Player.CamVel.X-10., -100.)
						cameraMoveX = true
					}
				case 2:
					if d.Physics.Grounded {
						d.faceLeft = false
						d.Physics.SetVelX(d.Speed, GroundAcceleration)
					} else {
						d.Physics.SetVelX(d.Speed, AirAcceleration)
					}
					if !d.Physics.RightBound {
						d.Player.CamVel.X = math.Min(d.Player.CamVel.X+10., 100.)
						cameraMoveX = true
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
					PlayStep(0.)
					d.distFell = 0.
					d.Physics.SetVelY(JumpVel, 0.)
				} else if d.climbing {
					if in.Get("jump").JustPressed() {
						d.climbing = false
					} else if d.Physics.WallBound {
						d.distFell = 0.
						if in.Get("up").Pressed() && !in.Get("down").Pressed() {
							d.Physics.SetVelY(d.ClimbSpeed, 0.)
						} else if in.Get("down").Pressed() && !in.Get("up").Pressed() {
							d.Physics.SetVelY(-d.ClimbSpeed, 0.)
						} else {
							d.Physics.SetVelY(0., 0.)
						}
						if d.Physics.RightBound && (!d.Physics.LeftBound || in.Get("right").Pressed()) {
							d.faceLeft = false
						} else if d.Physics.LeftBound && (!d.Physics.RightBound || in.Get("left").Pressed()) {
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
				} else if d.Physics.WallBound && !d.toJump && in.Get("up").Pressed() {
					d.climbing = true
					d.walking = false
					d.jumping = false
					d.toJump = false
					d.distFell = 0.
					d.Physics.SetVelY(d.ClimbSpeed, 0.)
					if d.Physics.RightBound && !d.Physics.LeftBound {
						d.faceLeft = false
					} else if d.Physics.LeftBound && !d.Physics.RightBound {
						d.faceLeft = true
					}
				} else if !d.jumping && !d.toJump && d.Physics.Grounded {
					if math.Abs(d.Physics.Velocity.X) < 20.0 {
						if in.Get("up").Pressed() && !in.Get("down").Pressed() {
							d.distFell = 0.
							d.Player.CamVel.Y = math.Min(d.Player.CamVel.Y + 10., 100.)
							cameraMoveY = true
						} else if in.Get("down").Pressed() && !in.Get("up").Pressed() {
							d.distFell = 0.
							d.Player.CamVel.Y = math.Max(d.Player.CamVel.Y - 10., -100.)
							cameraMoveY = true
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
				in.Get("prev").Consume()
				PrevItem(d.Player.Inventory)
			} else if in.Get("next").JustPressed() {
				in.Get("next").Consume()
				NextItem(d.Player.Inventory)
			} else if in.Get("use").JustPressed() {
				in.Get("use").Consume()
				UseEquipped(d.Player, d.Entity, d.Transform.Pos, d.Hovered.Transform.Pos)
				d.facing = util.Cardinal(d.Transform.Pos, d.Hovered.Transform.Pos)
				useItem = true
			} else if in.Get("interact").Pressed() {
				if in.Get("interact").JustPressed() {
					d.stopDrop = false
					d.dropTimer = timing.New(0.25)
				} else if in.Get("interact").Repeated() && d.dropTimer.UpdateDone() && !d.stopDrop {
					d.stopDrop = DropEquipped(d.Player.Inventory, d.Transform.Pos)
					d.dropTimer = timing.New(0.25)
				}
			}
		}
		d.Collider.Fallthrough = in.Get("down").Pressed() || d.climbing
	}
	if d.facingTimer != nil && d.facingTimer.UpdateDone() {
		if d.faceLeft {
			d.facing.X = -1
		} else {
			d.facing.X = 1
		}
		d.facing.Y = 0
	} else if d.facing.X < 0 && (useItem || d.digHold || d.flagHold || d.digging || d.flagging || d.attacking) {
		d.faceLeft = true
	} else if d.facing.X > 0 && (useItem || d.digHold || d.flagHold || d.digging || d.flagging || d.attacking) {
		d.faceLeft = false
	}
	if !cameraMoveX {
		if d.Player.CamVel.X > 0. {
			d.Player.CamVel.X -= 10.
			if d.Player.CamVel.X < 0. {
				d.Player.CamVel.X = 0.
			}
		} else if d.Player.CamVel.X < 0. {
			d.Player.CamVel.X += 10.
			if d.Player.CamVel.X > 0. {
				d.Player.CamVel.X = 0.
			}
		}
	}
	if !cameraMoveY {
		if d.Player.CamVel.Y > 0. {
			d.Player.CamVel.Y -= 10.
			if d.Player.CamVel.Y < 0. {
				d.Player.CamVel.Y = 0.
			}
		} else if d.Player.CamVel.Y < 0. {
			d.Player.CamVel.Y += 10.
			if d.Player.CamVel.Y > 0. {
				d.Player.CamVel.Y = 0.
			}
		}
	}
	d.Transform.Flip = d.faceLeft
}

func (d *Dwarf) Delete() {
	myecs.Manager.DisposeEntity(d.Entity)
}

func FlagTile(p *data.Player, tile *cave.Tile) {
	if tile != nil && tile.Solid() && !tile.Destroyed && tile.Breakable() {
		if !tile.Flagged {
			tile.Flagged = true
			CreateFlag(p, tile)
		} else {
			tile.Flagged = false
		}
	}
}

func Dig(tile *cave.Tile, p *data.Player) bool {
	if tile.Diggable() {
		if p != nil {
			profile.CurrentProfile.Stats.BlocksDug++
			p.Stats.BlocksDug++
		}
		tile.Destroy(p, true)
		return true
	}
	return false
}