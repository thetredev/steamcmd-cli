package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/daemon"
	"github.com/thetredev/steamcmd-cli/shared"
)

var daemonCmd = &cobra.Command{
	Use:     "daemon",
	Version: shared.Version,
	Short:   "Run the daemon socket",
	Run:     daemonCallback,
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}

func daemonCallback(cmd *cobra.Command, args []string) {
	daemon.StartSocket()
}
