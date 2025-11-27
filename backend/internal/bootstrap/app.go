package bootstrap

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/config"
	"github.com/hrygo/echomind/pkg/database"
	"github.com/hrygo/echomind/pkg/logger"
	"github.com/hrygo/echomind/pkg/telemetry"
	"gorm.io/gorm"
)

type App struct {
	Config      *configs.Config
	DB          *gorm.DB
	Logger      logger.Logger
	AsynqClient *asynq.Client
	Telemetry   *telemetry.Telemetry
}

func Init(configPath string, production bool) (*App, error) {
	// 1. Config
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	// 2. Logger - 使用新的日志框架
	var logConfig *logger.Config

	// 尝试从YAML文件加载配置
	if configPath != "" {
		loggerConfigPath := filepath.Join(filepath.Dir(configPath), "logger.yaml")
		if _, err := os.Stat(loggerConfigPath); err == nil {
			logConfig, err = logger.LoadConfigFromFile(loggerConfigPath)
			if err != nil {
				fmt.Printf("Warning: Failed to load logger config from %s: %v, using default config\n", loggerConfigPath, err)
				if production {
					logConfig = logger.ProductionConfig()
				} else {
					logConfig = logger.DevelopmentConfig()
				}
			}
		} else {
			// YAML文件不存在，使用默认配置
			if production {
				logConfig = logger.ProductionConfig()
			} else {
				logConfig = logger.DevelopmentConfig()
			}
		}
	} else {
		// 未指定配置路径，使用默认配置
		if production {
			logConfig = logger.ProductionConfig()
		} else {
			logConfig = logger.DevelopmentConfig()
		}
	}

	// 从环境变量加载配置（覆盖YAML配置）
	logConfig = logger.LoadConfigFromEnv()

	// 确保生产模式设置正确
	if production {
		logConfig.Production = true
	}

	if err := logger.Init(logConfig); err != nil {
		return nil, fmt.Errorf("logger init failed: %w", err)
	}
	log := logger.GetDefaultLogger()

	// 3. Database
	if cfg.Database.DSN == "" {
		return nil, fmt.Errorf("database DSN not found")
	}
	db, err := database.New(cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	// 4. Redis/Asynq
	var asynqClient *asynq.Client
	if cfg.Redis.Addr != "" {
		asynqClient = asynq.NewClient(asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
	}

	// 5. Telemetry (OpenTelemetry)
	var tel *telemetry.Telemetry
	if cfg.Telemetry.Enabled {
		telCfg := &telemetry.TelemetryConfig{
			ServiceName:     cfg.Telemetry.ServiceName,
			ServiceVersion:  cfg.Telemetry.ServiceVersion,
			Environment:     cfg.Telemetry.Environment,
			ExporterType:    cfg.Telemetry.Exporter.Type,
			TracesFilePath:  cfg.Telemetry.Exporter.File.TracesPath,
			MetricsFilePath: cfg.Telemetry.Exporter.File.MetricsPath,
			OTLPEndpoint:    cfg.Telemetry.Exporter.OTLP.Endpoint,
			OTLPInsecure:    cfg.Telemetry.Exporter.OTLP.Insecure,
			SamplingType:    cfg.Telemetry.Sampling.Type,
			SamplingRatio:   cfg.Telemetry.Sampling.Ratio,
		}

		tel, err = telemetry.InitTelemetry(context.Background(), telCfg)
		if err != nil {
			log.Warn("Failed to initialize telemetry, continuing without it",
				logger.Error(err))
		} else {
			log.Info("OpenTelemetry initialized",
				logger.String("service", telCfg.ServiceName),
				logger.String("version", telCfg.ServiceVersion),
				logger.String("exporter", telCfg.ExporterType))
		}
	}

	app := &App{
		Config:      cfg,
		DB:          db,
		Logger:      log,
		AsynqClient: asynqClient,
		Telemetry:   tel,
	}

	return app, nil
}

func (app *App) SetupDB() error {
	// Step 1: Create Extensions
	app.Logger.Info("Creating PostgreSQL extensions...")
	if err := app.DB.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		return fmt.Errorf("failed to create vector extension: %w", err)
	}
	if err := app.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		app.Logger.Warn("Failed to create uuid-ossp extension (may already exist or using gen_random_uuid)", logger.Error(err))
	}
	app.Logger.Info("Extensions created successfully")

	// Step 2: Create Custom Types (Enums)
	app.Logger.Info("Creating custom types...")
	customTypes := []string{
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'opportunity_type') THEN
				CREATE TYPE opportunity_type AS ENUM ('buying', 'partnership', 'renewal', 'strategic');
			END IF;
		END $$;`,
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'opportunity_status') THEN
				CREATE TYPE opportunity_status AS ENUM ('new', 'active', 'won', 'lost', 'on_hold');
			END IF;
		END $$;`,
	}
	for _, sql := range customTypes {
		if err := app.DB.Exec(sql).Error; err != nil {
			app.Logger.Warn("Failed to create custom type (may already exist)", logger.Error(err))
		}
	}
	app.Logger.Info("Custom types created successfully")

	// Step 3: Run Migrations for ALL Models
	// Note: We include ALL models here to ensure consistency across apps. GORM AutoMigrate handles column additions like SnoozedUntil.
	app.Logger.Info("Running database migrations...")
	models := []interface{}{
		// Core entities
		&model.User{},
		&model.Organization{},
		&model.OrganizationMember{},
		&model.Team{},
		&model.TeamMember{},
		// Email entities
		&model.Email{},
		&model.EmailAccount{},
		&model.EmailEmbedding{},
		// Context and relationship entities
		&model.Contact{},
		&model.Context{},
		&model.EmailContext{},
		&model.Task{},
		// Opportunity entities
		&model.Opportunity{},
		&model.OpportunityContact{},
		&model.Activity{},
	}

	if err := app.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}
	app.Logger.Info("Database migrations completed successfully")

	// Step 4: Create Indices
	app.Logger.Info("Creating database indices...")
	indices := []struct {
		name string
		sql  string
	}{
		{
			name: "email_embeddings_vector_idx",
			sql:  "CREATE INDEX IF NOT EXISTS email_embeddings_vector_idx ON email_embeddings USING hnsw (vector vector_cosine_ops)",
		},
		{
			name: "idx_emails_user_date",
			sql:  "CREATE INDEX IF NOT EXISTS idx_emails_user_date ON emails (user_id, date DESC)",
		},
		{
			name: "idx_opportunities_user_status",
			sql:  "CREATE INDEX IF NOT EXISTS idx_opportunities_user_status ON opportunities (user_id, status)",
		},
		{
			name: "idx_tasks_user_status",
			sql:  "CREATE INDEX IF NOT EXISTS idx_tasks_user_status ON tasks (user_id, status)",
		},
	}

	for _, idx := range indices {
		if err := app.DB.Exec(idx.sql).Error; err != nil {
			app.Logger.Warn("Failed to create index",
				logger.String("index", idx.name),
				logger.Error(err))
		} else {
			app.Logger.Info("Index created", logger.String("index", idx.name))
		}
	}

	// Step 5: Verify Tables
	app.Logger.Info("Verifying database tables...")
	var tableCount int64
	if err := app.DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'").Scan(&tableCount).Error; err != nil {
		return fmt.Errorf("failed to verify tables: %w", err)
	}
	app.Logger.Info("Database setup completed",
		logger.Int64("total_tables", tableCount))

	return nil
}

func (app *App) Close() {
	if app.AsynqClient != nil {
		app.AsynqClient.Close()
	}
	// Shutdown telemetry
	if app.Telemetry != nil {
		if err := app.Telemetry.Shutdown(context.Background()); err != nil {
			app.Logger.Error("Failed to shutdown telemetry", logger.Error(err))
		}
	}
	// 新日志框架会自动清理
	_ = logger.Close()
}
