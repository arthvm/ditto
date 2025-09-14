package git

import (
	"context"
	"os/exec"
)

var (
	Staged DiffOption = "--staged"
)

type DiffOption string

func Diff(ctx context.Context, options ...DiffOption) (string, error) {
	args := make([]string, len(options))
	for i, opt := range options {
		args[i] = string(opt)
	}
	gitArgs := append([]string{"diff"}, args...)

	cmd := exec.CommandContext(ctx, "git", gitArgs...)

	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(res), nil
}
