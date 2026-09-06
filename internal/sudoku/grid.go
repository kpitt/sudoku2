package sudoku

import (
	"fmt"
	"math/bits"
	"strings"
)

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

// Display prints the grid to standard output according to the specified style.
// Supported styles:
//   - "raw", "string", "line": 81 characters with '.' for empty cells, followed by a newline
//   - "grid": human-readable 9x9 grid with 3x3 box borders (default)
func (g *Grid) Display(style string) {
	fmt.Print(g.Format(style))
}

// Format returns the string representation of the grid for the specified style.
func (g *Grid) Format(style string) string {
	switch strings.ToLower(strings.TrimSpace(style)) {
	case "raw", "string", "line":
		var sb strings.Builder
		sb.Grow(82)
		for i := range 81 {
			if g.Values[i] == 0 {
				sb.WriteByte('.')
			} else {
				sb.WriteByte('0' + g.Values[i])
			}
		}
		sb.WriteByte('\n')
		return sb.String()

	case "grid", "":
		fallthrough
	default:
		var sb strings.Builder
		sb.Grow(350)
		border := "+-------+-------+-------+\n"
		for r := range 9 {
			if r%3 == 0 {
				sb.WriteString(border)
			}
			for c := range 9 {
				if c%3 == 0 {
					sb.WriteString("| ")
				}
				i := r*9 + c
				if g.Values[i] == 0 {
					sb.WriteByte('.')
				} else {
					sb.WriteByte('0' + g.Values[i])
				}
				sb.WriteByte(' ')
			}
			sb.WriteString("|\n")
		}
		sb.WriteString(border)
		return sb.String()
	}
}
