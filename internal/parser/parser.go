package parser

import (
	"fmt"
	"strings"

	"github.com/kpitt/sudoku2/internal/solver"
)

// Parse converts an 81-character string into a Board state.
// It supports '0' or '.' for empty cells and '1'-'9' for solved cells.
func Parse(input string) (*solver.Board, error) {
	// Simple normalization: keep only digits and dots
	var sb strings.Builder

	for _, r := range input {
		if (r >= '0' && r <= '9') || r == '.' {
			sb.WriteRune(r)
		}
	}

	normalized := sb.String()

	if len(normalized) != 81 {
		return nil, fmt.Errorf("invalid puzzle length: expected 81, got %d", len(normalized))
	}

	b := &solver.Board{}

	// Initialize all cells with all 9 candidates (0x01FF)
	// Initialize states with all 9 candidates
	for i := range 81 {
		b.Cells[i] = 0x01FF
	}

	for i := range 9 {
		b.RowState[i] = 0x01FF
		b.ColState[i] = 0x01FF
		b.BoxState[i] = 0x01FF
	}

	// Set given digits
	for i, r := range normalized {
		var digit int
		if r >= '1' && r <= '9' {
			digit = int(r - '0')
		} else {
			continue
		}

		mask := uint16(1 << (digit - 1))

		// If cell already has this candidate removed, it's an invalid puzzle
		if b.Cells[i]&mask == 0 {
			return nil, fmt.Errorf("invalid puzzle: conflict at cell %d with digit %d", i, digit)
		}

		// Set the cell to the solved value
		b.Cells[i] = mask
		b.Resolved++

		// Update states
		rowIdx := solver.RowLUT[i]
		colIdx := solver.ColLUT[i]
		boxIdx := solver.BoxLUT[i]

		b.RowState[rowIdx] &^= mask
		b.ColState[colIdx] &^= mask
		b.BoxState[boxIdx] &^= mask

		// Propagate: remove this candidate from all peers
		for _, peerIdx := range solver.PeersLUT[i] {
			b.Cells[peerIdx] &^= mask
			// Note: We don't increment b.Resolved here if a peer becomes solved
			// to avoid complex recursive propagation during parsing.
			// The solver loop will handle cells that were solved by propagation.
		}
	}

	// Re-calculate Resolved count just in case propagation solved more cells
	// but the solver loop should really be the one doing that.
	// Actually, let's keep Resolved count strictly to what we explicitly set.
	// Wait, it's better if Resolved count is accurate.
	// I'll do a quick pass to count solved cells.
	count := uint8(0)

	for i := range 81 {
		// Use bits.OnesCount16 if I had it here, but I can just check
		// if (c & (c-1)) == 0 && c != 0
		c := b.Cells[i]
		if c != 0 && (c&(c-1)) == 0 {
			count++
		}
	}

	b.Resolved = count

	return b, nil
}
