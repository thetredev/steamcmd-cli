package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var serverFilesListCmd = &cobra.Command{
	Use:     "list",
	Version: shared.Version,
	Short:   "List files on the server.",
	Run:     serverFilesListCallback,
}

func init() {
	serverFilesCmd.AddCommand(serverFilesListCmd)
}

func serverFilesListCallback(cmd *cobra.Command, args []string) {
	server.SendMessage(shared.MESSAGE_SERVER_FILES_LIST, args...)
}
