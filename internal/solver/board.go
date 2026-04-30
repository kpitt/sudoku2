package solver

import "errors"

// Board represents a 9x9 Sudoku board.
type Board [9][9]int

// ParseBoard converts an 81-character string into a Board.
func ParseBoard(input string) (Board, error) {
	return Board{}, errors.New("not implemented")
}

// IsValid checks if the current board state follows Sudoku rules.
func (b *Board) IsValid() bool {
	return false
}

