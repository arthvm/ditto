/*
Copyright © 2025 Arthur Mariano
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/arthvm/ditto/internal/git"
	"github.com/arthvm/ditto/internal/llm"
)

//TODO: Maybe I should refactor this for better readability...

const (
	noTemplateFlagName = "no-template"
	draftFlagName      = "draft"
	issuesFlagName     = "issues"
)

func findPRTemplate(root string) (string, error) {
	paths := []string{
		filepath.Join(root, ".github", "pull_request_template.md"),
		filepath.Join(root, "docs", "pull_request_template.md"),
		filepath.Join(root, "PULL_REQUEST_TEMPLATE.md"),
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

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Used to generated a pr title and body from commit diff between branches",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseBranch, err := cmd.Flags().GetString(baseBranchFlag)
		if err != nil {
			return fmt.Errorf("get base branch: %w", err)
		}

		headBranch, err := cmd.Flags().GetString(headBranchFlag)
		if err != nil {
			return fmt.Errorf("get head branch: %w", err)
		}

		if headBranch == "" {
			headBranch, err = git.CurrentBranch(cmd.Context())
			if err != nil {
				return fmt.Errorf("get current branch: %w", err)
			}
		}

		additionalPrompt, err := cmd.Flags().GetString(promptFlagName)
		if err != nil {
			return fmt.Errorf("get prompt flag: %w", err)
		}

		ignoreTemplate, err := cmd.Flags().GetBool(noTemplateFlagName)
		if err != nil {
			return fmt.Errorf("get ignore template flag: %w", err)
		}

		draft, err := cmd.Flags().GetBool(draftFlagName)
		if err != nil {
			return fmt.Errorf("get draft flag: %w", err)
		}

		issues, err := cmd.Flags().GetStringSlice(issuesFlagName)
		if err != nil {
			return fmt.Errorf("get issues flag: %w", err)
		}

		providerName, err := cmd.Flags().GetString(providerFlagName)
		if err != nil {
			return fmt.Errorf("get provider flag: %w", err)
		}

		provider, err := llm.GetProvider(providerName)
		if err != nil {
			return err
		}

		s := spinner.New(
			spinner.CharSets[14],
			time.Millisecond*100,
			spinner.WithColor("yellow"),
		)
		s.Suffix = " Generating PR..."

		s.Start()
		defer s.Stop()

		log, err := git.Log(cmd.Context(), git.Branches(baseBranch, headBranch))
		if err != nil {
			return fmt.Errorf("get log: %w", err)
		}

		diff, err := git.Diff(
			cmd.Context(),
			git.Stats,
			git.Branches(baseBranch, headBranch),
		)
		if err != nil {
			return fmt.Errorf("diff stats: %w", err)
		}

		root, err := git.Root(cmd.Context())
		if err != nil {
			return fmt.Errorf("get root dir: %w", err)
		}

		var template string
		if !ignoreTemplate {
			template, err = findPRTemplate(root)
			if err != nil {
				return fmt.Errorf("get pr template: %w", err)
			}
		}

		msg, err := provider.GeneratePr(cmd.Context(), llm.GeneratePrParams{
			HeadBranch:        headBranch,
			BaseBranch:        baseBranch,
			Log:               log,
			DiffStats:         diff,
			AdditionalContext: additionalPrompt,
			Issues:            issues,
			Template:          template,
		})
		if err != nil {
			return fmt.Errorf("generate pr: %w", err)
		}

		nLine := strings.Index(msg, "\n")
		if nLine == -1 {
			return fmt.Errorf("generate pr: failed to generate body")
		}
		title := strings.TrimSpace(msg[:nLine])
		body := strings.TrimSpace(msg[nLine+1:])

		s.Stop()
		if err := git.OpenPr(cmd.Context(), git.OpenPrParams{
			Title:     title,
			Body:      body,
			Head:      headBranch,
			Base:      baseBranch,
			UseEditor: true,
			Draft:     draft,
		}); err != nil {
			return fmt.Errorf("open pr: %w", err)
		}

		return nil
	},
}

func init() {
	prCmd.Flags().
		Bool(noTemplateFlagName, false, "Set this flag to ignore any template defined in the repo")

	prCmd.Flags().
		Bool(draftFlagName, false, "Set this flag to create the PR as a draft")

	prCmd.Flags().
		StringSlice(issuesFlagName, nil, "Specifies the issues that are addressed by the PR.")

	rootCmd.AddCommand(prCmd)
}
