package board

// Board represents a 9x9 Sudoku board.
type Board [9][9]int

// IsValid checks if the current board state follows Sudoku rules.
// It returns true if no rules are violated. Empty cells (0) are ignored.
func (b *Board) IsValid() bool {
	// Check rows
	for r := range 9 {
		if hasDuplicates(b.GetRow(r)) {
			return false
		}
	}

	// Check columns
	for c := range 9 {
		if hasDuplicates(b.GetCol(c)) {
			return false
		}
	}

	// Check 3x3 boxes
	for r := 0; r < 9; r += 3 {
		for c := 0; c < 9; c += 3 {
			if hasDuplicates(b.GetBox(r, c)) {
				return false
			}
		}
	}

	return true
}

// GetRow returns all values in the specified row.
func (b *Board) GetRow(r int) []int {
	return b[r][:]
}

// GetCol returns all values in the specified column.
func (b *Board) GetCol(c int) []int {
	col := make([]int, 9)
	for r := range 9 {
		col[r] = b[r][c]
	}
	return col
}

// GetBox returns all values in the 3x3 box starting at the specified row and column.
func (b *Board) GetBox(startRow, startCol int) []int {
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

// IsSolved checks if the board is completely filled and valid.
func (b *Board) IsSolved() bool {
	if !b.IsValid() {
		return false
	}

	for r := range 9 {
		for c := range 9 {
			if b[r][c] == 0 {
				return false
			}
		}
	}

	return true
}

// IsSafe checks if placing num at (row, col) violates any Sudoku rules.
func (b *Board) IsSafe(row, col, num int) bool {
	// Check row
	for c := range 9 {
		if b[row][c] == num {
			return false
		}
	}

	// Check column
	for r := range 9 {
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
