package systems

import (
	"dwarf-sweeper/internal/cave"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/transform"
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
			loc := cave.CurrCave.GetTile(tran.Pos)
			if loc != nil {
				if math.Abs(loc.Transform.Pos.X-tran.Pos.X) > world.TileSize || math.Abs(loc.Transform.Pos.Y-tran.Pos.Y) > world.TileSize {
					fmt.Println("Time to teleport")
				}
				up := cave.CurrCave.GetTile(pixel.V(tran.Pos.X, tran.Pos.Y+world.TileSize))
				upl := cave.CurrCave.GetTile(pixel.V(tran.Pos.X-world.TileSize*0.3, tran.Pos.Y+world.TileSize))
				upr := cave.CurrCave.GetTile(pixel.V(tran.Pos.X+world.TileSize*0.3, tran.Pos.Y+world.TileSize))
				dwn := cave.CurrCave.GetTile(pixel.V(tran.Pos.X, tran.Pos.Y-world.TileSize*0.6))
				dwnl := cave.CurrCave.GetTile(pixel.V(tran.Pos.X-world.TileSize*0.3, tran.Pos.Y-world.TileSize*0.6))
				dwnr := cave.CurrCave.GetTile(pixel.V(tran.Pos.X+world.TileSize*0.3, tran.Pos.Y-world.TileSize*0.6))
				if ((up != nil && up.Solid) || (upl != nil && upl.Solid) || (upr != nil && upr.Solid)) && tran.Pos.Y > loc.Transform.Pos.Y {
					tran.Pos.Y = loc.Transform.Pos.Y
					if phys.Velocity.Y > 0 {
						phys.Velocity.Y = 0
					}
				}
				if (dwn != nil && dwn.Solid) || (dwnr != nil && dwnr.Solid) || (dwnl != nil && dwnl.Solid) {
					tran.Pos.Y = loc.Transform.Pos.Y
					if phys.Velocity.Y < 0 {
						phys.Velocity.Y = 0
					}
					phys.Grounded = true
				} else {
					phys.Grounded = false
				}
				right := cave.CurrCave.GetTile(pixel.V(tran.Pos.X+world.TileSize, tran.Pos.Y))
				left := cave.CurrCave.GetTile(pixel.V(tran.Pos.X-world.TileSize, tran.Pos.Y))
				if right != nil && right.Solid && tran.Pos.X > loc.Transform.Pos.X {
					tran.Pos.X = loc.Transform.Pos.X
					if phys.Velocity.X > 0 {
						if phys.RicochetX {
							phys.Velocity.X = phys.Velocity.X * -0.6
						} else {
							phys.Velocity.X = 0
						}
					}
				}
				if left != nil && left.Solid && tran.Pos.X < loc.Transform.Pos.X {
					tran.Pos.X = loc.Transform.Pos.X
					if phys.Velocity.X < 0 {
						if phys.RicochetX {
							phys.Velocity.X = phys.Velocity.X * -0.6
						} else {
							phys.Velocity.X = 0
						}
					}
				}
				tran.Update()
				if debug.Debug {
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
}