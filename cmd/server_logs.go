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
	Long:    `A longer description`,
	Run:     logsCallback,
}

func init() {
	serverCmd.AddCommand(logsCmd)
}

func logsCallback(cmd *cobra.Command, args []string) {
	server.SendMessage(shared.ServerLogsMessage)
}
