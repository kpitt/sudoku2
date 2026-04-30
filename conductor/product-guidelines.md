# Product Guidelines

## Tone and Voice
- **Professional & Technical:** The application should communicate in a straightforward, concise manner. Avoid fluff or overly conversational language. Focus on delivering accurate data and clear instructions.

## UX Principles (CLI)
- **Human-Readable Output:** Prioritize formatting that is easy for humans to scan and understand within a terminal.
- **Rich Formatting:** Utilize ANSI colors and text styles to differentiate between puzzle data, hints, and application status. For example, use different colors for solved cells versus original cells.
- **Consistency:** Ensure that all CLI commands and flags follow a predictable pattern.

## Error Handling
- **Concise & Coded:** Errors should be reported with a brief, clear message and a unique error code. This allows for quick identification of issues while keeping the output clean.

## Documentation Style
- **Reference-Focused with Examples:** Documentation should provide a complete reference of all available flags and parameters.
- **Input/Output Clarity:** Crucially, provide clear examples of all supported Sudoku input and output formats (e.g., 81-character strings, multi-line grids, etc.) to ensure users can easily integrate the tool.