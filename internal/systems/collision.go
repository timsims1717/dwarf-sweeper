package systems

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"math"
)

const (
	collisionDistance = 0.85
	collisionPush     = 10.
	CollisionStep     = 3.
	NearGroundThresh  = 4.
	BounceThresholdX  = 70.
	BounceThresholdY  = 120.
)

func TileCollisionSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasCollision) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		coll, okC := result.Components[myecs.Collision].(*data.Collider)
		phys, okP := result.Components[myecs.Physics].(*physics.Physics)
		if okT && okC && okP && tran.Load {
			if coll.NoClip {
				debug.AddLine(colornames.Red, imdraw.RoundEndShape, tran.Pos, tran.Pos, 2.0)
				coll.BottomBound = false
				coll.TopBound = false
				coll.LeftBound = false
				coll.RightBound = false
			} else {
				coll.Collided = false
				stopped := false
				var hb pixel.Rect
				if math.Abs(tran.Rot) == 0.5 {
					hb = pixel.R(0., 0., coll.Hitbox.H(), coll.Hitbox.W())
				} else {
					hb = pixel.R(0., 0., coll.Hitbox.W(), coll.Hitbox.H())
				}
				hb = hb.Moved(tran.Pos).Moved(pixel.V(hb.W()*-0.5, hb.H()*-0.5))
				lastPos := tran.LastPos
				done := false
				var next pixel.Vec
				count := 0
				stepSize := math.Min(CollisionStep, math.Min(hb.H(), hb.W())*0.5)
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
					if debug.Debug && coll.Debug {
						col := colornames.Red
						if count == 1 {
							col = colornames.Blue
						} else if count == 2 {
							col = colornames.Green
						}
						count++
						debug.AddLine(col, imdraw.RoundEndShape, lastPos, next, 2.0)
					}
					nTile := descent.Descent.GetCave().GetTile(next)
					if nTile != nil {
						// collision rays init
						w := hb.W()
						h := hb.H()
						wcrf := w / (world.TileSize * 0.5)
						hcrf := h / (world.TileSize * 0.5)
						wcr := int(math.Round(wcrf))
						hcr := int(math.Round(hcrf))
						if wcr < 1 {
							wcr = 1
						}
						if hcr < 1 {
							hcr = 1
						}

						iw := w - CollisionStep*2.
						ih := h - CollisionStep*2.

						// collision rays up and down
						var dwn, up *cave.Tile
						for i := 0; i < wcr; i++ {
							dy := next.Y - h*0.51
							uy := next.Y + h*0.51
							var x float64
							if i == 0 {
								if wcr == 1 {
									x = next.X
								} else {
									x = next.X - iw*0.5
								}
							} else if i == wcr-1 {
								x = next.X + iw*0.5
							} else {
								x = next.X - w*0.5 + stepSize + float64(i)*iw/float64(wcr-1)
							}
							d := descent.Descent.GetCave().GetTile(pixel.V(x, dy))
							u := descent.Descent.GetCave().GetTile(pixel.V(x, uy))
							if d != nil && d.Solid() && d.Transform.Pos.Y+d.Type.Rect().Max.Y > dy {
								if i == 0 {
									coll.DL = true
								} else if i == wcr-1 {
									coll.DR = true
								}
								dwn = d
								if debug.Debug && coll.Debug {
									if d.Type == cave.Bridge {
										debug.AddLine(colornames.Yellow, imdraw.RoundEndShape, pixel.V(x, next.Y-ih*0.5), pixel.V(x, dy), 1.0)
									} else {
										debug.AddLine(colornames.Green, imdraw.RoundEndShape, pixel.V(x, next.Y-ih*0.5), pixel.V(x, dy), 1.0)
									}
								}
							} else if debug.Debug && coll.Debug {
								debug.AddLine(colornames.Red, imdraw.RoundEndShape, pixel.V(x, next.Y-ih*0.5), pixel.V(x, dy), 1.0)
							}
							if u != nil && u.Solid() && u.Type != cave.Bridge {
								if i == 0 {
									coll.UL = true
								} else if i == wcr-1 {
									coll.UR = true
								}
								up = u

								if debug.Debug && coll.Debug {
									debug.AddLine(colornames.Green, imdraw.RoundEndShape, pixel.V(x, next.Y+ih*0.5), pixel.V(x, uy), 1.0)
								}
							} else if debug.Debug && coll.Debug {
								debug.AddLine(colornames.Red, imdraw.RoundEndShape, pixel.V(x, next.Y+ih*0.5), pixel.V(x, uy), 1.0)
							}
						}

						// collision checks up and down
						uY := nTile.Transform.Pos.Y + (world.TileSize-h)*0.5
						if up != nil && up.Solid() {
							//if next.Y > uY {
							//	next.Y = uY
							//}
							next.Y = uY
							coll.TopBound = true
							if phys.Velocity.Y > 0 {
								if phys.RagDollY && math.Abs(phys.Velocity.Y) > BounceThresholdY/phys.Bounciness {
									phys.Velocity.Y *= -phys.Bounciness
									coll.TopBound = false
								} else {
									phys.Velocity.Y = 0
								}
								stopped = true
							}
						} else {
							coll.TopBound = false
						}
						if dwn != nil && dwn.Solid() && !(coll.Fallthrough && dwn.Type == cave.Bridge) {
							dY := dwn.Transform.Pos.Y + dwn.Type.Rect().Max.Y + h*0.5
							if next.Y < dY {
								next.Y = dY
							}
							phys.Grounded = true
							coll.BottomBound = true
							wasRDX := phys.RagDollX
							phys.RagDollX = false
							if phys.Velocity.Y < 0 {
								if phys.RagDollY && math.Abs(phys.Velocity.Y) > BounceThresholdY/phys.Bounciness {
									phys.Velocity.Y *= -phys.Bounciness
									phys.Grounded = false
									coll.BottomBound = false
									phys.RagDollX = wasRDX
								} else {
									phys.Velocity.Y = 0
									phys.RagDollY = false
								}
								stopped = true
							}
							phys.NearGround = next.Y-4. < dY
						} else {
							phys.Grounded = false
							coll.BottomBound = false
							phys.NearGround = false
						}

						if !coll.ThroughWalls {
							// collision rays left and right
							var left, right *cave.Tile
							for i := 0; i < hcr; i++ {
								lx := next.X - w*0.51
								rx := next.X + w*0.51
								var y float64
								if i == 0 {
									if hcr == 1 {
										y = next.Y
									} else {
										y = next.Y + ih*0.5
									}
								} else if i == hcr-1 {
									y = next.Y - ih*0.5
								} else {
									y = next.Y - h*0.5 + stepSize + float64(i)*ih/float64(hcr-1)
								}
								l := descent.Descent.GetCave().GetTile(pixel.V(lx, y))
								r := descent.Descent.GetCave().GetTile(pixel.V(rx, y))
								if l != nil && l.Solid() &&
									l.Transform.Pos.Y+l.Type.Rect().Max.Y > y &&
									l.Transform.Pos.Y+l.Type.Rect().Min.Y < y {
									if i == 0 {
										coll.LU = true
									} else if i == hcr-1 {
										coll.LD = true
									}
									left = l

									if debug.Debug && coll.Debug {
										debug.AddLine(colornames.Green, imdraw.RoundEndShape, pixel.V(next.X-iw*0.5, y), pixel.V(lx, y), 1.0)
									}
								} else if debug.Debug && coll.Debug {
									debug.AddLine(colornames.Red, imdraw.RoundEndShape, pixel.V(next.X-iw*0.5, y), pixel.V(lx, y), 1.0)
								}
								if r != nil && r.Solid() &&
									r.Transform.Pos.Y+r.Type.Rect().Max.Y > y &&
									r.Transform.Pos.Y+r.Type.Rect().Min.Y < y {
									if i == 0 {
										coll.RU = true
									} else if i == hcr-1 {
										coll.RD = true
									}
									right = r
									if debug.Debug && coll.Debug {
										debug.AddLine(colornames.Green, imdraw.RoundEndShape, pixel.V(next.X+iw*0.5, y), pixel.V(rx, y), 1.0)
									}
								} else if debug.Debug && coll.Debug {
									debug.AddLine(colornames.Red, imdraw.RoundEndShape, pixel.V(next.X+iw*0.5, y), pixel.V(rx, y), 1.0)
								}
							}

							// collision checks left and right
							rX := nTile.Transform.Pos.X + (world.TileSize-w)*0.5
							lX := nTile.Transform.Pos.X - (world.TileSize-w)*0.5
							if right != nil && right.Solid() {
								if next.X > rX {
									next.X = rX
								}
								coll.RightBound = true
								if phys.Velocity.X > 0 {
									if phys.RagDollX && math.Abs(phys.Velocity.X) > BounceThresholdX/phys.Bounciness {
										phys.Velocity.X *= -phys.Bounciness
										coll.RightBound = false
									} else {
										phys.Velocity.X = 0
									}
									stopped = true
								}
							} else {
								coll.RightBound = false
							}
							if left != nil && left.Solid() {
								if next.X < lX {
									next.X = lX
								}
								coll.LeftBound = true
								if phys.Velocity.X < 0 {
									if phys.RagDollX && math.Abs(phys.Velocity.X) > BounceThresholdX/phys.Bounciness {
										phys.Velocity.X *= -phys.Bounciness
										coll.LeftBound = false
									} else {
										phys.Velocity.X = 0
									}
									stopped = true
								}
							} else {
								coll.LeftBound = false
							}
							phys.WallBound = coll.LeftBound || coll.RightBound
						} else {
							coll.LeftBound = false
							coll.RightBound = false
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
						if debug.Debug && coll.Debug {
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
						}

						if stopped {
							done = true
							coll.Collided = true
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
