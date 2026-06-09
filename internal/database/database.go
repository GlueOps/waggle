package database

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	"github.com/glueops/waggle/internal/models/control"
	"github.com/glueops/waggle/internal/utils"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrTenantNotProvisioned = errors.New("tenant database not provisioned")
	ErrTenantNotActive      = errors.New("tenant not active")
)

type TenantManager struct {
	ControlDB *gorm.DB
	MasterKey []byte
	tenants   sync.Map
}

func OpenControlDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
}

func NewTenantManager(controlDB *gorm.DB, masterKeyB64 string) (*TenantManager, error) {
	if controlDB == nil {
		return nil, fmt.Errorf("control db is nil")
	}

	mk, err := base64.StdEncoding.DecodeString(masterKeyB64)
	if err != nil || len(mk) != 32 {
		return nil, fmt.Errorf("invalid master key: must be 32 bytes base64 encoded")
	}

	return &TenantManager{
		ControlDB: controlDB,
		MasterKey: mk,
	}, nil
}

// For resolves a tenant DB connection for the given organization, caching the
// result for subsequent calls. The Organization.ConnectionString is read as
// plaintext (see envelope-encryption design: the connection string is
// intentionally not encrypted so the migration runner can use it without the
// master key).
func (tm *TenantManager) For(ctx context.Context, orgID uuid.UUID) (*gorm.DB, error) {
	key := orgID.String()

	if cached, ok := tm.tenants.Load(key); ok {
		return cached.(*gorm.DB), nil
	}

	var org control.Organization
	if err := tm.ControlDB.WithContext(ctx).First(&org, "id = ?", orgID).Error; err != nil {
		return nil, fmt.Errorf("load organization %s: %w", orgID, err)
	}

	if org.ConnectionString == "" {
		return nil, ErrTenantNotProvisioned
	}

	if org.Status != "active" {
		return nil, fmt.Errorf("%w: org %s status=%q", ErrTenantNotActive, orgID, org.Status)
	}

	db, err := gorm.Open(postgres.Open(org.ConnectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open tenant db for org %s: %w", orgID, err)
	}

	if existing, loaded := tm.tenants.LoadOrStore(key, db); loaded {
		if sqlDB, e := db.DB(); e == nil {
			_ = sqlDB.Close()
		}
		return existing.(*gorm.DB), nil
	}

	return db, nil
}

// TenantDEK unwraps and returns an organization's 32-byte data encryption key.
// The DEK is stored wrapped (envelope encryption): the master key (KEK)
// decrypts Organization.EncryptedTenantKey. Callers use the returned DEK to
// encrypt/decrypt tenant field-level secrets (e.g. a datacenter's PVE token).
// The DEK is never persisted unwrapped and should not be cached by callers.
func (tm *TenantManager) TenantDEK(ctx context.Context, orgID uuid.UUID) ([]byte, error) {
	if tm == nil || len(tm.MasterKey) != 32 {
		return nil, errors.New("tenant manager: master key not configured")
	}
	var org control.Organization
	if err := tm.ControlDB.WithContext(ctx).First(&org, "id = ?", orgID).Error; err != nil {
		return nil, fmt.Errorf("load organization %s: %w", orgID, err)
	}
	if org.EncryptedTenantKey == "" || org.TenantKeyIV == "" || org.TenantKeyTag == "" {
		return nil, ErrTenantNotProvisioned
	}
	ct, err := base64.StdEncoding.DecodeString(org.EncryptedTenantKey)
	if err != nil {
		return nil, fmt.Errorf("decode tenant key: %w", err)
	}
	iv, err := base64.StdEncoding.DecodeString(org.TenantKeyIV)
	if err != nil {
		return nil, fmt.Errorf("decode tenant key iv: %w", err)
	}
	tag, err := base64.StdEncoding.DecodeString(org.TenantKeyTag)
	if err != nil {
		return nil, fmt.Errorf("decode tenant key tag: %w", err)
	}
	dek, err := utils.DecryptAESGCM(ct, tm.MasterKey, iv, tag)
	if err != nil {
		return nil, fmt.Errorf("unwrap tenant dek: %w", err)
	}
	if len(dek) != 32 {
		return nil, fmt.Errorf("tenant dek has wrong length %d", len(dek))
	}
	return dek, nil
}

// Forget evicts the cached tenant DB connection for the given organization and
// best-effort closes the underlying *sql.DB. Subsequent calls to For will
// reopen the connection.
func (tm *TenantManager) Forget(orgID uuid.UUID) {
	cached, loaded := tm.tenants.LoadAndDelete(orgID.String())
	if !loaded {
		return
	}
	db, ok := cached.(*gorm.DB)
	if !ok {
		return
	}
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}
