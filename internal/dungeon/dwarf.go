package dungeon

import (
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
	"time"
)

var (
	ClimbSpeed = 50.
	Speed = 75.
	JumpVel = 150.
	DigRange = 1.4
	MaxJump = 4
)

var Player1 *Dwarf

type Dwarf struct {
	Transform   *physics.Physics
	Reanimator  *reanimator.Tree
	faceLeft    bool
	selectLegal bool
	walkTimer   time.Time
	walking     bool
	jumping     bool
	jumpOrigY   float64
	jumpHeight  int
	toJump      bool
	toJumpTimer time.Time
	grounded    bool
	digging     bool
	marking     bool
	climbing    bool
	selected    *Tile
	distFell    float64
	cursorV     pixel.Vec
	relWorld    pixel.Vec
	Hurt        bool
	dmg         float64
	source      pixel.Vec
	knockback   float64
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
								reanimator.NewAnim("hit_front", dwarfSheet, []int{15}, reanimator.Hold, nil), // hit_front
								reanimator.NewAnim("hit_back", dwarfSheet, []int{16}, reanimator.Hold, nil), // hit_back
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
					reanimator.NewAnim("flat", dwarfSheet, []int{17}, reanimator.Hold, nil), // flat
				),
				Check: func() int {
					if !d.grounded || d.Transform.IsMovingX() {
						return 0
					} else {
						return 1
					}
				},
			}),
			reanimator.NewAnim("dig", dwarfSheet, []int{11, 12, 13}, reanimator.Tran, func() {
				d.digging = false
			}), // digging
			reanimator.NewAnim("mark", dwarfSheet, []int{14}, reanimator.Tran, func() {
				d.marking = false
			}), // marking
			reanimator.NewSwitch(&reanimator.Switch{
				Elements: reanimator.NewElements(
					reanimator.NewSwitch(&reanimator.Switch{
						Elements: reanimator.NewElements(
							reanimator.NewAnim("run", dwarfSheet, []int{4, 5, 6, 7}, reanimator.Loop, nil), // run
							reanimator.NewSwitch(&reanimator.Switch{
								Elements: reanimator.NewElements(
									reanimator.NewAnim("flat", dwarfSheet, []int{17}, reanimator.Hold, nil), // flat
									reanimator.NewAnim("idle", dwarfSheet, []int{0, 1, 2, 3}, reanimator.Loop, nil), // idle
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
							reanimator.NewAnim("climb_up", dwarfSheet, []int{18,19,20,21}, reanimator.Loop, nil), // climb_up
							reanimator.NewAnim("climb_dwn", dwarfSheet, []int{21,20,19,18}, reanimator.Loop, nil), // climb_dwn
							reanimator.NewAnim("climb_still", dwarfSheet, []int{18}, reanimator.Hold, nil), // climb_still
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
							reanimator.NewAnim("jump", dwarfSheet, []int{8, 9}, reanimator.Hold, nil), // jump
							reanimator.NewAnim("fall", dwarfSheet, []int{10}, reanimator.Hold, nil), // fall
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
					if d.grounded && !d.jumping {
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
	myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, tran).
		AddComponent(myecs.Physics, physicsT).
		AddComponent(myecs.Collision, myecs.Collider{}).
		AddComponent(myecs.Animation, d.Reanimator)
	return d
}

func (d *Dwarf) Update() {
	loc1 := CurrCave.GetTile(d.Transform.Pos)
	dwn1 := CurrCave.GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize))
	dwnlw := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.3, d.Transform.Pos.Y-world.TileSize))
	dwnrw := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.3, d.Transform.Pos.Y-world.TileSize))
	d.grounded = ((dwn1 != nil && dwn1.Solid) || (dwnlw != nil && dwnlw.Solid) || (dwnrw != nil && dwnrw.Solid)) && (loc1 != nil && d.Transform.Pos.Y <= loc1.Transform.Pos.Y+1.0)
	if d.Hurt {
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
		} else if d.grounded && !d.Transform.IsMovingX() {
			d.Transform.CancelMovement()
			d.distFell = 150.
			d.Dead = true
		}
		d.digging = false
		d.jumping = false
		d.walking = false
		d.climbing = false
	} else {
		d.selected = CurrCave.GetTile(input.Input.World)
		if d.selected != nil {
			d.selectLegal = math.Abs(d.Transform.Pos.X-d.selected.Transform.Pos.X) < world.TileSize*DigRange && math.Abs(d.Transform.Pos.Y-d.selected.Transform.Pos.Y) < world.TileSize*DigRange
			if input.Input.IsDig && !d.digging && !d.marking && d.selected.Solid && d.selectLegal {
				d.digging = true
				d.jumping = false
				d.walking = false
				d.climbing = false
				if d.selected.Transform.Pos.X < d.Transform.Pos.X {
					d.faceLeft = true
				} else if d.selected.Transform.Pos.X > d.Transform.Pos.X {
					d.faceLeft = false
				}
				BlocksDug += 1
				d.selected.Destroy()
				sfx.SoundPlayer.PlaySound("shovel", 1.0)
			} else if input.Input.IsMark && !d.digging && !d.marking && d.selected.Solid && d.selectLegal {
				if d.selected.Transform.Pos.X < d.Transform.Pos.X {
					d.faceLeft = true
				} else if d.selected.Transform.Pos.X > d.Transform.Pos.X {
					d.faceLeft = false
				}
				d.marking = true
				d.selected.Mark(d.Transform.Pos)
			}
		}
		if d.digging {
			d.Transform.Velocity = pixel.ZV
		} else if !d.marking {
			dwnlj := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.5, d.Transform.Pos.Y-world.TileSize))
			dwnrj := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.5, d.Transform.Pos.Y-world.TileSize))
			dwn2 := CurrCave.GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize*1.5))
			right := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.6, d.Transform.Pos.Y-world.TileSize*0.48))
			left := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.6, d.Transform.Pos.Y))
			canJump := (dwn1 != nil && dwn1.Solid) || (dwn2 != nil && dwn2.Solid) || (dwnlj != nil && dwnlj.Solid) || (dwnrj != nil && dwnrj.Solid)
			canClimb := (right != nil && right.Solid) || (left != nil && left.Solid)
			switch input.Input.XDir {
			case input.Left:
				if input.Input.XDirC || d.Transform.Velocity.X >= 0. {
					d.walkTimer = time.Now()
				}
				d.Transform.SetVelX(-Speed, 5.)
				if d.grounded {
					d.faceLeft = true
				}
			case input.Right:
				if input.Input.XDirC || d.Transform.Velocity.X <= 0. {
					d.walkTimer = time.Now()
				}
				d.Transform.SetVelX(Speed, 5.)
				if d.grounded {
					d.faceLeft = false
				}
			}
			input.Input.XDirC = false
			// Ground test, considered on the ground for jumping purposes until half a tile out
			if !d.jumping && loc1 != nil && canJump && input.Input.Jumping.JustPressed() {
				d.toJump = true
				d.climbing = false
				d.toJumpTimer = time.Now()
			} else if d.toJump && time.Since(d.toJumpTimer).Seconds() > 0.05 {
				d.climbing = false
				d.toJump = false
				d.walking = false
				d.jumping = true
				d.jumpOrigY = d.Transform.Pos.Y
				d.jumpHeight = -1
				sfx.SoundPlayer.PlaySound(fmt.Sprintf("step%d", rand.Intn(4)+1), 0.)
				d.distFell = 0.
				d.Transform.SetVelY(JumpVel, 0.)
			} else if d.climbing {
				if canClimb {
					if input.Input.ClimbUp.Pressed() && !input.Input.ClimbDown.Pressed() {
						d.Transform.SetVelY(ClimbSpeed, 0.)
					} else if input.Input.ClimbDown.Pressed() && !input.Input.ClimbUp.Pressed() {
						d.Transform.SetVelY(-ClimbSpeed, 0.)
					} else {
						d.Transform.SetVelY(0., 0.)
					}
				} else {
					d.climbing = false
				}
			} else if canClimb && input.Input.ClimbUp.Pressed() {
				d.climbing = true
				d.Transform.SetVelY(ClimbSpeed, 0.)
				if right != nil && right.Solid && (left == nil || !left.Solid) {
					d.faceLeft = false
				} else if left != nil && left.Solid && (right == nil || !right.Solid) {
					d.faceLeft = true
				}
			} else if !d.jumping && d.grounded {
				if math.Abs(d.Transform.Velocity.X) < 20.0 {
					if input.Input.LookUp.Pressed() && !input.Input.LookDown.Pressed() {
						camera.Cam.Up()
					} else if input.Input.LookDown.Pressed() && !input.Input.LookUp.Pressed() {
						camera.Cam.Down()
					}
					d.walking = false
					d.climbing = false
				} else if d.Transform.Velocity.X > 0. {
					d.faceLeft = false
					d.walking = true
					d.climbing = false
				} else if d.Transform.Velocity.X < 0. {
					d.faceLeft = true
					d.walking = true
					d.climbing = false
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
	debug.AddLine(colornames.White, imdraw.RoundEndShape, d.Transform.Pos, d.Transform.Pos, 2.0)
	if d.selected != nil {
		debug.AddLine(colornames.Yellow, imdraw.RoundEndShape, d.selected.Transform.Pos, d.selected.Transform.Pos, 3.0)
	}
	currLevel := int(-d.Transform.Pos.Y / world.TileSize)
	if LowestLevel < currLevel && !d.Hurt {
		LowestLevel = currLevel
	}
}

func (d *Dwarf) Draw(win *pixelgl.Window) {
	d.Reanimator.CurrentSprite().Draw(win, d.Transform.Mat)
	if d.walking && time.Since(d.walkTimer).Seconds() > 0.4 {
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("step%d", rand.Intn(4) + 1), 0.)
		d.walkTimer = time.Now()
	}
	if d.selected != nil && !d.Hurt {
		if d.selected.Solid && d.selectLegal {
			particles.CreateStaticParticle("target", d.selected.Transform.Pos)
		} else {
			particles.CreateStaticParticle("target_blank", d.selected.Transform.Pos)
		}
	}
	if debug.Debug {
		if d.selected != nil {
			debug.AddText(fmt.Sprintf("world coords: (%d,%d)", int(input.Input.World.X), int(input.Input.World.Y)))
			debug.AddText(fmt.Sprintf("chunk coords: (%d,%d)", d.selected.Chunk.Coords.X, d.selected.Chunk.Coords.Y))
			debug.AddText(fmt.Sprintf("tile coords: (%d,%d)", d.selected.SubCoords.X, d.selected.SubCoords.Y))
			debug.AddText(fmt.Sprintf("tile type: '%s'", d.selected.Type))
			debug.AddText(fmt.Sprintf("tile sprite: '%s'", d.selected.BGSpriteS))
		}
		debug.AddText(fmt.Sprintf("dwarf position: (%d,%d)", int(d.Transform.APos.X), int(d.Transform.APos.Y)))
		debug.AddText(fmt.Sprintf("dwarf velocity: (%d,%d)", int(d.Transform.Velocity.X), int(d.Transform.Velocity.Y)))
		debug.AddText(fmt.Sprintf("dwarf moving?: (%t,%t)", d.Transform.IsMovingX(), d.Transform.IsMovingY()))
		debug.AddText(fmt.Sprintf("jump pressed?: %t", input.Input.Jumping.Pressed()))
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