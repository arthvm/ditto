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
	"github.com/arthvm/ditto/internal/llm/gemini"
	"github.com/arthvm/ditto/internal/llm/ollama"
)

type ProviderFunc func(context.Context, string, string) (string, error)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Used to generated a git commit message from staged changes",
	Long: `Uses the git-commit-message llm model by default -
	but will be able to change providers in the future - to
	generate git commit messages for the stage changes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		diff, err := git.StagedDiff(cmd.Context())
		if err != nil {
			return fmt.Errorf("staged changes: %w", err)
		}

		additionalPrompt, err := cmd.Flags().GetString(promptFlagName)
		if err != nil {
			return fmt.Errorf("get prompt flag: %w", err)
		}

		provider, err := cmd.Flags().GetString(providerFlagName)
		if err != nil {
			return fmt.Errorf("get provider flag: %w", err)
		}

		var providerFunc ProviderFunc
		switch provider {
		case "gemini":
			providerFunc = gemini.GenerateCommitMessage
		case "ollama":
			providerFunc = ollama.GenerateCommitMessage
		default:
			return fmt.Errorf("invalid provider")
		}

		s := spinner.New(
			spinner.CharSets[14],
			time.Millisecond*100,
			spinner.WithColor("yellow"),
		)
		s.Suffix = " Generating commit messaging..."

		s.Start()
		defer s.Stop()

		ctx, cancel := context.WithTimeout(cmd.Context(), time.Second*30)
		defer cancel()

		msg, err := providerFunc(ctx, diff, additionalPrompt)
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
