package git

import (
	"context"
	"os"
	"os/exec"
)

type CommitOption string

const (
	Amend CommitOption = "--amend"
	All   CommitOption = "--all"
)

func CommitWithMessage(ctx context.Context, msg string, options ...CommitOption) error {
	args := make([]string, len(options))
	for i, opt := range options {
		args[i] = string(opt)
	}
	gitArgs := append([]string{"commit", "-em", msg}, args...)

	cmd := exec.CommandContext(ctx, "git", gitArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
