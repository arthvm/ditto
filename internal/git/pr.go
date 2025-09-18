package git

import (
	"context"
	"os"
	"os/exec"
)

type OpenPrParams struct {
	Title     string
	Head      string
	Base      string
	Body      string
	UseEditor bool
	Draft     bool
}

func OpenPr(ctx context.Context, params OpenPrParams) error {
	args := []string{
		"--title", params.Title,
		"--body", params.Body,
		"--base", params.Base,
		"--head", params.Head,
	}

	if params.UseEditor {
		args = append(args, "--editor")
	}

	if params.Draft {
		args = append(args, "--draft")
	}

	ghArgs := append([]string{"pr", "create"}, args...)

	cmd := exec.CommandContext(ctx, "gh", ghArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
