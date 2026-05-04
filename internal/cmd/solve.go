package cmd

import (
	"github.com/kpitt/sudoku2/internal/solver"
	"github.com/spf13/cobra"
)

var solveCmd = &cobra.Command{
	Use:   "solve [puzzle]",
	Short: "Solve a Sudoku puzzle",
	Long:  `Solve a Sudoku puzzle provided as an 81-character string using deductive human-like techniques.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		puzzle := args[0]

		board, err := solver.ParseBoard(puzzle)
		if err != nil {
			cmd.Printf("Error: %v\n", err)
			return
		}

		cmd.Println("Initial Puzzle:")
		cmd.Println(board.String())

		if board.SolveDeductive() {
			cmd.Println("Solved Puzzle (using deductive logic):")
			cmd.Println(board.String())
		} else {
			cmd.Println("Could not solve the puzzle using purely deductive techniques.")
			cmd.Println("Partial Solution:")
			cmd.Println(board.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(solveCmd)
}
