package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var startCmd = &cobra.Command{
	Use:     "start [server certificate file] [server key file]",
	Version: shared.Version,
	Short:   "Start the game server via the daemon socket",
	Run:     startCallback,
}

func init() {
	serverCmd.AddCommand(startCmd)
}

func startCallback(cmd *cobra.Command, args []string) {
	parseCertificateConfig(args)
	server.SendMessage(shared.MESSAGE_SERVER_START)
}
