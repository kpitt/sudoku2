package solver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoard_SolveBacktracking(t *testing.T) {
	t.Run("easy puzzle", func(t *testing.T) {
		input := "530070000600195000098000060800060003400803001700020006060000280000419005000080079"
		board, err := ParseBoard(input)
		require.NoError(t, err)

		solved := board.SolveBacktracking()
		assert.True(t, solved)
		assert.True(t, board.IsSolved())

		// Verify first row of solution
		expectedRow := [9]int{5, 3, 4, 6, 7, 8, 9, 1, 2}
		assert.Equal(t, expectedRow, board[0])
	})

	t.Run("unsolvable puzzle", func(t *testing.T) {
		// Puzzle with duplicate in same row (invalid start)
		board := Board{
			{5, 5, 0, 0, 0, 0, 0, 0, 0},
		}
		solved := board.SolveBacktracking()
		assert.False(t, solved)
	})
}

func TestBoard_SolveDeductive(t *testing.T) {
	t.Run("easy puzzle - naked singles", func(t *testing.T) {
		// A puzzle that can be solved by naked singles alone (or mostly)
		input := "003020600900305001001806400008102900700000008006708200002609500800203009005010300"
		board, err := ParseBoard(input)
		require.NoError(t, err)

		solved := board.SolveDeductive()
		assert.True(t, solved)
		assert.True(t, board.IsSolved())
	})
}
