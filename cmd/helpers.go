package cmd

import (
	"time"

	"github.com/briandowns/spinner"
)

func newSpinner(suffix string) *spinner.Spinner {
	s := spinner.New(
		spinner.CharSets[14],
		time.Millisecond*100,
		spinner.WithColor("yellow"),
	)
	s.Suffix = suffix
	return s
}
