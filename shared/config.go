package shared

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type SharedConfig struct {
	SocketIp   string `env:"STEAMCMD_CLI_SOCKET_IP,default=127.0.0.1"`
	SocketPort int    `env:"STEAMCMD_CLI_SOCKET_PORT,default=65000"`
}

var Config SharedConfig

func init() {
	ctx := context.Background()

	if err := envconfig.Process(ctx, &Config); err != nil {
		log.Fatal(err)
	}
}
