package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/daemon"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the daemon socket",
	Long:  `A longer description`,
	Run:   daemonCallback,
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}

func daemonCallback(cmd *cobra.Command, args []string) {
	daemon.StartSocket()
}
