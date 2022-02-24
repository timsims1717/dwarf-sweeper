package descent

import (
	player2 "dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/generate/builder"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
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
	Start      bool
	Type       cave.CaveType
	ExitPop    *menus.PopUp
	canExit    bool
	Builder    *builder.CaveBuilder

	Dwarves []*Dwarf

	FreeCam      bool
	DisableInput bool
	CoordsMap    map[string]world.Coords

	Builders [][]builder.CaveBuilder
	Timer    *timing.FrameTimer
}

func New() {
	Clear()
	Descent = &descent{
		Difficulty: Difficulty,
		Depth:      Depth,
		Start:      true,
		CoordsMap:  make(map[string]world.Coords),
		ExitPop:    menus.NewPopUp(""),
		Timer:      timing.New(0.),
	}
}

func Update() {
	if Descent != nil && Descent.Cave != nil {
		Descent.Timer.Update()
		Descent.Cave.Pivots = []pixel.Vec{}
		for _, d := range Descent.Dwarves {
			Descent.Cave.Pivots = append(Descent.Cave.Pivots, d.Transform.Pos)
		}
		if Descent.Cave.Type == cave.Infinite {
			var all []world.Coords
			for _, pivot := range Descent.Cave.Pivots {
				p := cave.WorldToChunk(pivot)
				all = append(all, p)
				all = append(all, p.Neighbors()...)
			}
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
			Descent.canExit = player2.OverallStats.CaveBombsFlagged == player2.CaveBombsLeft && player2.OverallStats.CaveWrongFlags < 1
			Descent.ExitPop.SetText("Flag all the remaining bombs to exit.")
		default:
			Descent.canExit = true
		}
		if Descent.canExit {
			Descent.ExitPop.SetText("{symbol:up}:Exit")
		}
	}
}

func UpdatePlayers() {
	for _, d := range Descent.Dwarves {
		UpdatePlayer(d)
	}
}

func UpdateViews() {
	for i, d := range Descent.Dwarves {
		UpdateView(d, i, len(Descent.Dwarves))
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
	d.ExitPop = menus.NewPopUp("")
	myecs.Manager.NewEntity().
		AddComponent(myecs.PopUp, d.ExitPop).
		AddComponent(myecs.Transform, d.GetCave().GetExit().Transform).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
}

func (d *descent) GetPlayers() []*Dwarf {
	return d.Dwarves
}

func (d *descent) GetClosestPlayer(pos pixel.Vec) *Dwarf {
	dist := -1.
	var pick *Dwarf
	for _, p := range d.GetPlayers() {
		newDist := util.Magnitude(p.Transform.Pos.Sub(pos))
		if dist < 0. || newDist < dist {
			dist = newDist
			pick = p
		}
	}
	return pick
}

func (d *descent) GetClosestPlayerTile(pos pixel.Vec) *cave.Tile {
	dwarf := d.GetClosestPlayer(pos)
	if dwarf != nil {
		return d.Cave.GetTile(dwarf.Transform.Pos)
	}
	return nil
}

func (d *descent) GetRandomPlayer() *Dwarf {
	if len(d.Dwarves) > 0 {
		return d.Dwarves[random.Effects.Intn(len(d.Dwarves))]
	}
	return nil
}

func (d *descent) GetRandomPlayerTile() *cave.Tile {
	dwarf := d.GetRandomPlayer()
	if dwarf != nil {
		return d.Cave.GetTile(dwarf.Transform.Pos)
	}
	return nil
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

func Clear() {
	if Descent != nil {
		for _, d := range Descent.Dwarves {
			d.Delete()
			if d.Player.Puzzle != nil {
				d.Player.Puzzle.Close()
				d.Player.Puzzle = nil
			}
		}
		Descent.Dwarves = []*Dwarf{}
	}
}