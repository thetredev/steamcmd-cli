package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var consoleCmd = &cobra.Command{
	Use:     "console [server certificate file] [server key file] <command list>",
	Version: shared.Version,
	Short:   "Send commands to the game server console via the daemon socket",
	Run:     consoleCallback,
}

func init() {
	serverCmd.AddCommand(consoleCmd)
}

func consoleCallback(cmd *cobra.Command, args []string) {
	parseCertificateConfig(args)
	server.SendMessage(shared.MESSAGE_SERVER_CONSOLE_COMMAND, args...)
}
