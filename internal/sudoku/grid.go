package sudoku

import "math/bits"

// Grid represents the 9x9 Sudoku board
type Grid struct {
	Values     [81]byte
	Candidates [81]uint16
	Clues      int
}

// CandidatesCount returns the number of valid candidates for a cell
func (g *Grid) CandidatesCount(i int) int {
	return bits.OnesCount16(g.Candidates[i])
}
