package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "waggle",
	Short: "Waggle CLI",
	Long:  "Named after the \"waggle dance\" bees use to communicate directions and timing to the rest of the hive",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
