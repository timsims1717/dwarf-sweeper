package systems

import (
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"math"
)

func PopUpSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasPopUp) {
		tran, okT := result.Components[myecs.Transform].(*transform.Transform)
		pop, okP := result.Components[myecs.PopUp].(*menus.PopUp)
		if okT && okP {
			pop.Tran.Pos = tran.Pos
			pop.Display1 = false
			pop.Display2 = false
			pop.Display3 = false
			pop.Display4 = false
		}
	}
	for _, d := range descent.Descent.GetPlayers() {
		pivot := d.Transform.Pos
		dist := -1.
		var toDisplay *menus.PopUp
		for _, result := range myecs.Manager.Query(myecs.HasPopUp) {
			tran, okT := result.Components[myecs.Transform].(*transform.Transform)
			pop, okP := result.Components[myecs.PopUp].(*menus.PopUp)
			if okT && okP {
				disp := math.Abs(pivot.X-tran.Pos.X) < pop.Dist && math.Abs(pivot.Y-tran.Pos.Y) < pop.Dist
				tDist := util.Magnitude(pivot.Sub(tran.Pos))
				if disp && dist == -1. || dist > tDist && pop.Raw != "" {
					dist = tDist
					toDisplay = pop
				}
			}
		}
		if toDisplay != nil {
			switch d.Player.Code {
			case "p1":
				toDisplay.Display1 = true
				toDisplay.Player1 = d.Player
			case "p2":
				toDisplay.Display2 = true
				toDisplay.Player2 = d.Player
			case "p3":
				toDisplay.Display3 = true
				toDisplay.Player3 = d.Player
			case "p4":
				toDisplay.Display4 = true
				toDisplay.Player4 = d.Player
			}
		}
	}
	for _, result := range myecs.Manager.Query(myecs.HasPopUp) {
		if pop, ok := result.Components[myecs.PopUp].(*menus.PopUp); ok {
			pop.Update()
		}
	}
}

func PopUpDraw() {
	for _, result := range myecs.Manager.Query(myecs.HasPopUp) {
		if pop, ok := result.Components[myecs.PopUp].(*menus.PopUp); ok {
			pop.Draw()
		}
	}
}
