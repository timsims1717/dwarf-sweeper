package systems

import (
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/dungeon"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"math"
)

func CollisionSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasCollision) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		_, okC := result.Components[myecs.Collision].(myecs.Collider)
		phys, okP := result.Components[myecs.Physics].(*physics.Physics)
		if okT && okC && okP {
			lastPos := tran.LastPos
			done := false
			var next pixel.Vec
			count := 0
			for !done {
				posChange := tran.Pos.Sub(lastPos)
				mag := util.Magnitude(posChange)
				if mag > world.TileSize {
					posChange = util.Normalize(posChange).Scaled(world.TileSize)
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
				loc := dungeon.Dungeon.GetCave().GetTile(next)
				if loc != nil {
					stopped := false
					if math.Abs(loc.Transform.Pos.X-next.X) > world.TileSize || math.Abs(loc.Transform.Pos.Y-next.Y) > world.TileSize {
						fmt.Println("Time to teleport")
					}
					up := dungeon.Dungeon.GetCave().GetTile(pixel.V(next.X, next.Y+world.TileSize))
					upl := dungeon.Dungeon.GetCave().GetTile(pixel.V(next.X-world.TileSize*0.3, next.Y+world.TileSize))
					upr := dungeon.Dungeon.GetCave().GetTile(pixel.V(next.X+world.TileSize*0.3, next.Y+world.TileSize))
					dwn := dungeon.Dungeon.GetCave().GetTile(pixel.V(next.X, next.Y-world.TileSize*0.6))
					dwnl := dungeon.Dungeon.GetCave().GetTile(pixel.V(next.X-world.TileSize*0.3, next.Y-world.TileSize*0.6))
					dwnr := dungeon.Dungeon.GetCave().GetTile(pixel.V(next.X+world.TileSize*0.3, next.Y-world.TileSize*0.6))
					if ((up != nil && up.Solid) || (upl != nil && upl.Solid) || (upr != nil && upr.Solid)) && next.Y >= loc.Transform.Pos.Y {
						next.Y = loc.Transform.Pos.Y
						if phys.Velocity.Y > 0 {
							phys.Velocity.Y = 0
							stopped = true
						}
					}
					if ((dwn != nil && dwn.Solid) || (dwnr != nil && dwnr.Solid) || (dwnl != nil && dwnl.Solid)) && next.Y <= loc.Transform.Pos.Y {
						next.Y = loc.Transform.Pos.Y
						if phys.Velocity.Y < 0 {
							phys.Velocity.Y = 0
							stopped = true
						}
						phys.Grounded = true
					} else {
						phys.Grounded = false
					}
					right := dungeon.Dungeon.GetCave().GetTile(pixel.V(next.X+world.TileSize, next.Y))
					left := dungeon.Dungeon.GetCave().GetTile(pixel.V(next.X-world.TileSize, next.Y))
					if right != nil && right.Solid && next.X >= loc.Transform.Pos.X {
						next.X = loc.Transform.Pos.X
						if phys.Velocity.X > 0 {
							if phys.RicochetX {
								phys.Velocity.X = phys.Velocity.X * -0.6
								stopped = true
							} else {
								phys.Velocity.X = 0
								stopped = true
							}
						}
					}
					if left != nil && left.Solid && next.X <= loc.Transform.Pos.X {
						next.X = loc.Transform.Pos.X
						if phys.Velocity.X < 0 {
							if phys.RicochetX {
								phys.Velocity.X = phys.Velocity.X * -0.6
								stopped = true
							} else {
								phys.Velocity.X = 0
								stopped = true
							}
						}
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
			tran.Update()
			if debug.Debug {
				up := dungeon.Dungeon.GetCave().GetTile(pixel.V(tran.Pos.X, tran.Pos.Y+world.TileSize))
				upl := dungeon.Dungeon.GetCave().GetTile(pixel.V(tran.Pos.X-world.TileSize*0.3, tran.Pos.Y+world.TileSize))
				upr := dungeon.Dungeon.GetCave().GetTile(pixel.V(tran.Pos.X+world.TileSize*0.3, tran.Pos.Y+world.TileSize))
				dwn := dungeon.Dungeon.GetCave().GetTile(pixel.V(tran.Pos.X, tran.Pos.Y-world.TileSize*0.6))
				dwnl := dungeon.Dungeon.GetCave().GetTile(pixel.V(tran.Pos.X-world.TileSize*0.3, tran.Pos.Y-world.TileSize*0.6))
				dwnr := dungeon.Dungeon.GetCave().GetTile(pixel.V(tran.Pos.X+world.TileSize*0.3, tran.Pos.Y-world.TileSize*0.6))
				right := dungeon.Dungeon.GetCave().GetTile(pixel.V(tran.Pos.X+world.TileSize, tran.Pos.Y))
				left := dungeon.Dungeon.GetCave().GetTile(pixel.V(tran.Pos.X-world.TileSize, tran.Pos.Y))
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
			}
		}
	}
}