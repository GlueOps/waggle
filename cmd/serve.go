package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/glueops/waggle/internal/api"
	"github.com/glueops/waggle/internal/app"
	"github.com/glueops/waggle/internal/config"
	"github.com/glueops/waggle/internal/database"
	"github.com/glueops/waggle/internal/jobs"
	"github.com/glueops/waggle/internal/repo"
	"github.com/glueops/waggle/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Nexus server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer()
	},
}

func runServer() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return err
	}
	httpAddr := fmt.Sprintf("%s:%s", cfg.BindHost, cfg.BindPort)

	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to pgx pool for River: %v", err)
	}
	defer dbPool.Close()

	controlDB, err := database.OpenControlDB(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("open control db: %w", err)
	}

	var tenants *database.TenantManager
	if strings.TrimSpace(cfg.EncryptionMasterKey) != "" {
		tenants, err = database.NewTenantManager(controlDB, cfg.EncryptionMasterKey)
		if err != nil {
			return fmt.Errorf("init tenant manager: %w", err)
		}
	}

	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{})
	if err != nil {
		return fmt.Errorf("init river client: %w", err)
	}

	jobEnqueuer, err := jobs.NewEnqueuer(riverClient)
	if err != nil {
		return fmt.Errorf("init job enqueuer: %w", err)
	}

	tokenSvc, err := service.NewTokenService(
		cfg.JWTSecret,
		cfg.JWTIssuer,
		time.Duration(cfg.JWTAccessTTLMin)*time.Minute,
		time.Duration(cfg.JWTRefreshTTLHour)*time.Hour,
		cfg.JWTAudience,
	)
	if err != nil {
		return fmt.Errorf("init token service: %w", err)
	}

	sessionRepo := repo.NewTokenSessionRepoGorm(controlDB)
	auditRepo := repo.NewAuthAuditRepoGorm(controlDB)
	userRepo := repo.NewUserRepoGorm(controlDB)
	accountRepo := repo.NewAccountRepoGorm(controlDB)
	accountEmailRepo := repo.NewAccountEmailRepoGorm(controlDB)
	tokenSessions := service.NewTokenSessionService(tokenSvc, sessionRepo, auditRepo, &userRepo)
	// Use SMTP when configured, otherwise log tokens (dev).
	var emailSender service.EmailSender = service.LogSender{}
	if strings.TrimSpace(cfg.SMTPServer) != "" {
		emailSender = service.NewSMTPSender(cfg.SMTPServer, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom, cfg.BaseURL)
		log.Printf("email: SMTP transport via %s:%d", cfg.SMTPServer, cfg.SMTPPort)
	}
	authSvc := service.NewAuthService(controlDB, tokenSvc, accountRepo, accountEmailRepo, &userRepo, sessionRepo, auditRepo, jobEnqueuer, emailSender)
	fleetSvc := service.NewFleetService(tenants, service.ReservationDefaults{
		CPU: cfg.ReserveCPU, RAMGB: cfg.ReserveRAMGB, DiskGB: cfg.ReserveDiskGB,
	})

	platformPolicyRepo := repo.NewPlatformPolicyRepoGorm(controlDB)
	orgPolicyRepo := repo.NewOrgPolicyRepoGorm(controlDB)
	passkeyRepo := repo.NewUserPasskeyRepoGorm(controlDB)
	policySvc := service.NewPolicyService(platformPolicyRepo, orgPolicyRepo, &userRepo, passkeyRepo)

	apiKeyRepo := repo.NewOrgAPIKeyRepoGorm(controlDB)
	apiKeySvc := service.NewAPIKeyService(apiKeyRepo)

	orgSvc := service.NewOrgService(controlDB, tokenSvc, jobEnqueuer, emailSender)

	deps := &app.Deps{
		Config:        cfg,
		ControlDB:     controlDB,
		DBPool:        dbPool,
		River:         riverClient,
		Tenants:       tenants,
		Jobs:          jobEnqueuer,
		Auth:          authSvc,
		Fleet:         fleetSvc,
		Tokens:        tokenSvc,
		TokenSessions: tokenSessions,
		Policies:      policySvc,
		APIKeys:       apiKeySvc,
		Orgs:          orgSvc,
	}

	apiServer, err := api.Build(*cfg, deps)
	if err != nil {
		return fmt.Errorf("api build: %w", err)
	}

	srv := &http.Server{
		Addr:           httpAddr,
		Handler:        apiServer.Router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	errCh := make(chan error, 1)
	go func() {
		urls := apiServer.PublicURLs()
		log.Printf("listening on %s (frontend_mode=%s)", httpAddr, cfg.FrontendMode)
		log.Printf("web: %s", urls["web"])
		log.Printf("docs: %s", urls["docs"])
		log.Printf("openapi: %s", urls["openapi"])

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case sig := <-stop:
		log.Printf("signal received: %s", sig)
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("listen: %w", err)
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("shutting down...")
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	if err := <-errCh; err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
