package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/glueops/waggle/internal/config"
	"github.com/glueops/waggle/internal/database"
	"github.com/glueops/waggle/internal/jobs"
	"github.com/glueops/waggle/internal/service"
	"github.com/glueops/waggle/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run the River job worker process",
	Long:  "Runs background workers (tenant provisioner, tenant destroyer, ...). Scale independently from `serve`.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWorker()
	},
}

func runWorker() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("open pgx pool: %w", err)
	}
	defer pool.Close()

	controlDB, err := database.OpenControlDB(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("open control db: %w", err)
	}

	masterKey, err := utils.DecodeB64(cfg.EncryptionMasterKey)
	if err != nil || len(masterKey) != 32 {
		return fmt.Errorf("invalid ENCRYPTION_MASTER_KEY (must be 32 bytes base64)")
	}

	tenants, err := database.NewTenantManager(controlDB, cfg.EncryptionMasterKey)
	if err != nil {
		return fmt.Errorf("init tenant manager: %w", err)
	}

	provisioner, err := jobs.NewTenantProvisionerService(controlDB, cfg.DatabaseURL, cfg.AdminDatabaseURL, masterKey)
	if err != nil {
		return fmt.Errorf("init provisioner: %w", err)
	}

	destroyer, err := jobs.NewTenantDestroyerService(controlDB, cfg.DatabaseURL, cfg.AdminDatabaseURL, tenants)
	if err != nil {
		return fmt.Errorf("init destroyer: %w", err)
	}

	discovery := jobs.NewHypervisorDiscoveryService(service.NewFleetService(tenants, service.ReservationDefaults{
		CPU: cfg.ReserveCPU, RAMGB: cfg.ReserveRAMGB, DiskGB: cfg.ReserveDiskGB,
	}))

	workers := river.NewWorkers()
	river.AddWorker(workers, jobs.NewTenantProvisionerWorker(provisioner))
	river.AddWorker(workers, jobs.NewTenantDestroyerWorker(destroyer))
	river.AddWorker(workers, jobs.NewHypervisorDiscoveryWorker(discovery))

	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: workers,
	})
	if err != nil {
		return fmt.Errorf("new river client: %w", err)
	}

	log.Printf("worker starting (admin_dsn_override=%t)", cfg.AdminDatabaseURL != "")
	if err := client.Start(ctx); err != nil {
		return fmt.Errorf("start river client: %w", err)
	}

	<-ctx.Done()
	log.Printf("worker stopping")
	stopCtx, stopCancel := context.WithCancel(context.Background())
	defer stopCancel()
	return client.Stop(stopCtx)
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
