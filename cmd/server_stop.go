package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var stopCmd = &cobra.Command{
	Use:     "stop",
	Version: shared.Version,
	Short:   "Stop the game server via the daemon socket",
	Run:     stopCallback,
}

func init() {
	serverCmd.AddCommand(stopCmd)
}

func stopCallback(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Help()
		return
	}

	server.SendMessage(args[0], args[1], shared.ServerStopMessage)
}
