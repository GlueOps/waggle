package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Eyrie",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("waggle %s (%s) %s\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
