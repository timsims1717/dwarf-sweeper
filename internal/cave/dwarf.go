package cave

import (
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/animation"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
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

const (
	Speed = 75.
	JumpVel = 225.
)

var Player1 *Dwarf

type Dwarf struct {
	Transform   *physics.Physics
	Animations  map[string]*animation.AnimationInstance
	currAnim    string
	faceLeft    bool
	selectLegal bool
	walkTimer   time.Time
	walking     bool
	toJumpT     time.Time
	toJump      bool
	jumping     bool
	digging     bool
	marking     bool
	selected    *Tile
	distFell    float64
	cursorV     pixel.Vec
	relWorld    pixel.Vec
	hurt        bool
	dmg         float64
	source      pixel.Vec
	Dead        bool
}

func NewDwarf() *Dwarf {
	dwarfSheet, err := img.LoadSpriteSheet("assets/img/dwarf.json")
	if err != nil {
		panic(err)
	}
	idle := animation.NewAnimation(dwarfSheet, 0, 2,true, false, 0.5)
	run := animation.NewAnimation(dwarfSheet, 2, 6,true, false, 0.4)
	jump := animation.NewAnimation(dwarfSheet, 6, 8,false, true, 0.2)
	fall := animation.NewAnimation(dwarfSheet, 8, 9,false, true, 0.1)
	dig := animation.NewAnimation(dwarfSheet, 9, 12,false, true, 0.4)
	flag := animation.NewAnimation(dwarfSheet, 12, 13,false, true, 0.2)
	hitfront := animation.NewAnimation(dwarfSheet, 13, 14,false, true, 0.1)
	hitback := animation.NewAnimation(dwarfSheet, 14, 15,false, true, 0.1)
	flat := animation.NewAnimation(dwarfSheet, 15, 16,false, true, 0.1)
	animations := make(map[string]*animation.AnimationInstance)
	animations["idle"] = idle.NewInstance()
	animations["run"] = run.NewInstance()
	animations["jump"] = jump.NewInstance()
	animations["fall"] = fall.NewInstance()
	animations["dig"] = dig.NewInstance()
	animations["flag"] = flag.NewInstance()
	animations["hit-front"] = hitfront.NewInstance()
	animations["hit-back"] = hitback.NewInstance()
	animations["flat"] = flat.NewInstance()
	transform := animation.NewTransform(true)
	physicsT := &physics.Physics{
		Transform: transform,
	}
	physicsT.Pos = pixel.V(16 * world.TileSize, -8 * world.TileSize)
	return &Dwarf{
		Transform:  physicsT,
		Animations: animations,
		currAnim:   "idle",
	}
}

func (d *Dwarf) Update() {
	newAnim := d.currAnim
	if d.hurt {
		if d.dmg > 0 {
			d.Transform.CancelMovement()
			dir := d.Transform.Pos.Sub(d.source)
			d.Transform.Velocity.X = dir.X * d.dmg * 10
			d.Transform.Velocity.Y = dir.Y * d.dmg * 10
			d.dmg = 0
		}
		d.digging = false
		d.toJump = false
		d.jumping = false
		d.walking = false
		if d.faceLeft {
			if d.Transform.Velocity.X > 0 {
				newAnim = "hit-front"
			} else {
				newAnim = "hit-back"
			}
		} else {
			if d.Transform.Velocity.X > 0 {
				newAnim = "hit-back"
			} else {
				newAnim = "hit-front"
			}
		}
		loc := CurrCave.GetTile(d.Transform.Pos)
		dwn := CurrCave.GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize))
		dwnl := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.40, d.Transform.Pos.Y-world.TileSize))
		dwnr := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.40, d.Transform.Pos.Y-world.TileSize))
		if ((dwn != nil && dwn.Solid) || (dwnr != nil && dwnr.Solid) || (dwnl != nil && dwnl.Solid)) && d.Transform.Pos.Y <= loc.Transform.Pos.Y && d.Transform.Velocity.Y < 0. && d.Transform.Velocity.X < 200. && d.Transform.Velocity.X > -200. {
			d.Transform.CancelMovement()
			d.distFell = 150.
			newAnim = "flat"
			d.Dead = true
		}
	} else {
		d.selected = CurrCave.GetTile(input.Input.World)
		if d.selected != nil {
			d.selectLegal = math.Abs(d.Transform.Pos.X-d.selected.Transform.Pos.X) < world.TileSize*1.3 && math.Abs(d.Transform.Pos.Y-d.selected.Transform.Pos.Y) < world.TileSize*1.3
			if input.Input.IsDig && !d.digging && !d.marking && d.selected.Solid && d.selectLegal {
				d.digging = true
				if d.selected.Transform.Pos.X < d.Transform.Pos.X {
					d.faceLeft = true
				} else if d.selected.Transform.Pos.X > d.Transform.Pos.X {
					d.faceLeft = false
				}
				newAnim = "dig"
				BlocksDug += 1
				d.selected.Destroy()
				sfx.SoundPlayer.PlaySound("shovel", 1.0)
			} else if input.Input.IsMark && !d.digging && !d.marking && d.selected.Solid && d.selectLegal {
				if d.selected.Transform.Pos.X < d.Transform.Pos.X {
					d.faceLeft = true
				} else if d.selected.Transform.Pos.X > d.Transform.Pos.X {
					d.faceLeft = false
				}
				newAnim = "flag"
				d.marking = true
				d.selected.Mark(d.Transform.Pos)
			}
		}
		if d.digging {
			d.Transform.Velocity = pixel.ZV
		} else if !d.marking {
			switch input.Input.XDir {
			case input.Left:
				if input.Input.XDirC || d.Transform.Velocity.X >= 0. {
					d.Transform.SetVelX(-Speed, 0.1) // Walk speed
					d.walkTimer = time.Now()
				}
				d.walking = true
			case input.Right:
				if input.Input.XDirC || d.Transform.Velocity.X <= 0. {
					d.Transform.SetVelX(Speed, 0.1)
					d.walkTimer = time.Now()
				}
				d.walking = true
			case input.None:
				if input.Input.XDirC {
					d.Transform.SetVelX(0., 0.1)
				}
				d.walking = false
			}
			input.Input.XDirC = false
			// Ground test, considered on the ground for jumping purposes until half a tile out
			loc1 := CurrCave.GetTile(d.Transform.Pos)
			dwn1 := CurrCave.GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize))
			dwnl := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.5, d.Transform.Pos.Y-world.TileSize))
			dwnr := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.5, d.Transform.Pos.Y-world.TileSize))
			if !d.jumping && loc1 != nil && ((dwn1 != nil && dwn1.Solid) || (dwnl != nil && dwnl.Solid) || (dwnr != nil && dwnr.Solid)) && d.Transform.Pos.Y < loc1.Transform.Pos.Y+1.0 {
				if input.Input.Jumping && d.toJump {
					s := time.Since(d.toJumpT).Seconds()
					if s >= 0.1 {
						sfx.SoundPlayer.PlaySound(fmt.Sprintf("step%d", rand.Intn(4) + 1), 0.)
						d.Transform.Velocity.Y = JumpVel
						d.toJump = false
						d.jumping = true
						d.walking = false
					}
				} else {
					if input.Input.Jumped {
						newAnim = "jump"
						d.toJump = true
						d.toJumpT = time.Now()
					} else {
						if math.Abs(d.Transform.Velocity.X) < 20.0 {
							if d.distFell > 100. {
								newAnim = "flat"
							} else {
								newAnim = "idle"
							}
						} else if d.Transform.Velocity.X > 0. {
							newAnim = "run"
							d.faceLeft = false
						} else if d.Transform.Velocity.X < 0. {
							newAnim = "run"
							d.faceLeft = true
						}
					}
				}
				if newAnim != "flat" {
					d.distFell = 0.
				}
			} else {
				d.walking = false
				if d.jumping {
					s := time.Since(d.toJumpT).Seconds()
					if s >= 0.1 && s <= 0.2 {
						d.Transform.Velocity.Y = JumpVel
						//} else if s >= 0.2 && s <= 0.4 && input.Input.Jumping {
						//	d.Transform.Velocity.Y = 200.
					} else {
						d.jumping = false
					}
				}
				if d.Transform.Velocity.Y <= 0. {
					newAnim = "fall"
					if d.currAnim != newAnim {
						d.distFell = 0.
					}
					d.distFell += math.Abs(d.Transform.Velocity.Y * timing.DT)
				}
			}
		}
	}
	d.Transform.Flip = d.faceLeft
	if newAnim != d.currAnim {
		d.Animations[d.currAnim].Reset()
		d.currAnim = newAnim
	}
	d.Transform.Update()
	loc := CurrCave.GetTile(d.Transform.Pos)
	if loc != nil {
		if math.Abs(loc.Transform.Pos.X-d.Transform.Pos.X) > world.TileSize || math.Abs(loc.Transform.Pos.Y-d.Transform.Pos.Y) > world.TileSize {
			fmt.Println("Time to teleport")
		}
		up := CurrCave.GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y+world.TileSize))
		upl := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.35, d.Transform.Pos.Y+world.TileSize))
		upr := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.35, d.Transform.Pos.Y+world.TileSize))
		dwn := CurrCave.GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize))
		dwnl := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.35, d.Transform.Pos.Y-world.TileSize))
		dwnr := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.35, d.Transform.Pos.Y-world.TileSize))
		right := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize, d.Transform.Pos.Y))
		left := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize, d.Transform.Pos.Y))
		if up != nil {
			debug.AddLine(colornames.Green, imdraw.RoundEndShape, up.Transform.Pos, up.Transform.Pos, 2.0)
		}
		if upl != nil {
			debug.AddLine(colornames.Green, imdraw.RoundEndShape, upl.Transform.Pos, upl.Transform.Pos, 2.0)
		}
		if upr != nil {
			debug.AddLine(colornames.Green, imdraw.RoundEndShape, upr.Transform.Pos, upr.Transform.Pos, 2.0)
		}
		if dwn != nil {
			debug.AddLine(colornames.Green, imdraw.RoundEndShape, dwn.Transform.Pos, dwn.Transform.Pos, 2.0)
		}
		if dwnl != nil {
			debug.AddLine(colornames.Green, imdraw.RoundEndShape, dwnl.Transform.Pos, dwnl.Transform.Pos, 2.0)
		}
		if dwnr != nil {
			debug.AddLine(colornames.Green, imdraw.RoundEndShape, dwnr.Transform.Pos, dwnr.Transform.Pos, 2.0)
		}
		if right != nil {
			debug.AddLine(colornames.Green, imdraw.RoundEndShape, right.Transform.Pos, right.Transform.Pos, 2.0)
		}
		if left != nil {
			debug.AddLine(colornames.Green, imdraw.RoundEndShape, left.Transform.Pos, left.Transform.Pos, 2.0)
		}
		if ((up != nil && up.Solid) || (upl != nil && upl.Solid) || (upr != nil && upr.Solid)) && d.Transform.Pos.Y > loc.Transform.Pos.Y {
			d.Transform.Pos.Y = loc.Transform.Pos.Y
			if d.Transform.Velocity.Y > 0 {
				d.Transform.Velocity.Y = 0
			}
		}
		if ((dwn != nil && dwn.Solid) || (dwnr != nil && dwnr.Solid) || (dwnl != nil && dwnl.Solid)) && d.Transform.Pos.Y < loc.Transform.Pos.Y {
			d.Transform.Pos.Y = loc.Transform.Pos.Y
			if d.Transform.Velocity.Y < 0 {
				d.Transform.Velocity.Y = 0
			}
		}
		if right != nil && right.Solid && d.Transform.Pos.X > loc.Transform.Pos.X {
			d.Transform.Pos.X = loc.Transform.Pos.X
			if d.Transform.Velocity.X > 0 {
				if d.hurt {
					d.Transform.Velocity.X = d.Transform.Velocity.X * -0.6
				} else {
					d.Transform.Velocity.X = 0
				}
			}
		}
		if left != nil && left.Solid && d.Transform.Pos.X < loc.Transform.Pos.X {
			d.Transform.Pos.X = loc.Transform.Pos.X
			if d.Transform.Velocity.X < 0 {
				if d.hurt {
					d.Transform.Velocity.X = d.Transform.Velocity.X * -0.6
				} else {
					d.Transform.Velocity.X = 0
				}
			}
		}
	}
	d.Animations[d.currAnim].Update()
	d.Animations[d.currAnim].SetMatrix(d.Transform.Mat)
	camera.Cam.Follow(d.Transform.Pos, 5.)
	debug.AddLine(colornames.White, imdraw.RoundEndShape, d.Transform.Pos, d.Transform.Pos, 2.0)
	if d.selected != nil {
		debug.AddLine(colornames.Yellow, imdraw.RoundEndShape, d.selected.Transform.Pos, d.selected.Transform.Pos, 3.0)
	}
	if d.digging && d.Animations[d.currAnim].Done {
		d.digging = false
	}
	if d.marking && d.Animations[d.currAnim].Done {
		d.marking = false
	}
	currLevel := int(-d.Transform.Pos.Y / world.TileSize)
	if LowestLevel < currLevel && !d.hurt {
		LowestLevel = currLevel
	}
}

func (d *Dwarf) Draw(win *pixelgl.Window) {
	d.Animations[d.currAnim].Draw(win)
	if d.walking && time.Since(d.walkTimer).Seconds() > 0.4 {
		sfx.SoundPlayer.PlaySound(fmt.Sprintf("step%d", rand.Intn(4) + 1), 0.)
		d.walkTimer = time.Now()
	}
	if d.selected != nil && !d.hurt {
		if d.selected.Solid && d.selectLegal {
			particles.CreateStaticParticle("target", d.selected.Transform.Pos)
		} else {
			particles.CreateStaticParticle("target_blank", d.selected.Transform.Pos)
		}
	}
}

func (d *Dwarf) Damage(dmg float64, source pixel.Vec) {
	if dmg > 0 {
		d.hurt = true
		d.dmg = dmg
		d.source = source
	}
}