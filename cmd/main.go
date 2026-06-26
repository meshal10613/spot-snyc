package main

import (
	"fmt"
	"log"
	"net/http"
	"spot-sync/config"
	"spot-sync/models"
	"spot-sync/server"
)

func main() {
	// ── Load Configuration ────────────────────────────────────────────────
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	// ── Connect Database ──────────────────────────────────────────────────
	db := config.ConnectDatabase(cfg)

	// ── Run Migrations ────────────────────────────────────────────────────
	fmt.Println("⏳ Running database migrations...")
	if err := models.Migrate(db); err != nil {
		log.Fatal("❌ Migration failed: ", err)
	}
	fmt.Println("✅ Migrations completed")

	// ── Start Server ──────────────────────────────────────────────────────
	e := server.NewHTTPServer(cfg, db)
	fmt.Printf("🚀 Server running at http://localhost:%s\n", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		log.Fatal("❌ Failed to start server: ", err)
	}
}
