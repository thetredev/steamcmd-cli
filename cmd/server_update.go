package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Version: shared.Version,
	Short:   "Update the game server via the daemon socket",
	Run:     updateCallback,
}

func init() {
	serverCmd.AddCommand(updateCmd)
}

func updateCallback(cmd *cobra.Command, args []string) {
	parseCertificateConfig(args)
	server.SendMessage(shared.ServerUpdateMessage)
}
