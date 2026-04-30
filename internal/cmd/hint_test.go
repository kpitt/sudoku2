package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHintCommand(t *testing.T) {
	output := new(bytes.Buffer)
	rootCmd.SetOut(output)

	t.Run("with puzzle needing hint", func(t *testing.T) {
		rootCmd.SetArgs([]string{"hint", "003020600900305001001806400008102900700000008006708200002609500800203009005010300"})
		err := rootCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, output.String(), "Hint:")
	})
}
