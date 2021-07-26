package dungeon

import (
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/input"
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
	"github.com/faiface/pixel/pixelgl"
	"math"
	"time"
)

var (
	ClimbSpeed = 50.
	Speed = 80.
	JumpVel = 150.
	DigRange = 1.4
	MaxJump = 4
	GroundAcceleration = 5.
	AirAcceleration = 10.
)

type Dwarf struct {
	Transform  *physics.Physics
	Reanimator *reanimator.Tree
	entity     *ecs.Entity
	faceLeft   bool

	selectLegal bool
	jpSelect    bool
	hovered     *Tile
	relative    pixel.Vec
	digTile     *Tile
	tileQueue   []struct{
		a int
		t *Tile
	}

	walkTimer time.Time
	walking   bool

	jumping    bool
	jumpOrigY  float64
	jumpTarget float64
	jumpTimer  time.Time
	toJump     bool
	jumpEnd    bool
	distFell   float64

	digging  bool
	marking  bool
	climbing bool

	Hurt      bool
	dmg       float64
	source    pixel.Vec
	knockback float64
	DeadStop  bool
	Dead      bool
	Inv       bool
}

func NewDwarf(start pixel.Vec) *Dwarf {
	dwarfSheet, err := img.LoadSpriteSheet("assets/img/dwarf.json")
	if err != nil {
		panic(err)
	}
	tran := transform.NewTransform()
	physicsT := &physics.Physics{
		Transform: tran,
	}
	tran.Pos = start
	d := &Dwarf{
		Transform:  physicsT,
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
								if d.Transform.Velocity.X > 0 {
									return 0
								} else {
									return 1
								}
							} else {
								if d.Transform.Velocity.X > 0 {
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
					if !d.Transform.Grounded {
						return 0
					} else {
						return 1
					}
				},
			}),
			reanimator.NewAnimFromSheet("dig", dwarfSheet, []int{11, 12, 13}, reanimator.Tran, map[int]func() {
				1: func() {
					BlocksDug += 1
					d.digTile.Destroy()
					sfx.SoundPlayer.PlaySound("shovel", 1.0)
				},
				3: func() {
					d.digging = false
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
							if d.Transform.IsMovingX() {
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
							if d.Transform.Velocity.Y > 0. {
								return 0
							} else if d.Transform.Velocity.Y < 0. {
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
							if d.Transform.Velocity.Y > 0. || d.jumping || d.toJump || d.jumpEnd {
								return 0
							} else {
								return 1
							}
						},
					}),
				),
				Check: func() int {
					if d.Transform.Grounded && !d.jumping && !d.toJump {
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
			if d.Hurt {
				return 0
			} else if d.digging {
				return 1
			} else if d.marking {
				return 2
			} else {
				return 3
			}
		},
	}, "idle")
	d.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, tran).
		AddComponent(myecs.Physics, physicsT).
		AddComponent(myecs.Collision, myecs.Collider{}).
		AddComponent(myecs.Animation, d.Reanimator)
	return d
}

