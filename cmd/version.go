package cmd

import (
	"fmt"

	"github.com/glueops/waggle/internal/buildinfo"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of waggle",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("waggle %s (%s) %s\n", buildinfo.Version, buildinfo.Commit, buildinfo.Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
