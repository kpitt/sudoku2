package main

import (
	"fmt"

	"github.com/kpitt/sudoku2/internal/formatter"
	"github.com/kpitt/sudoku2/internal/parser"
	"github.com/kpitt/sudoku2/internal/solver"
	"github.com/spf13/cobra"
)

var solveCmd = &cobra.Command{
	Use:   "solve [puzzle]",
	Short: "Solve a Sudoku puzzle",
	Long:  `Solve an 81-character Sudoku puzzle string. Use '0' or '.' for empty cells.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		puzzle := args[0]

		b, err := parser.Parse(puzzle)
		if err != nil {
			fmt.Printf("Error parsing puzzle: %v\n", err)
			return
		}

		fmt.Println("Initial Board:")
		fmt.Print(formatter.FormatSS(b))

		steps := solver.Solve(b)

		fmt.Printf("\nApplied %d steps.\n", len(steps))

		if b.Resolved == 81 {
			fmt.Println("\nPuzzle SOLVED!")
		} else {
			fmt.Printf("\nPuzzle stalled at %d/81 cells.\n", b.Resolved)
		}

		fmt.Println("\nFinal Board:")
		fmt.Print(formatter.FormatSS(b))
	},
}

func init() {
	rootCmd.AddCommand(solveCmd)
}
