package git

import (
	"context"
	"strings"
)

func Root(ctx context.Context) (string, error) {
	res, err := run(ctx, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res), nil
}
