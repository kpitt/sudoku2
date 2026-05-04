package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckCommand(t *testing.T) {
	output := new(bytes.Buffer)
	rootCmd.SetOut(output)

	t.Run("no args", func(t *testing.T) {
		rootCmd.SetArgs([]string{"check"})
		err := rootCmd.Execute()
		assert.Error(t, err) // Should error due to ExactArgs(1)
	})

	t.Run("with puzzle", func(t *testing.T) {
		rootCmd.SetArgs([]string{"check", "003020600900305001001806400008102900700000008006708200002609500800203009005010300"})
		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, output.String(), "The puzzle is valid and solvable.")
	})
}
