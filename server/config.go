package server

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type ServerCertificateConfig struct {
	CertificatePath    string `env:"STEAMCMD_CLI_CERTIFICATE_PATH,default=/certs/server/cert.pem"`
	CertificateKeyPath string `env:"STEAMCMD_CLI_CERTIFICATE_KEY_PATH,default=/certs/server/cert.key"`
}

var ServerCertificates ServerCertificateConfig

func init() {
	ctx := context.Background()

	if err := envconfig.Process(ctx, &ServerCertificates); err != nil {
		log.Fatal(err)
	}
}
