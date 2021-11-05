package descent

import (
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
)

type CaveType int

const (
	Normal = iota
	Infinite
	Minesweeper
)

var Descent = &descent{}

type descent struct {
	Cave    *cave.Cave
	Level   int
	Player  *Dwarf
	Start   bool
	Type    CaveType
	ExitPop *menus.PopUp
	canExit bool
	FreeCam bool
}

func Update() {
	if Descent != nil && Descent.Cave != nil {
		Descent.Cave.Pivot = Descent.GetPlayer().Transform.Pos
		if !Descent.FreeCam {
			camera.Cam.StayWithin(Descent.Cave.Pivot, world.TileSize*1.5)
		}
		if !Descent.Cave.Finite {
			p := cave.WorldToChunk(Descent.Cave.Pivot)
			all := append([]world.Coords{p}, p.Neighbors()...)
			for _, i := range all {
				if i.X >= 0 && i.Y >= 0 {
					if _, ok := Descent.Cave.RChunks[i]; !ok {
						Descent.Cave.RChunks[i] = cave.NewChunk(i, Descent.Cave, cave.BlockCollapse)
						Descent.Cave.FillChunk(Descent.Cave.RChunks[i])
						Descent.Cave.UpdateBatch = true
						IncreaseLevelInf()
					}
				} else if i.X < 0 && i.Y >= 0 {
					if _, ok := Descent.Cave.LChunks[i]; !ok {
						Descent.Cave.LChunks[i] = cave.NewChunk(i, Descent.Cave, cave.BlockCollapse)
						Descent.Cave.FillChunk(Descent.Cave.RChunks[i])
						Descent.Cave.UpdateBatch = true
						IncreaseLevelInf()
					}
				}
			}
		}
		Descent.Cave.Update()
		switch Descent.Type {
		case Minesweeper:
			Descent.canExit = CaveBombsMarked == CaveBombsLeft && CaveWrongMarks < 1
			Descent.ExitPop.Raw = "Flag all the remaining bombs to exit."
		case Infinite:
			Descent.canExit = false
		default:
			Descent.canExit = true
		}
		if Descent.canExit {
			Descent.ExitPop.Raw = fmt.Sprintf("%s to Exit", typeface.SymbolItem)
			Descent.ExitPop.Symbols = []string{data.GameInput.FirstKey("up")}
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
	d.Cave.UpdateAllTileSprites()
	d.Cave.UpdateBatch = true
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
	Descent.Level++
	Descent.Cave.BombPMin += 0.01
	Descent.Cave.BombPMax += 0.01
	if Descent.Cave.BombPMin > 0.3 {
		Descent.Cave.BombPMin = 0.3
	}
	if Descent.Cave.BombPMax > 0.4 {
		Descent.Cave.BombPMax = 0.4
	}
}