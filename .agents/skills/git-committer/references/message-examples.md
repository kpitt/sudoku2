# Git Commit Message Examples

## Good Examples

### Example 1: New Feature
Implement Hidden Singles detection logic

Adds the ability to identify numbers that can only be placed in a
single cell within a row, column, or block.

Assisted-by: Gemini CLI Agent <gemini-cli-agent@google.com>

### Example 2: Bug Fix
Fix slice bounds error in board parser

Resolves an off-by-one error when processing the final row of the
input string, which causes a panic on incomplete boards.

Assisted-by: Gemini CLI Agent <gemini-cli-agent@google.com>

### Example 3: Minimal Change
Update documentation for solver constraints

Assisted-by: Gemini CLI Agent <gemini-cli-agent@google.com>

## Bad Examples

### Example 1: Subject too long
Implement Hidden Singles detection logic for rows, columns and blocks

**Reason:** Subject line is 67 characters (exceeds 50 character limit).

### Example 2: Semantic prefix used
feat: add hidden singles logic

**Reason:** Uses "feat:" prefix, which is forbidden by project rules.

### Example 3: Past tense and filler phrasing
Fix slice bounds error in board parser

This change resolved an off-by-one error when processing the row.

**Reason:** Uses past tense ("resolved") and filler phrasing ("This change").

### Example 4: Missing footer
Implement board validation

Check that the initial board state follows Sudoku rules.

**Reason:** Missing mandatory "Assisted-by" footer.
