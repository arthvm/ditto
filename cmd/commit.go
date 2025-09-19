/*
Copyright © 2025 Arthur Mariano
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/arthvm/ditto/internal/git"
	"github.com/arthvm/ditto/internal/llm"
)

const (
	amendFlagName = "amend"
	allFlagName   = "all"
)

//TODO: Yeah, this *needs* a refactor. I don't really like how I'm checking
// the flags, especially --all and --amend

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

		diffOpt := []git.DiffArg{}

		switch {
		case amend && all:
			diffOpt = append(diffOpt, git.Target("HEAD^"))
		case amend:
			diffOpt = append(diffOpt, git.Staged, git.Target("HEAD^"))
		case all:
			diffOpt = append(diffOpt, git.Target("HEAD"))
		default:
			diffOpt = append(diffOpt, git.Staged)
		}

		diff, err := git.Diff(cmd.Context(), diffOpt...)
		if err != nil {
			return fmt.Errorf("staged changes: %w", err)
		}

		if strings.TrimSpace(diff) == "" {
			msg := "no staged changes"

			if amend || all {
				msg = "no changes to commit"
			}

			return errors.New(msg)
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
		commitOpts := []git.CommitOption{}

		if amend {
			commitOpts = append(commitOpts, git.Amend)
		}

		if all {
			commitOpts = append(commitOpts, git.All)
		}

		if err := git.CommitWithMessage(cmd.Context(), msg, commitOpts...); err != nil {
			return fmt.Errorf("execute commit: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().
		Bool(amendFlagName, false, "Used to edit the past commit with the current changes")

	commitCmd.Flags().
		BoolP(allFlagName, "a", false, "Used to commit all tracked files")
}
