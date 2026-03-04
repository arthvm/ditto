package workflow

// Progress reports long-running operation status to the user.
// Implementations control how progress is displayed: a CLI spinner,
// a TUI progress bar, or a no-op for non-interactive use.
type Progress interface {
	StartSpinner(label string)
	StopSpinner()
}
