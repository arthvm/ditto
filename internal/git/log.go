package git

import (
	"context"
	"os/exec"
)

type logArg interface {
	String() string
	isLogArg()
}

type LogOption string

func (o LogOption) String() string { return string(o) }
func (o LogOption) isLogArg()      {}

func Log(ctx context.Context, options ...logArg) (string, error) {
	args := make([]string, len(options))
	for i, opt := range options {
		args[i] = opt.String()
	}
	gitArgs := append([]string{"log", "--pretty=format:%h %s%n%b%n"}, args...)

	cmd := exec.CommandContext(ctx, "git", gitArgs...)

	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(res), nil
}
