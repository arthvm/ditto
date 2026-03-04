/*
Copyright © 2025 Arthur Mariano
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/arthvm/ditto/internal/config"
	"github.com/arthvm/ditto/internal/llm"
	"github.com/arthvm/ditto/internal/llm/gemini"
	"github.com/arthvm/ditto/internal/llm/ollama"
)

const (
	promptFlagName   = "prompt"
	providerFlagName = "provider"
	modelFlagName    = "model"
	issuesFlagName   = "issues"
)

// Resolved at startup by PersistentPreRunE, available to all subcommands.
var (
	appConfig config.Config
	provider  llm.Provider
)

var rootCmd = &cobra.Command{
	Use:   "ditto",
	Short: "An ai-application to simplify git workflows",
	Long: `An AI application that allow devs to work faster by
	delegating some of the tedious git related operations to an AI agent.`,
	PersistentPreRunE: setup,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().
		String(promptFlagName, "", "Used to provide additional context to the model")

	rootCmd.PersistentFlags().
		String(providerFlagName, "", "LLM provider to use (gemini, ollama)")

	rootCmd.PersistentFlags().
		String(modelFlagName, "", "Model name to use with the selected provider")

	rootCmd.PersistentFlags().
		StringSlice(issuesFlagName, nil, "Specifies the issues that are addressed by the operation.")
}

func setup(cmd *cobra.Command, _ []string) error {
	repoRoot, _ := repoRootDir()

	cfg, err := config.Load(repoRoot)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if cmd.Flags().Changed(providerFlagName) {
		cfg.Provider, _ = cmd.Flags().GetString(providerFlagName)
	}
	if cmd.Flags().Changed(modelFlagName) {
		model, _ := cmd.Flags().GetString(modelFlagName)
		cfg.SetModelForProvider(model)
	}

	appConfig = cfg
	provider, err = buildProvider(cfg)
	if err != nil {
		return err
	}

	return nil
}

func buildProvider(cfg config.Config) (llm.Provider, error) {
	switch {
	case cfg.Provider == "ollama":
		return ollama.New(cfg.Ollama.Host, cfg.Ollama.Model, cfg.LLM.Temperature), nil

	case strings.HasPrefix(cfg.Provider, "gemini"):
		if cfg.Gemini.APIKey != "" {
			os.Setenv("GEMINI_API_KEY", cfg.Gemini.APIKey)
		}
		return gemini.New(cfg.Gemini.Model, cfg.LLM.Temperature), nil

	default:
		return nil, fmt.Errorf("unknown provider: %q", cfg.Provider)
	}
}

func repoRootDir() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
