package cmd

import (
	"github.com/kpitt/sudoku2/internal/io"
	"github.com/kpitt/sudoku2/internal/solver"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [puzzle]",
	Short: "Check a Sudoku puzzle",
	Long:  `Check if a Sudoku puzzle is valid and has a unique solution using backtracking.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		puzzle := args[0]

		board, err := io.ParseBoard(puzzle)
		if err != nil {
			cmd.Printf("Error: %v\n", err)
			return
		}

		if !board.IsValid() {
			cmd.Println("The puzzle is invalid (violates Sudoku rules).")
			return
		}

		if solver.SolveBacktracking(&board) {
			cmd.Println("The puzzle is valid and solvable.")
			// In the future, we can add logic to check for multiple solutions.
		} else {
			cmd.Println("The puzzle is unsolvable.")
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
