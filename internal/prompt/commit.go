package prompt

import (
	"fmt"
	"strings"
)

type CommitParams struct {
	Diff              string
	Issues            []string
	AdditionalContext string
}

func CommitSystem(customPrompt, additionalContext string) string {
	convention := defaultCommitConvention
	if customPrompt != "" {
		convention = customPrompt
	}

	return fmt.Sprintf(`You are a Git commit message expert. Your task is to analyze a Git diff and generate a commit message following the convention below.

%s

## Instructions:
1. Carefully analyze the provided diff
2. Generate a commit message that follows the convention above
3. Add a footer to reference the provided issues if any are given

## Response format:
Provide only the final commit message, without additional explanations.

%s

---
`, convention, wrapAdditionalContext(additionalContext))
}

const defaultCommitConvention = `## Conventional Commits Rules:
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
- Create a concise description (max 50 characters)
- If needed, add explanatory body after blank line
- For breaking changes, add '!' after type/scope and 'BREAKING CHANGE:' in footer
- Issue references in footer:
  - Use Closes #<issue_number> for fix or feat
  - Use Refs #<issue_number> for all other types
  - If multiple issues, group by keyword: Closes #1, #2
  - BREAKING CHANGE footer must be separated by a blank line from other metadata`

func CommitUser(params CommitParams) string {
	return fmt.Sprintf(`--- DIFF START ---
%s
--- DIFF END ---
--- RELATED ISSUES START ---
%s
--- RELATED ISSUES END ---
`, params.Diff, strings.Join(params.Issues, "\n"))
}
