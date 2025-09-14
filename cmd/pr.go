/*
Copyright © 2025 Arthur Mariano
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arthvm/ditto/internal/git"
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

		if headBranch == "" {
			headBranch, err = git.CurrentBranch(cmd.Context())
			if err != nil {
				return fmt.Errorf("get current branch: %w", err)
			}
		}

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

		fmt.Println(log)
		fmt.Println(diff)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
}
