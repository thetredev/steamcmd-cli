package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/client"
)

var commandCmd = &cobra.Command{
	Use:   "command",
	Short: "Client subcommand send commands to the game server console via the daemon socket",
	Long:  `A longer description`,
	Run:   commandCallback,
}

func init() {
	clientCmd.AddCommand(commandCmd)
}

func commandCallback(cmd *cobra.Command, args []string) {
	client.SendMessageToSocket("command", args...)
}
