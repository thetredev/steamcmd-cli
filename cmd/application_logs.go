package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/client"
)

var applicationLogsCmd = &cobra.Command{
	Use:   "application",
	Short: "Client subcommand retrieve application logs (SteamCMD) from the daemon socket",
	Long:  `A longer description`,
	Run:   applicationLogsCallback,
}

func init() {
	logsCmd.AddCommand(applicationLogsCmd)
}

func applicationLogsCallback(cmd *cobra.Command, args []string) {
	client.SendMessageToSocket("logs/application")
}
