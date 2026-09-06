package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var solveCmd = &cobra.Command{
	Use:   "solve [puzzle]",
	Short: "Solve a sudoku puzzle",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Solving:", args[0])
	},
}

func init() {
	rootCmd.AddCommand(solveCmd)
}
