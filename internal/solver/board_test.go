package solver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBoard(t *testing.T) {
	t.Run("valid 81-char string", func(t *testing.T) {
		input := "003020600900305001001806400008102900700000008006708200002609500800203009005010300"
		board, err := ParseBoard(input)
		assert.NoError(t, err)
		assert.Equal(t, 3, board[0][2])
		assert.Equal(t, 9, board[1][0])
	})

	t.Run("invalid length", func(t *testing.T) {
		input := "123"
		_, err := ParseBoard(input)
		assert.Error(t, err)
	})

	t.Run("invalid characters", func(t *testing.T) {
		input := "00302060090030500100180640000810290070000000800670820000260950080020300900501030X"
		_, err := ParseBoard(input)
		assert.Error(t, err)
	})
}

func TestBoard_IsValid(t *testing.T) {
	t.Run("valid board", func(t *testing.T) {
		board := Board{
			{5, 3, 0, 0, 7, 0, 0, 0, 0},
			{6, 0, 0, 1, 9, 5, 0, 0, 0},
			{0, 9, 8, 0, 0, 0, 0, 6, 0},
			{8, 0, 0, 0, 6, 0, 0, 0, 3},
			{4, 0, 0, 8, 0, 3, 0, 0, 1},
			{7, 0, 0, 0, 2, 0, 0, 0, 6},
			{0, 6, 0, 0, 0, 0, 2, 8, 0},
			{0, 0, 0, 4, 1, 9, 0, 0, 5},
			{0, 0, 0, 0, 8, 0, 0, 7, 9},
		}
		assert.True(t, board.IsValid())
	})

	t.Run("duplicate in row", func(t *testing.T) {
		board := Board{
			{5, 5, 0, 0, 7, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			// ... rest zero
		}
		assert.False(t, board.IsValid())
	})

	t.Run("duplicate in column", func(t *testing.T) {
		board := Board{
			{5, 0, 0, 0, 0, 0, 0, 0, 0},
			{5, 0, 0, 0, 0, 0, 0, 0, 0},
			// ... rest zero
		}
		assert.False(t, board.IsValid())
	})

	t.Run("duplicate in box", func(t *testing.T) {
		board := Board{
			{5, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 5, 0, 0, 0, 0, 0, 0, 0},
			// ... rest zero
		}
		assert.False(t, board.IsValid())
	})
}
