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
	Speed = 75.
	JumpVel = 150.
	DigRange = 1.4
	MaxJump = 4
)

type Dwarf struct {
	Transform   *physics.Physics
	Reanimator  *reanimator.Tree
	entity      *ecs.Entity
	faceLeft    bool
	selectLegal bool
	walkTimer   time.Time
	walking     bool
	jumping     bool
	jumpOrigY   float64
	jumpHeight  int
	toJump      bool
	toJumpTimer time.Time
	digging     bool
	marking     bool
	climbing    bool
	digTile     *Tile
	hovered     *Tile
	distFell    float64
	cursorV     pixel.Vec
	relWorld    pixel.Vec
	Hurt        bool
	dmg         float64
	source      pixel.Vec
	knockback   float64
	DeadStop    bool
	Dead        bool
	Inv         bool
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
							if d.Transform.Velocity.Y > 0. {
								return 0
							} else {
								return 1
							}
						},
					}),
				),
				Check: func() int {
					if d.Transform.Grounded && !d.jumping {
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
	})
	d.entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, tran).
		AddComponent(myecs.Physics, physicsT).
		AddComponent(myecs.Collision, myecs.Collider{}).
		AddComponent(myecs.Animation, d.Reanimator)
	return d
}

func (d *Dwarf) Update() {
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
		d.hovered = Dungeon.GetCave().GetTile(input.Input.World)
		if d.hovered != nil {
			d.selectLegal = math.Abs(d.Transform.Pos.X-d.hovered.Transform.Pos.X) < world.TileSize*DigRange && math.Abs(d.Transform.Pos.Y-d.hovered.Transform.Pos.Y) < world.TileSize*DigRange
			if input.Input.IsDig && !d.digging && !d.marking && d.hovered.Solid && d.selectLegal {
				d.digging = true
				d.jumping = false
				d.walking = false
				d.climbing = false
				d.distFell = 0.
				d.digTile = d.hovered
				if d.digTile.Transform.Pos.X < d.Transform.Pos.X {
					d.faceLeft = true
				} else if d.digTile.Transform.Pos.X > d.Transform.Pos.X {
					d.faceLeft = false
				}
			} else if input.Input.IsMark && !d.digging && !d.marking && d.hovered.Solid && d.selectLegal {
				if d.hovered.Transform.Pos.X < d.Transform.Pos.X {
					d.faceLeft = true
				} else if d.hovered.Transform.Pos.X > d.Transform.Pos.X {
					d.faceLeft = false
				}
				d.marking = true
				d.distFell = 0.
				d.hovered.Mark(d.Transform.Pos)
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
			switch input.Input.XDir {
			case input.Left:
				if input.Input.XDirC || d.Transform.Velocity.X >= 0. {
					d.walkTimer = time.Now()
				}
				d.Transform.SetVelX(-Speed, 5.)
				if d.Transform.Grounded {
					d.faceLeft = true
				}
			case input.Right:
				if input.Input.XDirC || d.Transform.Velocity.X <= 0. {
					d.walkTimer = time.Now()
				}
				d.Transform.SetVelX(Speed, 5.)
				if d.Transform.Grounded {
					d.faceLeft = false
				}
			}
			input.Input.XDirC = false
			// Ground test, considered on the ground for jumping purposes until half a tile out
			if !d.jumping && loc1 != nil && canJump && input.Input.Jumping.JustPressed() {
				d.toJump = true
				d.climbing = false
				d.walking = false
				d.distFell = 0.
				d.toJumpTimer = time.Now()
			} else if d.toJump && time.Since(d.toJumpTimer).Seconds() > 0.05 {
				d.climbing = false
				d.toJump = false
				d.walking = false
				d.jumping = true
				d.jumpOrigY = d.Transform.Pos.Y
				d.jumpHeight = -1
				sfx.SoundPlayer.PlaySound(fmt.Sprintf("step%d", random.Effects.Intn(4)+1), 0.)
				d.distFell = 0.
				d.Transform.SetVelY(JumpVel, 0.)
			} else if d.climbing {
				if canClimb {
					d.distFell = 0.
					if input.Input.ClimbUp.Pressed() && !input.Input.ClimbDown.Pressed() {
						d.Transform.SetVelY(ClimbSpeed, 0.)
					} else if input.Input.ClimbDown.Pressed() && !input.Input.ClimbUp.Pressed() {
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
			} else if canClimb && input.Input.ClimbUp.Pressed() {
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
			} else if !d.jumping && d.Transform.Grounded {
				if math.Abs(d.Transform.Velocity.X) < 20.0 {
					if input.Input.LookUp.Pressed() && !input.Input.LookDown.Pressed() {
						d.distFell = 0.
						camera.Cam.Up()
					} else if input.Input.LookDown.Pressed() && !input.Input.LookUp.Pressed() {
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
				}
			} else {
				d.walking = false
				d.climbing = false
				if d.jumping {
					dist := int((d.Transform.Pos.Y - d.jumpOrigY) / world.TileSize)
					if ((dist < MaxJump - 2 && input.Input.Jumping.Pressed()) || dist == d.jumpHeight) && d.Transform.Velocity.Y > 0. {
						d.Transform.SetVelY(JumpVel, 0.)
						d.jumpHeight = dist
					} else {
						input.Input.Jumping.Consume()
						d.jumping = false
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

func (d *Dwarf) Draw(win *pixelgl.Window) {
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
		debug.AddText(fmt.Sprintf("world coords: (%d,%d)", int(input.Input.World.X), int(input.Input.World.Y)))
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