package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertCommand(t *testing.T) {
	output := new(bytes.Buffer)
	rootCmd.SetOut(output)

	t.Run("to pretty", func(t *testing.T) {
		rootCmd.SetArgs([]string{"convert", "003020600900305001001806400008102900700000008006708200002609500800203009005010300", "--format", "pretty"})
		err := rootCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, output.String(), "|")
	})

	t.Run("to raw", func(t *testing.T) {
		output.Reset()
		rootCmd.SetArgs([]string{"convert", "003020600900305001001806400008102900700000008006708200002609500800203009005010300", "--format", "raw"})
		err := rootCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, output.String(), "003020600900305001001806400008102900700000008006708200002609500800203009005010300")
	})
}
