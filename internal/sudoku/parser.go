package sudoku

import (
	"fmt"
	"strings"
)

// Parse normalizes an input string, extracting the first 81 cell characters.
// '.' and '0' represent empty cells; '1'-'9' represent given clues.
// Non-cell characters (whitespace, borders, comments) are stripped.
// Returns an error if fewer than 81 valid cell characters are present.
func Parse(input string) (*Grid, error) {
	clean := strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' {
			return r
		}
		return -1
	}, input)

	if len(clean) < 81 {
		return nil, fmt.Errorf("invalid input: expected at least 81 cell characters, found %d", len(clean))
	}

	g := &Grid{}
	for idx := 0; idx < 81; idx++ {
		ch := clean[idx]
		if ch >= '1' && ch <= '9' {
			g.Values[idx] = byte(ch - '0')
			g.Clues++
		}
	}
	return g, nil
}

// InitializeCandidates populates candidates for empty cells and eliminates
// candidates that conflict with given clues in peer cells.
func (g *Grid) InitializeCandidates() {
	// Start with all candidates (bits 1-9) for empty cells
	for i := range 81 {
		if g.Values[i] == 0 {
			g.Candidates[i] = 0b1111111110 // Bits 1-9 set
		} else {
			g.Candidates[i] = 0
		}
	}

	// Eliminate based on givens
	for i := range 81 {
		val := g.Values[i]
		if val > 0 {
			mask := uint16(1 << val)
			for _, p := range Peers[i] {
				g.Candidates[p] &^= mask
			}
		}
	}
}
