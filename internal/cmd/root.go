package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sudoku",
	Short: "Sudoku2 is a high-performance Sudoku solver and educational tool",
	Long: `A fast and reliable Sudoku solver built with Go.
It supports solving puzzles, validating uniqueness, and providing educational hints.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Sudoku2! Use 'sudoku --help' for more information.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
