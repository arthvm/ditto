package git

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

type gitArg interface {
	String() string
	isDiffArg()
	isLogArg()
}

type GitOption string

func (o GitOption) String() string { return string(o) }
func (o GitOption) isDiffArg()     {}
func (o GitOption) isLogArg()      {}

func Branches(base string, head string) gitArg {
	return GitOption(fmt.Sprintf("%s..%s", base, head))
}

func run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)

	res, err := cmd.Output()
	if err != nil {
		var exitErr exec.ExitError
		if errors.Is(err, &exitErr) && len(exitErr.Stderr) > 0 {
			return "", fmt.Errorf("git %s: %s", args[0], string(exitErr.Stderr))
		}
		return "", fmt.Errorf("git %s: %w", args[0], err)
	}

	return string(res), nil
}
