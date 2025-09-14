package git

import (
	"context"
	"os/exec"
)

var (
	Staged DiffOption = "--staged"
	Stats  DiffOption = "--stat"
)

type diffArg interface {
	String() string
	isDiffArg()
}

type DiffOption string

func (o DiffOption) String() string { return string(o) }
func (o DiffOption) isDiffArg()     {}

func Diff(ctx context.Context, options ...diffArg) (string, error) {
	args := make([]string, len(options))
	for i, opt := range options {
		args[i] = opt.String()
	}
	gitArgs := append([]string{"diff"}, args...)

	cmd := exec.CommandContext(ctx, "git", gitArgs...)

	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(res), nil
}
