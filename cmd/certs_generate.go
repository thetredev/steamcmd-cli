package cmd

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var certsGenerateCmd = &cobra.Command{
	Use:     "generate",
	Version: shared.Version,
	Short:   "Generate daemon/server certificate and private key",
	Long:    `A longer description`,
	Run:     certsGenerateCallback,
}

func init() {
	certsCmd.AddCommand(certsGenerateCmd)
}

func certsGenerateCallback(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		log.Fatal("Please provide CA pem, CA key and output path!")
	}

	caCertPath := args[0]

	if _, err := os.Stat(caCertPath); os.IsNotExist(err) {
		log.Fatal(err)
	}

	caKeyPath := args[1]

	if _, err := os.Stat(caKeyPath); os.IsNotExist(err) {
		log.Fatal(err)
	}

	outputPath := args[2]

	err := os.MkdirAll(outputPath, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	caPemBytes, err := os.ReadFile(caCertPath)

	if err != nil {
		log.Fatal(err)
	}

	caCertPemBlock, _ := pem.Decode(caPemBytes)
	caCert, err := x509.ParseCertificate(caCertPemBlock.Bytes)

	if err != nil {
		log.Fatal(err)
	}

	caKeyBytes, err := os.ReadFile(caKeyPath)

	if err != nil {
		log.Fatal(err)
	}

	caKeyPemBlock, _ := pem.Decode(caKeyBytes)
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyPemBlock.Bytes)

	if err != nil {
		log.Fatal(err)
	}

	ca := &server.Certificate{
		X509:       *caCert,
		PrivateKey: caKey,
	}

	issued, err := server.IssueCertificate(ca, net.ParseIP("127.0.0.1"))

	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(fmt.Sprintf("%s/cert.pem", outputPath), issued.Certificate.PEM.Bytes(), os.ModePerm)
	os.WriteFile(fmt.Sprintf("%s/cert.key", outputPath), issued.Certificate.PrivateKeyPEM.Bytes(), os.ModePerm)
}
