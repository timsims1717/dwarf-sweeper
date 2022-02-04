package minesweeper

type Board struct {
	Board [][]Cell
}

type Cell struct {
	Num  int
	Bomb bool
	Flag bool
	Ex   bool
	Rev  bool
}

func (b *Board) cellIsReal(x, y int) bool {
	return !(x < 0 || y < 0 || x > len(b.Board[0])-1 || y > len(b.Board)-1)
}

func (b *Board) cellIsNotBomb(x, y int) bool {
	if b.cellIsReal(x, y) {
		return b.Board[y][x].Rev && !b.Board[y][x].Bomb
	}
	return false
}

func (b *Board) cellIsBomb(x, y int) bool {
	if b.cellIsReal(x, y) {
		return b.Board[y][x].Bomb
	}
	return false
}

func (b *Board) ex(x, y int) bool {
	if b.cellIsReal(x, y) {
		return b.Board[y][x].Ex
	}
	return false
}

func (b *Board) reveal(x, y int) {
	if b.cellIsReal(x, y) && !b.Board[y][x].Rev {
		b.Board[y][x].Rev = true
		b.Board[y][x].Ex = true
	}
}

func (b *Board) hide(x, y int) {
	if b.cellIsReal(x, y) {
		b.Board[y][x].Rev = false
		b.Board[y][x].Ex = false
	}
}

func (b *Board) neighborsAreRevealed(x, y int) bool {
	if !b.cellIsRevealed(x-1, y-1) {
		return false
	}
	if !b.cellIsRevealed(x, y-1) {
		return false
	}
	if !b.cellIsRevealed(x+1, y-1) {
		return false
	}
	if !b.cellIsRevealed(x-1, y) {
		return false
	}
	if !b.cellIsRevealed(x+1, y) {
		return false
	}
	if !b.cellIsRevealed(x-1, y+1) {
		return false
	}
	if !b.cellIsRevealed(x, y+1) {
		return false
	}
	if !b.cellIsRevealed(x+1, y+1) {
		return false
	}
	return true
}

func (b *Board) neighborsRevealed(x, y int) (int, int) {
	bc := 0
	rc := 0
	if b.cellIsRevealed(x-1, y-1) {
		if b.cellIsReal(x-1, y-1) && b.Board[y-1][x-1].Bomb {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x, y-1) {
		if b.cellIsReal(x, y-1) && b.Board[y-1][x].Bomb {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x+1, y-1) {
		if b.cellIsReal(x+1, y-1) && b.Board[y-1][x+1].Bomb {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x-1, y) {
		if b.cellIsReal(x-1, y) && b.Board[y][x-1].Bomb {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x+1, y) {
		if b.cellIsReal(x+1, y) && b.Board[y][x+1].Bomb {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x-1, y+1) {
		if b.cellIsReal(x-1, y+1) && b.Board[y+1][x-1].Bomb {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x, y+1) {
		if b.cellIsReal(x, y+1) && b.Board[y+1][x].Bomb {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x+1, y+1) {
		if b.cellIsReal(x+1, y+1) && b.Board[y+1][x+1].Bomb {
			bc++
		}
		rc++
	}
	return rc, bc
}

func (b *Board) cellIsRevealed(x, y int) bool {
	if x < 0 || y < 0 || x > len(b.Board[0])-1 || y > len(b.Board)-1 {
		return true
	}
	return b.Board[y][x].Rev
}

func (b *Board) revealNeighbors(x, y int) {
	b.reveal(x-1, y-1)
	b.reveal(x, y-1)
	b.reveal(x+1, y-1)
	b.reveal(x-1, y)
	b.reveal(x+1, y)
	b.reveal(x-1, y+1)
	b.reveal(x, y+1)
	b.reveal(x+1, y+1)
}