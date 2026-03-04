package git

import (
	"context"
)

type logArg interface {
	String() string
	isLogArg()
}

type LogOption string

func (o LogOption) String() string { return string(o) }
func (o LogOption) isLogArg()      {}

func LogRange(ctx context.Context, options ...logArg) (string, error) {
	args := make([]string, len(options))
	for i, opt := range options {
		args[i] = opt.String()
	}
	gitArgs := append([]string{"log", "--pretty=format:%h %s%n%b%n"}, args...)

	return run(ctx, gitArgs...)
}
