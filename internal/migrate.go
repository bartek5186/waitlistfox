package internal

import (
	"fmt"

	"github.com/bartek5186/waitlistfox/models"
	"gorm.io/gorm"
)

func MigrateRun(db *gorm.DB) error {
	if err := db.Transaction(func(tx *gorm.DB) error {
		tx = tx.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")

		if err := tx.AutoMigrate(
			&models.WaitlistSubscriber{},
			&models.WaitlistUnsubscribe{},
		); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}
