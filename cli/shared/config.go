package shared

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type SocketConfigType struct {
	SocketIp   string `env:"STEAMCMD_CLI_SOCKET_IP,default=127.0.0.1"`
	SocketPort int    `env:"STEAMCMD_CLI_SOCKET_PORT,default=65000"`
}

var SocketConfig SocketConfigType

func init() {
	ctx := context.Background()

	if err := envconfig.Process(ctx, &SocketConfig); err != nil {
		log.Fatal(err)
	}
}
