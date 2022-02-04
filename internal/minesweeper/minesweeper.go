package minesweeper

import (
	"fmt"
	"math/rand"
)

func CreateBoard(w, h, amt int, rando *rand.Rand) *Board {
	list := make([]bool, w*h)
	for i := 0; i < amt; i++ {
		list[i] = true
	}
	// randomize list
	for i := len(list) - 1; i > 0; i-- {
		j := rando.Intn(i)
		list[i], list[j] = list[j], list[i]
	}
	var board Board
	t := 0
	for i := 0; i < h; i++ {
		var row []Cell
		for j := 0; j < w; j++ {
			row = append(row, Cell{
				Bomb: list[t],
			})
			t++
		}
		board.Board = append(board.Board, row)
	}
	for y, row := range board.Board {
		for x, cell := range row {
			if cell.Bomb {
				board.Board[y][x].Num = -1
			} else {
				c := 0
				if board.cellIsBomb(x, y+1) {
					c++
				}
				if board.cellIsBomb(x+1, y+1) {
					c++
				}
				if board.cellIsBomb(x+1, y) {
					c++
				}
				if board.cellIsBomb(x+1, y-1) {
					c++
				}
				if board.cellIsBomb(x, y-1) {
					c++
				}
				if board.cellIsBomb(x-1, y-1) {
					c++
				}
				if board.cellIsBomb(x-1, y) {
					c++
				}
				if board.cellIsBomb(x-1, y+1) {
					c++
				}
				board.Board[y][x].Num = c
			}
		}
	}
	return &board
}

func (b *Board) AsArray() []bool {
	var result []bool
	for _, row := range b.Board {
		for _, cell := range row {
			result = append(result, cell.Bomb)
		}
	}
	return result
}

func (b *Board) PrintToTerminal() {
	fmt.Println("Printing board ... ")
	for _, row := range b.Board {
		for _, cell := range row {
			if cell.Rev {
				if cell.Bomb {
					fmt.Print("ó")
				} else if cell.Ex {
					fmt.Print("X")
				} else {
					fmt.Print(cell.Num)
				}
			} else {
				fmt.Print("□")
			}
		}
		fmt.Print("\n")
	}
	fmt.Println()
}

func (b *Board) PrintToTerminalFull() {
	fmt.Println("Printing full board ... ")
	for _, row := range b.Board {
		for _, cell := range row {
			if cell.Bomb {
				fmt.Print("ó")
			} else {
				fmt.Print(cell.Num)
			}
		}
		fmt.Print("\n")
	}
	fmt.Println()
}

func (b *Board) Copy() *Board {
	c := &Board{}
	for i, row := range b.Board {
		c.Board = append(c.Board, []Cell{})
		for _, cell := range row {
			c.Board[i] = append(c.Board[i], Cell{
				Num:  cell.Num,
				Bomb: cell.Bomb,
				Ex:   cell.Ex,
				Rev:  cell.Rev,
			})
		}
	}
	return c
}