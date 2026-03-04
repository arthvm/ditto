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
	Edit  CommitOption = "--edit"
)

func CommitWithMsg(ctx context.Context, msg string, options ...CommitOption) error {
	useEditor := false
	var extraArgs []string
	for _, opt := range options {
		if opt == Edit {
			useEditor = true
		} else {
			extraArgs = append(extraArgs, string(opt))
		}
	}

	var flag string
	if useEditor {
		flag = "-em"
	} else {
		flag = "-m"
	}

	gitArgs := append([]string{"commit", flag, msg}, extraArgs...)

	cmd := exec.CommandContext(ctx, "git", gitArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
