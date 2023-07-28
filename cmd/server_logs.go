package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/client"
)

var serverLogsCmd = &cobra.Command{
	Use:   "server",
	Short: "Client subcommand retrieve server logs (HLDS/SRCDS) from the daemon socket",
	Long:  `A longer description`,
	Run:   serverLogsCallback,
}

func init() {
	logsCmd.AddCommand(serverLogsCmd)
}

func serverLogsCallback(cmd *cobra.Command, args []string) {
	client.SendMessageToSocket("logs/server")
}
