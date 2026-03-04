/*
Copyright © 2025 Arthur Mariano
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arthvm/ditto/internal/ui"
	"github.com/arthvm/ditto/internal/vcs"
	"github.com/arthvm/ditto/internal/workflow"
)

const (
	amendFlagName = "amend"
	allFlagName   = "all"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Used to generated a git commit message from staged changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		amend, err := cmd.Flags().GetBool(amendFlagName)
		if err != nil {
			return fmt.Errorf("get amend flag: %w", err)
		}

		all, err := cmd.Flags().GetBool(allFlagName)
		if err != nil {
			return fmt.Errorf("get all flag: %w", err)
		}

		additionalPrompt, err := cmd.Flags().GetString(promptFlagName)
		if err != nil {
			return fmt.Errorf("get prompt flag: %w", err)
		}

		issues, err := cmd.Flags().GetStringSlice(issuesFlagName)
		if err != nil {
			return fmt.Errorf("get issues flag: %w", err)
		}

		return workflow.Commit(cmd.Context(), workflow.CommitDeps{
			VCS:      vcs.Git{},
			Provider: provider,
			Progress: ui.Default(),
		}, workflow.CommitParams{
			Amend:             amend,
			All:               all,
			AdditionalContext: additionalPrompt,
			Issues:            issues,
		})
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().
		Bool(amendFlagName, false, "Used to edit the past commit with the current changes")

	commitCmd.Flags().
		BoolP(allFlagName, "a", false, "Used to commit all tracked files")
}
