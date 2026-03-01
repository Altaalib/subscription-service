package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"subscription-service/config"
	"subscription-service/internal/handler"
	"subscription-service/internal/repository"
	"subscription-service/internal/service"
	"subscription-service/pkg/logger"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Логгер
	zapLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatal("failed to init logger", err)
	}
	defer zapLogger.Sync()

	// Конфиг
	cfg, err := config.Load()
	if err != nil {
		zapLogger.Fatal("failed to load config", zap.Error(err))
	}

	// База данных
	db, err := sqlx.Connect("pgx", cfg.Database.DSN())
	if err != nil {
		zapLogger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()
	zapLogger.Info("connected to database")

	// Слои
	repo := repository.NewSubscriptionRepository(db)
	svc := service.NewSubscriptionService(repo)
	h := handler.NewSubscriptionHandler(svc)

	// Echo
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Роуты
	api := e.Group("/api/v1")
	api.POST("/subscriptions", h.Create)
	api.GET("/subscriptions", h.GetAll)
	api.GET("/subscriptions/:id", h.GetByID)
	api.PUT("/subscriptions/:id", h.Update)
	api.DELETE("/subscriptions/:id", h.Delete)
	api.GET("/subscriptions/total", h.GetTotalCost)

	// Запуск
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		zapLogger.Info("starting server", zap.String("addr", addr))
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		zapLogger.Fatal("shutdown error", zap.Error(err))
	}
	zapLogger.Info("server stopped")
}
