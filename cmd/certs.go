package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/shared"
)

var certsCmd = &cobra.Command{
	Use:     "certs",
	Version: shared.Version,
	Short:   "Cert stuff",
	Long:    `A longer description`,
	Run:     certsCallback,
}

func init() {
	serverCmd.AddCommand(certsCmd)
}

func certsCallback(cmd *cobra.Command, args []string) {
	cmd.Help()
}
