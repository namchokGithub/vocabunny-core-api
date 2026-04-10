package protocol

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/namchokGithub/vocabunny-core-api/configs"
	"github.com/namchokGithub/vocabunny-core-api/infrastructure"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/service"
	"github.com/namchokGithub/vocabunny-core-api/internal/handler"
	appmiddleware "github.com/namchokGithub/vocabunny-core-api/internal/middleware"
	"github.com/namchokGithub/vocabunny-core-api/internal/repository"
	"github.com/namchokGithub/vocabunny-core-api/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type App struct {
	Config     configs.Config
	Echo       *echo.Echo
	Logger     *slog.Logger
	DB         *gorm.DB
	Redis      *redis.Client
	Cron       *cron.Cron
	Handlers   *handler.Handler
	Middleware *appmiddleware.Middleware
	Services   *service.Service
	Repo       *repository.Repository
	JWT        *infrastructure.JWTManager
	Websocket  *infrastructure.WebsocketManager
	shutdown   func(context.Context) error
}

func Initialize(ctx context.Context, cfg configs.Config) (*App, error) {
	appLogger := logger.New()

	db, err := infrastructure.NewDatabase(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	redisClient := infrastructure.NewRedis(cfg.Redis)
	if err := infrastructure.PingRedis(ctx, redisClient); err != nil {
		appLogger.Warn("redis unavailable", slog.String("error", err.Error()))
	}

	storage, err := infrastructure.NewFileStorage(cfg.Storage)
	if err != nil {
		return nil, fmt.Errorf("init storage: %w", err)
	}

	jwtManager := infrastructure.NewJWTManager(cfg.JWT)

	repositories := repository.NewRepository(repository.Dependencies{
		DB: db.Gorm,
	})

	services := service.NewService(service.Dependencies{
		Repositories: &service.RepositoryPorts{
			User:           repositories.User,
			Role:           repositories.Role,
			AuthIdentity:   repositories.AuthIdentity,
			Section:        repositories.Section,
			Lesson:         repositories.Lesson,
			Unit:           repositories.Unit,
			QuestionSet:    repositories.QuestionSet,
			Question:       repositories.Question,
			QuestionChoice: repositories.QuestionChoice,
			Tag:            repositories.Tag,
			MediaAsset:     repositories.MediaAsset,
		},
		TxManager:    infrastructure.NewTransactionManager(db.Gorm),
		Storage:      storage,
		TokenManager: jwtManager,
	})

	validator := helper.NewRequestValidator()
	handlers := handler.NewHandler(handler.Dependencies{
		Services:  services,
		Validator: validator,
	})
	middlewares := appmiddleware.New(appmiddleware.Dependencies{
		JWTManager:  jwtManager,
		UserService: services.User,
	})

	e := echo.New()
	e.Validator = validator
	e.HideBanner = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		_ = helper.RespondError(c, err)
	}
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	e.Use(echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Timeout: cfg.HTTP.WriteTimeout,
	}))

	cronScheduler := infrastructure.NewCron()
	cronScheduler.Start()

	websocketManager := infrastructure.NewWebsocketManager()

	app := &App{
		Config:     cfg,
		Echo:       e,
		Logger:     appLogger,
		DB:         db.Gorm,
		Redis:      redisClient,
		Cron:       cronScheduler,
		Handlers:   handlers,
		Middleware: middlewares,
		Services:   services,
		Repo:       repositories,
		JWT:        jwtManager,
		Websocket:  websocketManager,
	}

	RegisterHTTP(app)

	app.shutdown = func(shutdownCtx context.Context) error {
		cronScheduler.Stop()

		if redisClient != nil {
			if err := redisClient.Close(); err != nil {
				appLogger.Warn("close redis", slog.String("error", err.Error()))
			}
		}

		sqlDB, err := db.Gorm.DB()
		if err == nil {
			if closeErr := sqlDB.Close(); closeErr != nil {
				appLogger.Warn("close database", slog.String("error", closeErr.Error()))
			}
		}

		return nil
	}

	return app, nil
}

func (a *App) Start() error {
	server := &http.Server{
		Addr:         ":" + a.Config.HTTP.Port,
		ReadTimeout:  a.Config.HTTP.ReadTimeout,
		WriteTimeout: a.Config.HTTP.WriteTimeout,
	}

	a.Echo.Server = server
	return a.Echo.StartServer(server)
}

func (a *App) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, a.Config.HTTP.ShutdownTimeout)
	defer cancel()

	if err := a.Echo.Shutdown(shutdownCtx); err != nil {
		return err
	}

	if a.shutdown != nil {
		if err := a.shutdown(shutdownCtx); err != nil {
			return err
		}
	}

	return nil
}
