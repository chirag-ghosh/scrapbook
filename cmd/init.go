package cmd

import (
	"log"

	"github.com/chirag-ghosh/scrapbook/db"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new scrapbook",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.Initialize(); err != nil {
			log.Fatalf("failed to initialize scrapbook: %v", err)
		}
	},
}
