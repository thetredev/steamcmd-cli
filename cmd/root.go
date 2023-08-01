package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/shared"
)

var rootCmd = &cobra.Command{
	Use:     "steamcmd-cli",
	Version: shared.Version,
	Short:   "Custom SteamCMD client and game server manager implementation written in Go",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
