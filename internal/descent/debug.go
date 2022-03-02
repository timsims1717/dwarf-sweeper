package descent

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/physics"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/transform"
	"fmt"
)

func Debug(in *input.Input) {
	d := Descent.GetPlayers()[0]
	if debug.Text {
		debug.AddText(fmt.Sprintf("world coords: (%d,%d)", int(in.World.X), int(in.World.Y)))
		for i, p := range Descent.GetPlayers() {
			debug.AddText(fmt.Sprintf("P%d canvas pos: (%d,%d)", i+1, int(p.Player.CanvasPos.X), int(p.Player.CanvasPos.Y)))
			debug.AddText(fmt.Sprintf("P%d cam pos: (%d,%d)", i+1, int(p.Player.CamPos.X), int(p.Player.CamPos.Y)))
		}
		if d.Hovered != nil {
			debug.AddText(fmt.Sprintf("chunk coords: (%d,%d)", d.Hovered.Chunk.Coords.X, d.Hovered.Chunk.Coords.Y))
			debug.AddText(fmt.Sprintf("tile coords: (%d,%d)", d.Hovered.RCoords.X, d.Hovered.RCoords.Y))
			//debug.AddText(fmt.Sprintf("chunk coords: (%d,%d)", d.hovered.Chunk.Coords.X, d.hovered.Chunk.Coords.Y))
			//debug.AddText(fmt.Sprintf("tile sub coords: (%d,%d)", d.hovered.SubCoords.X, d.hovered.SubCoords.Y))
			debug.AddText(fmt.Sprintf("tile type: '%s'", d.Hovered.Type))
			if d.Hovered.Bomb {
				debug.AddText(fmt.Sprintf("tile has a bomb"))
			}
			for i, spr := range d.Hovered.Sprites {
				debug.AddText(fmt.Sprintf("tile sprite %d: '%s'", i, spr.K))
			}
			debug.AddText(fmt.Sprintf("tile smart str: '%s'", d.Hovered.SmartStr))
			debug.AddText(fmt.Sprintf("tile fg smart str: '%s'", d.Hovered.FogSmartStr))
			debug.AddText(fmt.Sprintf("tile fg sprite: '%s'", d.Hovered.FogSpriteS))
			//debug.AddText(fmt.Sprintf("tile s, d: %t, %t", d.Hovered.Surrounded, d.Hovered.DSurrounded))
			for _, result := range myecs.Manager.Query(myecs.HasCollision) {
				tran, okT := result.Components[myecs.Transform].(*transform.Transform)
				_, okC := result.Components[myecs.Collision].(*data.Collider)
				phys, okP := result.Components[myecs.Physics].(*physics.Physics)
				if okT && okC && okP {
					t := Descent.GetCave().GetTile(tran.Pos)
					if t.RCoords.X == d.Hovered.RCoords.X && t.RCoords.Y == d.Hovered.RCoords.Y {
						debug.AddText(fmt.Sprintf("el position: (%d,%d)", int(tran.APos.X), int(tran.APos.Y)))
						debug.AddText(fmt.Sprintf("el velocity: (%d,%d)", int(phys.Velocity.X), int(phys.Velocity.Y)))
						debug.AddText(fmt.Sprintf("el bounded: (l:%t,r:%t,u:%t,d:%t)", phys.LeftBound, phys.RightBound, phys.TopBound, phys.BottomBound))
					}
				}
			}
		}
		debug.AddText(fmt.Sprintf("dwarf position: (%d,%d)", int(d.Transform.APos.X), int(d.Transform.APos.Y)))
		//debug.AddText(fmt.Sprintf("dwarf actual position: (%f,%f)", d.Transform.Pos.X, d.Transform.Pos.Y))
		debug.AddText(fmt.Sprintf("dwarf velocity: (%d,%d)", int(d.Physics.Velocity.X), int(d.Physics.Velocity.Y)))
		//debug.AddText(fmt.Sprintf("dwarf moving?: (%t,%t)", d.Physics.IsMovingX(), d.Physics.IsMovingY()))
		//debug.AddText(fmt.Sprintf("dwarf grounded?: %t", d.Physics.Grounded))
		//debug.AddText(fmt.Sprintf("tile queue len: %d", len(d.tileQueue)))
	}
}
