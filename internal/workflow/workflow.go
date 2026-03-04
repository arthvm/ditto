package workflow

import (
	"context"
	"time"
)

// generateTimeout bounds the LLM generation call. Interactive phases
// (editor, gh pr create) run without a timeout so they can't be killed
// by an expired deadline.
const generateTimeout = 2 * time.Minute

// Provider generates text from a system prompt and user prompt.
type Provider interface {
	Generate(ctx context.Context, system, user string) (string, error)
}

// Progress reports long-running operation status to the user.
// Implementations control how progress is displayed: a CLI spinner,
// a TUI progress bar, or a no-op for non-interactive use.
type Progress interface {
	StartSpinner(label string)
	StopSpinner()
}

// VCS abstracts version control operations (git, jj, etc.) for
// testability and to decouple workflows from a specific VCS CLI.
type VCS interface {
	// CommitDiff returns the diff for changes that would be committed
	// given the amend and all flags.
	CommitDiff(ctx context.Context, amend, all bool) (string, error)

	// DiffStats returns the diffstat between two branches.
	DiffStats(ctx context.Context, base, head string) (string, error)

	// Log returns the commit log between two branches.
	Log(ctx context.Context, base, head string) (string, error)

	// CurrentBranch returns the name of the currently checked-out branch.
	CurrentBranch(ctx context.Context) (string, error)

	// Root returns the absolute path to the repository root.
	Root(ctx context.Context) (string, error)

	// CommitWithMessage creates a commit with the given message, opening
	// the user's editor for final editing.
	CommitWithMessage(ctx context.Context, msg string, amend, all bool) error
}

// OpenPRParams holds the parameters for opening a pull request.
type OpenPRParams struct {
	Title     string
	Body      string
	Head      string
	Base      string
	UseEditor bool
	Draft     bool
}

// Platform abstracts hosting platform operations (GitHub, GitLab, etc.)
// for testability and to allow swapping platforms independently of the VCS.
type Platform interface {
	// FindPRTemplate looks for a pull request template file in the
	// repository and returns its contents, or empty string if none found.
	// If customPath is non-empty it is checked first (relative to repoRoot
	// if not absolute) before falling back to well-known default locations.
	FindPRTemplate(repoRoot, customPath string) (string, error)

	// OpenPR creates a pull request via the platform CLI (e.g. gh, glab).
	OpenPR(ctx context.Context, params OpenPRParams) error
}