func (d *Dwarf) Update(in *input.Input) {
	loc1 := Dungeon.GetCave().GetTile(d.Transform.Pos)
	if d.Hurt {
		d.Dead = true
		if d.dmg > 0 {
			// todo: damage to health
			if d.knockback > 0.1 {
				d.Transform.CancelMovement()
				dir := util.Normalize(d.Transform.Pos.Sub(d.source))
				d.Transform.SetVelX(dir.X * d.knockback, 0.)
				d.Transform.SetVelY(dir.Y * d.knockback, 0.)
			}
			d.knockback = 0.
			d.dmg = 0
			d.Transform.RicochetX = true
		} else if d.Transform.Grounded && !d.Transform.IsMovingX() {
			d.Transform.CancelMovement()
			d.distFell = 150.
			d.DeadStop = true
		}
		d.digging = false
		d.jumping = false
		d.walking = false
		d.climbing = false
	} else {
		jpSelecting := in.Axes["targetX"].F > 0. || in.Axes["targetX"].F < 0. || in.Axes["targetY"].F > 0. || in.Axes["targetY"].F < 0.
		if jpSelecting {
			d.jpSelect = true
		} else if in.MouseMoved {
			d.jpSelect = false
		}
		if d.jpSelect {
			if jpSelecting {
				x := in.Axes["targetX"].F
				y := in.Axes["targetY"].F
				if x > input.Deadzone || x < -input.Deadzone {
					if x > 0. {
						x = world.TileSize
					} else {
						x = -world.TileSize
					}
				}
				if y > input.Deadzone || y < -input.Deadzone {
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
			d.hovered = Dungeon.GetCave().GetTile(p)
		} else if !d.jpSelect {
			d.hovered = Dungeon.GetCave().GetTile(in.World)
		}
		if d.hovered != nil {
			d.selectLegal = math.Abs(d.Transform.Pos.X-d.hovered.Transform.Pos.X) < world.TileSize*DigRange && math.Abs(d.Transform.Pos.Y-d.hovered.Transform.Pos.Y) < world.TileSize*DigRange
			if in.Get("dig").JustPressed() && d.hovered.Solid && d.selectLegal {
				d.tileQueue = append(d.tileQueue, struct{
					a int
					t *Tile
				}{
					a: 0,
					t: d.hovered,
				})
			} else if in.Get("mark").JustPressed() && d.hovered.Solid && d.selectLegal {
				d.tileQueue = append(d.tileQueue, struct{
					a int
					t *Tile
				}{
					a: 1,
					t: d.hovered,
				})
			}
		}
		if len(d.tileQueue) > 0 && !d.digging && !d.marking {
			next := d.tileQueue[0]
			d.tileQueue = d.tileQueue[1:]
			if math.Abs(d.Transform.Pos.X-next.t.Transform.Pos.X) < world.TileSize*DigRange && math.Abs(d.Transform.Pos.Y-next.t.Transform.Pos.Y) < world.TileSize*DigRange && next.t.Solid {
				if next.a == 0 {
					d.digging = true
					d.jumping = false
					d.walking = false
					d.climbing = false
					d.distFell = 0.
					d.digTile = next.t
					if d.digTile.Transform.Pos.X < d.Transform.Pos.X {
						d.faceLeft = true
					} else if d.digTile.Transform.Pos.X > d.Transform.Pos.X {
						d.faceLeft = false
					}
				} else if next.a == 1 {
					if next.t.Transform.Pos.X < d.Transform.Pos.X {
						d.faceLeft = true
					} else if next.t.Transform.Pos.X > d.Transform.Pos.X {
						d.faceLeft = false
					}
					d.marking = true
					d.distFell = 0.
					next.t.Mark(d.Transform.Pos)
				}
			}
		}
		if d.digging {
			d.Transform.Velocity = pixel.ZV
		} else if !d.marking {
			dwnlj := Dungeon.GetCave().GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.4, d.Transform.Pos.Y-world.TileSize))
			dwnrj := Dungeon.GetCave().GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.4, d.Transform.Pos.Y-world.TileSize))
			dwn1 := Dungeon.GetCave().GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize))
			dwn2 := Dungeon.GetCave().GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize*1.5))
			right := Dungeon.GetCave().GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.6, d.Transform.Pos.Y-world.TileSize*0.48))
			left := Dungeon.GetCave().GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.6, d.Transform.Pos.Y))
			canJump := (dwn1 != nil && dwn1.Solid) || (dwn2 != nil && dwn2.Solid) || (dwnlj != nil && dwnlj.Solid) || (dwnrj != nil && dwnrj.Solid)
			canClimb := (right != nil && right.Solid) || (left != nil && left.Solid)

			xDir := 0
			if in.Get("left").Pressed() && !in.Get("right").Pressed() {
				xDir = 1
			} else if in.Get("right").Pressed() && !in.Get("left").Pressed() {
				xDir = 2
			}

			switch xDir {
			case 1:
				if d.Transform.Grounded {
					d.faceLeft = true
					d.Transform.SetVelX(-Speed, GroundAcceleration)
				} else {
					d.Transform.SetVelX(-Speed, AirAcceleration)
				}
			case 2:
				if d.Transform.Grounded {
					d.faceLeft = false
					d.Transform.SetVelX(Speed, GroundAcceleration)
				} else {
					d.Transform.SetVelX(Speed, AirAcceleration)
				}
			}
			// Ground test, considered on the ground for jumping purposes until half a tile out
			if !d.jumping && loc1 != nil && canJump && in.Get("jump").JustPressed() {
				d.toJump = true
				d.climbing = false
				d.walking = false
				d.distFell = 0.
				d.jumpTimer = time.Now()
			} else if d.toJump && time.Since(d.jumpTimer).Seconds() > 0.05 {
				d.climbing = false
				d.toJump = false
				d.walking = false
				d.jumping = true
				d.jumpOrigY = d.Transform.Pos.Y
				sfx.SoundPlayer.PlaySound(fmt.Sprintf("step%d", random.Effects.Intn(4)+1), 0.)
				d.distFell = 0.
				d.Transform.SetVelY(JumpVel, 0.)
			} else if d.climbing {
				if canClimb {
					d.distFell = 0.
					if in.Get("climbUp").Pressed() && !in.Get("climbDown").Pressed() {
						d.Transform.SetVelY(ClimbSpeed, 0.)
					} else if in.Get("climbDown").Pressed() && !in.Get("climbUp").Pressed() {
						d.Transform.SetVelY(-ClimbSpeed, 0.)
					} else {
						d.Transform.SetVelY(0., 0.)
					}
					if right != nil && right.Solid && (left == nil || !left.Solid) {
						d.faceLeft = false
					} else if left != nil && left.Solid && (right == nil || !right.Solid) {
						d.faceLeft = true
					}
				} else {
					d.climbing = false
				}
			} else if canClimb && !d.toJump && in.Get("climbUp").Pressed() {
				d.climbing = true
				d.walking = false
				d.jumping = false
				d.toJump = false
				d.distFell = 0.
				d.Transform.SetVelY(ClimbSpeed, 0.)
				if right != nil && right.Solid && (left == nil || !left.Solid) {
					d.faceLeft = false
				} else if left != nil && left.Solid && (right == nil || !right.Solid) {
					d.faceLeft = true
				}
			} else if !d.jumping && !d.toJump && d.Transform.Grounded {
				wasWalking := d.walking
				if math.Abs(d.Transform.Velocity.X) < 20.0 {
					if in.Get("lookUp").Pressed() && !in.Get("lookDown").Pressed() {
						d.distFell = 0.
						camera.Cam.Up()
					} else if in.Get("lookDown").Pressed() && !in.Get("lookUp").Pressed() {
						d.distFell = 0.
						camera.Cam.Down()
					}
					d.walking = false
					d.climbing = false
				} else if d.Transform.Velocity.X > 0. {
					d.faceLeft = false
					d.walking = true
				} else if d.Transform.Velocity.X < 0. {
					d.faceLeft = true
					d.walking = true
				}
				if d.walking {
					d.distFell = 0.
					if !wasWalking {
						d.walkTimer = time.Now()
					}
				}
			} else {
				d.walking = false
				d.climbing = false
				if d.jumping || d.jumpEnd {
					height := int(((d.Transform.Pos.Y - d.jumpOrigY) + world.TileSize * 1.0) / world.TileSize)
					if d.Transform.Velocity.Y <= 0. {
						d.jumping = false
						d.jumpEnd = false
					} else if height < MaxJump - 1 && in.Get("jump").Pressed() {
						d.Transform.SetVelY(JumpVel, 0.)
					} else if !d.jumpEnd {
						in.Get("jump").Consume()
						d.jumping = false
						d.jumpEnd = true
						d.jumpTarget = Dungeon.GetCave().GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y+world.TileSize*1.0)).Transform.Pos.Y
					}
					if d.jumpEnd {
						if d.jumpTarget > d.Transform.Pos.Y {
							d.Transform.SetVelY(0., 0.5)
						} else {
							d.Transform.SetVelY(0., 0.)
							d.jumpEnd = false
						}
					}
				}
				if d.Transform.Velocity.Y < 0. {
					d.distFell += math.Abs(d.Transform.Velocity.Y * timing.DT)
				}
			}
		}
	}
	d.Transform.Flip = d.faceLeft
	camera.Cam.StayWithin(d.Transform.Pos, world.TileSize * 1.5)
	currLevel := int(-d.Transform.Pos.Y / world.TileSize)
	if LowestLevel < currLevel && !d.Hurt {
		LowestLevel = currLevel
	}
}

