package workflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/arthvm/ditto/internal/git"
	"github.com/arthvm/ditto/internal/llm"
)

type PRParams struct {
	BaseBranch        string
	HeadBranch        string
	ProviderName      string
	AdditionalContext string
	Issues            []string
	IgnoreTemplate    bool
	Draft             bool
}

func CreatePR(ctx context.Context, progress Progress, params PRParams) error {
	headBranch := params.HeadBranch
	if headBranch == "" {
		var err error
		headBranch, err = git.CurrentBranch(ctx)
		if err != nil {
			return fmt.Errorf("get current branch: %w", err)
		}
	}

	provider, err := llm.GetProvider(params.ProviderName)
	if err != nil {
		return fmt.Errorf("get provider: %w", err)
	}

	log, err := git.Log(ctx, git.Branches(params.BaseBranch, headBranch))
	if err != nil {
		return fmt.Errorf("get log: %w", err)
	}

	diff, err := git.Diff(ctx, git.Stats, git.Branches(params.BaseBranch, headBranch))
	if err != nil {
		return fmt.Errorf("diff stats: %w", err)
	}

	root, err := git.Root(ctx)
	if err != nil {
		return fmt.Errorf("get root dir: %w", err)
	}

	var template string
	if !params.IgnoreTemplate {
		template, err = git.FindPRTemplate(root)
		if err != nil {
			return fmt.Errorf("get pr template: %w", err)
		}
	}

	progress.StartSpinner(" Generating PR...")

	msg, err := provider.GeneratePR(ctx, llm.GeneratePRParams{
		HeadBranch:        headBranch,
		BaseBranch:        params.BaseBranch,
		Log:               log,
		DiffStats:         diff,
		AdditionalContext: params.AdditionalContext,
		Issues:            params.Issues,
		Template:          template,
	})

	progress.StopSpinner()

	if err != nil {
		return fmt.Errorf("generate pr: %w", err)
	}

	title, body, err := parsePRMessage(msg)
	if err != nil {
		return err
	}

	return git.OpenPR(ctx, git.OpenPRParams{
		Title:     title,
		Body:      body,
		Head:      headBranch,
		Base:      params.BaseBranch,
		UseEditor: true,
		Draft:     params.Draft,
	})
}

func parsePRMessage(msg string) (title string, body string, err error) {
	before, after, ok := strings.Cut(msg, "\n")
	if !ok {
		return "", "", fmt.Errorf("generate pr: failed to generate body")
	}

	return strings.TrimSpace(before), strings.TrimSpace(after), nil
}
