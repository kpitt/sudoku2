package formatter

import (
	"math/bits"
	"strings"

	"github.com/kpitt/sudoku2/internal/solver"
)

// FormatSS returns the board in "Simple Sudoku" ASCII format.
func FormatSS(b *solver.Board) string {
	var sb strings.Builder

	border := "*-----------*"
	divider := "|---+---+---|"

	sb.WriteString(border + "\n")

	for r := range 9 {
		if r > 0 && r%3 == 0 {
			sb.WriteString(divider + "\n")
		}

		sb.WriteString("|")

		for c := range 9 {
			if c > 0 && c%3 == 0 {
				sb.WriteString("|")
			}

			cellIdx := r*9 + c
			mask := b.Cells[cellIdx]

			if bits.OnesCount16(mask) == 1 {
				digit := bits.TrailingZeros16(mask) + 1
				sb.WriteByte(byte('0' + digit))
			} else {
				sb.WriteString(".")
			}
		}

		sb.WriteString("|\n")
	}

	sb.WriteString(border + "\n")

	return sb.String()
}
