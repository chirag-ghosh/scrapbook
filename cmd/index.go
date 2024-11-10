package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chirag-ghosh/scrapbook/indexer"
	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index a directory of photos",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// Ask the user for the directory to index
		fmt.Print("Enter the directory to index: ")
		dirPath, _ := reader.ReadString('\n')
		dirPath = strings.TrimSpace(dirPath)

		// Check if the directory exists
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			fmt.Println("Directory does not exist")
			return
		}

		// Get default name
		defaultName := filepath.Base(dirPath)

		// Ask the user for the name of the directory
		fmt.Printf("Enter the name of the directory [%s]: ", defaultName)
		dirName, _ := reader.ReadString('\n')
		dirName = strings.TrimSpace(dirName)

		if dirName == "" {
			dirName = defaultName
		}

		// Index the directory
		if err := indexer.IndexRootDirectory(dirName, dirPath); err != nil {
			log.Fatalf("failed to index directory: %v", err)
		}
	},
}
