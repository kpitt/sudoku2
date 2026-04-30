package cmd

import (
	"github.com/kpitt/sudoku2/internal/solver"
	"github.com/spf13/cobra"
)

var format string

var convertCmd = &cobra.Command{
	Use:   "convert [puzzle]",
	Short: "Convert a Sudoku puzzle to a different format",
	Long:  `Convert an 81-character Sudoku string to other text representations like pretty-print or raw strings.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		puzzle := args[0]
		board, err := solver.ParseBoard(puzzle)
		if err != nil {
			cmd.Printf("Error: %v\n", err)
			return
		}

		switch format {
		case "pretty":
			cmd.Println(board.String())
		case "raw":
			cmd.Println(board.RawString())
		default:
			cmd.Printf("Unknown format: %s. Use 'pretty' or 'raw'.\n", format)
		}
	},
}

func init() {
	convertCmd.Flags().StringVarP(&format, "format", "f", "pretty", "Target format (pretty, raw)")
	rootCmd.AddCommand(convertCmd)
}
