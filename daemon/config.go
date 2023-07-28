package daemon

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type DaemonConfig struct {
	Application     string `env:"STEAMCMD_SH"`
	ServerAppConfig string `env:"STEAMCMD_SERVER_APP_CONFIG"`
	ServerAppId     int    `env:"STEAMCMD_SERVER_APPID"`
	ServerHome      string `env:"STEAMCMD_SERVER_HOME"`
	ServerPort      int    `env:"STEAMCMD_SERVER_PORT,default=27015"`
}

var Config DaemonConfig

func init() {
	ctx := context.Background()

	if err := envconfig.Process(ctx, &Config); err != nil {
		log.Fatal(err)
	}
}
