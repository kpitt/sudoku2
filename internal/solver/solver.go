package solver

import (
	"errors"
	"fmt"

	"github.com/kpitt/sudoku2/internal/board"
)

// SolveBacktracking attempts to solve the Sudoku puzzle using a backtracking algorithm.
// It modifies the board in-place and returns true if a solution is found.
func SolveBacktracking(b *board.Board) bool {
	// If the board is initially invalid, it's unsolvable.
	if !b.IsValid() {
		return false
	}
	return backtrack(b)
}

// SolveDeductive attempts to solve the Sudoku puzzle using deductive techniques.
// It mimics human solving strategies (Naked Singles, Hidden Singles, etc.).
func SolveDeductive(b *board.Board) bool {
	if !b.IsValid() {
		return false
	}

	for {
		changed := false

		// 1. Try Naked Singles
		if applyNakedSingles(b) {
			changed = true
		}

		// 2. Try Hidden Singles
		if applyHiddenSingles(b) {
			changed = true
		}

		if !changed {
			break // No more progress can be made with current strategies
		}

		if b.IsSolved() {
			return true
		}
	}

	return b.IsSolved()
}

// PossibleValues returns a slice of numbers (1-9) that could legally be placed at (row, col).
func PossibleValues(b *board.Board, row, col int) []int {
	if b[row][col] != 0 {
		return nil
	}

	possible := make([]int, 0, 9)
	for num := 1; num <= 9; num++ {
		if b.IsSafe(row, col, num) {
			possible = append(possible, num)
		}
	}

	return possible
}

// Hint represents a single deductive move and its explanation.
type Hint struct {
	Row     int
	Col     int
	Value   int
	Message string
}

// GetHint identifies a single move that can be made using deductive logic.
func GetHint(b *board.Board) (*Hint, error) {
	if !b.IsValid() {
		return nil, errors.New("the current board is invalid")
	}

	// 1. Check for Naked Singles
	for r := range 9 {
		for c := range 9 {
			if b[r][c] == 0 {
				p := PossibleValues(b, r, c)
				if len(p) == 1 {
					return &Hint{
						Row:     r,
						Col:     c,
						Value:   p[0],
						Message: fmt.Sprintf("Naked Single: Cell (%d, %d) has only one possible value: %d", r+1, c+1, p[0]),
					}, nil
				}
			}
		}
	}

	// 2. Check for Hidden Singles (Rows)
	for r := range 9 {
		for num := 1; num <= 9; num++ {
			count := 0
			lastCol := -1
			for c := range 9 {
				if b[r][c] == 0 && b.IsSafe(r, c, num) {
					count++
					lastCol = c
				}
			}

			if count == 1 {
				return &Hint{
					Row:     r,
					Col:     lastCol,
					Value:   num,
					Message: fmt.Sprintf("Hidden Single in Row %d: %d can only go in cell (%d, %d)", r+1, num, r+1, lastCol+1),
				}, nil
			}
		}
	}

	// (Additional hidden single checks for columns/boxes can be added here)

	return nil, errors.New("no simple deductive hints found")
}

// applyNakedSingles finds cells that have only one possible candidate.
func applyNakedSingles(b *board.Board) bool {
	changed := false
	for r := range 9 {
		for c := range 9 {
			if b[r][c] == 0 {
				p := PossibleValues(b, r, c)
				if len(p) == 1 {
					b[r][c] = p[0]
					changed = true
				}
			}
		}
	}

	return changed
}

// applyHiddenSingles finds cells that are the only possible location for a number within a row, column, or box.
func applyHiddenSingles(b *board.Board) bool {
	changed := false

	// Check rows
	for r := range 9 {
		for num := 1; num <= 9; num++ {
			count := 0
			lastCol := -1
			for c := range 9 {
				if b[r][c] == 0 && b.IsSafe(r, c, num) {
					count++
					lastCol = c
				}
			}

			if count == 1 {
				b[r][lastCol] = num
				changed = true
			}
		}
	}

	// Check columns
	for c := range 9 {
		for num := 1; num <= 9; num++ {
			count := 0
			lastRow := -1
			for r := range 9 {
				if b[r][c] == 0 && b.IsSafe(r, c, num) {
					count++
					lastRow = r
				}
			}

			if count == 1 {
				b[lastRow][c] = num
				changed = true
			}
		}
	}

	// Check boxes
	for boxRow := 0; boxRow < 9; boxRow += 3 {
		for boxCol := 0; boxCol < 9; boxCol += 3 {
			for num := 1; num <= 9; num++ {
				count := 0
				lastR, lastC := -1, -1
				for r := boxRow; r < boxRow+3; r++ {
					for c := boxCol; c < boxCol+3; c++ {
						if b[r][c] == 0 && b.IsSafe(r, c, num) {
							count++
							lastR, lastC = r, c
						}
					}
				}

				if count == 1 {
					b[lastR][lastC] = num
					changed = true
				}
			}
		}
	}

	return changed
}

func backtrack(b *board.Board) bool {
	row, col, found := findEmptyCell(b)
	if !found {
		return true // No empty cells left, solved!
	}

	for num := 1; num <= 9; num++ {
		if b.IsSafe(row, col, num) {
			b[row][col] = num
			if backtrack(b) {
				return true
			}

			b[row][col] = 0 // Backtrack
		}
	}

	return false
}

func findEmptyCell(b *board.Board) (int, int, bool) {
	for r := range 9 {
		for c := range 9 {
			if b[r][c] == 0 {
				return r, c, true
			}
		}
	}

	return 0, 0, false
}
