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
	// Extensions
	if err := app.DB.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		return fmt.Errorf("failed to create vector extension: %w", err)
	}

	// Migrations
	// Note: We include ALL models here to ensure consistency across apps. GORM AutoMigrate handles column additions like SnoozedUntil.
	models := []interface{}{
		&model.Email{},
		&model.User{},
		&model.Contact{},
		&model.EmailAccount{},
		&model.EmailEmbedding{},
		&model.Organization{},
		&model.OrganizationMember{},
		&model.Team{},
		&model.TeamMember{},
		&model.Task{},
		&model.Context{},
		&model.EmailContext{},
	}

	if err := app.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	// Indices
	if err := app.DB.Exec("CREATE INDEX IF NOT EXISTS email_embeddings_vector_idx ON email_embeddings USING hnsw (vector vector_cosine_ops)").Error; err != nil {
		app.Logger.Warn("Failed to create HNSW index", logger.Error(err))
	}

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
