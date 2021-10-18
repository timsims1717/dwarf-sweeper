package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/typeface"
	"dwarf-sweeper/pkg/world"
	"fmt"
)

type CaveType int

const (
	Normal = iota
	Infinite
	Minesweeper
)

var Descent = &descent{}

type descent struct {
	Cave     *cave.Cave
	Level    int
	Player   *Dwarf
	Start    bool
	Type     CaveType
	ExitPop  *menus.PopUp
	canExit  bool
}

func Update() {
	if Descent != nil && Descent.Cave != nil {
		Descent.Cave.Pivot = Descent.GetPlayer().Transform.Pos
		if !Descent.Cave.Finite {
			p := cave.WorldToChunk(Descent.Cave.Pivot)
			all := append([]world.Coords{p}, p.Neighbors()...)
			for _, i := range all {
				if i.X >= 0 && i.Y >= 0 {
					if _, ok := Descent.Cave.RChunks[i]; !ok {
						Descent.Cave.RChunks[i] = cave.NewChunk(i, Descent.Cave)
						FillChunk(Descent.Cave.RChunks[i])
						Descent.Cave.UpdateBatch = true
						IncreaseLevelInf()
					}
				} else if i.X < 0 && i.Y >= 0 {
					if _, ok := Descent.Cave.LChunks[i]; !ok {
						Descent.Cave.LChunks[i] = cave.NewChunk(i, Descent.Cave)
						FillChunk(Descent.Cave.LChunks[i])
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
			Descent.ExitPop.Raw = "Mark all the remaining bombs to exit."
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
	Descent.Cave.FuseLen -= 0.1
	if Descent.Cave.FuseLen < 0.4 {
		Descent.Cave.FuseLen = 0.4
	}
}

func FillChunk(ch *cave.Chunk) {
	for _, row := range ch.Rows {
		for _, tile := range row {
			if tile.Solid() && tile.Breakable() && (tile.Fillable || !ch.Cave.Finite) {
				if tile.Bomb {
					tile.Entity = &Bomb{
						Tile:       tile,
						FuseLength: tile.Chunk.Cave.FuseLen,
					}
					CaveTotalBombs++
					CaveBombsLeft++
					tile.XRay = img.Batchers[constants.EntityKey].Sprites["bomb_fuse"]
				} else if random.CaveGen.Intn(ch.Cave.GemRate) == 0 && tile.Solid() && tile.Breakable() {
					collect := Collectibles[GemDiamond]
					tile.Entity = &CollectibleItem{
						Collect: collect,
					}
					tile.XRay = collect.Sprite
				} else if random.CaveGen.Intn(75) == 0 && tile.Solid() && tile.Breakable() {
					tile.Entity = &MadMonk{}
				}
			} else if tile.Bomb {
				tile.Bomb = false
			}
		}
	}
}

func StartMinesweeper(c *cave.Cave, t *cave.Tile) {
	nb := false
	first := true
	for !nb || first {
		for _, ch := range c.LChunks {
			nb = FillMinesweeper(ch, t, nb)
			if !first && !nb {
				return
			}
		}
		for _, ch := range c.RChunks {
			nb = FillMinesweeper(ch, t, nb)
			if !first && !nb {
				return
			}
		}
		first = false
	}
}

func FillMinesweeper(ch *cave.Chunk, t *cave.Tile, nb bool) bool {
	needBomb := nb
	for _, row := range ch.Rows {
		for _, tile := range row {
			if tile.Solid() && tile.Breakable() && tile.Fillable && t.RCoords != tile.RCoords {
				if tile.Bomb {
					tile.Entity = &Mine{
						Tile:       tile,
					}
					tile.XRay = img.Batchers[constants.EntityKey].Sprites["mine_1"]
				} else if needBomb {
					tile.Bomb = true
					tile.Entity = &Mine{
						Tile:       tile,
					}
					tile.XRay = img.Batchers[constants.EntityKey].Sprites["mine_1"]
					needBomb = false
				}
			} else if tile.Bomb {
				tile.Bomb = false
				needBomb = true
			}
		}
	}
	return needBomb
}