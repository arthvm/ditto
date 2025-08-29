package gemini

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var systemPrompt = `
You are a Git and Conventional Commits expert. Your task is to analyze a Git diff and generate a commit message strictly following the Conventional Commits standard.

## Conventional Commits Rules:
- **Format**: '<type>(<scope>): <description>'
- **Valid types**:
  - 'feat': new feature
  - 'fix': bug fix
  - 'docs': documentation
  - 'style': formatting (no logic changes)
  - 'refactor': code refactoring
  - 'test': adding/fixing tests
  - 'chore': maintenance tasks, build, dependencies
  - 'perf': performance improvement
  - 'ci': CI/CD changes
  - 'revert': revert previous commit

## Instructions:
1. Carefully analyze the provided diff
2. Identify the predominant type of change
3. Determine if there's a relevant scope (optional)
4. Create a concise description (max 50 characters)
5. If needed, add explanatory body after blank line
6. For breaking changes, add '!' after type/scope and 'BREAKING CHANGE:' in footer

## Response format:
Provide only the final commit message, without additional explanations.

---
`

func GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	promptFile, err := os.CreateTemp("", "gemini-prompt-*.md")
	if err != nil {
		return "", fmt.Errorf("create prompt file: %w", err)
	}
	defer os.Remove(promptFile.Name())
	defer promptFile.Close()

	if _, err := promptFile.WriteString(systemPrompt); err != nil {
		return "", fmt.Errorf("write to prompt file: %w", err)
	}

	cmd := exec.CommandContext(ctx, "gemini")
	cmd.Env = append(os.Environ(), fmt.Sprintf("GEMINI_SYSTEM_MD=%s", promptFile.Name()))

	cmd.Stdin = strings.NewReader(diff)

	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(res), nil
}
