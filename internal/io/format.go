package io

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kpitt/sudoku2/internal/board"
)

// FormatRaw returns the board as a simple 81-character string.
func FormatRaw(b *board.Board) string {
	var out strings.Builder
	for r := range 9 {
		for c := range 9 {
			out.WriteString(strconv.Itoa(b[r][c]))
		}
	}

	return out.String()
}

// FormatPretty returns a human-readable representation of the board.
func FormatPretty(b *board.Board) string {
	var out strings.Builder
	for r := range 9 {
		if r > 0 && r%3 == 0 {
			out.WriteString("------+-------+------\n")
		}

		for c := range 9 {
			if c > 0 && c%3 == 0 {
				out.WriteString("| ")
			}

			if b[r][c] == 0 {
				out.WriteString(". ")
			} else {
				fmt.Fprintf(&out, "%d ", b[r][c])
			}
		}

		out.WriteString("\n")
	}

	return out.String()
}
