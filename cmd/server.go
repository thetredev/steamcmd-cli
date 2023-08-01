package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Version: shared.Version,
	Short:   "Subcommands to communicate with the game server via the daemon socket",
	Run:     serverCallback,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func serverCallback(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func parseCertificateConfig(args []string) {
	if len(args) >= 2 {
		server.ServerCertificates.CertificatePath = args[0]
		server.ServerCertificates.CertificateKeyPath = args[1]
	}
}
