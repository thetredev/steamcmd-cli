package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var stopCmd = &cobra.Command{
	Use:     "stop [server certificate file] [server key file]",
	Version: shared.Version,
	Short:   "Stop the game server via the daemon socket",
	Run:     stopCallback,
}

func init() {
	serverCmd.AddCommand(stopCmd)
}

func stopCallback(cmd *cobra.Command, args []string) {
	parseCertificateConfig(args)
	server.SendMessage(shared.MESSAGE_SERVER_STOP)
}
