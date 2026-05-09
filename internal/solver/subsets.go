package solver

import "math/bits"

// FindNakedPairs looks for two cells in a house that contain the same two candidates.
func FindNakedPairs(b *Board) (Step, bool) {
	for hIdx := range 27 {
		house := HouseLUT[hIdx]

		for i := range 9 {
			c1 := house[i]
			mask1 := b.Cells[c1]

			if bits.OnesCount16(mask1) != 2 {
				continue
			}

			for j := i + 1; j < 9; j++ {
				c2 := house[j]
				mask2 := b.Cells[c2]

				if mask1 == mask2 {
					s := Step{
						Technique: "Naked Pair",
						Action:    ActionEliminate,
					}

					for k := range 9 {
						c3 := house[k]
						if c3 == c1 || c3 == c2 {
							continue
						}

						if (b.Cells[c3] & mask1) != 0 {
							s.Target[s.TargetLen] = c3
							s.TargetLen++
						}
					}

					if s.TargetLen > 0 {
						tempMask := mask1
						for tempMask != 0 {
							digit := bits.TrailingZeros16(tempMask) + 1
							s.Values[s.ValuesLen] = digit
							s.ValuesLen++
							tempMask &^= (1 << (digit - 1))
						}

						return s, true
					}
				}
			}
		}
	}

	return Step{}, false
}

// FindHiddenPairs looks for two candidates that appear only in the same two cells in a house.
func FindHiddenPairs(b *Board) (Step, bool) {
	for hIdx := range 27 {
		house := HouseLUT[hIdx]

		var candPos [9]uint16
		for digit := range 9 {
			mask := uint16(1 << digit)
			for pos, cellIdx := range house {
				if (b.Cells[cellIdx] & mask) != 0 {
					candPos[digit] |= (1 << pos)
				}
			}
		}

		for d1 := range 9 {
			p1 := candPos[d1]
			if bits.OnesCount16(p1) != 2 {
				continue
			}

			for d2 := d1 + 1; d2 < 9; d2++ {
				p2 := candPos[d2]
				if p1 == p2 {
					pairMask := uint16((1 << d1) | (1 << d2))
					s := Step{
						Technique: "Hidden Pair",
						Action:    ActionEliminate,
					}

					tempPos := p1
					for tempPos != 0 {
						pos := bits.TrailingZeros16(tempPos)
						cellIdx := house[pos]

						if (b.Cells[cellIdx] &^ pairMask) != 0 {
							s.Target[s.TargetLen] = cellIdx
							s.TargetLen++
						}

						tempPos &^= (1 << pos)
					}

					if s.TargetLen > 0 {
						removedMask := uint16(0)
						for i := range s.TargetLen {
							t := s.Target[i]
							removedMask |= (b.Cells[t] &^ pairMask)
						}

						tempRemove := removedMask
						for tempRemove != 0 {
							digit := bits.TrailingZeros16(tempRemove) + 1
							s.Values[s.ValuesLen] = digit
							s.ValuesLen++
							tempRemove &^= (1 << (digit - 1))
						}

						return s, true
					}
				}
			}
		}
	}

	return Step{}, false
}
