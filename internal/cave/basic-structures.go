package cave

func CarveEntrance(chunk *Chunk) {
	for y, row := range chunk.Rows {
		for x, tile := range row {
			// entrance room
			if x > 11 && x < 21 && y > 4 && y < 10 {
				tile.Solid = false
				tile.destroyed = true
				tile.bomb = false
				tile.BGSprite = nil
			} else if x > 10 && x < 22 && y > 5 && y < 11 {
				tile.bomb = false
			}
		}
	}
}