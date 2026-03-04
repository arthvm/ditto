package ui

import (
	"io"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer

	spinner *spinner.Spinner
}

func Default() *IOStreams {
	return &IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

func (s *IOStreams) StartSpinner(label string) {
	sp := spinner.New(
		spinner.CharSets[14],
		time.Millisecond*100,
		spinner.WithColor("yellow"),
		spinner.WithWriter(s.ErrOut),
	)
	sp.Suffix = label
	sp.Start()

	s.spinner = sp
}

func (s *IOStreams) StopSpinner() {
	if s.spinner != nil {
		s.spinner.Stop()
		s.spinner = nil
	}
}
