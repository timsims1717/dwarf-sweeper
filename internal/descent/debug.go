package descent

import (
	"dwarf-sweeper/internal/debug"
	"dwarf-sweeper/pkg/input"
	"fmt"
)

func Debug(in *input.Input) {
	d := Descent.GetPlayer()
	if debug.Text {
		debug.AddText(fmt.Sprintf("world coords: (%d,%d)", int(in.World.X), int(in.World.Y)))
		if d.Hovered != nil {
			debug.AddText(fmt.Sprintf("tile coords: (%d,%d)", d.Hovered.RCoords.X, d.Hovered.RCoords.Y))
			//debug.AddText(fmt.Sprintf("chunk coords: (%d,%d)", d.hovered.Chunk.Coords.X, d.hovered.Chunk.Coords.Y))
			//debug.AddText(fmt.Sprintf("tile sub coords: (%d,%d)", d.hovered.SubCoords.X, d.hovered.SubCoords.Y))
			debug.AddText(fmt.Sprintf("tile type: '%s'", d.Hovered.Type))
			debug.AddText(fmt.Sprintf("tile bg sprite: '%s'", d.Hovered.BGSpriteS))
			debug.AddText(fmt.Sprintf("tile fg sprite: '%s'", d.Hovered.FGSpriteS))
			debug.AddText(fmt.Sprintf("tile smart str: '%s'", d.Hovered.SmartStr))
			debug.AddText(fmt.Sprintf("tile s, d: %t, %t", d.Hovered.Surrounded, d.Hovered.DSurrounded))
		}
		debug.AddText(fmt.Sprintf("dwarf position: (%d,%d)", int(d.Transform.APos.X), int(d.Transform.APos.Y)))
		//debug.AddText(fmt.Sprintf("dwarf actual position: (%f,%f)", d.Transform.Pos.X, d.Transform.Pos.Y))
		debug.AddText(fmt.Sprintf("dwarf velocity: (%d,%d)", int(d.Physics.Velocity.X), int(d.Physics.Velocity.Y)))
		//debug.AddText(fmt.Sprintf("dwarf moving?: (%t,%t)", d.Physics.IsMovingX(), d.Physics.IsMovingY()))
		//debug.AddText(fmt.Sprintf("dwarf grounded?: %t", d.Physics.Grounded))
		//debug.AddText(fmt.Sprintf("tile queue len: %d", len(d.tileQueue)))
	}
}