package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sudoku",
	Short: "A high-performance Sudoku solver CLI",
	Long:  `A zero-allocation, high-performance Sudoku solver CLI using deductive techniques.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Root flags can be defined here
}
