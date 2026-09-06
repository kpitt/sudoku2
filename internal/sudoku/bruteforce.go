package sudoku

// SolveBruteForce solves the Sudoku grid using recursive backtracking with
// the Minimum Remaining Values (MRV) heuristic. It operates with zero heap
// allocations during recursion by backing up candidate bitmasks on the stack.
// Returns true if a valid solution was found, false otherwise.
func (g *Grid) SolveBruteForce() bool {
	if g.Clues == 81 {
		return true
	}

	// Find cell with fewest candidates (MRV heuristic)
	minCand := 10
	bestCell := -1
	for i := range 81 {
		if g.Values[i] == 0 {
			c := g.CandidatesCount(i)
			if c < minCand {
				minCand = c
				bestCell = i
				if minCand == 0 {
					break
				}
			}
		}
	}

	if bestCell == -1 || minCand == 0 {
		return false // Unsolvable or dead end
	}

	// Backup state (stored on stack, zero heap alloc)
	var backup [81]uint16
	copy(backup[:], g.Candidates[:])

	cands := g.Candidates[bestCell]
	for d := 1; d <= 9; d++ {
		if cands&(1<<d) != 0 {
			g.Values[bestCell] = byte(d)
			g.Clues++

			// Eliminate from peers
			mask := uint16(1 << d)
			for _, p := range Peers[bestCell] {
				g.Candidates[p] &^= mask
			}

			if g.SolveBruteForce() {
				return true
			}

			// Backtrack
			g.Values[bestCell] = 0
			g.Clues--
			copy(g.Candidates[:], backup[:])
		}
	}
	return false
}
