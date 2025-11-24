package main

import (
	"log"

	"github.com/hrygo/echomind/internal/app"
	"github.com/hrygo/echomind/pkg/logger"
)

func main() {
	// Parse CLI configuration to get config path
	cli := app.ParseCLI()

	// Initialize application container (connects to DB)
	container, err := app.NewContainer(cli.ConfigPath, cli.IsProduction)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer container.Close()

	container.Logger.Info("Starting database initialization...")

	// Setup Database (Migrations & Extensions)
	if err := container.SetupDB(); err != nil {
		container.Logger.Fatal("Database setup failed", logger.Error(err))
	}

	container.Logger.Info("Database initialization completed successfully.")
}
