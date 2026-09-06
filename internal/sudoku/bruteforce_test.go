package sudoku

import "testing"

func isValidSolution(g *Grid) bool {
	if g.Clues != 81 {
		return false
	}
	for h := range 27 {
		var seen uint16
		for c := range 9 {
			cellIdx := Houses[h][c]
			val := g.Values[cellIdx]
			if val < 1 || val > 9 {
				return false
			}
			mask := uint16(1 << val)
			if seen&mask != 0 {
				return false
			}
			seen |= mask
		}
	}
	return true
}

func TestBruteForce(t *testing.T) {
	input := ".....6....59.....82....8....45........3........6..3.54...325..6.................."
	g, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	g.InitializeCandidates()
	solved := g.SolveBruteForce()
	if !solved || g.Clues != 81 {
		t.Errorf("Expected board to be solved, got solved=%v clues=%d", solved, g.Clues)
	}
	if !isValidSolution(g) {
		t.Errorf("Solved board is not a valid Sudoku solution")
	}
}

func TestBruteForce_AlreadySolved(t *testing.T) {
	input := ".....6....59.....82....8....45........3........6..3.54...325..6.................."
	g, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	g.InitializeCandidates()
	if !g.SolveBruteForce() {
		t.Fatalf("failed initial solve")
	}

	// Calling SolveBruteForce again on already-solved board should return true immediately
	if !g.SolveBruteForce() {
		t.Errorf("expected SolveBruteForce on solved board to return true")
	}
}

func TestBruteForce_Unsolvable(t *testing.T) {
	// Cell (2, 2) has 1-8 in its box and 9 in its row -> 0 candidates
	input := "123......456......78.9............................................................"
	g, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	g.InitializeCandidates()
	if g.SolveBruteForce() {
		t.Errorf("expected unsolvable board to return false")
	}
}

func TestBruteForce_ZeroAllocations(t *testing.T) {
	input := "53..7....6..195....98....6.8...6...34..8.3..17...2...6.6....28....419..5....8..79"
	base, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	base.InitializeCandidates()

	var g Grid
	allocs := testing.AllocsPerRun(10, func() {
		g = *base
		if !g.SolveBruteForce() {
			panic("expected solve to succeed")
		}
	})

	if allocs > 0 {
		t.Errorf("expected 0 heap allocations, got %f", allocs)
	}
}

func BenchmarkSolveBruteForce_Typical(b *testing.B) {
	input := "53..7....6..195....98....6.8...6...34..8.3..17...2...6.6....28....419..5....8..79"
	base, err := Parse(input)
	if err != nil {
		b.Fatalf("unexpected parse error: %v", err)
	}
	base.InitializeCandidates()

	var g Grid
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g = *base
		g.SolveBruteForce()
	}
}

func BenchmarkSolveBruteForce_Hard17(b *testing.B) {
	input := ".....6....59.....82....8....45........3........6..3.54...325..6.................."
	base, err := Parse(input)
	if err != nil {
		b.Fatalf("unexpected parse error: %v", err)
	}
	base.InitializeCandidates()

	var g Grid
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g = *base
		g.SolveBruteForce()
	}
}
