package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/client"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the game server via the daemon socket",
	Long:  `A longer description`,
	Run:   updateCallback,
}

func init() {
	clientCmd.AddCommand(updateCmd)
}

func updateCallback(cmd *cobra.Command, args []string) {
	client.SendMessageToSocket("update")
}
