package workflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/arthvm/ditto/internal/llm"
	"github.com/arthvm/ditto/internal/prompt"
)

type PRDeps struct {
	VCS      VCS
	Platform Platform
	Progress Progress
}

type PRParams struct {
	BaseBranch        string
	HeadBranch        string
	ProviderName      string
	AdditionalContext string
	Issues            []string
	IgnoreTemplate    bool
	Draft             bool
}

func CreatePR(ctx context.Context, deps PRDeps, params PRParams) error {
	headBranch := params.HeadBranch
	if headBranch == "" {
		var err error
		headBranch, err = deps.VCS.CurrentBranch(ctx)
		if err != nil {
			return fmt.Errorf("get current branch: %w", err)
		}
	}

	provider, err := llm.GetProvider(params.ProviderName)
	if err != nil {
		return fmt.Errorf("get provider: %w", err)
	}

	log, err := deps.VCS.Log(ctx, params.BaseBranch, headBranch)
	if err != nil {
		return fmt.Errorf("get log: %w", err)
	}

	diff, err := deps.VCS.DiffStats(ctx, params.BaseBranch, headBranch)
	if err != nil {
		return fmt.Errorf("diff stats: %w", err)
	}

	root, err := deps.VCS.Root(ctx)
	if err != nil {
		return fmt.Errorf("get root dir: %w", err)
	}

	var template string
	if !params.IgnoreTemplate {
		template, err = deps.Platform.FindPRTemplate(root)
		if err != nil {
			return fmt.Errorf("get pr template: %w", err)
		}
	}

	system := prompt.PRSystem(template, params.AdditionalContext)
	user := prompt.PRUser(prompt.PRParams{
		HeadBranch: headBranch,
		BaseBranch: params.BaseBranch,
		Log:        log,
		DiffStats:  diff,
		Issues:     params.Issues,
	})

	deps.Progress.StartSpinner(" Generating PR...")

	msg, err := provider.Generate(ctx, system, user)

	deps.Progress.StopSpinner()

	if err != nil {
		return fmt.Errorf("generate pr: %w", err)
	}

	title, body, err := parsePRMessage(msg)
	if err != nil {
		return err
	}

	return deps.Platform.OpenPR(ctx, OpenPRParams{
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
