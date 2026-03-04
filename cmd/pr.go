/*
Copyright © 2025 Arthur Mariano
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arthvm/ditto/internal/platform"
	"github.com/arthvm/ditto/internal/ui"
	"github.com/arthvm/ditto/internal/vcs"
	"github.com/arthvm/ditto/internal/workflow"
)

const (
	noTemplateFlagName = "no-template"
	draftFlagName      = "draft"
)

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

		return workflow.CreatePR(cmd.Context(), workflow.PRDeps{
			VCS:      vcs.Git{},
			Platform: platform.GitHub{},
			Progress: ui.Default(),
		}, workflow.PRParams{
			BaseBranch:        baseBranch,
			HeadBranch:        headBranch,
			ProviderName:      providerName,
			AdditionalContext: additionalPrompt,
			Issues:            issues,
			IgnoreTemplate:    ignoreTemplate,
			Draft:             draft,
		})
	},
}

func init() {
	prCmd.Flags().
		Bool(noTemplateFlagName, false, "Set this flag to ignore any template defined in the repo")

	prCmd.Flags().
		Bool(draftFlagName, false, "Set this flag to create the PR as a draft")

	rootCmd.AddCommand(prCmd)
}
