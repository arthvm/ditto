/*
Copyright © 2025 Arthur Mariano
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	promptFlagName = "prompt"

// providerFlagName = "provider"
)

var rootCmd = &cobra.Command{
	Use:   "ditto",
	Short: "An ai-application to simplify git workflows",
	Long: `An AI application that allow devs to work faster by
	delegating some of the tedious git related operations to an AI agent.`,
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

	// rootCmd.PersistentFlags().
	// 	String(providerFlagName, "", "Used to select the model to be used")
}
