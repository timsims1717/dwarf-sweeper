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
	angleDiff          = 0.15
)

var (
	ClimbSpeed      = 50.
	Speed           = 80.
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
	Collider   data.Collider
	Reanimator *reanimator.Tree
	Entity     *ecs.Entity
	Enchants   []string
	EnchantMax int
	faceLeft   bool

	selectLegal bool
	gpSelect    bool
	mouseSelect bool
	selectTimer *timing.FrameTimer
	angleTimer  *timing.FrameTimer

	hovered     *cave.Tile
	relative    pixel.Vec
	digTile     *cave.Tile
	tileQueue   []struct{
		a int
		t *cave.Tile
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

	digging   bool
	attacking bool
	marking   bool
	climbing  bool
	airDig    bool

	Health    *data.Health
	DeadStop  bool

	Bubble *Bubble
}

func NewDwarf(start pixel.Vec) *Dwarf {
	dwarfSheet, err := img.LoadSpriteSheet("assets/img/dwarf.json")
	if err != nil {
		panic(err)
	}
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
	d.Reanimator = reanimator.New(&reanimator.Switch{
		Elements: reanimator.NewElements(
			reanimator.NewSwitch(&reanimator.Switch{
				Elements: reanimator.NewElements(
					reanimator.NewSwitch(&reanimator.Switch{
						Elements: reanimator.NewElements(
								reanimator.NewAnimFromSheet("hit_front", dwarfSheet, []int{15}, reanimator.Hold, nil), // hit_front
								reanimator.NewAnimFromSheet("hit_back", dwarfSheet, []int{16}, reanimator.Hold, nil),  // hit_back
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
					reanimator.NewAnimFromSheet("flat", dwarfSheet, []int{17}, reanimator.Hold, nil), // flat
				),
				Check: func() int {
					if !d.Physics.Grounded {
						return 0
					} else {
						return 1
					}
				},
			}),
			reanimator.NewAnimFromSheet("dig", dwarfSheet, []int{11, 12, 13}, reanimator.Tran, map[int]func() {
				1: func() {
					if d.digTile != nil {
						if d.digTile.Solid() {
							CaveBlocksDug++
							d.digTile.Destroy(true)
						} else {
							var x, y float64
							if d.Transform.Pos.X > d.digTile.Transform.Pos.X {
								x = d.Transform.Pos.X - world.TileSize * 0.75
							} else {
								x = d.Transform.Pos.X + world.TileSize * 0.75
							}
							if d.Transform.Pos.Y > d.digTile.Transform.Pos.Y {
								y = d.Transform.Pos.Y - world.TileSize * 0.75
							} else {
								y = d.Transform.Pos.Y + world.TileSize * 0.75
							}
							myecs.Manager.NewEntity().
								AddComponent(myecs.AreaDmg, &data.AreaDamage{
									Center:    pixel.V(x, y),
									Radius:    world.TileSize * 1.5,
									Amount:    d.ShovelDamage,
									Dazed:     d.ShovelDazed,
									Knockback: d.ShovelKnockback,
									Source:    d.Transform.Pos,
									Type:      data.Shovel,
								})
						}
						sfx.SoundPlayer.PlaySound("shovel", 1.0)
						d.digTile = nil
					}
				},
				3: func() {
					d.digging = false
					d.attacking = false
				},
			}), // digging
			reanimator.NewAnimFromSheet("mark", dwarfSheet, []int{14}, reanimator.Tran, map[int]func() {
				1: func() {
					d.marking = false
				},
			}), // marking
			reanimator.NewSwitch(&reanimator.Switch{
				Elements: reanimator.NewElements(
					reanimator.NewSwitch(&reanimator.Switch{
						Elements: reanimator.NewElements(
							reanimator.NewAnimFromSheet("run", dwarfSheet, []int{4, 5, 6, 7}, reanimator.Loop, nil), // run
							reanimator.NewSwitch(&reanimator.Switch{
								Elements: reanimator.NewElements(
									reanimator.NewAnimFromSheet("flat", dwarfSheet, []int{17}, reanimator.Hold, nil),         // flat
									reanimator.NewAnimFromSheet("idle", dwarfSheet, []int{0, 1, 2, 3}, reanimator.Loop, nil), // idle
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
							reanimator.NewAnimFromSheet("climb_up", dwarfSheet, []int{18,19,20,21}, reanimator.Loop, nil),  // climb_up
							reanimator.NewAnimFromSheet("climb_dwn", dwarfSheet, []int{21,20,19,18}, reanimator.Loop, nil), // climb_dwn
							reanimator.NewAnimFromSheet("climb_still", dwarfSheet, []int{18}, reanimator.Hold, nil),        // climb_still
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
							reanimator.NewAnimFromSheet("jump", dwarfSheet, []int{8, 9}, reanimator.Hold, nil), // jump
							reanimator.NewAnimFromSheet("fall", dwarfSheet, []int{10}, reanimator.Hold, nil),   // fall
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
			} else if d.digging || d.attacking {
				return 1
			} else if d.marking {
				return 2
			} else {
				return 3
			}
		},
	}, "idle")
	d.Collider = data.Collider{
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
	loc1 := Descent.GetCave().GetTile(d.Transform.Pos)
	if d.Physics.Grounded || d.climbing || d.Bubble != nil {
		d.airDig = false
	}
	if d.Health.Dazed || d.Health.Dead {
		if d.Bubble != nil {
			d.Bubble.Pop()
		}
		d.tileQueue = []struct{
			a int
			t *cave.Tile
		}{}
		d.digging = false
		d.attacking = false
		d.jumping = false
		d.walking = false
		d.climbing = false
		if d.Health.Dead {
			if d.Physics.Grounded && !d.Physics.IsMovingX() && !d.DeadStop {
				d.Physics.CancelMovement()
				d.distFell = 150.
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
		jpSelecting := in.Axes["targetX"].F > 0. || in.Axes["targetX"].F < 0. || in.Axes["targetY"].F > 0. || in.Axes["targetY"].F < 0.
		moveSelecting := in.Get("left").Pressed() || in.Get("right").Pressed() || in.Get("up").Pressed() || in.Get("down").Pressed()
		jpMS := in.Get("left").JustPressed() || in.Get("right").JustPressed() || in.Get("up").JustPressed() || in.Get("down").JustPressed()
		if in.Get("left").JustReleased() || in.Get("right").JustReleased() || in.Get("up").JustReleased() || in.Get("down").JustReleased() {
			d.angleTimer = timing.New(AngleSec)
		}
		if jpSelecting && constants.DigMode != data.Movement {
			d.gpSelect = true
			d.mouseSelect = false
			d.selectTimer = timing.New(selectTimerSec)
		} else if in.MouseMoved && constants.DigMode != data.Movement {
			d.mouseSelect = true
			d.gpSelect = false
			d.selectTimer = timing.New(selectTimerSec)
		} else if d.selectTimer.UpdateDone() && constants.DigMode != data.Dedicated {
			d.gpSelect = false
			d.mouseSelect = false
		}
		if d.gpSelect {
			if jpSelecting {
				x := in.Axes["targetX"].R
				y := in.Axes["targetY"].R
				xA := math.Abs(x)
				yA := math.Abs(y)
				diff := math.Abs(x-y)
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
			}
			p := d.Transform.Pos
			p.X += d.relative.X
			p.Y += d.relative.Y
			d.hovered = Descent.GetCave().GetTile(p)
		} else if d.mouseSelect {
			d.hovered = Descent.GetCave().GetTile(in.World)
		} else if constants.DigMode != data.Dedicated {
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
			}
			p := d.Transform.Pos
			p.X += d.relative.X
			p.Y += d.relative.Y
			d.hovered = Descent.GetCave().GetTile(p)
		}
		if d.hovered != nil && !d.airDig && len(d.tileQueue) < 3 {
			d.selectLegal = math.Abs(d.Transform.Pos.X-d.hovered.Transform.Pos.X) < world.TileSize*DigRange && math.Abs(d.Transform.Pos.Y-d.hovered.Transform.Pos.Y) < world.TileSize*DigRange
			if in.Get("dig").JustPressed() && d.selectLegal {
				if d.hovered.Solid() {
					d.tileQueue = append(d.tileQueue, struct {
						a int
						t *cave.Tile
					}{
						a: 0,
						t: d.hovered,
					})
				} else {
					d.tileQueue = append(d.tileQueue, struct{
						a int
						t *cave.Tile
					}{
						a: 2,
						t: d.hovered,
					})
				}
			} else if in.Get("mark").JustPressed() && d.hovered.Solid() && d.selectLegal {
				d.tileQueue = append(d.tileQueue, struct{
					a int
					t *cave.Tile
				}{
					a: 1,
					t: d.hovered,
				})
			}
		}
		if len(d.tileQueue) > 0 && !d.digging && !d.attacking && !d.marking {
			next := d.tileQueue[0]
			d.tileQueue = d.tileQueue[1:]
			if next.t.Transform.Pos.X < d.Transform.Pos.X {
				d.faceLeft = true
			} else if next.t.Transform.Pos.X > d.Transform.Pos.X {
				d.faceLeft = false
			}
			if math.Abs(d.Transform.Pos.X-next.t.Transform.Pos.X) < world.TileSize*DigRange && math.Abs(d.Transform.Pos.Y-next.t.Transform.Pos.Y) < world.TileSize*DigRange && !cave.TileInTile(d.Transform.Pos, next.t.Transform.Pos) {
				if next.a == 0 && next.t.Solid() {
					d.digging = true
					d.attacking = false
					d.jumping = false
					d.walking = false
					d.climbing = false
					d.distFell = 0.
					d.digTile = next.t
				} else if next.a == 1 && next.t.Solid() {
					d.marking = true
					d.distFell = 0.
					Mark(next.t)
				} else if next.a == 2 {
					d.digging = false
					d.attacking = true
					d.jumping = false
					d.walking = false
					d.climbing = false
					d.distFell = 0.
					d.digTile = next.t
				}
			}
		}
		if d.digging || d.attacking {
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
		} else if !d.marking {
			dwnlj := Descent.GetCave().GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.4, d.Transform.Pos.Y-world.TileSize))
			dwnrj := Descent.GetCave().GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.4, d.Transform.Pos.Y-world.TileSize))
			dwn1 := Descent.GetCave().GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize))
			dwn2 := Descent.GetCave().GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize*1.25))
			right := Descent.GetCave().GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.6, d.Transform.Pos.Y-world.TileSize*0.48))
			left := Descent.GetCave().GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.6, d.Transform.Pos.Y))
			canJump := (dwn1 != nil && dwn1.Solid()) || (dwn2 != nil && dwn2.Solid()) || (dwnlj != nil && dwnlj.Solid()) || (dwnrj != nil && dwnrj.Solid())
			canClimb := (right != nil && right.Solid()) || (left != nil && left.Solid())

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
			if !d.jumping && loc1 != nil && canJump && in.Get("jump").JustPressed() {
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
				if canClimb {
					d.distFell = 0.
					if in.Get("up").Pressed() && !in.Get("down").Pressed() {
						d.Physics.SetVelY(d.ClimbSpeed, 0.)
					} else if in.Get("down").Pressed() && !in.Get("up").Pressed() {
						d.Physics.SetVelY(-d.ClimbSpeed, 0.)
					} else {
						d.Physics.SetVelY(0., 0.)
					}
					if right != nil && right.Solid() && (left == nil || !left.Solid()) {
						d.faceLeft = false
					} else if left != nil && left.Solid() && (right == nil || !right.Solid()) {
						d.faceLeft = true
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
				if right != nil && right.Solid() && (left == nil || !left.Solid()) {
					d.faceLeft = false
				} else if left != nil && left.Solid() && (right == nil || !right.Solid()) {
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
					height := int(((d.Transform.Pos.Y - d.jumpOrigY) + world.TileSize * 1.0) / world.TileSize)
					if d.Physics.Velocity.Y <= 0. {
						d.jumping = false
						d.jumpEnd = false
					} else if height < d.MaxJump - 1 && in.Get("jump").Pressed() {
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
			particles.CreateStaticParticle("target", d.hovered.Transform.Pos)
		} else {
			particles.CreateStaticParticle("target_blank", d.hovered.Transform.Pos)
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