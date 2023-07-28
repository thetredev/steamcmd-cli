package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/client"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the game server via the daemon socket",
	Long:  `A longer description`,
	Run:   stopCallback,
}

func init() {
	clientCmd.AddCommand(stopCmd)
}

func stopCallback(cmd *cobra.Command, args []string) {
	client.SendMessageToSocket("stop")
}
