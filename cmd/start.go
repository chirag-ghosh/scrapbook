package cmd

import (
	"github.com/chirag-ghosh/scrapbook/server"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the scrapbook web server",
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer()
	},
}
