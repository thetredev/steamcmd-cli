package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var commandCmd = &cobra.Command{
	Use:   "command",
	Short: "Send commands to the game server console via the daemon socket",
	Long:  `A longer description`,
	Run:   commandCallback,
}

func init() {
	serverCmd.AddCommand(commandCmd)
}

func commandCallback(cmd *cobra.Command, args []string) {
	server.SendMessage(shared.ServerCommandMessage, args...)
}
