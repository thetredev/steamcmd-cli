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
	Use:     "ca <output dir>",
	Version: shared.Version,
	Short:   "Generate CA certificate and private key",
	Run:     certsCaCmdCallback,
}

func init() {
	certsCmd.AddCommand(certsCaCmd)
	certCmdAddFlags(certsCaCmd)
}

func certsCaCmdCallback(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Help()
		return
	}

	certPath := args[0]
	err := os.MkdirAll(certPath, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	ca, err := server.NewCertificateAuthority(certCmdParseFlags(cmd))

	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(fmt.Sprintf("%s/cert.pem", certPath), ca.PEM.Bytes(), os.ModePerm)
	os.WriteFile(fmt.Sprintf("%s/cert.key", certPath), ca.PrivateKeyPEM.Bytes(), os.ModePerm)
}
