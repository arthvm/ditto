package git

import (
	"context"
	"os/exec"
	"strings"
)

func Root(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")

	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(res)), nil
}
