package cmd

import (
	"github.com/spf13/cobra"
)

var solveCmd = &cobra.Command{
	Use:   "solve [puzzle]",
	Short: "Solve a Sudoku puzzle",
	Long:  `Solve a Sudoku puzzle provided as an 81-character string.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Println("Please provide a Sudoku puzzle string.")
			return
		}
		puzzle := args[0]
		cmd.Printf("Attempting to solve puzzle: %s\n", puzzle)
		// Logic will be added in later tasks
	},
}

func init() {
	rootCmd.AddCommand(solveCmd)
}
