package io

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBoard(t *testing.T) {
	t.Run("valid 81-char string", func(t *testing.T) {
		input := "003020600900305001001806400008102900700000008006708200002609500800203009005010300"
		board, err := ParseBoard(input)
		require.NoError(t, err)
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
