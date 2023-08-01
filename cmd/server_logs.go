package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var logsCmd = &cobra.Command{
	Use:     "logs",
	Version: shared.Version,
	Short:   "Retrieve game server logs from the daemon socket",
	Run:     logsCallback,
}

func init() {
	serverCmd.AddCommand(logsCmd)
}

func logsCallback(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Help()
		return
	}

	server.SendMessage(args[0], args[1], shared.ServerLogsMessage)
}
