package descent

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
)

var (
	Difficulty = 1
	Depth      = 6
)

var Descent = &descent{
	Difficulty: 1,
	CoordsMap:  make(map[string]world.Coords),
	ExitPop:    nil,
}

type descent struct {
	Cave       *cave.Cave
	CurrDepth  int
	Depth      int
	Difficulty int
	Player     *Dwarf
	Start      bool
	Type       cave.CaveType
	ExitPop    *menus.PopUp
	canExit    bool
	Builder    *builder.CaveBuilder

	FreeCam      bool
	DisableInput bool
	CoordsMap    map[string]world.Coords

	Builders [][]builder.CaveBuilder
}

func New() {
	Descent = &descent{
		Difficulty: Difficulty,
		Depth:      Depth,
		Start:      true,
		CoordsMap:  make(map[string]world.Coords),
		ExitPop:    menus.NewPopUp("", nil),
	}
}

func Update() {
	if Descent != nil && Descent.Cave != nil {
		Descent.Cave.Pivot = Descent.GetPlayer().Transform.Pos
		if !Descent.FreeCam {
			camera.Cam.StayWithin(Descent.Cave.Pivot, world.TileSize*1.5)
		}
		if Descent.Cave.Type == cave.Infinite {
			p := cave.WorldToChunk(Descent.Cave.Pivot)
			all := append([]world.Coords{p}, p.Neighbors()...)
			for _, i := range all {
				if _, ok := Descent.Cave.Chunks[i]; !ok {
					Descent.Cave.Chunks[i] = cave.NewChunk(i, Descent.Cave, cave.BlockCollapse)
					Descent.Cave.FillChunk(Descent.Cave.Chunks[i])
					Descent.Cave.UpdateBatch = true
					IncreaseLevelInf()
				}
			}
		}
		Descent.Cave.Update()
		switch Descent.Type {
		case cave.Minesweeper:
			Descent.canExit = CaveBombsMarked == CaveBombsLeft && CaveWrongMarks < 1
			Descent.ExitPop.Raw = "Flag all the remaining bombs to exit."
		default:
			Descent.canExit = true
		}
		if Descent.canExit {
			Descent.ExitPop.Raw = fmt.Sprintf("%s to Exit", typeface.SymbolItem)
			Descent.ExitPop.Symbols = []string{data.GameInput.FirstKey("up")}
		}
	}
}

func UpdatePlayer(in *input.Input) {
	if Descent.Player != nil {
		if Descent.DisableInput {
			Descent.Player.Update(nil)
		} else {
			Descent.Player.Update(in)
		}
	}
}

func (d *descent) CanExit() bool {
	return d.canExit
}

func (d *descent) GetCave() *cave.Cave {
	return d.Cave
}

func (d *descent) SetCave(cave *cave.Cave) {
	d.Cave = cave
	d.Type = cave.Type
	d.Cave.UpdateAllTileSprites()
	d.Cave.UpdateBatch = true
	d.SetExitPopup()
}

func (d *descent) SetExitPopup() {
	d.ExitPop = menus.NewPopUp("", nil)
	myecs.Manager.NewEntity().
		AddComponent(myecs.PopUp, d.ExitPop).
		AddComponent(myecs.Transform, d.GetCave().GetExit().Transform).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
}

func (d *descent) GetPlayer() *Dwarf {
	return d.Player
}

func (d *descent) SetPlayer(dwarf *Dwarf) {
	d.Player = dwarf
}

func (d *descent) GetPlayerTile() *cave.Tile {
	return d.Cave.GetTile(d.Player.Transform.Pos)
}

func (d *descent) GetTile(pos pixel.Vec) *cave.Tile {
	return d.Cave.GetTile(pos)
}

func IncreaseLevelInf() {
	Descent.CurrDepth++
	Descent.Cave.BombPMin += 0.01
	Descent.Cave.BombPMax += 0.01
	if Descent.Cave.BombPMin > 0.3 {
		Descent.Cave.BombPMin = 0.3
	}
	if Descent.Cave.BombPMax > 0.4 {
		Descent.Cave.BombPMax = 0.4
	}
}
