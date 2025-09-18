/*
Copyright © 2025 Arthur Mariano
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/arthvm/ditto/internal/git"
	"github.com/arthvm/ditto/internal/llm"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Used to generated a git commit message from staged changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		diff, err := git.Diff(cmd.Context(), git.Staged)
		if err != nil {
			return fmt.Errorf("staged changes: %w", err)
		}

		additionalPrompt, err := cmd.Flags().GetString(promptFlagName)
		if err != nil {
			return fmt.Errorf("get prompt flag: %w", err)
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
		s.Suffix = " Generating commit messaging..."

		s.Start()
		defer s.Stop()

		ctx, cancel := context.WithTimeout(cmd.Context(), time.Minute*1)
		defer cancel()

		msg, err := provider.GenerateCommitMessage(ctx, llm.GenerateCommitParams{
			Diff:              diff,
			Issues:            issues,
			AdditionalContext: additionalPrompt,
		})
		if err != nil {
			return fmt.Errorf("generate git commit: %w", err)
		}

		s.Stop()
		if err := git.CommitWithMessage(cmd.Context(), msg); err != nil {
			return fmt.Errorf("execute commit: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
