package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand(t *testing.T) {
	t.Run("help flag", func(t *testing.T) {
		output := new(bytes.Buffer)
		rootCmd.SetOut(output)
		rootCmd.SetArgs([]string{"--help"})

		err := rootCmd.Execute()

		require.NoError(t, err)
		assert.Contains(t, output.String(), "A fast and reliable Sudoku solver built with Go")
	})

	t.Run("welcome message", func(t *testing.T) {
		output := new(bytes.Buffer)
		rootCmd.SetOut(output)
		rootCmd.SetArgs([]string{})

		err := rootCmd.Execute()

		require.NoError(t, err)

		outStr := output.String()
		assert.NotEmpty(t, outStr)
	})

	t.Run("Execute function", func(t *testing.T) {
		output := new(bytes.Buffer)
		rootCmd.SetOut(output)
		rootCmd.SetArgs([]string{"--help"})
		Execute()
	})
}
