package structures

import (
	"dwarf-sweeper/internal/descent"
	"dwarf-sweeper/internal/descent/cave"
	"dwarf-sweeper/internal/descent/player"
	"dwarf-sweeper/internal/random"
)

func BasicDestroy(t *cave.Tile) {
	gemRate := t.Chunk.Cave.GemRate*t.GemRate*descent.Descent.Player.GemRate
	if t.Bomb {
		descent.CreateBomb(t.Transform.Pos)
	} else if random.CaveGen.Float64() < gemRate {
		descent.CreateGem(t.Transform.Pos)
	}
}

func UpdateTiles(tiles []*cave.Tile) {
	for _, tile := range tiles {
		if tile.Solid() && tile.Breakable() {
			tile.DestroyTrigger = BasicDestroy
			if tile.Bomb {
				player.CaveTotalBombs++
				player.CaveBombsLeft++
				tile.XRay = "bomb"
			}
		}
	}
}

func FillCave(c *cave.Cave) {
	for _, ch := range c.Chunks {
		FillBasic(ch)
	}
}

func FillBasic(ch *cave.Chunk) {
	for _, row := range ch.Rows {
		for _, tile := range row {
			if tile.Solid() && tile.Breakable() {
				tile.DestroyTrigger = BasicDestroy
				if tile.Bomb {
					player.CaveTotalBombs++
					player.CaveBombsLeft++
					tile.XRay = "bomb"
				} else if random.CaveGen.Intn(80) == 0 {
					if ch.Cave.Biome == "mine" {
						//tile.Entity = &descent.Slug{}
					} else {
						//tile.Entity = &descent.MadMonk{}
					}
				} else if random.CaveGen.Intn(50) == 0 && tile.Solid() && tile.Breakable() {
					//tile.Entity = &descent.Bat{}
				}
			}
		}
	}
}

func StartMinesweeper(c *cave.Cave, t *cave.Tile) {
	nb := false
	first := true
	for !nb || first {
		for _, ch := range c.Chunks {
			nb = FillMinesweeper(ch, t, nb)
			if !first && !nb {
				return
			}
		}
		first = false
	}
}

func MineDestroy(t *cave.Tile) {
	descent.CreateMine(t.Transform.Pos)
}

func FillMinesweeper(ch *cave.Chunk, t *cave.Tile, nb bool) bool {
	needBomb := nb
	for _, row := range ch.Rows {
		for _, tile := range row {
			if tile.Solid() && tile.Breakable() {
				if t.RCoords != tile.RCoords {
					if tile.Bomb {
						tile.DestroyTrigger = MineDestroy
						tile.XRay = "mine"
					} else if needBomb {
						tile.Bomb = true
						tile.DestroyTrigger = MineDestroy
						tile.XRay = "mine"
						needBomb = false
					}
				} else if tile.Bomb {
					tile.Bomb = false
					needBomb = true
				}
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
