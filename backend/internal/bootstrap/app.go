package bootstrap

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/config"
	"github.com/hrygo/echomind/pkg/database"
	"github.com/hrygo/echomind/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Config      *configs.Config
	DB          *gorm.DB
	Logger      *zap.Logger
	AsynqClient *asynq.Client
	Sugar       *zap.SugaredLogger
}

func Init(configPath string, production bool) (*App, error) {
	// 1. Config
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	// 2. Logger
	log, err := logger.New(production)
	if err != nil {
		return nil, fmt.Errorf("logger init failed: %w", err)
	}
	sugar := log.Sugar()

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

	app := &App{
		Config:      cfg,
		DB:          db,
		Logger:      log,
		Sugar:       sugar,
		AsynqClient: asynqClient,
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
		app.Sugar.Warnf("Failed to create HNSW index: %v", err)
	}

	return nil
}

func (app *App) Close() {
	if app.AsynqClient != nil {
		app.AsynqClient.Close()
	}
	if app.Logger != nil {
		app.Logger.Sync()
	}
}
