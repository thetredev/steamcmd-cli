package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var serverFilesDownloadCmd = &cobra.Command{
	Use:     "download",
	Version: shared.Version,
	Short:   "Download files from the server.",
	Run:     serverFilesDownloadCallback,
}

func init() {
	serverFilesCmd.AddCommand(serverFilesDownloadCmd)
}

func serverFilesDownloadCallback(cmd *cobra.Command, args []string) {
	rootPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	if len(args) > 1 {
		rootPath = args[1]
	}

	server.SendMessage(shared.MESSAGE_SERVER_FILES_DOWNLOAD, args[0], rootPath)
}
