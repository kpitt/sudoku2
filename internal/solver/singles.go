package solver

import "math/bits"

// FindNakedSingles looks for cells that have only one candidate remaining.
func FindNakedSingles(b *Board) (Step, bool) {
	for i := range 81 {
		mask := b.Cells[i]

		if bits.OnesCount16(mask) == 1 {
			r := RowLUT[i]
			c := ColLUT[i]
			box := BoxLUT[i]

			// Only return a step if the cell hasn't been "officially" resolved in its houses yet
			if (b.RowState[r]&mask) != 0 || (b.ColState[c]&mask) != 0 || (b.BoxState[box]&mask) != 0 {
				digit := bits.TrailingZeros16(mask) + 1
				s := Step{
					Technique: "Naked Single",
					TargetLen: 1,
					ValuesLen: 1,
					Action:    ActionPlacement,
				}
				s.Target[0] = i
				s.Values[0] = digit

				return s, true
			}
		}
	}

	return Step{}, false
}

// FindHiddenSingles looks for candidates that can only be placed in one cell within a house.
func FindHiddenSingles(b *Board) (Step, bool) {
	for hIdx := range 27 {
		house := HouseLUT[hIdx]

		var seenOnce, seenMultiple uint16
		for _, cellIdx := range house {
			mask := b.Cells[cellIdx]
			seenMultiple |= (seenOnce & mask)
			seenOnce |= mask
		}

		hiddenSingles := seenOnce &^ seenMultiple

		if hiddenSingles != 0 {
			digitMask := uint16(1 << bits.TrailingZeros16(hiddenSingles))

			for _, cellIdx := range house {
				if (b.Cells[cellIdx] & digitMask) != 0 {
					// Check if it's actually hidden (i.e. cell has multiple candidates)
					if bits.OnesCount16(b.Cells[cellIdx]) > 1 {
						digit := bits.TrailingZeros16(digitMask) + 1
						s := Step{
							Technique: "Hidden Single",
							TargetLen: 1,
							ValuesLen: 1,
							Action:    ActionPlacement,
						}
						s.Target[0] = cellIdx
						s.Values[0] = digit

						return s, true
					}
				}
			}
		}
	}

	return Step{}, false
}
