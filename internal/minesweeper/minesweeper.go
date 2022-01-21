package minesweeper

import (
	"fmt"
	"math/rand"
)

type Board struct {
	board    [][]bool
	nums     [][]int
	revealed [][]bool
}

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
		var row []bool
		for j := 0; j < w; j++ {
			row = append(row, list[t])
			t++
		}
		board.board = append(board.board, row)
	}
	for y, row := range board.board {
		var numRow []int
		for x, bomb := range row {
			if bomb {
				numRow = append(numRow, -1)
			} else {
				c := 0
				if x > 0 && y > 0 && board.board[y-1][x-1] {
					c++
				}
				if y > 0 && board.board[y-1][x] {
					c++
				}
				if x < w-1 && y > 0 && board.board[y-1][x+1] {
					c++
				}
				if x > 0 && board.board[y][x-1] {
					c++
				}
				if x < w-1 && board.board[y][x+1] {
					c++
				}
				if x > 0 && y < h-1 && board.board[y+1][x-1] {
					c++
				}
				if y < h-1 && board.board[y+1][x] {
					c++
				}
				if x < w-1 && y < h-1 && board.board[y+1][x+1] {
					c++
				}
				numRow = append(numRow, c)
			}
		}
		board.nums = append(board.nums, numRow)
	}
	board.revealed = make([][]bool, h)
	for i := range board.revealed {
		board.revealed[i] = make([]bool, w)
	}
	return &board
}

func (b *Board) AsArray() []bool {
	var result []bool
	for _, row := range b.board {
		for _, bomb := range row {
			result = append(result, bomb)
		}
	}
	return result
}

func (b *Board) PrintToTerminal() {
	fmt.Println("Printing board ... ")
	fmt.Println()
	for y, row := range b.nums {
		for x, cell := range row {
			if b.revealed[y][x] {
				if cell == -1 {
					fmt.Print("ó")
				} else {
					fmt.Print(cell)
				}
			} else {
				fmt.Print("□")
			}
		}
		fmt.Print("\n")
	}
}

func (b *Board) PrintToTerminalFull() {
	fmt.Println("Printing full board ... ")
	fmt.Println()
	for _, row := range b.nums {
		for _, cell := range row {
			if cell == -1 {
				fmt.Print("ó")
			} else {
				fmt.Print(cell)
			}
		}
		fmt.Print("\n")
	}
}
