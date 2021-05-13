package cave

import (
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/input"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
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
	Speed = 75.
	JumpVel = 170.
	DigRange = 1.4
	MaxJump = 4
)

var Player1 *Dwarf

type Dwarf struct {
	Transform   *physics.Physics
	Animations  map[string]*img.Instance
	currAnim    string
	faceLeft    bool
	selectLegal bool
	walkTimer   time.Time
	walking     bool
	toJumpT     time.Time
	toJump      bool
	jumping     bool
	jumpOrigY   float64
	jumpHeight  int
	grounded    bool
	digging     bool
	marking     bool
	selected    *Tile
	distFell    float64
	cursorV     pixel.Vec
	relWorld    pixel.Vec
	hurt        bool
	dmg         float64
	source      pixel.Vec
	knockback   float64
	Dead        bool
}

func NewDwarf() *Dwarf {
	dwarfSheet, err := img.LoadSpriteSheet("assets/img/dwarf.json")
	if err != nil {
		panic(err)
	}
	idle := img.NewAnimation(dwarfSheet, []pixel.Rect{dwarfSheet.Sprites[0], dwarfSheet.Sprites[1]},true, false, 0.5)
	run := img.NewAnimation(dwarfSheet, []pixel.Rect{dwarfSheet.Sprites[2], dwarfSheet.Sprites[3], dwarfSheet.Sprites[4], dwarfSheet.Sprites[5]},true, false, 0.4)
	jump := img.NewAnimation(dwarfSheet, []pixel.Rect{dwarfSheet.Sprites[6], dwarfSheet.Sprites[7]},false, true, 0.2)
	fall := img.NewAnimation(dwarfSheet, []pixel.Rect{dwarfSheet.Sprites[8]},false, true, 0.1)
	dig := img.NewAnimation(dwarfSheet, []pixel.Rect{dwarfSheet.Sprites[9], dwarfSheet.Sprites[10], dwarfSheet.Sprites[11]},false, true, 0.4)
	flag := img.NewAnimation(dwarfSheet, []pixel.Rect{dwarfSheet.Sprites[12]},false, true, 0.2)
	hitfront := img.NewAnimation(dwarfSheet, []pixel.Rect{dwarfSheet.Sprites[13]},false, true, 0.1)
	hitback := img.NewAnimation(dwarfSheet, []pixel.Rect{dwarfSheet.Sprites[14]},false, true, 0.1)
	flat := img.NewAnimation(dwarfSheet, []pixel.Rect{dwarfSheet.Sprites[15]},false, true, 0.1)
	animations := make(map[string]*img.Instance)
	animations["idle"] = idle.NewInstance()
	animations["run"] = run.NewInstance()
	animations["jump"] = jump.NewInstance()
	animations["fall"] = fall.NewInstance()
	animations["dig"] = dig.NewInstance()
	animations["flag"] = flag.NewInstance()
	animations["hit-front"] = hitfront.NewInstance()
	animations["hit-back"] = hitback.NewInstance()
	animations["flat"] = flat.NewInstance()
	tran := transform.NewTransform()
	physicsT := &physics.Physics{
		Transform: tran,
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
	d.Transform.YOff = false
	d.Transform.XOff = false
	if d.hurt {
		if d.dmg > 0 {
			// todo: damage to health
			if d.knockback > 0.1 {
				d.Transform.CancelMovement()
				dir := util.Normalize(d.Transform.Pos.Sub(d.source))
				d.Transform.Velocity.X = dir.X * d.knockback
				d.Transform.Velocity.Y = dir.Y * d.knockback
			}
			d.knockback = 0.
			d.dmg = 0
		}
		d.digging = false
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
			if debug.Debug {
				debug.AddText(fmt.Sprintf("world coords: (%d,%d)", int(input.Input.World.X), int(input.Input.World.Y)))
				debug.AddText(fmt.Sprintf("tile coords: (%d,%d)", d.selected.Coords.X, d.selected.Coords.Y))
				debug.AddText(fmt.Sprintf("tile code: '%s'", d.selected.GetTileCode()))
			}
			d.selectLegal = math.Abs(d.Transform.Pos.X-d.selected.Transform.Pos.X) < world.TileSize*DigRange && math.Abs(d.Transform.Pos.Y-d.selected.Transform.Pos.Y) < world.TileSize*DigRange
			if input.Input.IsDig && !d.digging && !d.marking && d.selected.Solid && d.selectLegal {
				d.digging = true
				d.toJump = false
				d.jumping = false
				d.walking = false
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
			loc1 := CurrCave.GetTile(d.Transform.Pos)
			dwn1 := CurrCave.GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize))
			dwnlj := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.5, d.Transform.Pos.Y-world.TileSize))
			dwnrj := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.5, d.Transform.Pos.Y-world.TileSize))
			dwnlw := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.3, d.Transform.Pos.Y-world.TileSize))
			dwnrw := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.3, d.Transform.Pos.Y-world.TileSize))
			canJump := ((dwn1 != nil && dwn1.Solid) || (dwnlj != nil && dwnlj.Solid) || (dwnrj != nil && dwnrj.Solid)) && d.Transform.Pos.Y <= loc1.Transform.Pos.Y+1.0 && d.grounded
			onGround := ((dwn1 != nil && dwn1.Solid) || (dwnlw != nil && dwnlw.Solid) || (dwnrw != nil && dwnrw.Solid)) && d.Transform.Pos.Y <= loc1.Transform.Pos.Y+1.0
			switch input.Input.XDir {
			case input.Left:
				if input.Input.XDirC || d.Transform.Velocity.X >= 0. {
					d.Transform.SetVelX(-Speed, 0.1) // Walk speed
					d.walkTimer = time.Now()
				}
				if onGround {
					d.faceLeft = true
				}
			case input.Right:
				if input.Input.XDirC || d.Transform.Velocity.X <= 0. {
					d.Transform.SetVelX(Speed, 0.1)
					d.walkTimer = time.Now()
				}
				if onGround {
					d.faceLeft = false
				}
			case input.None:
				if input.Input.XDirC {
					d.Transform.SetVelX(0., 0.1)
				}
			}
			input.Input.XDirC = false
			// Ground test, considered on the ground for jumping purposes until half a tile out
			if (!d.jumping && loc1 != nil && canJump && input.Input.Jumping.Pressed()) || d.toJump {
				d.walking = false
				if d.toJump && time.Since(d.toJumpT).Seconds() > 0.1 {
					newAnim = "jump"
					d.toJump = false
					d.jumping = true
					d.grounded = false
					d.jumpOrigY = d.Transform.Pos.Y
					d.jumpHeight = -1
					sfx.SoundPlayer.PlaySound(fmt.Sprintf("step%d", rand.Intn(4)+1), 0.)
				} else if !d.toJump {
					newAnim = "jump"
					d.toJumpT = time.Now()
					d.toJump = true
				}
				d.distFell = 0.
			} else if !d.jumping && onGround {
				d.grounded = true
				if math.Abs(d.Transform.Velocity.X) < 20.0 {
					if d.distFell > 100. {
						newAnim = "flat"
					} else {
						newAnim = "idle"
					}
					d.walking = false
				} else if d.Transform.Velocity.X > 0. {
					newAnim = "run"
					d.faceLeft = false
					d.walking = true
				} else if d.Transform.Velocity.X < 0. {
					newAnim = "run"
					d.faceLeft = true
					d.walking = true
				}
				if newAnim != "flat" {
					d.distFell = 0.
				}
			} else {
				d.grounded = false
				d.walking = false
				if d.jumping {
					newAnim = "jump"
					dist := int((d.Transform.Pos.Y - d.jumpOrigY) / world.TileSize)
					if (dist < MaxJump - 2 && input.Input.Jumping.Pressed()) || dist == d.jumpHeight {
						d.Transform.Velocity.Y = JumpVel
						d.jumpHeight = dist
						d.Transform.YOff = true
					} else {
						input.Input.Jumping.Consume()
						d.Transform.Velocity.Y = JumpVel
						d.jumping = false
					}
				}
				if d.Transform.Velocity.Y < 0. {
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
	d.Transform.Update()
	loc := CurrCave.GetTile(d.Transform.Pos)
	if loc != nil {
		if math.Abs(loc.Transform.Pos.X-d.Transform.Pos.X) > world.TileSize || math.Abs(loc.Transform.Pos.Y-d.Transform.Pos.Y) > world.TileSize {
			fmt.Println("Time to teleport")
		}
		up := CurrCave.GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y+world.TileSize))
		upl := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.3, d.Transform.Pos.Y+world.TileSize))
		upr := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.3, d.Transform.Pos.Y+world.TileSize))
		dwn := CurrCave.GetTile(pixel.V(d.Transform.Pos.X, d.Transform.Pos.Y-world.TileSize))
		dwnl := CurrCave.GetTile(pixel.V(d.Transform.Pos.X-world.TileSize*0.3, d.Transform.Pos.Y-world.TileSize))
		dwnr := CurrCave.GetTile(pixel.V(d.Transform.Pos.X+world.TileSize*0.3, d.Transform.Pos.Y-world.TileSize))
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
			d.jumping = false
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
		d.Transform.Transform.Update()
	}
	if newAnim != d.currAnim {
		d.Animations[d.currAnim].Reset()
		d.currAnim = newAnim
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

func (d *Dwarf) Damage(dmg float64, source pixel.Vec, knockback float64) {
	if dmg > 0 {
		d.hurt = true
		d.dmg = dmg
		d.source = source
		d.knockback = knockback
	}
}