package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/particles"
	"dwarf-sweeper/internal/pathfinding"
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
	gnomeSpeed    = 75.
	gnomeAcc      = 10.
	gnomeJumpX    = 100.
	gnomeJumpY    = 50.
	gnomeDigSpeed = 15.
	gnomeColSpeed = 25.
	gnomeDigAcc   = 5.
)

type GnomeState int

const (
	GnomeWait = iota
	GnomeSeek
	GnomeToEmerge
	GnomeEmerge
	GnomeEmergeSide
	GnomeEmergeDown
	GnomeDig
	GnomeDigSide
	GnomeTarget
	GnomeChase
	GnomePrep
	GnomeLeap
	GnomeIdle
	GnomeDazed
	GnomeDead
)

type Gnome struct {
	Transform *transform.Transform
	Physics   *physics.Physics
	Collider  *data.Collider
	Animation *reanimator.Tree
	Entity    *ecs.Entity
	Health    *data.Health

	timer   *timing.Timer
	State   GnomeState
	counter int

	direction data.Direction
	path      []*cave.Tile
	target    *cave.Tile
	atkTar    *cave.Tile
	atkTran   *transform.Transform

	effectTimer *timing.Timer
}

func CreateGnome(c *cave.Cave, pos pixel.Vec) *Gnome {
	g := &Gnome{
		State: GnomeIdle,
	}
	g.Transform = transform.New().WithID("gnome")
	g.Transform.Pos = pos
	g.Physics = physics.New()
	g.Health = &data.Health{
		Max:    1,
		Curr:   1,
		Immune: data.EnemyImmunity,
	}
	g.Collider = data.NewCollider(pixel.R(0., 0., 16., 16.), data.Critter)
	g.Collider.Debug = true
	batch := img.Batchers[constants.EntityKey]
	g.Animation = reanimator.New(reanimator.NewSwitch().
		AddAnimation(reanimator.NewAnimFromSprite("idle", batch.GetSprite("gnome_idle"), reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprite("idle_hat", batch.GetSprite("gnome_idle_hat"), reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprite("idle_foot", batch.GetSprite("gnome_idle_foot"), reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprite("idle_nose", batch.GetSprite("gnome_idle_nose"), reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("run", batch.GetAnimation("gnome_run").S, reanimator.Loop).
			SetTrigger(4, func() {
				PlayStep(-1.)
			})).
		AddAnimation(reanimator.NewAnimFromSprites("leap", batch.GetAnimation("gnome_run").S[:2], reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("leap_start", batch.GetAnimation("gnome_leap").S, reanimator.Tran).
			SetTrigger(4, func() {
				if g.Transform.Pos.X > g.atkTran.Pos.X {
					g.Transform.Flip = true
					g.atkTran = transform.New()
					g.atkTran.Pos.X = g.Transform.Pos.X - world.TileSize * 50.
				} else {
					g.Transform.Flip = false
					g.atkTran = transform.New()
					g.atkTran.Pos.X = g.Transform.Pos.X + world.TileSize * 50.
				}
				g.State = GnomeLeap
				g.timer = timing.New(0.3)
				PlaySqueak()
				PlayStep(-1.)
			})).
		AddAnimation(reanimator.NewAnimFromSprites("dazed", batch.GetAnimation("gnome_dazed").S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("dead", batch.GetAnimation("gnome_dead").S, reanimator.Hold)).
		AddAnimation(reanimator.NewAnimFromSprites("emerge", batch.GetAnimation("gnome_emerge").S, reanimator.Hold).
			SetTriggerAll(func() {
				orig := g.Transform.Pos
				orig.Y -= world.TileSize * 0.49
				particles.BiomeParticles(orig, Descent.Cave.Biome, 2, 3, world.TileSize * 0.25, 0., math.Pi*0.5, math.Pi, 80., 10., 0.75, 0.1, true)
			}).
			SetTrigger(4, func() {
				g.State = GnomeTarget
			})).
		AddAnimation(reanimator.NewAnimFromSprites("dig", batch.GetAnimation("gnome_dig").S, reanimator.Hold).
			SetTriggerAll(func() {
				orig := g.Transform.Pos
				orig.Y -= world.TileSize * 0.49
				particles.BiomeParticles(orig, Descent.Cave.Biome, 2, 3, world.TileSize * 0.25, 0., math.Pi*0.5, math.Pi, 80., 10., 0.75, 0.1, true)
			}).
			SetTrigger(4, func() {
				g.State = GnomeWait
				g.timer = timing.New(4.)
				g.Transform.Pos.Y -= world.TileSize
			})).
		AddAnimation(reanimator.NewAnimFromSprites("dig_side", batch.GetAnimation("gnome_dig_side").S, reanimator.Hold).
			SetTriggerAll(func() {
				orig := g.Transform.Pos
				var rot float64
				if g.direction == data.Left {
					orig.X -= world.TileSize * 0.49
					rot = 0.
				} else {
					orig.X += world.TileSize * 0.49
					rot = math.Pi
				}
				particles.BiomeParticles(orig, Descent.Cave.Biome, 2, 3, world.TileSize * 0.25, 0., rot, math.Pi, 80., 10., 0.75, 0.1, true)
			}).
			SetTrigger(4, func() {
				g.State = GnomeWait
				g.timer = timing.New(4.)
				if g.direction == data.Left {
					g.Transform.Pos.X -= world.TileSize
				} else {
					g.Transform.Pos.X += world.TileSize
				}
			})).
		AddAnimation(reanimator.NewAnimFromSprites("emerge_side", batch.GetAnimation("gnome_emerge_side").S, reanimator.Tran).
			SetTriggerAll(func() {
				orig := g.Transform.Pos
				var rot float64
				if g.direction == data.Left {
					orig.X += world.TileSize * 0.49
					rot = math.Pi
				} else {
					orig.X -= world.TileSize * 0.49
					rot = 0.
				}
				particles.BiomeParticles(orig, Descent.Cave.Biome, 2, 3, world.TileSize * 0.25, 0., rot, math.Pi, 80., 10., 0.75, 0.1, true)
			}).
			SetTrigger(4, func() {
				g.State = GnomeTarget
			})).
		AddAnimation(reanimator.NewAnimFromSprites("emerge_top", batch.GetAnimation("gnome_emerge_top").S, reanimator.Tran).
			SetTriggerAll(func() {
				orig := g.Transform.Pos
				orig.Y += world.TileSize * 0.49
				particles.BiomeParticles(orig, Descent.Cave.Biome, 2, 3, world.TileSize * 0.25, 0., math.Pi * -0.5, math.Pi, 80., 10., 0.75, 0.1, true)
			}).
			SetTrigger(2, func() {
				g.State = GnomeTarget
			})).
		AddAnimation(reanimator.NewAnimFromSprite("fall", batch.GetFrame("gnome_run", 3), reanimator.Hold)).
		AddNull().SetChooseFn(func() int {
			switch g.State {
			case GnomeIdle, GnomeTarget, GnomeChase:
				if !g.Physics.NearGround && !g.Physics.Grounded {
					return 14
				} else if g.State == GnomeChase && g.Physics.IsMovingX() {
					return 4
				} else {
					if random.Effects.Intn(25) == 0 {
						r := random.Effects.Intn(3) + 1
						switch r {
						case 2:
							if g.Transform.Load {
								PlayStep(-1.)
							}
						case 3:
							if g.Transform.Load {
								PlaySqueak()
							}
						}
						return r
					}
					return 0
				}
			case GnomeLeap:
				return 5
			case GnomePrep:
				return 6
			case GnomeDazed, GnomeDead:
				if g.State == GnomeDead && (g.Physics.NearGround || g.Physics.Grounded) {
					return 8
				} else {
					return 7
				}
			case GnomeEmerge:
				return 9
			case GnomeDig:
				return 10
			case GnomeDigSide:
				return 11
			case GnomeEmergeSide:
				return 12
			case GnomeEmergeDown:
				return 13
			default:
				return 15
			}
	}), "idle")
	g.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Transform, g.Transform).
		AddComponent(myecs.Physics, g.Physics).
		AddComponent(myecs.Collision, g.Collider).
		AddComponent(myecs.Health, g.Health).
		AddComponent(myecs.Update, data.NewFrameFunc(g.Update)).
		AddComponent(myecs.Animation, g.Animation).
		AddComponent(myecs.Drawable, g.Animation).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	return g
}

func (g *Gnome) Update() bool {
	g.timer.Update()
	g.SwitchState()
actionSwitch:
	switch g.State {
	case GnomeSeek:
		// pathfinding
		if g.target == nil || g.path == nil || !g.target.Solid() {
			cd := Descent.GetClosestPlayer(g.Transform.Pos)
			dt := Descent.GetTile(cd.Transform.Pos)
			outline := Descent.Cave.GetOutline(dt.RCoords, 5.5)
			inline := Descent.Cave.GetOutline(dt.RCoords, 2.25)
			candidates := world.NotIn(outline, inline)
			if len(candidates) > 0 {
				i := random.Effects.Intn(len(candidates))
				tc := candidates[i]
				t := Descent.Cave.GetTileInt(tc.X, tc.Y)
				if t == nil || !t.Solid() || !t.Diggable() {
					break actionSwitch
				}
				start := Descent.GetTile(g.Transform.Pos)
				cave.Origin = start.RCoords
				cave.NeighborsFn = pathfinding.DigNeighbors
				cave.CostFn = pathfinding.DigCost
				path, dist, found := astar.Path(t, start)
				if found && dist < 20. {
					g.target = t
					g.path = []*cave.Tile{}
					for _, ip := range path {
						it := ip.(*cave.Tile)
						g.path = append(g.path, it)
					}
				} else {
					break actionSwitch
				}
			}
		}
		// movement
		if len(g.path) > 0 {
			next := g.path[0]
			mag := util.Magnitude(next.Transform.Pos.Sub(g.Transform.Pos))
			if mag < world.TileSize * 0.25 || mag > world.TileSize * 3. {
				if len(g.path) > 1 {
					g.path = g.path[1:]
					next = g.path[0]
				} else {
					g.Physics.CancelMovement()
					g.Transform.Pos = next.Transform.Pos
					g.path = nil
					next = nil
					g.State = GnomeToEmerge
					break actionSwitch
				}
			}
			if next.Solid() {
				move := util.Normalize(next.Transform.Pos.Sub(g.Transform.Pos))
				spd := gnomeColSpeed
				if next.Type == cave.Dig {
					spd = gnomeDigSpeed
				}
				g.Physics.SetVelX(move.X*spd, gnomeDigAcc)
				g.Physics.SetVelY(move.Y*spd, gnomeDigAcc)
				if g.effectTimer == nil || g.effectTimer.UpdateDone() {
					currT := Descent.GetTile(g.Transform.Pos)
					DigParticle(g.Transform.Pos, currT.Biome)
					g.effectTimer = timing.New(digTimer)
				}
			} else {
				g.Physics.CancelMovement()
				g.path = nil
				next = nil
			}
		} else {
			g.Physics.CancelMovement()
		}
	case GnomeChase:
		if g.atkTran != nil {
			// pathfinding
			dist := util.Magnitude(g.Transform.Pos.Sub(g.atkTran.Pos))
			if g.timer.Done() && dist > 4. * world.TileSize &&
				random.Effects.Intn(6*int(1/timing.DT)) == 0 {
				g.atkTran = nil
				break actionSwitch
			}
			if g.target == nil || g.path == nil || g.atkTar == nil ||
				util.Magnitude(g.atkTar.Transform.Pos.Sub(g.atkTran.Pos)) > world.TileSize {
				// update pathfinding
				g.atkTar = Descent.GetTile(g.atkTran.Pos)
				g.target = nil
				below := g.atkTar.RCoords
				tries := 0
				for g.target == nil {
					if tries > 3 {
						g.atkTran = nil
						break actionSwitch
					}
					below.Y++
					below.X += random.Effects.Intn(3) - 1
					tt := Descent.Cave.GetTileInt(below.X, below.Y-1)
					tb := Descent.Cave.GetTileInt(below.X, below.Y)
					if tb != nil && tb.Solid() && tt != nil && !tt.Solid() {
						g.target = tt
					}
					tries++
				}
				start := Descent.GetTile(g.Transform.Pos)
				cave.Origin = start.RCoords
				cave.NeighborsFn = pathfinding.RunNeighbors
				cave.CostFn = pathfinding.RunCost
				path, dist, found := astar.Path(g.target, start)
				if !found || dist > 20. {
					g.atkTran = nil
					break
				}
				g.path = []*cave.Tile{}
				for _, ip := range path {
					i := ip.(*cave.Tile)
					g.path = append(g.path, i)
				}
			}
			// movement
			if len(g.path) > 0 {
				next := g.path[0]
				mag := util.Magnitude(next.Transform.Pos.Sub(g.Transform.Pos))
				if mag < world.TileSize * 0.5 || mag > world.TileSize * 3. ||
					next.Transform.Pos.Y > g.Transform.Pos.Y + world.TileSize * 0.75 {
					if len(g.path) > 1 {
						g.path = g.path[1:]
						next = g.path[0]
					} else {
						g.target = nil
						g.path = nil
					}
				}
				if g.Transform.Pos.X > next.Transform.Pos.X {
					g.Transform.Flip = true
					g.Physics.SetVelX(-gnomeSpeed, gnomeAcc)
				} else {
					g.Transform.Flip = false
					g.Physics.SetVelX(gnomeSpeed, gnomeAcc)
				}
			}
		}
	case GnomeLeap:
		if !g.timer.Done() {
			g.Physics.SetVelY(gnomeJumpY, 0.)
			if g.Transform.Pos.X > g.atkTran.Pos.X {
				g.Transform.Flip = true
				g.Physics.SetVelX(-gnomeJumpX, 0.)
			} else {
				g.Transform.Flip = false
				g.Physics.SetVelX(gnomeJumpX, 0.)
			}
		} else if g.Physics.Grounded || g.Physics.NearGround {
			g.State = GnomeIdle
			g.timer = timing.New(1.2)
			g.atkTran = nil
			break
		}
		myecs.Manager.NewEntity().AddComponent(myecs.AreaDmg, &data.AreaDamage{
			SourceID:  g.Transform.ID,
			Center:    g.Transform.Pos,
			Radius:    4.,
			Amount:    1,
			Dazed:     1.,
			Knockback: 8.,
			Type:      data.Enemy,
		})
	}
	if debug.Debug && len(g.path) > 0 {
		var n *cave.Tile
		col := colornames.Green
		for _, t := range g.path {
			if n != nil {
				debug.AddLine(col, imdraw.SharpEndShape, n.Transform.Pos, t.Transform.Pos, 2.)
				col.R += 25
			}
			n = t
		}
		col = colornames.Aquamarine
		for _, t := range g.path {
			debug.AddLine(col, imdraw.RoundEndShape, t.Transform.Pos, t.Transform.Pos, 2.)
			col.G -= 25
		}
	}
	return false
}

func (g *Gnome) SwitchState() {
	t := Descent.GetTile(g.Transform.Pos)
	if t == nil || t.Solid() && !t.Diggable() {
		g.State = GnomeDead
	} else if t.Solid() && (g.State != GnomeWait && g.State != GnomeSeek && g.State != GnomeToEmerge) {
		g.State = GnomeWait
	}
	if g.Health.Dead {
		if g.State != GnomeDead {
			PlaySqueak()
			g.Entity.RemoveComponent(myecs.Update)
			g.Entity.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
				myecs.AddEffect(g.Entity, data.NewBlink(2.))
				return true
			}, 2.))
			g.Entity.AddComponent(myecs.Temp, timing.New(4.))
		}
		g.State = GnomeDead
	} else if g.Health.Dazed {
		if g.State != GnomeDazed {
			PlaySqueak()
		}
		g.State = GnomeDazed
	} else if g.State == GnomeDazed || g.State == GnomeDead {
		g.State = GnomeIdle
		g.timer = timing.New(1.2)
		g.counter = 0
	} else {
		switch g.State {
		case GnomeTarget:
			g.atkTran = Descent.GetClosestPlayer(g.Transform.Pos).Transform
			dist := util.Magnitude(g.Transform.Pos.Sub(g.atkTran.Pos))
			if dist < world.TileSize * 1.5 {
				g.State = GnomePrep
				g.counter = 0
			} else {
				g.State = GnomeChase
				g.timer = timing.New(1.2)
				g.counter++
			}
		case GnomeIdle:
			if g.timer.Done() && random.Effects.Intn(int(1/timing.DT)) == 0 && g.Physics.Grounded {
				switch random.Effects.Intn(4) {
				case 0:
					g.atkTran = Descent.GetClosestPlayer(g.Transform.Pos).Transform
					dist := util.Magnitude(g.Transform.Pos.Sub(g.atkTran.Pos))
					if dist < world.TileSize * 1.5 {
						g.State = GnomePrep
						g.counter = 0
					} else {
						g.State = GnomeChase
						g.timer = timing.New(1.2)
						g.counter++
					}
				case 1:
					// dig
					pos := g.Transform.Pos
					pos.Y -= world.TileSize
					if Descent.Cave.GetTile(pos).Diggable() {
						g.State = GnomeDig
						g.direction = data.Down
						PlayRocks(-2.0)
						break
					}
					// couldn't dig down, try right and left
					lpos := g.Transform.Pos
					rpos := g.Transform.Pos
					lpos.X -= world.TileSize
					rpos.X += world.TileSize
					ldig := Descent.Cave.GetTile(lpos).Diggable() && g.Physics.LeftBound
					rdig := Descent.Cave.GetTile(rpos).Diggable() && g.Physics.RightBound
					if ldig || rdig {
						g.State = GnomeDigSide
						if ldig && !rdig {
							g.direction = data.Left
						} else if rdig && !ldig {
							g.direction = data.Right
						} else {
							if random.Effects.Intn(2) == 0 {
								g.direction = data.Left
							} else {
								g.direction = data.Right
							}
						}
						g.Transform.Flip = g.direction == data.Left
						PlayRocks(-2.0)
						break
					}
					g.counter++
				case 3:
					if g.counter > 0 {
						mpos := g.Transform.Pos
						mpos.X += (random.Effects.Float64() - 0.5) * 10. * world.TileSize
						g.atkTran = Descent.GetTile(mpos).Transform
						g.State = GnomeChase
						g.timer = timing.New(1.2)
						break
					}
					fallthrough
				default:
					g.Transform.Flip = !g.Transform.Flip
					g.timer = timing.New(0.5)
				}
			}
		case GnomeChase:
			if g.atkTran == nil {
				g.State = GnomeIdle
				g.timer = timing.New(1.2)
				break
			}
			dist := util.Magnitude(g.Transform.Pos.Sub(g.atkTran.Pos))
			if dist < world.TileSize * 1.5 {
				g.State = GnomePrep
				g.counter = 0
			}
		case GnomeSeek, GnomeWait:
			currT := Descent.GetTile(g.Transform.Pos)
			if currT != nil {
				if !currT.Solid() && !g.Entity.HasComponent(myecs.Damage) {
					g.Physics.CancelMovement()
					g.State = GnomeDazed
					g.Entity.AddComponent(myecs.Damage, &data.Damage{
						Dazed:     1.5,
						Knockback: 8.,
						Angle:     &upAngle,
						Source:    g.Transform.Pos,
						Type:      data.Blast,
					})
				}
				if g.timer.Done() {
					d := Descent.GetClosestPlayer(g.Transform.Pos)
					dt := Descent.GetTile(d.Transform.Pos)
					dist := util.Magnitude(g.Transform.Pos.Sub(d.Transform.Pos))
					if g.State == GnomeWait {
						outline := Descent.Cave.GetOutline(dt.RCoords, 5.5)
						emerge := world.CoordsIn(currT.RCoords, outline) &&
							(Descent.Cave.DestroyedWithin(currT.RCoords, 4, 2) || dist < 2. * world.TileSize)
						seek := Descent.Cave.DestroyedWithin(currT.RCoords, 5, 3) ||
							dist < 3. * world.TileSize || random.Effects.Intn(int(1/timing.DT)) == 0
						if emerge {
							g.State = GnomeToEmerge
						} else if seek {
							g.State = GnomeSeek
							g.timer = timing.New(4.)
						}
					} else if g.State == GnomeSeek &&
						(dist > 16. * world.TileSize || random.Effects.Intn(12*int(1/timing.DT)) == 0) {
						g.State = GnomeWait
						g.timer = timing.New(4.)
					}
				}
			}
		case GnomeToEmerge:
			g.Physics.CancelMovement()
			currT := Descent.GetTile(g.Transform.Pos)
			if currT != nil {
				d := Descent.GetClosestPlayer(g.Transform.Pos)
				dt := Descent.GetTile(d.Transform.Pos)
				blob := Descent.Cave.GetBlob(dt.RCoords, 6.)
				above := Descent.Cave.GetTileInt(currT.RCoords.X, currT.RCoords.Y-1)
				if !above.Solid() && world.CoordsIn(above.RCoords, blob) {
					g.State = GnomeEmerge
					g.Transform.Pos.Y += world.TileSize
					break
				}
				left := Descent.Cave.GetTileInt(currT.RCoords.X-1, currT.RCoords.Y)
				if !left.Solid() && world.CoordsIn(left.RCoords, blob) {
					g.State = GnomeEmergeSide
					g.Transform.Pos.X -= world.TileSize
					g.Transform.Flip = true
					g.direction = data.Left
					break
				}
				right := Descent.Cave.GetTileInt(currT.RCoords.X+1, currT.RCoords.Y)
				if !right.Solid() && world.CoordsIn(right.RCoords, blob) {
					g.State = GnomeEmergeSide
					g.Transform.Pos.X += world.TileSize
					g.Transform.Flip = false
					g.direction = data.Right
					break
				}
				below := Descent.Cave.GetTileInt(currT.RCoords.X, currT.RCoords.Y+1)
				if !below.Solid() && world.CoordsIn(below.RCoords, blob) {
					g.State = GnomeEmergeDown
					g.Transform.Pos.Y -= world.TileSize
					break
				}
				g.State = GnomeSeek
			}
		}
	}
	if g.State == GnomeWait || g.State == GnomeSeek ||
		g.State == GnomeEmerge || g.State == GnomeEmergeSide || g.State == GnomeEmergeDown {
		if g.State == GnomeWait {
			g.Physics.CancelMovement()
		}
		if g.State == GnomeWait || g.State == GnomeSeek {
			g.Health.Immune = data.UndergroundImmunity
		} else {
			g.Health.Immune = data.EnemyImmunity
		}
		g.Collider.NoClip = true
		g.Physics.GravityOff = true
	} else {
		g.Collider.NoClip = false
		g.Physics.GravityOff = false
		g.Health.Immune = data.EnemyImmunity
	}
	if g.State != GnomeSeek && g.State != GnomeChase {
		g.path = nil
		g.target = nil
	}
}