package app

import (
	"flag"
	"os"
)

// CLIConfig holds configuration from CLI flags and environment variables
type CLIConfig struct {
	ConfigPath   string
	IsProduction bool
}

// ParseCLI parses command-line flags and environment variables
// Priority: CLI flags > Environment variables > Default values
func ParseCLI() *CLIConfig {
	cfg := &CLIConfig{}

	// Define flags with environment variable fallback
	flag.StringVar(
		&cfg.ConfigPath,
		"config",
		getEnv("CONFIG_PATH", "configs/config.yaml"),
		"Path to configuration file (env: CONFIG_PATH)",
	)

	flag.BoolVar(
		&cfg.IsProduction,
		"production",
		getEnv("PRODUCTION", "false") == "true",
		"Run in production mode (env: PRODUCTION)",
	)

	flag.Parse()
	return cfg
}

// getEnv retrieves an environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
