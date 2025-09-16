package git

import (
	"context"
	"os/exec"
	"strings"
)

func CurrentBranch(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "branch", "--show-current")

	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(res)), nil
}
