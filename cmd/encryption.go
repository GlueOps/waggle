package cmd

import (
	"fmt"
	"log"

	"github.com/glueops/waggle/internal/config"
	"github.com/glueops/waggle/internal/models/control"
	"github.com/glueops/waggle/internal/utils"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Manage encryption keys",
}

var generateMasterCmd = &cobra.Command{
	Use:   "generate-master",
	Short: "Generate a new random 32-byte master key (for initial setup)",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := utils.RandomBytes(32)
		if err != nil {
			return err
		}
		fmt.Printf("ENCRYPTION_MASTER_KEY=%s\n", utils.EncodeB64(key))
		return nil
	},
}

var rotateMasterCmd = &cobra.Command{
	Use:   "rotate-master",
	Short: "Rotate the master encryption key (KEK)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		oldMasterKey, err := utils.DecodeB64(cfg.EncryptionMasterKey)
		if err != nil || len(oldMasterKey) != 32 {
			return fmt.Errorf("invalid current master key in config (must be 32 bytes base64)")
		}

		log.Println("Generating new Master Key...")
		newMasterKey, err := utils.RandomBytes(32)
		if err != nil {
			return err
		}

		controlDB, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
		if err != nil {
			return err
		}

		err = controlDB.Transaction(func(tx *gorm.DB) error {
			var orgs []control.Organization
			if err := tx.Find(&orgs).Error; err != nil {
				return err
			}

			log.Printf("Found %d organizations. Rotating wrapped keys...", len(orgs))

			for _, org := range orgs {
				if org.EncryptedTenantKey == "" {
					continue // Skip if this org somehow has no key
				}

				ct, _ := utils.DecodeB64(org.EncryptedTenantKey)
				iv, _ := utils.DecodeB64(org.TenantKeyIV)
				tag, _ := utils.DecodeB64(org.TenantKeyTag)

				// A. Decrypt the Tenant Key using the OLD Master Key
				tenantKey, err := utils.DecryptAESGCM(ct, oldMasterKey, iv, tag)
				if err != nil {
					return fmt.Errorf("failed to decrypt tenant key for org %s: %w", org.ID, err)
				}

				// B. Re-encrypt the Tenant Key using the NEW Master Key
				newCt, newIv, newTag, err := utils.EncryptAESGCM(tenantKey, newMasterKey)
				if err != nil {
					return fmt.Errorf("failed to encrypt tenant key for org %s: %w", org.ID, err)
				}

				// C. Save the updated fields back to the Org record
				if err := tx.Model(&org).Updates(map[string]interface{}{
					"encrypted_tenant_key": utils.EncodeB64(newCt),
					"tenant_key_iv":        utils.EncodeB64(newIv),
					"tenant_key_tag":       utils.EncodeB64(newTag),
				}).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("rotation failed: %w", err)
		}

		fmt.Println("\n✅ Master key rotated successfully in the database.")
		fmt.Println("⚠️  IMPORTANT: Update your environment configuration IMMEDIATELY.")
		fmt.Println("Replace your current ENCRYPTION_MASTER_KEY in your .env with this new key:")
		fmt.Printf("ENCRYPTION_MASTER_KEY=%s\n\n", utils.EncodeB64(newMasterKey))
		fmt.Println("Then restart your API and Worker processes.")

		return nil
	},
}

func init() {
	encryptCmd.AddCommand(generateMasterCmd, rotateMasterCmd)
	rootCmd.AddCommand(encryptCmd)
}
