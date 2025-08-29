package git

import (
	"context"
	"os"
	"os/exec"
)

func CommitWithMessage(ctx context.Context, msg string) error {
	cmd := exec.CommandContext(ctx, "git", "commit", "-em", msg)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
