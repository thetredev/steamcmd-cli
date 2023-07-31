package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var certsCaCmd = &cobra.Command{
	Use:     "ca",
	Version: shared.Version,
	Short:   "Generate CA certificate and private key",
	Long:    `A longer description`,
	Run:     certsCaCmdCallback,
}

func init() {
	certsCmd.AddCommand(certsCaCmd)
}

func certsCaCmdCallback(cmd *cobra.Command, args []string) {
	const certPath = "certs/ca"
	err := os.MkdirAll(certPath, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	ca, err := server.NewCertificateAuthority()

	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(fmt.Sprintf("%s/cert.pem", certPath), ca.PEM.Bytes(), os.ModePerm)
	os.WriteFile(fmt.Sprintf("%s/cert.key", certPath), ca.PrivateKeyPEM.Bytes(), os.ModePerm)
}
