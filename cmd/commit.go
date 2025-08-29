/*
Copyright © 2025 Arthur Mariano
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arthvm/ditto/internal/git"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Used to generated a git commit message from staged changes",
	Long: `Uses the git-commit-message llm model by default -
	but will be able to change providers in the future - to
	generate git commit messages for the stage changes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := git.StagedDiff(cmd.Context())
		if err != nil {
			return fmt.Errorf("staged changes: %w", err)
		}

		fmt.Println(res)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
