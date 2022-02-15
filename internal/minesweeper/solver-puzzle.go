package minesweeper

import (
	"fmt"
	"math/rand"
)

// can the board be solved without guessing?
func SolvableP(b *Board, print bool) bool {
	if print {
		fmt.Println("State of the board:")
		b.PrintToTerminal()
	}
	for {
		found := false
		for y, row := range b.Board {
			for x := range row {
				if b.cellIsNotBomb(x, y) && !b.neighborsAreRevealed(x, y) && !b.ex(x, y) {
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
			if print {
				fmt.Println("failed to reveal new cells")
			}
			break
		}
		if print {
			fmt.Println("After a pass:")
			b.PrintToTerminal()
		}
	}
	allRevealed := true
outer:
	for _, row := range b.Board {
		for _, cell := range row {
			if !cell.Rev && !cell.Bomb {
				allRevealed = false
				break outer
			}
		}
	}
	return allRevealed
}

func RevealTilSolvableP(b *Board, rando *rand.Rand, print bool) {
	tries := 0
	for !SolvableP(b.Copy(), print) && tries < 200 {
		x := rando.Intn(len(b.Board[0]))
		y := rando.Intn(len(b.Board))
		for b.Board[y][x].Rev || b.Board[y][x].Bomb {
			x = rando.Intn(len(b.Board[0]))
			y = rando.Intn(len(b.Board))
			tries++
			if tries >= 200 {
				return
			}
		}
		b.Board[y][x].Rev = true
		b.Board[y][x].Ex = false
	}
	if print {
		fmt.Println("Final state:")
		b.PrintToTerminal()
	}
}

func UnRevealWhileSolvableP(b *Board, rando *rand.Rand, print bool) {
	tries := 0
	for tries < 50 {
		x := rando.Intn(len(b.Board[0]))
		y := rando.Intn(len(b.Board))
		for !b.Board[y][x].Rev || b.Board[y][x].Bomb {
			x = rando.Intn(len(b.Board[0]))
			y = rando.Intn(len(b.Board))
			tries++
			if tries >= 50 {
				return
			}
		}
		b.Board[y][x].Rev = false
		b.Board[y][x].Ex = false
		if !SolvableP(b.Copy(), print) {
			b.Board[y][x].Rev = true
			b.Board[y][x].Ex = false
			break
		}
	}
	if print {
		fmt.Println("Final state:")
		b.PrintToTerminal()
	}
}
