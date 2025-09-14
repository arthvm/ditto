/*
Copyright © 2025 Arthur Mariano
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/arthvm/ditto/internal/git"
	"github.com/arthvm/ditto/internal/llm"
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

		additionalPrompt, err := cmd.Flags().GetString(promptFlagName)
		if err != nil {
			return fmt.Errorf("get prompt flag: %w", err)
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

		msg, err := provider.GeneratePr(cmd.Context(), llm.GeneratePrParams{
			HeadBranch: headBranch,
			BaseBranch: baseBranch,
			Log:        log,
			DiffStats:  diff,
		}, additionalPrompt)
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
		}); err != nil {
			return fmt.Errorf("open pr: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
}
