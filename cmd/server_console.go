package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var consoleCmd = &cobra.Command{
	Use:     "console",
	Version: shared.Version,
	Short:   "Send commands to the game server console via the daemon socket",
	Long:    `A longer description`,
	Run:     consoleCallback,
}

func init() {
	serverCmd.AddCommand(consoleCmd)
}

func consoleCallback(cmd *cobra.Command, args []string) {
	server.SendMessage(shared.ServerConsoleCommandMessage, args...)
}
