package io

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatBoard(t *testing.T) {
	input := "530070000600195000098000060800060003400803001700020006060000280000419005000080079"
	b, _ := ParseBoard(input)

	t.Run("FormatRaw", func(t *testing.T) {
		assert.Equal(t, input, FormatRaw(&b))
	})

	t.Run("FormatPretty", func(t *testing.T) {
		out := FormatPretty(&b)
		// Check for an exact horizontal separator line
		assert.Contains(t, out, "------+-------+------\n")
		// Check for an exact data row line (Row 0 of the provided puzzle)
		assert.Contains(t, out, "5 3 . | . 7 . | . . . \n")
	})
}
