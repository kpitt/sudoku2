package solver

import "math/bits"

// Board represents the state of a Sudoku grid.
// It uses bitmasks (uint16) to represent candidates for each cell.
// Bit n (0-8) represents candidate n+1.
type Board struct {
	Cells    [81]uint16 // Cell candidate masks
	RowState [9]uint16  // Candidates remaining in each row
	ColState [9]uint16  // Candidates remaining in each col
	BoxState [9]uint16  // Candidates remaining in each box
	Resolved uint8      // Number of solved cells
}

// ApplyStep applies a Step to the Board and handles propagation.
func (b *Board) ApplyStep(s Step) {
	switch s.Action {
	case ActionPlacement:
		cellIdx := s.Target[0]
		digit := s.Values[0]
		mask := uint16(1 << (digit - 1))

		// Set the cell
		b.Cells[cellIdx] = mask

		// Update states
		r := RowLUT[cellIdx]
		c := ColLUT[cellIdx]
		box := BoxLUT[cellIdx]
		b.RowState[r] &^= mask
		b.ColState[c] &^= mask
		b.BoxState[box] &^= mask

		// Propagate to peers
		for _, peerIdx := range PeersLUT[cellIdx] {
			if (b.Cells[peerIdx] & mask) != 0 {
				b.Cells[peerIdx] &^= mask
			}
		}

		// Recalculate Resolved
		b.SyncResolved()

	case ActionEliminate:
		// Used for Pairs/Subsets
		digitMask := uint16(0)
		for i := range s.ValuesLen {
			val := s.Values[i]
			digitMask |= (1 << (val - 1))
		}

		for i := range s.TargetLen {
			cellIdx := s.Target[i]
			b.Cells[cellIdx] &^= digitMask
		}

		// Recalculate Resolved
		b.SyncResolved()
	}
}

// SyncResolved recalculates the Resolved count by scanning all cells.
func (b *Board) SyncResolved() {
	count := uint8(0)
	for i := range 81 {
		if bits.OnesCount16(b.Cells[i]) == 1 {
			count++
		}
	}

	b.Resolved = count
}
