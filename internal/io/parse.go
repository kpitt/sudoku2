package io

import (
	"fmt"
	"strconv"

	"github.com/kpitt/sudoku2/internal/board"
)

// ParseBoard converts an 81-character string into a board.Board.
func ParseBoard(input string) (board.Board, error) {
	if len(input) != 81 {
		return board.Board{}, fmt.Errorf("invalid input length: expected 81, got %d", len(input))
	}

	var b board.Board
	for i, char := range input {
		if char < '0' || char > '9' {
			return board.Board{}, fmt.Errorf("invalid character at index %d: %c", i, char)
		}

		val, _ := strconv.Atoi(string(char))
		b[i/9][i%9] = val
	}

	return b, nil
}
