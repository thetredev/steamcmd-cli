package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/daemon"
	"github.com/thetredev/steamcmd-cli/shared"
)

var daemonCmd = &cobra.Command{
	Use:     "daemon <CA certificate input file path> <Daemon certificate input key path>",
	Version: shared.Version,
	Short:   "Run the daemon socket",
	Run:     daemonCallback,
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}

func daemonCallback(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		cmd.Help()
		return
	}

	daemon.StartSocket(args[0], args[1], args[2])
}
