package seed

import (
	"fmt"
	"spot-sync/internal/config"
	"spot-sync/internal/models"

	"gorm.io/gorm"
)

// AdminSeeder checks if the admin user exists. If not, it creates one
// using the credentials from environment variables. If the admin already
// exists, it silently skips.
func AdminSeeder(db *gorm.DB, cfg *config.Config) error {
	var existing models.User
	result := db.Where("email = ?", cfg.AdminEmail).First(&existing)

	// Admin already exists — skip
	if result.Error == nil {
		fmt.Printf("✅ Admin already exists (%s) — skipping seed\n", cfg.AdminEmail)
		return nil
	}

	// Unexpected DB error (not "record not found")
	if result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check admin user: %w", result.Error)
	}

	// Create admin user
	admin := &models.User{
		Name:  cfg.AdminName,
		Email: cfg.AdminEmail,
		Role:  models.RoleAdmin,
	}

	if err := admin.HashPassword(cfg.AdminPassword); err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}

	fmt.Printf("✅ Admin seeded successfully (%s)\n", cfg.AdminEmail)
	return nil
}
