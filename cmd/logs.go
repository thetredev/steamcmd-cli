package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Client subcommand retrieve certain logs from the daemon socket",
	Long:  `A longer description`,
	Run:   logsCallback,
}

func init() {
	clientCmd.AddCommand(logsCmd)
}

func logsCallback(cmd *cobra.Command, args []string) {
	client.SendMessageToSocket("logs")
}
