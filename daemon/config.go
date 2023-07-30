package daemon

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type DaemonConfig struct {
	Application      string `env:"STEAMCMD_SH"`
	ServerAppId      int    `env:"STEAMCMD_SERVER_APPID"`
	ServerMod        string `env:"STEAMCMD_SERVER_MOD"`
	ServerHome       string `env:"STEAMCMD_SERVER_HOME"`
	ServerPort       int    `env:"STEAMCMD_SERVER_PORT,default=27015"`
	ServerMaxPlayers int    `env:"STEAMCMD_SERVER_MAXPLAYERS"`
	ServerMap        string `env:"STEAMCMD_SERVER_MAP"`
	ServerTickrate   int    `env:"STEAMCMD_SERVER_TICKRATE,default=128"`
	ServerThreads    int    `env:"STEAMCMD_SERVER_THREADS,default=3"`
	ServerFpsMax     int    `env:"STEAMCMD_SERVER_FPSMAX,default=300"`
}

var Config DaemonConfig

func init() {
	ctx := context.Background()

	if err := envconfig.Process(ctx, &Config); err != nil {
		log.Fatal(err)
	}
}
