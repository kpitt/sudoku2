package solver

import (
	"fmt"
	"strconv"
)

// Board represents a 9x9 Sudoku board.
type Board [9][9]int

// ParseBoard converts an 81-character string into a Board.
func ParseBoard(input string) (Board, error) {
	if len(input) != 81 {
		return Board{}, fmt.Errorf("invalid input length: expected 81, got %d", len(input))
	}

	var board Board
	for i, char := range input {
		if char < '0' || char > '9' {
			return Board{}, fmt.Errorf("invalid character at index %d: %c", i, char)
		}
		val, _ := strconv.Atoi(string(char))
		board[i/9][i%9] = val
	}

	return board, nil
}

// IsValid checks if the current board state follows Sudoku rules.
// It returns true if no rules are violated. Empty cells (0) are ignored.
func (b *Board) IsValid() bool {
	// Check rows
	for r := 0; r < 9; r++ {
		if hasDuplicates(b.getRow(r)) {
			return false
		}
	}

	// Check columns
	for c := 0; c < 9; c++ {
		if hasDuplicates(b.getCol(c)) {
			return false
		}
	}

	// Check 3x3 boxes
	for r := 0; r < 9; r += 3 {
		for c := 0; c < 9; c += 3 {
			if hasDuplicates(b.getBox(r, c)) {
				return false
			}
		}
	}

	return true
}

// IsSolved checks if the board is completely filled and valid.
func (b *Board) IsSolved() bool {
	if !b.IsValid() {
		return false
	}
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b[r][c] == 0 {
				return false
			}
		}
	}
	return true
}

// SolveBacktracking attempts to solve the Sudoku puzzle using a backtracking algorithm.
// It modifies the board in-place and returns true if a solution is found.
// This is used for the 'check' command to verify solutions and uniqueness.
func (b *Board) SolveBacktracking() bool {
	// If the board is initially invalid, it's unsolvable.
	if !b.IsValid() {
		return false
	}
	return b.backtrack()
}

// SolveDeductive attempts to solve the Sudoku puzzle using deductive techniques.
// It mimics human solving strategies (Naked Singles, Hidden Singles, etc.).
func (b *Board) SolveDeductive() bool {
	if !b.IsValid() {
		return false
	}

	for {
		changed := false

		// 1. Try Naked Singles
		if b.applyNakedSingles() {
			changed = true
		}

		// 2. Try Hidden Singles
		if b.applyHiddenSingles() {
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
func (b *Board) PossibleValues(row, col int) []int {
	if b[row][col] != 0 {
		return nil
	}

	possible := make([]int, 0, 9)
	for num := 1; num <= 9; num++ {
		if b.isSafe(row, col, num) {
			possible = append(possible, num)
		}
	}
	return possible
}

// applyNakedSingles finds cells that have only one possible candidate.
func (b *Board) applyNakedSingles() bool {
	changed := false
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b[r][c] == 0 {
				p := b.PossibleValues(r, c)
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
func (b *Board) applyHiddenSingles() bool {
	changed := false

	// Check rows
	for r := 0; r < 9; r++ {
		for num := 1; num <= 9; num++ {
			count := 0
			lastCol := -1
			for c := 0; c < 9; c++ {
				if b[r][c] == 0 && b.isSafe(r, c, num) {
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
	for c := 0; c < 9; c++ {
		for num := 1; num <= 9; num++ {
			count := 0
			lastRow := -1
			for r := 0; r < 9; r++ {
				if b[r][c] == 0 && b.isSafe(r, c, num) {
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
						if b[r][c] == 0 && b.isSafe(r, c, num) {
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

func (b *Board) backtrack() bool {
	row, col, found := b.findEmptyCell()
	if !found {
		return true // No empty cells left, solved!
	}

	for num := 1; num <= 9; num++ {
		if b.isSafe(row, col, num) {
			b[row][col] = num
			if b.backtrack() {
				return true
			}
			b[row][col] = 0 // Backtrack
		}
	}

	return false
}

func (b *Board) findEmptyCell() (int, int, bool) {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b[r][c] == 0 {
				return r, c, true
			}
		}
	}
	return 0, 0, false
}

func (b *Board) isSafe(row, col, num int) bool {
	// Check row
	for c := 0; c < 9; c++ {
		if b[row][c] == num {
			return false
		}
	}

	// Check column
	for r := 0; r < 9; r++ {
		if b[r][col] == num {
			return false
		}
	}

	// Check box
	startRow, startCol := row-row%3, col-col%3
	for r := startRow; r < startRow+3; r++ {
		for c := startCol; c < startCol+3; c++ {
			if b[r][c] == num {
				return false
			}
		}
	}

	return true
}

func (b *Board) getRow(r int) []int {
	return b[r][:]
}

func (b *Board) getCol(c int) []int {
	col := make([]int, 9)
	for r := 0; r < 9; r++ {
		col[r] = b[r][c]
	}
	return col
}

func (b *Board) getBox(startRow, startCol int) []int {
	box := make([]int, 0, 9)
	for r := startRow; r < startRow+3; r++ {
		for c := startCol; c < startCol+3; c++ {
			box = append(box, b[r][c])
		}
	}
	return box
}

func hasDuplicates(nums []int) bool {
	seen := make(map[int]bool)
	for _, n := range nums {
		if n == 0 {
			continue
		}
		if seen[n] {
			return true
		}
		seen[n] = true
	}
	return false
}

// String returns a human-readable representation of the board.
func (b *Board) String() string {
	var out string
	for r := 0; r < 9; r++ {
		if r > 0 && r%3 == 0 {
			out += "------+-------+------\n"
		}
		for c := 0; c < 9; c++ {
			if c > 0 && c%3 == 0 {
				out += "| "
			}
			if b[r][c] == 0 {
				out += ". "
			} else {
				out += fmt.Sprintf("%d ", b[r][c])
			}
		}
		out += "\n"
	}
	return out
}
