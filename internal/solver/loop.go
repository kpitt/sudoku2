package solver

// Solve attempts to solve the Sudoku puzzle using the implemented techniques.
// It returns a slice of steps taken.
func Solve(b *Board) []Step {
	var steps []Step

	for {
		// Priority 1: Naked Singles
		if s, ok := FindNakedSingles(b); ok {
			b.ApplyStep(s)
			steps = append(steps, s)
			continue
		}

		// Priority 2: Hidden Singles
		if s, ok := FindHiddenSingles(b); ok {
			b.ApplyStep(s)
			steps = append(steps, s)
			continue
		}

		// Priority 3: Naked Pairs
		if s, ok := FindNakedPairs(b); ok {
			b.ApplyStep(s)
			steps = append(steps, s)
			continue
		}

		// Priority 4: Hidden Pairs
		if s, ok := FindHiddenPairs(b); ok {
			b.ApplyStep(s)
			steps = append(steps, s)
			continue
		}

		// No more steps found
		break
	}

	return steps
}
