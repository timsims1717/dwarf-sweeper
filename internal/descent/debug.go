package descent

import (
	"dwarf-sweeper/internal/debug"
	"fmt"
	pxginput "github.com/timsims1717/pixel-go-input"
)

func Debug(in *pxginput.Input) {
	d := Descent.GetPlayers()[0]
	if debug.Text {
		debug.AddText(fmt.Sprintf("world coords: (%d,%d)", int(in.World.X), int(in.World.Y)))
		for i, p := range Descent.GetPlayers() {
			debug.AddText(fmt.Sprintf("P%d pos: (%d,%d)", i+1, int(p.Transform.APos.X), int(p.Transform.APos.Y)))
			debug.AddText(fmt.Sprintf("P%d vel: (%d,%d)", i+1, int(d.Physics.Velocity.X), int(d.Physics.Velocity.Y)))
			//debug.AddText(fmt.Sprintf("P%d canvas pos: (%d,%d)", i+1, int(p.Player.CanvasPos.X), int(p.Player.CanvasPos.Y)))
			//debug.AddText(fmt.Sprintf("P%d cam pos: (%d,%d)", i+1, int(p.Player.CamPos.X), int(p.Player.CamPos.Y)))
			//debug.AddText(fmt.Sprintf("P%d cam target: (%d,%d)", i+1, int(p.Player.CamTar.X), int(p.Player.CamTar.Y)))
			//debug.AddText(fmt.Sprintf("P%d cam vel: (%d,%d)", i+1, int(p.Player.CamVel.X), int(p.Player.CamVel.Y)))
			//debug.AddText(fmt.Sprintf("P%d cam rel X: %d", i+1, int(p.Player.RelX)))
		}
		if d.Hovered != nil {
			debug.AddText(fmt.Sprintf("tile world: (%d,%d)", int(d.Hovered.Transform.Pos.X), int(d.Hovered.Transform.Pos.Y)))
			debug.AddText(fmt.Sprintf("chunk coords: (%d,%d)", d.Hovered.Chunk.Coords.X, d.Hovered.Chunk.Coords.Y))
			debug.AddText(fmt.Sprintf("tile coords: (%d,%d)", d.Hovered.RCoords.X, d.Hovered.RCoords.Y))
			debug.AddText(fmt.Sprintf("tile type: '%s'", d.Hovered.Type))
			debug.AddText(fmt.Sprintf("tile biome: '%s'", d.Hovered.Biome))
			if d.Hovered.Bomb {
				debug.AddText(fmt.Sprintf("tile has a bomb"))
			}
			for i, spr := range d.Hovered.Sprites {
				debug.AddText(fmt.Sprintf("tile sprite %d: '%s'", i, spr.K))
			}
			debug.AddText(fmt.Sprintf("tile smart str: '%s'", d.Hovered.SmartStr))
			debug.AddText(fmt.Sprintf("tile bg smart str: '%s'", d.Hovered.BGSmartStr))
			debug.AddText(fmt.Sprintf("tile fog smart str: '%s'", d.Hovered.FogSmartStr))
			debug.AddText(fmt.Sprintf("tile fog sprite: '%s'", d.Hovered.FogSpriteS))
		}
		debug.AddText(fmt.Sprintf("cave biome: %s", Descent.Cave.Biome))
		debug.AddText(fmt.Sprintf("cave level: %d", Descent.Cave.Level))
		debug.AddText(fmt.Sprintf("cave depth: %d", Descent.CurrDepth))
		debug.AddText(fmt.Sprintf("descent depth: %d", Descent.Depth))
		debug.AddText(fmt.Sprintf("bombs in cave: %d", Descent.GetCave().BombsLeft))
	}
}
