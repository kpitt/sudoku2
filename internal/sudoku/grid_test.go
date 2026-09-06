package sudoku

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestGrid_Initialization(t *testing.T) {
	g := &Grid{}
	if len(g.Values) != 81 || len(g.Candidates) != 81 {
		t.Errorf("Grid must have 81 cells")
	}

	g.Candidates[0] = 0b1111111110
	if g.CandidatesCount(0) != 9 {
		t.Errorf("Expected 9 candidates")
	}
}

func TestGrid_CandidatesCount_Boundaries(t *testing.T) {
	tests := []struct {
		name     string
		mask     uint16
		expected int
	}{
		{name: "0 candidates (mask 0)", mask: 0, expected: 0},
		{name: "1 candidate (digit 1)", mask: 1 << 1, expected: 1},
		{name: "1 candidate (digit 9)", mask: 1 << 9, expected: 1},
		{name: "all 9 candidates", mask: 0b1111111110, expected: 9},
		{name: "arbitrary 3 candidates (digits 2, 5, 8)", mask: (1 << 2) | (1 << 5) | (1 << 8), expected: 3},
		{name: "even candidates (digits 2, 4, 6, 8)", mask: (1 << 2) | (1 << 4) | (1 << 6) | (1 << 8), expected: 4},
		{name: "odd candidates (digits 1, 3, 5, 7, 9)", mask: (1 << 1) | (1 << 3) | (1 << 5) | (1 << 7) | (1 << 9), expected: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Grid{}
			g.Candidates[0] = tt.mask
			if got := g.CandidatesCount(0); got != tt.expected {
				t.Errorf("CandidatesCount(0) with mask %#x = %d, expected %d", tt.mask, got, tt.expected)
			}
		})
	}
}

// captureStdout captures os.Stdout during the execution of f.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	oldStdout := os.Stdout
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	f()

	_ = w.Close()
	os.Stdout = oldStdout
	return <-outC
}

func TestGrid_Display_Raw(t *testing.T) {
	// Empty grid
	emptyGrid := &Grid{}
	for _, style := range []string{"raw", "string", "line", "RAW", "Line"} {
		output := captureStdout(t, func() {
			emptyGrid.Display(style)
		})
		expected := strings.Repeat(".", 81) + "\n"
		if output != expected {
			t.Errorf("Display(%q) on empty grid = %q, expected %q", style, output, expected)
		}
	}

	// Partially filled grid
	puzzle := "1" + strings.Repeat(".", 80)
	g, err := Parse(puzzle)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	output := captureStdout(t, func() {
		g.Display("line")
	})
	expected := puzzle + "\n"
	if output != expected {
		t.Errorf("Display(\"line\") = %q, expected %q", output, expected)
	}
}

func TestGrid_Display_Grid(t *testing.T) {
	emptyGrid := &Grid{}
	expectedEmptyGrid := "+-------+-------+-------+\n" +
		"| . . . | . . . | . . . |\n" +
		"| . . . | . . . | . . . |\n" +
		"| . . . | . . . | . . . |\n" +
		"+-------+-------+-------+\n" +
		"| . . . | . . . | . . . |\n" +
		"| . . . | . . . | . . . |\n" +
		"| . . . | . . . | . . . |\n" +
		"+-------+-------+-------+\n" +
		"| . . . | . . . | . . . |\n" +
		"| . . . | . . . | . . . |\n" +
		"| . . . | . . . | . . . |\n" +
		"+-------+-------+-------+\n"

	for _, style := range []string{"grid", "GRID", "", "unknown"} {
		output := captureStdout(t, func() {
			emptyGrid.Display(style)
		})
		if output != expectedEmptyGrid {
			t.Errorf("Display(%q) = %q, expected %q", style, output, expectedEmptyGrid)
		}
	}
}

func TestGrid_Display_RoundTripWithParser(t *testing.T) {
	input := ".....6....59.....82....8....45........3........6..3.54...325..6.................."
	g, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Round-trip with "raw"
	rawOutput := captureStdout(t, func() {
		g.Display("raw")
	})
	gParsedRaw, err := Parse(rawOutput)
	if err != nil {
		t.Fatalf("Parse error on raw output: %v", err)
	}
	if gParsedRaw.Values != g.Values || gParsedRaw.Clues != g.Clues {
		t.Errorf("Round-trip failed for raw output")
	}

	// Round-trip with "grid"
	gridOutput := captureStdout(t, func() {
		g.Display("grid")
	})
	gParsedGrid, err := Parse(gridOutput)
	if err != nil {
		t.Fatalf("Parse error on grid output: %v", err)
	}
	if gParsedGrid.Values != g.Values || gParsedGrid.Clues != g.Clues {
		t.Errorf("Round-trip failed for grid output")
	}
}
