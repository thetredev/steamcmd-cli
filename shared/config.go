package shared

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type SharedConfig struct {
	SocketPort int `env:"STEAMCMD_CLI_SOCKET_PORT,default=65000"`
}

var Config SharedConfig

func init() {
	ctx := context.Background()

	if err := envconfig.Process(ctx, &Config); err != nil {
		log.Fatal(err)
	}
}
