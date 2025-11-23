package main

import (
	"log"

	"github.com/hrygo/echomind/internal/app"
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

	container.Sugar.Info("Starting database initialization...")

	// Setup Database (Migrations & Extensions)
	if err := container.SetupDB(); err != nil {
		container.Sugar.Fatalf("Database setup failed: %v", err)
	}

	container.Sugar.Info("Database initialization completed successfully.")
}
