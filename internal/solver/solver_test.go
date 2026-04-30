package solver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoard_Solve(t *testing.T) {
	t.Run("easy puzzle", func(t *testing.T) {
		input := "530070000600195000098000060800060003400803001700020006060000280000419005000080079"
		board, err := ParseBoard(input)
		assert.NoError(t, err)

		solved := board.Solve()
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
		solved := board.Solve()
		assert.False(t, solved)
	})
}
