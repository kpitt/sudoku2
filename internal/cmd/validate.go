package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [puzzle]",
	Short: "Validate a Sudoku puzzle",
	Long:  `Check if a Sudoku puzzle is valid and has a unique solution.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a Sudoku puzzle string.")
			return
		}
		puzzle := args[0]
		fmt.Printf("Attempting to validate puzzle: %s\n", puzzle)
		// Logic will be added in later tasks
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
