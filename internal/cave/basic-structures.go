package cave

func CarveEntrance(chunk *Chunk) {
	for y, row := range chunk.Rows {
		for x, tile := range row {
			// entrance room
			if x > 11 && x < 21 && y > 4 && y < 10 {
				tile.Solid = false
				tile.bomb = false
				if x == 16 && y == 9 {
					tile.Type = Deco
					tile.BGSpriteS = "door"
					tile.BGSprite = chunk.Cave.batcher.Sprites["door"]
					tile.breakable = false
				} else if x == 15 && y == 9 {
					tile.Type = Deco
					tile.BGSpriteS = "door_l"
					tile.BGSprite = chunk.Cave.batcher.Sprites["door_l"]
					tile.breakable = false
				} else if x == 17 && y == 9 {
					tile.Type = Deco
					tile.BGSpriteS = "door_r"
					tile.BGSprite = chunk.Cave.batcher.Sprites["door_r"]
					tile.breakable = false
				} else if x == 16 && y == 8 {
					tile.Type = Deco
					tile.BGSpriteS = "door_t"
					tile.BGSprite = chunk.Cave.batcher.Sprites["door_t"]
					tile.breakable = false
				} else if x == 15 && y == 8 {
					tile.Type = Deco
					tile.BGSpriteS = "door_tl"
					tile.BGSprite = chunk.Cave.batcher.Sprites["door_tl"]
					tile.breakable = false
				} else if x == 17 && y == 8 {
					tile.Type = Deco
					tile.BGSpriteS = "door_tr"
					tile.BGSprite = chunk.Cave.batcher.Sprites["door_tr"]
					tile.breakable = false
				} else {
					tile.destroyed = true
					tile.Type = Empty
					tile.BGSprite = nil
				}
				tile.Entities = []Entity{}
			} else if x > 10 && x < 22 && y > 3 && y < 11 {
				if x > 14 && x < 18 && y == 10 {
					tile.Type = Wall
					tile.breakable = false
				}
				tile.bomb = false
				tile.Entities = []Entity{}
			}
		}
	}
}