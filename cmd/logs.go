package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/client"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve game server logs from the daemon socket",
	Long:  `A longer description`,
	Run:   logsCallback,
}

func init() {
	clientCmd.AddCommand(logsCmd)
}

func logsCallback(cmd *cobra.Command, args []string) {
	client.SendMessageToSocket("logs")
}
