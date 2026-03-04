package workflow

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/arthvm/ditto/internal/prompt"
)

type CommitDeps struct {
	VCS             VCS
	Provider        Provider
	Progress        Progress
	GenerateTimeout time.Duration
}

type CommitParams struct {
	Amend             bool
	All               bool
	Edit              bool
	SystemPrompt      string
	AdditionalContext string
	Issues            []string
}

func Commit(ctx context.Context, deps CommitDeps, params CommitParams) error {
	diff, err := deps.VCS.CommitDiff(ctx, params.Amend, params.All)
	if err != nil {
		return fmt.Errorf("staged changes: %w", err)
	}

	if strings.TrimSpace(diff) == "" {
		if params.Amend || params.All {
			return errors.New("no changes to commit")
		}
		return errors.New("no staged changes")
	}

	system := prompt.CommitSystem(params.SystemPrompt, params.AdditionalContext)
	user := prompt.CommitUser(prompt.CommitParams{
		Diff:   diff,
		Issues: params.Issues,
	})

	deps.Progress.StartSpinner(" Generating commit message...")

	timeout := deps.GenerateTimeout
	if timeout == 0 {
		timeout = generateTimeout
	}

	genCtx, genCancel := context.WithTimeout(ctx, timeout)
	defer genCancel()

	msg, err := deps.Provider.Generate(genCtx, system, user)

	deps.Progress.StopSpinner()

	if err != nil {
		return fmt.Errorf("generate git commit: %w", err)
	}

	return deps.VCS.CommitWithMessage(ctx, msg, params.Amend, params.All, params.Edit)
}
