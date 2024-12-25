package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "wmn",
	Short: "A CLI to generate weekly reports for student training",
	Long:  `Write My Naita helps generate weekly training reports using LLMs.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use --help to see available commands.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
