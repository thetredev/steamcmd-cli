package shared

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type SocketConfigType struct {
	SocketIp   string `env:"STEAMCMD_CLI_SOCKET_IP,default=0.0.0.0"`
	SocketPort int    `env:"STEAMCMD_CLI_SOCKET_PORT,default=27015"`
}

var SocketConfig SocketConfigType

func init() {
	ctx := context.Background()

	if err := envconfig.Process(ctx, &SocketConfig); err != nil {
		log.Fatal(err)
	}
}
