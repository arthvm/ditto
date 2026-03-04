package git

import (
	"context"
	"strings"
)

var (
	Staged DiffOption = "--staged"
	Stats  DiffOption = "--stat"
)

type DiffArg interface {
	String() string
	isDiffArg()
}

type DiffOption string

func (o DiffOption) String() string { return string(o) }
func (o DiffOption) isDiffArg()     {}

func Target(target string) DiffOption {
	return DiffOption(target)
}

func Diff(ctx context.Context, options ...DiffArg) (string, error) {
	var args []string

	for _, opt := range options {
		parts := strings.Fields(opt.String())
		args = append(args, parts...)
	}
	gitArgs := append([]string{"diff"}, args...)

	return run(ctx, gitArgs...)
}
