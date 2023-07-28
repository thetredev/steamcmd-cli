package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "steamcmd-cli",
	Short: "Custom SteamCMD client implementation written in Go",
	Long:  `Some multiline description.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
