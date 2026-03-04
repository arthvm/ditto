package vcs

import (
	"context"

	"github.com/arthvm/ditto/internal/git"
)

// Git implements the workflow.VCS interface using the git CLI.
type Git struct{}

func (g Git) CommitDiff(ctx context.Context, amend, all bool) (string, error) {
	opts := buildDiffOptions(amend, all)
	return git.Diff(ctx, opts...)
}

func (g Git) DiffStats(ctx context.Context, base, head string) (string, error) {
	return git.Diff(ctx, git.Stats, git.Branches(base, head))
}

func (g Git) Log(ctx context.Context, base, head string) (string, error) {
	return git.LogRange(ctx, git.Branches(base, head))
}

func (g Git) CurrentBranch(ctx context.Context) (string, error) {
	return git.CurrentBranch(ctx)
}

func (g Git) Root(ctx context.Context) (string, error) {
	return git.Root(ctx)
}

func (g Git) CommitWithMessage(ctx context.Context, msg string, amend, all bool) error {
	var opts []git.CommitOption
	if amend {
		opts = append(opts, git.Amend)
	}
	if all {
		opts = append(opts, git.All)
	}
	return git.CommitWithMsg(ctx, msg, opts...)
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
