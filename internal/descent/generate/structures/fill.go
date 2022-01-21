package structures

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
)

func BasicDestroy(t *cave.Tile) {
	gemRate := t.Chunk.Cave.GemRate*t.GemRate*descent.Descent.Player.GemRate
	if t.Bomb {
		descent.CreateBomb(t.Transform.Pos)
	} else if random.CaveGen.Float64() < gemRate {
		descent.CreateCollectible(t.Transform.Pos, descent.GemDiamond)
	}
}

func UpdateTiles(tiles []*cave.Tile) {
	for _, tile := range tiles {
		if tile.Solid() && tile.Breakable() {
			tile.DestroyTrigger = BasicDestroy
			if tile.Bomb {
				descent.CaveTotalBombs++
				descent.CaveBombsLeft++
				tile.XRay = img.Batchers[constants.EntityKey].Sprites["bomb_fuse"]
			}
		}
	}
}

func FillCave(c *cave.Cave) {
	for _, ch := range c.LChunks {
		FillBasic(ch)
	}
	for _, ch := range c.RChunks {
		FillBasic(ch)
	}
}

func FillBasic(ch *cave.Chunk) {
	for _, row := range ch.Rows {
		for _, tile := range row {
			if tile.Solid() && tile.Breakable() {
				tile.DestroyTrigger = BasicDestroy
				if tile.Bomb {
					descent.CaveTotalBombs++
					descent.CaveBombsLeft++
					tile.XRay = img.Batchers[constants.EntityKey].Sprites["bomb_fuse"]
				} else if random.CaveGen.Intn(80) == 0 {
					switch random.CaveGen.Intn(2) {
					case 0:
						tile.Entity = &descent.Slug{}
					case 1:
						tile.Entity = &descent.MadMonk{}
					}
				} else if random.CaveGen.Intn(50) == 0 && tile.Solid() && tile.Breakable() {
					tile.Entity = &descent.Bat{}
				}
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
			if tile.Solid() && tile.Breakable() && t.RCoords != tile.RCoords {
				if tile.Bomb {
					tile.Entity = &descent.Mine{
						Tile: tile,
					}
					tile.XRay = img.Batchers[constants.EntityKey].Sprites["mine_1"]
				} else if needBomb {
					tile.Bomb = true
					tile.Entity = &descent.Mine{
						Tile: tile,
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

func FillChunkWall(ch *cave.Chunk) {
	for _, row := range ch.Rows {
		for _, tile := range row {
			ToBlock(tile, cave.Wall, false, false)
		}
	}
}