func (d *Dwarf) Draw(win *pixelgl.Window, in *input.Input) {
	d.Reanimator.CurrentSprite().Draw(win, d.Transform.Mat)
	if d.walking && time.Since(d.walkTimer).Seconds() > 0.4 {
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("step%d", random.Effects.Intn(4) + 1), 0.)
		d.walkTimer = time.Now()
	}
	if d.hovered != nil && !d.Hurt {
		if d.hovered.Solid && d.selectLegal {
			particles.CreateStaticParticle("target", d.hovered.Transform.Pos)
		} else {
			particles.CreateStaticParticle("target_blank", d.hovered.Transform.Pos)
		}
	}
	if debug.Debug {
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
		debug.AddText(fmt.Sprintf("dwarf velocity: (%d,%d)", int(d.Transform.Velocity.X), int(d.Transform.Velocity.Y)))
		debug.AddText(fmt.Sprintf("dwarf moving?: (%t,%t)", d.Transform.IsMovingX(), d.Transform.IsMovingY()))
		//debug.AddText(fmt.Sprintf("jump pressed?: %t", input.Input.Jumping.Pressed()))
		debug.AddText(fmt.Sprintf("dwarf grounded?: %t", d.Transform.Grounded))
		debug.AddText(fmt.Sprintf("tile queue len: %d", len(d.tileQueue)))
	}
}

func (d *Dwarf) Damage(dmg float64, source pixel.Vec, knockback float64) {
	if dmg > 0 && !d.Inv {
		d.Hurt = true
		d.dmg = dmg
		d.source = source
		d.knockback = knockback
	}
}

func (d *Dwarf) Delete() {
	myecs.Manager.DisposeEntity(d.entity)
}