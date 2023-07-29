package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/shared"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Version: shared.Version,
	Short:   "Subcommands to communicate with the game server via the daemon socket",
	Long:    `A longer description`,
	Run:     serverCallback,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func serverCallback(cmd *cobra.Command, args []string) {
	cmd.Help()
}
