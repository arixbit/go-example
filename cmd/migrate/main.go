package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"go-skeleton/config"
	"go-skeleton/internal/model"
	"go-skeleton/pkg/database"
	applog "go-skeleton/pkg/log"
)

func main() {
	config.LoadEnv("cmd/migrate/.env")
	cfg := config.Load()
	if err := config.InitRuntime(cfg, "migrate"); err != nil {
		panic(fmt.Sprintf("init runtime: %v", err))
	}
	defer func() { _ = applog.Sync() }()

	dbMgr, err := database.Init(database.Config{
		DSN:             cfg.Postgres.DSN,
		LogLevel:        cfg.Postgres.LogLevel,
		MaxIdleConns:    cfg.Postgres.MaxIdleConns,
		MaxOpenConns:    cfg.Postgres.MaxOpenConns,
		ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Postgres.ConnMaxIdleTime,
	})
	if err != nil {
		applog.L().Fatal("initialize database", zap.Error(err))
	}
	defer func() {
		if err := dbMgr.Close(); err != nil {
			applog.L().Warn("close database", zap.Error(err))
		}
	}()
	if dbMgr.DB() == nil {
		applog.L().Fatal("database is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := dbMgr.DB().WithContext(ctx).AutoMigrate(&model.Example{}); err != nil {
		applog.L().Fatal("run migrations", zap.Error(err))
	}
	applog.L().Info("migrations completed")
}
