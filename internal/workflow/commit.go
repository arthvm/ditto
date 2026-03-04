package workflow

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/arthvm/ditto/internal/git"
	"github.com/arthvm/ditto/internal/llm"
)

type CommitParams struct {
	Amend             bool
	All               bool
	ProviderName      string
	AdditionalContext string
	Issues            []string
}

func Commit(ctx context.Context, progress Progress, params CommitParams) error {
	diffOpts := buildDiffOptions(params.Amend, params.All)

	diff, err := git.Diff(ctx, diffOpts...)
	if err != nil {
		return fmt.Errorf("staged changes: %w", err)
	}

	if strings.TrimSpace(diff) == "" {
		if params.Amend || params.All {
			return errors.New("no changes to commit")
		}
		return errors.New("no staged changes")
	}

	provider, err := llm.GetProvider(params.ProviderName)
	if err != nil {
		return fmt.Errorf("get provider: %w", err)
	}

	progress.StartSpinner(" Generating commit message...")

	msg, err := provider.GenerateCommitMessage(ctx, llm.GenerateCommitParams{
		Diff:              diff,
		Issues:            params.Issues,
		AdditionalContext: params.AdditionalContext,
	})

	progress.StopSpinner()

	if err != nil {
		return fmt.Errorf("generate git commit: %w", err)
	}

	var opts []git.CommitOption
	if params.Amend {
		opts = append(opts, git.Amend)
	}
	if params.All {
		opts = append(opts, git.All)
	}

	return git.CommitWithMessage(ctx, msg, opts...)
}

func buildDiffOptions(amend, all bool) []git.DiffArg {
	switch {
	case amend && all:
		return []git.DiffArg{git.Target("HEAD^")}
	case amend:
		return []git.DiffArg{git.Staged, git.Target("HEAD^")}
	case all:
		return []git.DiffArg{git.Target("HEAD")}
	default:
		return []git.DiffArg{git.Staged}
	}
}
