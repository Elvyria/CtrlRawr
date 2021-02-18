package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dbPath = "data.db"

var rootCmd = &cobra.Command{
	Use:   "ctrlrawr",
	Short: "Short Description",
	Long:  "Long Description",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
