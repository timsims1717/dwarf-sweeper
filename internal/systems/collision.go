package systems

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"math"
)

const (
	collisionDistance = 0.8
	collisionPush     = 10.
	CollisionThresh   = 3.
	NearGroundThresh  = 4.
)

func CollisionSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasCollision) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		coll, okC := result.Components[myecs.Collision].(*data.Collider)
		phys, okP := result.Components[myecs.Physics].(*physics.Physics)
		if okT && okC && okP {
			dist := camera.Cam.Pos.Sub(tran.Pos)
			if math.Abs(dist.X) < constants.DrawDistance && math.Abs(dist.Y) < constants.DrawDistance {
				w := coll.Hitbox.W()
				h := coll.Hitbox.H()
				if math.Abs(tran.Rot) == 0.5 {
					w = coll.Hitbox.H()
					h = coll.Hitbox.W()
				}
				if !coll.GroundOnly {
					for _, result1 := range myecs.Manager.Query(myecs.HasCollision) {
						tran1, okT1 := result1.Components[myecs.Transform].(*transform.Transform)
						coll1, okC1 := result1.Components[myecs.Collision].(*data.Collider)
						phys1, okP1 := result1.Components[myecs.Physics].(*physics.Physics)
						distX := math.Abs(tran.Pos.X - tran1.Pos.X)
						collDist := (w + coll1.Hitbox.W()) * 0.5 * collisionDistance
						if okT1 && okC1 && okP1 {
							h1 := coll1.Hitbox.H()
							if math.Abs(tran1.Rot) == 0.5 {
								h1 = coll1.Hitbox.W()
							}
							if !coll1.GroundOnly && distX < collDist && math.Abs(tran1.Pos.Y-tran.Pos.Y) < (h + h1) * 0.5 {
								if tran.Pos.X < tran1.Pos.X {
									tran.Pos.X -= math.Min(collisionPush*timing.DT, math.Abs(distX-collDist)*0.5)
									tran1.Pos.X += math.Min(collisionPush*timing.DT, math.Abs(distX-collDist)*0.5)
									if phys.Velocity.X > 0. && !coll.CanPass {
										phys.Velocity.X = 0.
									}
									if phys1.Velocity.X < 0. && !coll1.CanPass {
										phys1.Velocity.X = 0.
									}
								} else if tran.Pos.X > tran1.Pos.X {
									tran.Pos.X += math.Min(collisionPush*timing.DT, math.Abs(distX-collDist)*0.5)
									tran1.Pos.X -= math.Min(collisionPush*timing.DT, math.Abs(distX-collDist)*0.5)
									if phys.Velocity.X < 0. && !coll.CanPass {
										phys.Velocity.X = 0.
									}
									if phys1.Velocity.X > 0. && !coll1.CanPass {
										phys1.Velocity.X = 0.
									}
								}
							}
						}
					}
				}
				lastPos := tran.LastPos
				done := false
				var next pixel.Vec
				count := 0
				stepSize := CollisionThresh
				for !done {
					posChange := tran.Pos.Sub(lastPos)
					mag := util.Magnitude(posChange)
					if mag > stepSize {
						posChange = util.Normalize(posChange).Scaled(stepSize)
						next = lastPos.Add(posChange)
					} else {
						next = tran.Pos
						done = true
					}
					if debug.Debug {
						col := colornames.Red
						if count == 1 {
							col = colornames.Blue
						} else if count == 2 {
							col = colornames.Green
						}
						count++
						debug.AddLine(col, imdraw.RoundEndShape, lastPos, next, 2.0)
					}
					loc := descent.Descent.GetCave().GetTile(next)
					if loc != nil {
						stopped := false

						// collision rays init
						wcr := int(w) / int(world.TileSize * 0.5)
						hcr := int(h) / int(world.TileSize * 0.5)
						if wcr < 2 {
							wcr = 2
						}
						if hcr < 2 {
							hcr = 2
						}
						iw := w - CollisionThresh * 2.
						ih := h - CollisionThresh * 2.

						// collision rays up and down
						var dwn, up, gr *cave.Tile
						for i := 0; i < wcr; i++ {
							dy := next.Y-h*0.51
							uy := next.Y+h*0.51
							gy := dy - NearGroundThresh
							var x float64
							if i == 0 {
								x = next.X - iw * 0.5
							} else if i == wcr - 1 {
								x = next.X + iw * 0.5
							} else {
								x = next.X + iw / float64(wcr) * float64(i) - iw * 0.5
							}
							d := descent.Descent.GetCave().GetTile(pixel.V(x, dy))
							u := descent.Descent.GetCave().GetTile(pixel.V(x, uy))
							g := descent.Descent.GetCave().GetTile(pixel.V(x, gy))
							if d != nil && d.Solid() {
								if i == 0 {
									coll.DL = true
								} else if i == wcr - 1 {
									coll.DR = true
								}
								dwn = d
								debug.AddLine(colornames.Green, imdraw.RoundEndShape, pixel.V(x, next.Y-ih*0.5), pixel.V(x, dy), 1.0)
							} else {
								debug.AddLine(colornames.Red, imdraw.RoundEndShape, pixel.V(x, next.Y-ih*0.5), pixel.V(x, dy), 1.0)
							}
							if u != nil && u.Solid() {
								if i == 0 {
									coll.UL = true
								} else if i == wcr - 1 {
									coll.UR = true
								}
								up = u
								debug.AddLine(colornames.Green, imdraw.RoundEndShape, pixel.V(x, next.Y+ih*0.5), pixel.V(x, uy), 1.0)
							} else {
								debug.AddLine(colornames.Red, imdraw.RoundEndShape, pixel.V(x, next.Y+ih*0.5), pixel.V(x, uy), 1.0)
							}
							if g != nil && g.Solid() {
								gr = g
							}
						}

						// collision checks up and down
						uY := loc.Transform.Pos.Y + (world.TileSize - h) * 0.5
						dY := loc.Transform.Pos.Y - (world.TileSize - h) * 0.5
						if up != nil && up.Solid() {
							if next.Y > uY {
								next.Y = uY
							}
							next.Y = uY
							phys.Ceilinged = true
							if phys.Velocity.Y > 0 {
								phys.Velocity.Y = 0
								stopped = true
							}
						} else {
							phys.Ceilinged = false
						}
						phys.Grounded = false
						if dwn != nil && dwn.Solid() {
							if next.Y < dY {
								next.Y = dY
							}
							phys.Grounded = true
							wasRDX := phys.RagDollX
							phys.RagDollX = false
							if phys.Velocity.Y < 0 {
								if phys.RagDollY {
									phys.Velocity.Y = phys.Velocity.Y * -phys.Bounciness
									phys.RagDollY = false
									phys.Grounded = false
									phys.RagDollX = wasRDX
								} else {
									phys.Velocity.Y = 0
								}
								stopped = true
							}
						} else {
							phys.Grounded = false
						}
						phys.NearGround = gr != nil && gr.Solid()

						// collision rays left and right
						var left, right *cave.Tile
						for i := 0; i < hcr; i++ {
							lx := next.X-w*0.51
							rx := next.X+w*0.51
							var y float64
							if i == 0 {
								y = next.Y + ih * 0.5
							} else if i == hcr - 1 {
								y = next.Y - ih * 0.5
							} else {
								y = next.Y - ih / float64(hcr) * float64(i) - ih * 0.5
							}
							l := descent.Descent.GetCave().GetTile(pixel.V(lx, y))
							r := descent.Descent.GetCave().GetTile(pixel.V(rx, y))
							if l != nil && l.Solid() {
								if i == 0 {
									coll.LU = true
								} else if i == hcr - 1 {
									coll.LD = true
								}
								left = l
								debug.AddLine(colornames.Green, imdraw.RoundEndShape, pixel.V(next.X-iw*0.5, y), pixel.V(lx, y), 1.0)
							} else {
								debug.AddLine(colornames.Red, imdraw.RoundEndShape, pixel.V(next.X-iw*0.5, y), pixel.V(lx, y), 1.0)
							}
							if r != nil && r.Solid() {
								if i == 0 {
									coll.RU = true
								} else if i == hcr - 1 {
									coll.RD = true
								}
								right = r
								debug.AddLine(colornames.Green, imdraw.RoundEndShape, pixel.V(next.X+iw*0.5, y), pixel.V(rx, y), 1.0)
							} else {
								debug.AddLine(colornames.Red, imdraw.RoundEndShape, pixel.V(next.X+iw*0.5, y), pixel.V(rx, y), 1.0)
							}
						}

						// collision checks left and right
						rX := loc.Transform.Pos.X + (world.TileSize - w) * 0.5
						lX := loc.Transform.Pos.X - (world.TileSize - w) * 0.5
						if right != nil && right.Solid() {
							if next.X > rX {
								next.X = rX
							}
							phys.RightWalled = true
							if phys.Velocity.X > 0 {
								if phys.RagDollX {
									phys.Velocity.X = phys.Velocity.X * -phys.Bounciness
									phys.RightWalled = false
								} else {
									phys.Velocity.X = 0
								}
								stopped = true
							}
						} else {
							phys.RightWalled = false
						}
						if left != nil && left.Solid() {
							if next.X < lX {
								next.X = lX
							}
							phys.LeftWalled = true
							if phys.Velocity.X < 0 {
								if phys.RagDollX {
									phys.Velocity.X = phys.Velocity.X * -phys.Bounciness
									phys.LeftWalled = false
								} else {
									phys.Velocity.X = 0
								}
								stopped = true
							}
						} else {
							phys.LeftWalled = false
						}

						// corner collision check
						vul := pixel.V(next.X-w*0.51, next.Y+h*0.51)
						vur := pixel.V(next.X+w*0.51, next.Y+h*0.51)
						vdl := pixel.V(next.X-w*0.51, next.Y-h*0.51)
						vdr := pixel.V(next.X+w*0.51, next.Y-h*0.51)
						ul := descent.Descent.GetCave().GetTile(vul)
						ur := descent.Descent.GetCave().GetTile(vur)
						dl := descent.Descent.GetCave().GetTile(vdl)
						dr := descent.Descent.GetCave().GetTile(vdr)
						coll.CUL = ul != nil && ul.Solid()
						coll.CUR = ur != nil && ur.Solid()
						coll.CDL = dl != nil && dl.Solid()
						coll.CDR = dr != nil && dr.Solid()
						if coll.CUL {
							debug.AddLine(colornames.Green, imdraw.RoundEndShape, vul, vul, 1.0)
						} else {
							debug.AddLine(colornames.Red, imdraw.RoundEndShape, vul, vul, 1.0)
						}
						if coll.CUR {
							debug.AddLine(colornames.Green, imdraw.RoundEndShape, vur, vur, 1.0)
						} else {
							debug.AddLine(colornames.Red, imdraw.RoundEndShape, vur, vur, 1.0)
						}
						if coll.CDL {
							debug.AddLine(colornames.Green, imdraw.RoundEndShape, vdl, vdl, 1.0)
						} else {
							debug.AddLine(colornames.Red, imdraw.RoundEndShape, vdl, vdl, 1.0)
						}
						if coll.CDR {
							debug.AddLine(colornames.Green, imdraw.RoundEndShape, vdr, vdr, 1.0)
						} else {
							debug.AddLine(colornames.Red, imdraw.RoundEndShape, vdr, vdr, 1.0)
						}

						if stopped {
							done = true
						}
						lastPos = next
					} else {
						done = true
					}
				}
				tran.Pos = next
			}
		}
	}
}