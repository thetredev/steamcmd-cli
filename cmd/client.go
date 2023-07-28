package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client subcommands to communicate with the daemon socket",
	Long:  `A longer description`,
	Run:   clientCallback,
}

func init() {
	rootCmd.AddCommand(clientCmd)
}

func clientCallback(cmd *cobra.Command, args []string) {
	fmt.Println("client called")
}
