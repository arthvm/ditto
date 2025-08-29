/*
Copyright © 2025 Arthur Mariano
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arthvm/ditto/git"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
