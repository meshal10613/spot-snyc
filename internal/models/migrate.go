package models

import "gorm.io/gorm"

// Migrate runs auto-migration for all database models.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&ParkingZone{},
		&Reservation{},
	)
}
