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
	Run:     consoleCallback,
}

func init() {
	serverCmd.AddCommand(consoleCmd)
}

func consoleCallback(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Help()
		return
	}

	server.SendMessage(args[0], args[1], shared.ServerConsoleCommandMessage, args...)
}
