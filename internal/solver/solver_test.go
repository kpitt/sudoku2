package solver

import (
	"testing"
)

func TestFindNakedSingles(t *testing.T) {
	b := &Board{}

	for i := range 81 {
		b.Cells[i] = 0x01FF
	}

	for i := range 9 {
		b.RowState[i] = 0x01FF
		b.ColState[i] = 0x01FF
		b.BoxState[i] = 0x01FF
	}

	// Set a naked single at cell 0
	b.Cells[0] = 0x0001 // Candidate 1

	s, ok := FindNakedSingles(b)
	if !ok {
		t.Fatal("Expected to find a naked single")
	}

	if s.Target[0] != 0 || s.Values[0] != 1 {
		t.Errorf("Unexpected naked single step: %+v", s)
	}
}

func TestFindHiddenSingles(t *testing.T) {
	b := &Board{}

	for i := range 81 {
		b.Cells[i] = 0x01FF
	}

	for i := range 9 {
		b.RowState[i] = 0x01FF
		b.ColState[i] = 0x01FF
		b.BoxState[i] = 0x01FF
	}

	// Hidden single: Candidate 1 only possible in cell 0 within Row 0
	for i := 1; i < 9; i++ {
		b.Cells[i] &^= 0x0001
	}

	s, ok := FindHiddenSingles(b)
	if !ok {
		t.Fatal("Expected to find a hidden single")
	}

	if s.Target[0] != 0 || s.Values[0] != 1 {
		t.Errorf("Unexpected hidden single step: %+v", s)
	}
}

func BenchmarkSolve(b *testing.B) {
	input := "003020600900305001001806400008102900700000008006708200002609500800203009005010300"

	b.ResetTimer()

	for b.Loop() {
		board := Board{}

		for i := range 81 {
			board.Cells[i] = 0x01FF
		}

		for i := range 9 {
			board.RowState[i] = 0x01FF
			board.ColState[i] = 0x01FF
			board.BoxState[i] = 0x01FF
		}

		for i, r := range input {
			if r >= '1' && r <= '9' {
				digit := int(r - '0')
				mask := uint16(1 << (digit - 1))
				board.Cells[i] = mask
			}
		}

		Solve(&board)
	}
}
