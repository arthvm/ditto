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
		diff, err := git.Diff(cmd.Context(), git.Stats, git.Branches("feat/pr-generation", "main"))
		if err != nil {
			return fmt.Errorf("diff stats: %w", err)
		}

		fmt.Println(diff)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
}
