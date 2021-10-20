package generate

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
)

func FillChunk(ch *cave.Chunk) {
	for _, row := range ch.Rows {
		for _, tile := range row {
			if tile.Solid() && tile.Breakable() && (tile.Fillable || !ch.Cave.Finite) {
				if tile.Bomb {
					tile.Entity = &descent.Bomb{
						Tile:       tile,
						FuseLength: tile.Chunk.Cave.FuseLen,
					}
					descent.CaveTotalBombs++
					descent.CaveBombsLeft++
					tile.XRay = img.Batchers[constants.EntityKey].Sprites["bomb_fuse"]
				} else if random.CaveGen.Intn(ch.Cave.GemRate) == 0 && tile.Solid() && tile.Breakable() {
					collect := descent.Collectibles[descent.GemDiamond]
					tile.Entity = &descent.CollectibleItem{
						Collect: collect,
					}
					tile.XRay = collect.Sprite
				} else if random.CaveGen.Intn(75) == 0 && tile.Solid() && tile.Breakable() {
					switch random.CaveGen.Intn(2) {
					case 0:
						tile.Entity = &descent.Slug{}
					case 1:
						tile.Entity = &descent.MadMonk{}
					}
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
					tile.Entity = &descent.Mine{
						Tile:       tile,
					}
					tile.XRay = img.Batchers[constants.EntityKey].Sprites["mine_1"]
				} else if needBomb {
					tile.Bomb = true
					tile.Entity = &descent.Mine{
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