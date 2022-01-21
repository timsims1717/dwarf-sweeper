package minesweeper

import (
	"fmt"
	"math/rand"
)

// can the board be solved without guessing?
func (b *Board) Solvable() bool {
	revealed := b.revealed
	for {
		found := false
		for y, row := range b.board {
			for x := range row {
				if b.cellIsNotBomb(x, y) && !b.neighborsAreRevealed(x, y) {
					// candidate for examination
					num := b.nums[y][x]
					if num == 0 {
						b.revealNeighbors(x, y)
						found = true
					} else {
						rc, bc := b.neighborsRevealed(x, y)
						if rc-bc == 8-num {
							b.revealNeighbors(x, y)
							found = true
						} else if bc == num {
							b.revealNeighbors(x, y)
							found = true
						}
					}
				}
			}
		}
		if !found {
			fmt.Println("failed to reveal new cells")
			break
		}
		b.PrintToTerminal()
	}
	allRevealed := true
outer:
	for _, row := range b.revealed {
		for _, cell := range row {
			if !cell {
				allRevealed = false
				break outer
			}
		}
	}
	b.revealed = revealed
	return allRevealed
}

func (b *Board) neighborsRevealed(x, y int) (int, int) {
	bc := 0
	rc := 0
	if b.cellIsRevealed(x-1, y-1) {
		if b.cellIsReal(x-1, y-1) && b.board[y-1][x-1] {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x, y-1) {
		if b.cellIsReal(x, y-1) && b.board[y-1][x] {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x+1, y-1) {
		if b.cellIsReal(x+1, y-1) && b.board[y-1][x+1] {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x-1, y) {
		if b.cellIsReal(x-1, y) && b.board[y][x-1] {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x+1, y) {
		if b.cellIsReal(x+1, y) && b.board[y][x+1] {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x-1, y+1) {
		if b.cellIsReal(x-1, y+1) && b.board[y+1][x-1] {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x, y+1) {
		if b.cellIsReal(x, y+1) && b.board[y+1][x] {
			bc++
		}
		rc++
	}
	if b.cellIsRevealed(x+1, y+1) {
		if b.cellIsReal(x+1, y+1) && b.board[y+1][x+1] {
			bc++
		}
		rc++
	}
	return rc, bc
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

func (b *Board) cellIsReal(x, y int) bool {
	return !(x < 0 || y < 0 || x > len(b.board[0])-1 || y > len(b.board)-1)
}

func (b *Board) cellIsRevealed(x, y int) bool {
	if x < 0 || y < 0 || x > len(b.board[0])-1 || y > len(b.board)-1 {
		return true
	}
	return b.revealed[y][x]
}

func (b *Board) cellIsNotBomb(x, y int) bool {
	if x < 0 || y < 0 || x > len(b.board[0])-1 || y > len(b.board)-1 {
		return true
	}
	return b.revealed[y][x] && !b.board[y][x]
}

func (b *Board) reveal(x, y int) {
	if !(x < 0 || y < 0 || x > len(b.board[0])-1 || y > len(b.board)-1) {
		b.revealed[y][x] = true
	}
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

func (b *Board) RevealTilSolvable(rando *rand.Rand) {
	tries := 0
	for !b.Solvable() && tries < 200 {
		x := rando.Intn(len(b.board[0]))
		y := rando.Intn(len(b.board))
		for b.revealed[y][x] && !b.board[y][x] && tries < 200 {
			x = rando.Intn(len(b.board[0]))
			y = rando.Intn(len(b.board))
			tries++
		}
		b.revealed[y][x] = true
	}
}
