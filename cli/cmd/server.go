package cmd

import (
	"os"
	"strings"

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
		cert := args[0]
		key := args[1]

		if _, errCert := os.Stat(cert); errCert != nil {
			return
		}

		if _, errKey := os.Stat(key); errKey != nil {
			return
		}

		if !strings.HasSuffix(cert, ".pem") || !strings.HasSuffix(cert, ".crt") {
			return
		}

		if !strings.HasSuffix(key, ".key") {
			return
		}

		server.ServerCertificates.CertificatePath = args[0]
		server.ServerCertificates.CertificateKeyPath = args[1]
	}
}
