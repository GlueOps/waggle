package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/glueops/waggle/internal/config"
	"github.com/spf13/cobra"
)

const psqlSessionTimeout = 72 * time.Hour

var pagerMode string

var dbPsqlCmd = &cobra.Command{
	Use:   "psql",
	Short: "Open a psql session to the app database",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.DatabaseURL == "" {
			return errors.New("database.url is empty")
		}

		psql := "psql"
		if runtime.GOOS == "windows" {
			psql = "psql.exe"
		}

		ctx, cancel := context.WithTimeout(context.Background(), psqlSessionTimeout)
		defer cancel()

		args = []string{"-P", "pager=" + pagerMode, cfg.DatabaseURL}
		psqlCmd := exec.CommandContext(ctx, psql, args...)
		psqlCmd.Stdin = os.Stdin
		psqlCmd.Stdout = os.Stdout
		psqlCmd.Stderr = os.Stderr

		fmt.Println("Launching psql…")
		return psqlCmd.Run()
	},
}

func init() {
	dbPsqlCmd.Flags().StringVar(&pagerMode, "pager", "off", "psql pager mode: off, on, always")

	rootCmd.AddCommand(dbPsqlCmd)
}
