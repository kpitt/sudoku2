package cmd

import (
	"github.com/kpitt/sudoku2/internal/io"
	"github.com/kpitt/sudoku2/internal/solver"
	"github.com/spf13/cobra"
)

var hintCmd = &cobra.Command{
	Use:   "hint [puzzle]",
	Short: "Provide a hint for a Sudoku puzzle",
	Long:  `Provide a single deductive hint for a Sudoku puzzle provided as an 81-character string.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		puzzle := args[0]

		board, err := io.ParseBoard(puzzle)
		if err != nil {
			cmd.Printf("Error: %v\n", err)
			return
		}

		hint, err := solver.GetHint(&board)
		if err != nil {
			cmd.Printf("Hint Error: %v\n", err)
			return
		}

		cmd.Printf("Hint: %s\n", hint.Message)
	},
}

func init() {
	rootCmd.AddCommand(hintCmd)
}
