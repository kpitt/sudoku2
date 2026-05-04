package solver

import (
	"testing"

	"github.com/kpitt/sudoku2/internal/board"
	"github.com/kpitt/sudoku2/internal/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSolveBacktracking(t *testing.T) {
	t.Run("easy puzzle", func(t *testing.T) {
		input := "530070000600195000098000060800060003400803001700020006060000280000419005000080079"
		b, err := io.ParseBoard(input)
		require.NoError(t, err)

		solved := SolveBacktracking(&b)
		assert.True(t, solved)
		assert.True(t, b.IsSolved())

		// Verify first row of solution
		expectedRow := [9]int{5, 3, 4, 6, 7, 8, 9, 1, 2}
		assert.Equal(t, expectedRow, b[0])
	})

	t.Run("unsolvable puzzle", func(t *testing.T) {
		// Puzzle with duplicate in same row (invalid start)
		b := board.Board{
			{5, 5, 0, 0, 0, 0, 0, 0, 0},
		}
		solved := SolveBacktracking(&b)
		assert.False(t, solved)
	})
}

func TestSolveDeductive(t *testing.T) {
	t.Run("easy puzzle - naked singles", func(t *testing.T) {
		// A puzzle that can be solved by naked singles alone (or mostly)
		input := "003020600900305001001806400008102900700000008006708200002609500800203009005010300"
		b, err := io.ParseBoard(input)
		require.NoError(t, err)

		solved := SolveDeductive(&b)
		assert.True(t, solved)
		assert.True(t, b.IsSolved())
	})
}

func TestGetHint(t *testing.T) {
	t.Run("naked single hint", func(t *testing.T) {
		input := "003020600900305001001806400008102900700000008006708200002609500800203009005010300"
		b, _ := io.ParseBoard(input)
		hint, err := GetHint(&b)
		require.NoError(t, err)
		assert.Contains(t, hint.Message, "Naked Single")
	})

	t.Run("invalid board hint", func(t *testing.T) {
		b := board.Board{{5, 5, 0}}
		_, err := GetHint(&b)
		assert.Error(t, err)
	})
}
