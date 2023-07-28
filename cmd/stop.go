package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the game server via the daemon socket",
	Long:  `A longer description`,
	Run:   stopCallback,
}

func init() {
	serverCmd.AddCommand(stopCmd)
}

func stopCallback(cmd *cobra.Command, args []string) {
	server.SendMessageToSocket("stop")
}
