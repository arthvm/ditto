package git

import (
	"context"
	"strings"
)

func CurrentBranch(ctx context.Context) (string, error) {
	res, err := run(ctx, "branch", "--show-current")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(res), nil
}
