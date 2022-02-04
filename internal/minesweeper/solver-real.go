package minesweeper

import (
	"fmt"
	"math/rand"
)

// can the board be solved without guessing?
func Solvable(b Board) bool {
	for {
		found := false
		for y, row := range b.Board {
			for x := range row {
				if b.cellIsNotBomb(x, y) && !b.neighborsAreRevealed(x, y) {
					// candidate for examination
					num := b.Board[y][x].Num
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
	for _, row := range b.Board {
		for _, cell := range row {
			if !cell.Rev {
				allRevealed = false
				break outer
			}
		}
	}
	return allRevealed
}

func RevealTilSolvable(b *Board, rando *rand.Rand) {
	tries := 0
	for !Solvable(*b) && tries < 200 {
		x := rando.Intn(len(b.Board[0]))
		y := rando.Intn(len(b.Board))
		for b.Board[y][x].Rev && !b.Board[y][x].Bomb && tries < 200 {
			x = rando.Intn(len(b.Board[0]))
			y = rando.Intn(len(b.Board))
			tries++
		}
		b.Board[y][x].Rev = true
	}
}
