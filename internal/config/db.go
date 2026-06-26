package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase opens a PostgreSQL connection using GORM.
func ConnectDatabase(cfg *Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	// Configure connection pooling for production
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get underlying sql.DB: %v", err))
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	fmt.Println("✅ Database connected successfully")
	return db
}
