package solver

// ActionType defines what kind of action a step represents.
type ActionType int

const (
	// ActionPlacement indicates a cell is solved.
	ActionPlacement ActionType = iota
	// ActionEliminate indicates candidates are removed from a cell.
	ActionEliminate
)

// Step represents a single deductive step in the solving process.
// It uses fixed-size arrays to avoid heap allocations.
type Step struct {
	Technique string
	Target    [81]int // Enough for any technique (e.g. eliminating from all other cells)
	TargetLen int
	Values    [9]int
	ValuesLen int
	Action    ActionType
}
