package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

var (
	Staged DiffOption = "--staged"
	Stats  DiffOption = "--stat"
)

type DiffArg interface {
	String() string
	isDiffArg()
}

type DiffOption string

func (o DiffOption) String() string { return string(o) }
func (o DiffOption) isDiffArg()     {}

func Cached(target string) DiffOption {
	return DiffOption(fmt.Sprintf("--cached %s", target))
}

func Diff(ctx context.Context, options ...DiffArg) (string, error) {
	var args []string

	for _, opt := range options {
		parts := strings.Fields(opt.String())
		args = append(args, parts...)
	}
	gitArgs := append([]string{"diff"}, args...)

	cmd := exec.CommandContext(ctx, "git", gitArgs...)

	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(res), nil
}
