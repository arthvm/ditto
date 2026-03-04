package platform

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/arthvm/ditto/internal/workflow"
)

// GitHub implements the workflow.Platform interface using the gh CLI
// and GitHub-specific conventions (e.g. .github/pull_request_template.md).
type GitHub struct{}

func (g GitHub) FindPRTemplate(repoRoot, customPath string) (string, error) {
	if customPath != "" {
		p := customPath
		if !filepath.IsAbs(p) {
			p = filepath.Join(repoRoot, p)
		}
		content, err := os.ReadFile(p)
		if err == nil {
			return string(content), nil
		}

		if !os.IsNotExist(err) {
			return "", fmt.Errorf("pr template: %w", err)
		}
	}

	paths := []string{
		filepath.Join(repoRoot, ".github", "pull_request_template.md"),
		filepath.Join(repoRoot, "docs", "pull_request_template.md"),
		filepath.Join(repoRoot, "PULL_REQUEST_TEMPLATE.md"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			content, err := os.ReadFile(p)
			if err != nil {
				return "", err
			}

			return string(content), nil
		}
	}

	return "", nil
}

func (g GitHub) OpenPR(ctx context.Context, params workflow.OpenPRParams) error {
	args := []string{
		"pr", "create",
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

	cmd := exec.CommandContext(ctx, "gh", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
