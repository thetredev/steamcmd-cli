package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/shared"
)

var serverFilesCmd = &cobra.Command{
	Use:     "files",
	Version: shared.Version,
	Short:   "Subcommands to manage files on the server.",
	Run:     serverFilesCallback,
}

func init() {
	serverCmd.AddCommand(serverFilesCmd)
}

func serverFilesCallback(cmd *cobra.Command, args []string) {
	cmd.Help()
}
