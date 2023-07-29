package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the game server via the daemon socket",
	Long:  `A longer description`,
	Run:   updateCallback,
}

func init() {
	serverCmd.AddCommand(updateCmd)
}

func updateCallback(cmd *cobra.Command, args []string) {
	server.SendMessage("update")
}
