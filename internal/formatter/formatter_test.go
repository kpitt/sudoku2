package formatter

import (
	"strings"
	"testing"

	"github.com/kpitt/sudoku2/internal/parser"
)

func TestFormatSS(t *testing.T) {
	input := "310004069000000200008005040000000005006000017807030000590700006600003050000100002"
	b, _ := parser.Parse(input)
	output := FormatSS(b)

	expected := `*-----------*
|31.|..4|.69|
|...|...|2..|
|..8|..5|.4.|
|---+---+---|
|...|...|..5|
|..6|...|.17|
|8.7|.3.|..4|
|---+---+---|
|59.|7..|..6|
|6..|..3|.5.|
|...|1..|..2|
*-----------*
`
	if strings.TrimSpace(output) != strings.TrimSpace(expected) {
		t.Errorf("FormatSS mismatch.\nGot:\n%s\nExpected:\n%s", output, expected)
	}
}
